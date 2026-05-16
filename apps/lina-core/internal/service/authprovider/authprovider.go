// Package authprovider implements host-governed external authentication
// provider metadata, source-plugin provider registration, and identity binding.
package authprovider

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"lina-core/internal/dao"
	"lina-core/internal/model/do"
	"lina-core/internal/model/entity"
	configsvc "lina-core/internal/service/config"
	"lina-core/internal/service/kvcache"
	pluginsvc "lina-core/internal/service/plugin"
	"lina-core/pkg/bizerr"
	"lina-core/pkg/pluginhost"
)

// Provider enabled status values.
const (
	// ProviderStatusDisabled marks a provider as hidden and unusable.
	ProviderStatusDisabled = 0
	// ProviderStatusEnabled marks a provider as visible and usable.
	ProviderStatusEnabled = 1
)

// ProviderConfigModeStatic references config.yaml or environment-backed settings.
const ProviderConfigModeStatic = "static_config"

const (
	authProviderCacheOwner          = "auth-provider"
	authProviderStateNamespace      = "oauth-state"
	authProviderStateConsumeNS      = "oauth-state-consume"
	authProviderStateSchema         = 1
	authProviderStateRandomBytes    = 32
	authProviderPurposeLogin        = "login"
	authProviderPurposeBind         = "bind"
	authProviderPKCEChallengeMethod = "S256"
	authProviderUserStatusDisabled  = 0
)

// Service defines pluggable authentication provider operations.
type Service interface {
	// ListPublicProviders returns enabled providers that can be shown before login.
	ListPublicProviders(ctx context.Context) ([]ProviderItem, error)
	// Authorize creates a temporary authorization state and returns a redirect target.
	Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error)
	// Callback consumes an authorization state and returns a normalized external identity.
	Callback(ctx context.Context, input CallbackInput) (*CallbackOutput, error)
	// ListCurrentUserIdentities returns external identity bindings for one user.
	ListCurrentUserIdentities(ctx context.Context, userID int) ([]IdentityItem, error)
	// BindCurrentUserIdentity binds an external identity to the current user.
	BindCurrentUserIdentity(ctx context.Context, userID int, providerKey string, identity *pluginhost.ExternalIdentity) error
	// UnbindCurrentUserProvider removes the current user's binding for one provider.
	UnbindCurrentUserProvider(ctx context.Context, userID int, providerKey string) error
	// ResolveRuntimeProvider returns enabled metadata plus the plugin handler.
	ResolveRuntimeProvider(ctx context.Context, providerKey string) (*RuntimeProvider, error)
	// InvalidateProviderCache invalidates process-local provider snapshots.
	InvalidateProviderCache(ctx context.Context, scope string) error
}

// ProviderItem is a public-safe provider projection.
type ProviderItem struct {
	ProviderKey  string // ProviderKey is the stable provider instance key.
	ProviderType string // ProviderType identifies the protocol family.
	Name         string // Name is the display name.
	Icon         string // Icon is the icon name or URL.
	Sort         int    // Sort is the display order.
}

// AuthorizeInput describes one host authorization request.
type AuthorizeInput struct {
	ProviderKey string // ProviderKey is the stable provider instance key.
	Purpose     string // Purpose identifies login or bind flow.
	RedirectURI string // RedirectURI is the frontend return URI after callback.
	UserID      int    // UserID is captured for current-user binding flows.
}

// AuthorizeOutput contains the browser redirect target and host state.
type AuthorizeOutput struct {
	RedirectURL string // RedirectURL is the external provider URL.
	State       string // State is the host-generated CSRF state value.
}

// CallbackInput describes one provider callback request.
type CallbackInput struct {
	ProviderKey string            // ProviderKey is the stable provider instance key.
	State       string            // State is the host-generated CSRF state value.
	Query       map[string]string // Query contains GET callback parameters.
	Form        map[string]string // Form contains POST callback parameters.
}

// CallbackOutput contains the normalized identity returned by a provider plugin.
type CallbackOutput struct {
	ProviderKey  string                       // ProviderKey is the stable provider instance key.
	ProviderType string                       // ProviderType identifies the provider family.
	Purpose      string                       // Purpose identifies login or bind flow.
	RedirectURI  string                       // RedirectURI is the frontend return URI after callback.
	UserID       int                          // UserID is the captured user for binding flows.
	Identity     *pluginhost.ExternalIdentity // Identity is the normalized external identity.
}

// IdentityItem is a current-user external identity projection.
type IdentityItem struct {
	ProviderKey      string      // ProviderKey is the bound provider instance key.
	ProviderType     string      // ProviderType identifies the protocol family.
	Subject          string      // Subject is the provider subject identifier.
	ExternalTenantID string      // ExternalTenantID is the provider tenant/corp identifier.
	Email            string      // Email is the external email address.
	Mobile           string      // Mobile is the external mobile phone number.
	DisplayName      string      // DisplayName is the external display name.
	Avatar           string      // Avatar is the external avatar URL.
	LastLoginAt      *gtime.Time // LastLoginAt records the last successful external login time.
	BoundAt          *gtime.Time // BoundAt records when the binding was created.
}

// RuntimeProvider combines host metadata with a plugin-owned handler.
type RuntimeProvider struct {
	Metadata ProviderItem            // Metadata is the host-governed provider projection.
	Handler  pluginhost.AuthProvider // Handler is the plugin-owned protocol adapter.
	Record   *entity.SysAuthProvider // Record is the database metadata row.
}

// pluginService is the narrow plugin facade required by this service.
type pluginService interface {
	IsEnabled(ctx context.Context, pluginID string) bool
	ListAuthProviders(ctx context.Context) ([]pluginsvc.AuthProviderRegistration, error)
}

// authProviderConfigReader is the narrow static config surface used here.
type authProviderConfigReader interface {
	GetAuthProviders(ctx context.Context) *configsvc.AuthProvidersConfig
}

// serviceImpl implements Service.
type serviceImpl struct {
	configSvc  authProviderConfigReader
	pluginSvc  pluginService
	kvCacheSvc kvcache.Service
}

// New creates an auth provider service.
func New(configSvc authProviderConfigReader, pluginSvc pluginService, kvCacheSvc kvcache.Service) Service {
	return &serviceImpl{
		configSvc:  configSvc,
		pluginSvc:  pluginSvc,
		kvCacheSvc: kvCacheSvc,
	}
}

// ListPublicProviders returns enabled providers that can be shown before login.
func (s *serviceImpl) ListPublicProviders(ctx context.Context) ([]ProviderItem, error) {
	var rows []*entity.SysAuthProvider
	err := dao.SysAuthProvider.Ctx(ctx).
		Where(do.SysAuthProvider{Enabled: ProviderStatusEnabled}).
		OrderAsc(dao.SysAuthProvider.Columns().Sort).
		OrderAsc(dao.SysAuthProvider.Columns().ProviderKey).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	registrations, err := s.providerRegistrations(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]ProviderItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		registration, ok := registrations.resolve(row)
		if !ok || registration.Provider == nil {
			continue
		}
		if s.pluginSvc != nil && !s.pluginSvc.IsEnabled(ctx, row.PluginId) {
			continue
		}
		items = append(items, providerItemFromRow(row, registration.Provider))
	}
	return items, nil
}

// Authorize creates a temporary authorization state and returns a redirect target.
func (s *serviceImpl) Authorize(ctx context.Context, input AuthorizeInput) (*AuthorizeOutput, error) {
	runtimeProvider, err := s.ResolveRuntimeProvider(ctx, input.ProviderKey)
	if err != nil {
		return nil, err
	}
	state, err := randomURLToken(authProviderStateRandomBytes)
	if err != nil {
		return nil, err
	}
	nonce, err := randomURLToken(authProviderStateRandomBytes)
	if err != nil {
		return nil, err
	}
	pkceVerifier, err := randomURLToken(authProviderStateRandomBytes)
	if err != nil {
		return nil, err
	}
	staticConfig := s.staticProviderConfig(ctx, runtimeProvider.Metadata.ProviderKey)
	record := authStateRecord{
		Schema:              authProviderStateSchema,
		ProviderKey:         runtimeProvider.Metadata.ProviderKey,
		ProviderType:        runtimeProvider.Metadata.ProviderType,
		State:               state,
		Nonce:               nonce,
		PKCEVerifier:        pkceVerifier,
		Purpose:             normalizePurpose(input.Purpose),
		RedirectURI:         strings.TrimSpace(input.RedirectURI),
		UserID:              input.UserID,
		ProviderRedirectURL: staticConfig.RedirectURL,
		CreatedAt:           time.Now(),
		ExpiresAt:           time.Now().Add(s.stateTTL(ctx)),
	}
	if err = s.storeState(ctx, record); err != nil {
		return nil, err
	}
	output, err := runtimeProvider.Handler.BuildAuthorizeURL(ctx, pluginhost.AuthProviderAuthorizeInput{
		ProviderKey:         record.ProviderKey,
		State:               record.State,
		Nonce:               record.Nonce,
		PKCEChallenge:       pkceChallenge(record.PKCEVerifier),
		PKCEChallengeMethod: authProviderPKCEChallengeMethod,
		RedirectURL:         record.ProviderRedirectURL,
		Scopes:              staticConfig.Scopes,
		Purpose:             record.Purpose,
	})
	if err != nil {
		if deleteErr := s.deleteState(ctx, record.State); deleteErr != nil {
			return nil, deleteErr
		}
		return nil, bizerr.WrapCode(err, CodeAuthProviderUnavailable)
	}
	if output == nil || strings.TrimSpace(output.RedirectURL) == "" {
		if deleteErr := s.deleteState(ctx, record.State); deleteErr != nil {
			return nil, deleteErr
		}
		return nil, bizerr.NewCode(CodeAuthProviderUnavailable)
	}
	return &AuthorizeOutput{RedirectURL: output.RedirectURL, State: record.State}, nil
}

// Callback consumes an authorization state and returns a normalized external identity.
func (s *serviceImpl) Callback(ctx context.Context, input CallbackInput) (*CallbackOutput, error) {
	record, ok, err := s.consumeState(ctx, input.ProviderKey, input.State)
	if err != nil {
		return nil, bizerr.WrapCode(err, CodeAuthProviderUnavailable)
	}
	if !ok {
		return nil, bizerr.NewCode(CodeAuthStateInvalid)
	}
	runtimeProvider, err := s.ResolveRuntimeProvider(ctx, record.ProviderKey)
	if err != nil {
		return nil, err
	}
	identity, err := runtimeProvider.Handler.ExchangeCallback(ctx, pluginhost.AuthProviderCallbackInput{
		ProviderKey:  record.ProviderKey,
		Query:        cloneStringMap(input.Query),
		Form:         cloneStringMap(input.Form),
		RedirectURL:  record.ProviderRedirectURL,
		Nonce:        record.Nonce,
		PKCEVerifier: record.PKCEVerifier,
		Purpose:      record.Purpose,
	})
	if err != nil {
		return nil, bizerr.WrapCode(err, CodeAuthProviderUnavailable)
	}
	if identity == nil || strings.TrimSpace(identity.Subject) == "" {
		return nil, bizerr.NewCode(CodeAuthIdentityNotBound)
	}
	if identity.ProviderKey == "" {
		identity.ProviderKey = record.ProviderKey
	}
	if identity.ProviderType == "" {
		identity.ProviderType = record.ProviderType
	}
	return &CallbackOutput{
		ProviderKey:  record.ProviderKey,
		ProviderType: record.ProviderType,
		Purpose:      record.Purpose,
		RedirectURI:  record.RedirectURI,
		UserID:       record.UserID,
		Identity:     identity,
	}, nil
}

// ResolveRuntimeProvider returns enabled metadata plus the plugin handler.
func (s *serviceImpl) ResolveRuntimeProvider(ctx context.Context, providerKey string) (*RuntimeProvider, error) {
	providerKey = strings.TrimSpace(providerKey)
	if providerKey == "" {
		return nil, bizerr.NewCode(CodeAuthProviderNotFound)
	}
	var row *entity.SysAuthProvider
	err := dao.SysAuthProvider.Ctx(ctx).
		Where(do.SysAuthProvider{ProviderKey: providerKey, Enabled: ProviderStatusEnabled}).
		Scan(&row)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, bizerr.NewCode(CodeAuthProviderNotFound)
	}
	if s.pluginSvc != nil && !s.pluginSvc.IsEnabled(ctx, row.PluginId) {
		return nil, bizerr.NewCode(CodeAuthProviderUnavailable)
	}
	registrations, err := s.providerRegistrations(ctx)
	if err != nil {
		return nil, err
	}
	registration, ok := registrations.resolve(row)
	if !ok || registration.Provider == nil {
		return nil, bizerr.NewCode(CodeAuthProviderUnavailable)
	}
	return &RuntimeProvider{
		Metadata: providerItemFromRow(row, registration.Provider),
		Handler:  registration.Provider,
		Record:   row,
	}, nil
}

// ListCurrentUserIdentities returns external identity bindings for one user.
func (s *serviceImpl) ListCurrentUserIdentities(ctx context.Context, userID int) ([]IdentityItem, error) {
	if userID <= 0 {
		return nil, nil
	}
	var rows []*entity.SysAuthIdentity
	err := dao.SysAuthIdentity.Ctx(ctx).
		Where(do.SysAuthIdentity{UserId: userID}).
		OrderAsc(dao.SysAuthIdentity.Columns().ProviderKey).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	items := make([]IdentityItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		items = append(items, IdentityItem{
			ProviderKey:      row.ProviderKey,
			ProviderType:     row.ProviderType,
			Subject:          row.Subject,
			ExternalTenantID: row.ExternalTenantId,
			Email:            row.Email,
			Mobile:           row.Mobile,
			DisplayName:      row.DisplayName,
			Avatar:           row.Avatar,
			LastLoginAt:      row.LastLoginAt,
			BoundAt:          row.BoundAt,
		})
	}
	return items, nil
}

// BindCurrentUserIdentity binds an external identity to the current user.
func (s *serviceImpl) BindCurrentUserIdentity(ctx context.Context, userID int, providerKey string, identity *pluginhost.ExternalIdentity) error {
	providerKey = strings.TrimSpace(providerKey)
	if userID <= 0 || providerKey == "" || identity == nil || strings.TrimSpace(identity.Subject) == "" {
		return bizerr.NewCode(CodeAuthIdentityNotBound)
	}
	runtimeProvider, err := s.ResolveRuntimeProvider(ctx, providerKey)
	if err != nil {
		return err
	}
	var user *entity.SysUser
	err = dao.SysUser.Ctx(ctx).
		Where(do.SysUser{Id: userID}).
		Scan(&user)
	if err != nil {
		return err
	}
	if user == nil || user.Status == authProviderUserStatusDisabled {
		return bizerr.NewCode(CodeAuthIdentityNotBound)
	}
	var existing *entity.SysAuthIdentity
	err = dao.SysAuthIdentity.Ctx(ctx).
		Where(do.SysAuthIdentity{ProviderKey: providerKey, Subject: identity.Subject}).
		Scan(&existing)
	if err != nil {
		return err
	}
	if existing != nil {
		return bizerr.NewCode(CodeAuthIdentityAlreadyBound)
	}
	err = dao.SysAuthIdentity.Ctx(ctx).
		Where(do.SysAuthIdentity{UserId: userID, ProviderKey: providerKey}).
		Scan(&existing)
	if err != nil {
		return err
	}
	if existing != nil {
		return bizerr.NewCode(CodeAuthIdentityAlreadyBound)
	}
	emailVerified := 0
	if identity.EmailVerified {
		emailVerified = 1
	}
	_, err = dao.SysAuthIdentity.Ctx(ctx).Data(do.SysAuthIdentity{
		TenantId:         user.TenantId,
		UserId:           userID,
		ProviderKey:      providerKey,
		ProviderType:     runtimeProvider.Metadata.ProviderType,
		Subject:          strings.TrimSpace(identity.Subject),
		UnionId:          strings.TrimSpace(identity.UnionID),
		OpenId:           strings.TrimSpace(identity.OpenID),
		ExternalTenantId: strings.TrimSpace(identity.ExternalTenantID),
		ExternalDeptIds:  marshalJSONString(identity.ExternalDeptIDs),
		Email:            strings.TrimSpace(identity.Email),
		EmailVerified:    emailVerified,
		Mobile:           strings.TrimSpace(identity.Mobile),
		DisplayName:      strings.TrimSpace(identity.DisplayName),
		Avatar:           strings.TrimSpace(identity.Avatar),
		RawProfile:       marshalJSONString(identity.RawProfile),
		BoundAt:          gtime.Now(),
	}).Insert()
	if err != nil {
		return err
	}
	return nil
}

// UnbindCurrentUserProvider removes the current user's binding for one provider.
func (s *serviceImpl) UnbindCurrentUserProvider(ctx context.Context, userID int, providerKey string) error {
	providerKey = strings.TrimSpace(providerKey)
	if userID <= 0 || providerKey == "" {
		return bizerr.NewCode(CodeAuthIdentityNotFound)
	}
	result, err := dao.SysAuthIdentity.Ctx(ctx).
		Where(do.SysAuthIdentity{UserId: userID, ProviderKey: providerKey}).
		Delete()
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return bizerr.NewCode(CodeAuthIdentityNotFound)
	}
	return nil
}

// InvalidateProviderCache invalidates process-local provider snapshots.
func (s *serviceImpl) InvalidateProviderCache(_ context.Context, _ string) error {
	// Provider metadata is database-backed in this iteration. Runtime plugin
	// enablement caches are already invalidated by plugin lifecycle flows; the
	// authority source remains the SQL provider table plus current plugin state.
	return nil
}

// authStateRecord is the persisted single-use OAuth/OIDC state envelope.
type authStateRecord struct {
	Schema              int       `json:"schema"`
	ProviderKey         string    `json:"providerKey"`
	ProviderType        string    `json:"providerType"`
	State               string    `json:"state"`
	Nonce               string    `json:"nonce"`
	PKCEVerifier        string    `json:"pkceVerifier"`
	Purpose             string    `json:"purpose"`
	RedirectURI         string    `json:"redirectUri"`
	UserID              int       `json:"userId"`
	ProviderRedirectURL string    `json:"providerRedirectUrl"`
	CreatedAt           time.Time `json:"createdAt"`
	ExpiresAt           time.Time `json:"expiresAt"`
}

// staticProviderConfig contains non-secret provider settings from config.yaml.
type staticProviderConfig struct {
	RedirectURL string
	Scopes      []string
}

// providerRegistrationIndex resolves provider handlers by concrete key first
// and by protocol type second for multi-instance providers such as OIDC.
type providerRegistrationIndex struct {
	byKey  map[string]pluginsvc.AuthProviderRegistration
	byType map[string][]pluginsvc.AuthProviderRegistration
}

// storeState stores a temporary state in the runtime-selected kvcache backend.
func (s *serviceImpl) storeState(ctx context.Context, record authStateRecord) error {
	if s == nil || s.kvCacheSvc == nil {
		return bizerr.NewCode(CodeAuthStateInvalid)
	}
	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}
	_, err = s.kvCacheSvc.Set(ctx, kvcache.OwnerTypeModule, authStateCacheKey(record.State), string(payload), s.stateTTL(ctx))
	return err
}

// consumeState reads and deletes a state once, rejecting repeats and expired records.
func (s *serviceImpl) consumeState(ctx context.Context, providerKey string, state string) (authStateRecord, bool, error) {
	providerKey = strings.TrimSpace(providerKey)
	state = strings.TrimSpace(state)
	if s == nil || s.kvCacheSvc == nil || providerKey == "" || state == "" {
		return authStateRecord{}, false, nil
	}
	consumeItem, err := s.kvCacheSvc.Incr(ctx, kvcache.OwnerTypeModule, authStateConsumeCacheKey(state), 1, s.stateTTL(ctx))
	if err != nil {
		return authStateRecord{}, false, err
	}
	if consumeItem.IntValue != 1 {
		return authStateRecord{}, false, nil
	}
	item, ok, err := s.kvCacheSvc.Get(ctx, kvcache.OwnerTypeModule, authStateCacheKey(state))
	if err != nil || !ok {
		return authStateRecord{}, ok, err
	}
	if err = s.deleteState(ctx, state); err != nil {
		return authStateRecord{}, false, err
	}
	var record authStateRecord
	if err = json.Unmarshal([]byte(item.Value), &record); err != nil {
		return authStateRecord{}, false, err
	}
	if record.Schema != authProviderStateSchema ||
		record.State != state ||
		record.ProviderKey != providerKey ||
		time.Now().After(record.ExpiresAt) {
		return authStateRecord{}, false, nil
	}
	return record, true, nil
}

// deleteState removes a temporary state record while keeping consume markers intact.
func (s *serviceImpl) deleteState(ctx context.Context, state string) error {
	if s == nil || s.kvCacheSvc == nil {
		return nil
	}
	return s.kvCacheSvc.Delete(ctx, kvcache.OwnerTypeModule, authStateCacheKey(state))
}

// stateTTL returns the configured temporary state TTL.
func (s *serviceImpl) stateTTL(ctx context.Context) time.Duration {
	if s == nil || s.configSvc == nil {
		return configsvc.DefaultAuthProviderStateTTL
	}
	cfg := s.configSvc.GetAuthProviders(ctx)
	if cfg == nil || cfg.StateTTL <= 0 {
		return configsvc.DefaultAuthProviderStateTTL
	}
	return cfg.StateTTL
}

// staticProviderConfig returns non-secret static settings for one provider.
func (s *serviceImpl) staticProviderConfig(ctx context.Context, providerKey string) staticProviderConfig {
	if s == nil || s.configSvc == nil {
		return staticProviderConfig{}
	}
	cfg := s.configSvc.GetAuthProviders(ctx)
	if cfg == nil {
		return staticProviderConfig{}
	}
	for _, item := range cfg.OIDC {
		if item.ProviderKey == providerKey {
			return staticProviderConfig{
				RedirectURL: item.RedirectURL,
				Scopes:      append([]string(nil), item.Scopes...),
			}
		}
	}
	for _, item := range cfg.WeCom {
		if item.ProviderKey == providerKey {
			return staticProviderConfig{RedirectURL: item.RedirectURL}
		}
	}
	return staticProviderConfig{}
}

// providerRegistrations collects enabled source-plugin auth provider handlers.
func (s *serviceImpl) providerRegistrations(ctx context.Context) (*providerRegistrationIndex, error) {
	out := &providerRegistrationIndex{
		byKey:  make(map[string]pluginsvc.AuthProviderRegistration),
		byType: make(map[string][]pluginsvc.AuthProviderRegistration),
	}
	if s == nil || s.pluginSvc == nil {
		return out, nil
	}
	items, err := s.pluginSvc.ListAuthProviders(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.Provider == nil {
			continue
		}
		providerKey := strings.TrimSpace(item.Provider.ProviderKey())
		if providerKey != "" {
			out.byKey[providerKey] = item
		}
		providerType := strings.TrimSpace(item.Provider.ProviderType())
		if providerType != "" {
			out.byType[providerType] = append(out.byType[providerType], item)
		}
	}
	return out, nil
}

// resolve returns the plugin handler for one enabled provider metadata row.
func (idx *providerRegistrationIndex) resolve(row *entity.SysAuthProvider) (pluginsvc.AuthProviderRegistration, bool) {
	if idx == nil || row == nil {
		return pluginsvc.AuthProviderRegistration{}, false
	}
	providerKey := strings.TrimSpace(row.ProviderKey)
	if providerKey != "" {
		if registration, ok := idx.byKey[providerKey]; ok && registrationMatchesRow(registration, row) {
			return registration, true
		}
	}
	providerType := strings.TrimSpace(row.ProviderType)
	if providerType == "" {
		return pluginsvc.AuthProviderRegistration{}, false
	}
	for _, registration := range idx.byType[providerType] {
		if registrationMatchesRow(registration, row) {
			return registration, true
		}
	}
	return pluginsvc.AuthProviderRegistration{}, false
}

// registrationMatchesRow confirms that a plugin-owned handler belongs to the
// provider metadata row and supports its declared protocol type.
func registrationMatchesRow(registration pluginsvc.AuthProviderRegistration, row *entity.SysAuthProvider) bool {
	if row == nil || registration.Provider == nil {
		return false
	}
	if strings.TrimSpace(registration.PluginID) != strings.TrimSpace(row.PluginId) {
		return false
	}
	providerType := strings.TrimSpace(row.ProviderType)
	return providerType == "" || strings.TrimSpace(registration.Provider.ProviderType()) == providerType
}

// providerItemFromRow builds one public projection from DB metadata and plugin defaults.
func providerItemFromRow(row *entity.SysAuthProvider, provider pluginhost.AuthProvider) ProviderItem {
	item := ProviderItem{
		ProviderKey:  row.ProviderKey,
		ProviderType: row.ProviderType,
		Name:         row.Name,
		Icon:         row.Icon,
		Sort:         row.Sort,
	}
	if provider != nil {
		if item.ProviderType == "" {
			item.ProviderType = provider.ProviderType()
		}
		if item.Name == "" {
			item.Name = provider.DisplayName()
		}
		if item.Icon == "" {
			item.Icon = provider.Icon()
		}
	}
	return item
}

// normalizePurpose returns a supported callback purpose.
func normalizePurpose(purpose string) string {
	switch strings.TrimSpace(purpose) {
	case authProviderPurposeBind:
		return authProviderPurposeBind
	default:
		return authProviderPurposeLogin
	}
}

// authStateCacheKey builds the scoped kvcache key for one OAuth state.
func authStateCacheKey(state string) string {
	return kvcache.BuildCacheKey(authProviderCacheOwner, authProviderStateNamespace, state)
}

// authStateConsumeCacheKey builds the scoped kvcache key for one state consume marker.
func authStateConsumeCacheKey(state string) string {
	return kvcache.BuildCacheKey(authProviderCacheOwner, authProviderStateConsumeNS, state)
}

// randomURLToken returns a URL-safe random token.
func randomURLToken(size int) (string, error) {
	data := make([]byte, size)
	if _, err := rand.Read(data); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(data), nil
}

// pkceChallenge returns the RFC 7636 S256 code challenge.
func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// cloneStringMap returns a detached copy of a callback parameter map.
func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

// marshalJSONString returns a JSON string for generated JSONB string fields.
func marshalJSONString(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(data)
}
