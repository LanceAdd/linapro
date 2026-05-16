// This file defines static configuration loading for pluggable authentication
// providers and their non-database secret references.

package config

import (
	"context"
	"strings"
	"time"
)

// DefaultAuthProviderStateTTL is the fallback OAuth state lifetime.
const DefaultAuthProviderStateTTL = 10 * time.Minute

// AuthProvidersConfig holds static pluggable authentication provider settings.
type AuthProvidersConfig struct {
	StateTTL time.Duration         `json:"stateTTL"` // StateTTL controls authorization temporary state lifetime.
	OIDC     []OIDCProviderConfig  `json:"oidc"`     // OIDC lists configured OIDC provider instances.
	WeCom    []WeComProviderConfig `json:"wecom"`    // WeCom lists configured Enterprise WeChat provider instances.
}

// OIDCProviderConfig describes one static OIDC provider instance.
type OIDCProviderConfig struct {
	ProviderKey     string   `json:"providerKey"`     // ProviderKey is the stable provider instance key.
	DisplayName     string   `json:"displayName"`     // DisplayName is the default provider display name.
	Issuer          string   `json:"issuer"`          // Issuer is the OIDC issuer URL.
	ClientID        string   `json:"clientId"`        // ClientID is the public OIDC client identifier.
	ClientSecretEnv string   `json:"clientSecretEnv"` // ClientSecretEnv names the environment variable containing the client secret.
	RedirectURL     string   `json:"redirectUrl"`     // RedirectURL is the host callback URL.
	Scopes          []string `json:"scopes"`          // Scopes are requested OIDC scopes.
}

// WeComProviderConfig describes one Enterprise WeChat provider instance.
type WeComProviderConfig struct {
	ProviderKey   string `json:"providerKey"`   // ProviderKey is the stable provider instance key.
	DisplayName   string `json:"displayName"`   // DisplayName is the default provider display name.
	CorpID        string `json:"corpId"`        // CorpID is the Enterprise WeChat corp identifier.
	AgentID       string `json:"agentId"`       // AgentID is the Enterprise WeChat application identifier.
	CorpSecretEnv string `json:"corpSecretEnv"` // CorpSecretEnv names the environment variable containing the corp secret.
	RedirectURL   string `json:"redirectUrl"`   // RedirectURL is the host callback URL.
}

// GetAuthProviders reads static pluggable authentication provider settings.
func (s *serviceImpl) GetAuthProviders(ctx context.Context) *AuthProvidersConfig {
	return cloneAuthProvidersConfig(processStaticConfigCaches.authProviders.load(func() *AuthProvidersConfig {
		cfg := &AuthProvidersConfig{
			StateTTL: DefaultAuthProviderStateTTL,
		}
		mustScanConfig(ctx, "authProviders", cfg)
		cfg.StateTTL = mustLoadDurationConfig(ctx, "authProviders.stateTTL", cfg.StateTTL)
		normalizeAuthProvidersConfig(cfg)
		return cfg
	}))
}

// normalizeAuthProvidersConfig trims provider configuration fields and applies
// OIDC defaults without reading any secret values.
func normalizeAuthProvidersConfig(cfg *AuthProvidersConfig) {
	if cfg == nil {
		return
	}
	for index := range cfg.OIDC {
		item := &cfg.OIDC[index]
		item.ProviderKey = strings.TrimSpace(item.ProviderKey)
		item.DisplayName = strings.TrimSpace(item.DisplayName)
		item.Issuer = strings.TrimSpace(item.Issuer)
		item.ClientID = strings.TrimSpace(item.ClientID)
		item.ClientSecretEnv = strings.TrimSpace(item.ClientSecretEnv)
		item.RedirectURL = strings.TrimSpace(item.RedirectURL)
		item.Scopes = normalizeProviderScopes(item.Scopes)
		if len(item.Scopes) == 0 {
			item.Scopes = []string{"openid", "profile", "email"}
		}
	}
	for index := range cfg.WeCom {
		item := &cfg.WeCom[index]
		item.ProviderKey = strings.TrimSpace(item.ProviderKey)
		item.DisplayName = strings.TrimSpace(item.DisplayName)
		item.CorpID = strings.TrimSpace(item.CorpID)
		item.AgentID = strings.TrimSpace(item.AgentID)
		item.CorpSecretEnv = strings.TrimSpace(item.CorpSecretEnv)
		item.RedirectURL = strings.TrimSpace(item.RedirectURL)
	}
}

// normalizeProviderScopes trims and de-duplicates OAuth scopes while preserving
// declaration order.
func normalizeProviderScopes(scopes []string) []string {
	if len(scopes) == 0 {
		return nil
	}
	out := make([]string, 0, len(scopes))
	seen := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	return out
}
