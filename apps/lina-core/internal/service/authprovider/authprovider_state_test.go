// This file tests temporary authorization state storage and consumption.

package authprovider

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"

	"lina-core/internal/model/entity"
	configsvc "lina-core/internal/service/config"
	"lina-core/internal/service/kvcache"
	pluginsvc "lina-core/internal/service/plugin"
	"lina-core/pkg/pluginhost"
)

// testAuthProviderConfigReader serves static auth-provider settings to tests.
type testAuthProviderConfigReader struct {
	cfg *configsvc.AuthProvidersConfig
}

// GetAuthProviders returns test auth-provider configuration.
func (r testAuthProviderConfigReader) GetAuthProviders(context.Context) *configsvc.AuthProvidersConfig {
	return r.cfg
}

// memoryKVCache stores kvcache items in memory for authprovider tests.
type memoryKVCache struct {
	items map[string]*kvcache.Item
}

// newMemoryKVCache creates an empty memory-backed kvcache test double.
func newMemoryKVCache() *memoryKVCache {
	return &memoryKVCache{items: make(map[string]*kvcache.Item)}
}

// BackendName returns the test backend name.
func (s *memoryKVCache) BackendName() kvcache.BackendName {
	return kvcache.BackendName("memory-test")
}

// RequiresExpiredCleanup reports no external cleanup requirement.
func (s *memoryKVCache) RequiresExpiredCleanup() bool {
	return false
}

// Get returns one unexpired item.
func (s *memoryKVCache) Get(_ context.Context, ownerType kvcache.OwnerType, cacheKey string) (*kvcache.Item, bool, error) {
	key := testKVKey(ownerType, cacheKey)
	item, ok := s.items[key]
	if !ok {
		return nil, false, nil
	}
	if item.ExpireAt != nil && time.Now().After(item.ExpireAt.Time) {
		delete(s.items, key)
		return nil, false, nil
	}
	cloned := *item
	return &cloned, true, nil
}

// GetInt returns one unexpired integer item.
func (s *memoryKVCache) GetInt(ctx context.Context, ownerType kvcache.OwnerType, cacheKey string) (int64, bool, error) {
	item, ok, err := s.Get(ctx, ownerType, cacheKey)
	if err != nil || !ok {
		return 0, ok, err
	}
	return item.IntValue, true, nil
}

// Set stores one string item.
func (s *memoryKVCache) Set(_ context.Context, ownerType kvcache.OwnerType, cacheKey string, value string, ttl time.Duration) (*kvcache.Item, error) {
	item := &kvcache.Item{Key: cacheKey, ValueKind: kvcache.ValueKindString, Value: value}
	if ttl > 0 {
		item.ExpireAt = gtime.New(time.Now().Add(ttl))
	}
	s.items[testKVKey(ownerType, cacheKey)] = item
	return item, nil
}

// Delete removes one item.
func (s *memoryKVCache) Delete(_ context.Context, ownerType kvcache.OwnerType, cacheKey string) error {
	delete(s.items, testKVKey(ownerType, cacheKey))
	return nil
}

// Incr increments one integer item.
func (s *memoryKVCache) Incr(_ context.Context, ownerType kvcache.OwnerType, cacheKey string, delta int64, ttl time.Duration) (*kvcache.Item, error) {
	key := testKVKey(ownerType, cacheKey)
	item, ok := s.items[key]
	if !ok || item.ValueKind != kvcache.ValueKindInt {
		item = &kvcache.Item{Key: cacheKey, ValueKind: kvcache.ValueKindInt}
	}
	item.IntValue += delta
	if ttl > 0 {
		item.ExpireAt = gtime.New(time.Now().Add(ttl))
	}
	s.items[key] = item
	return item, nil
}

// Expire updates one item TTL.
func (s *memoryKVCache) Expire(_ context.Context, ownerType kvcache.OwnerType, cacheKey string, ttl time.Duration) (bool, *gtime.Time, error) {
	item, ok := s.items[testKVKey(ownerType, cacheKey)]
	if !ok {
		return false, nil, nil
	}
	if ttl <= 0 {
		item.ExpireAt = nil
		return true, nil, nil
	}
	item.ExpireAt = gtime.New(time.Now().Add(ttl))
	return true, item.ExpireAt, nil
}

// CleanupExpired removes expired test entries.
func (s *memoryKVCache) CleanupExpired(context.Context) error {
	for key, item := range s.items {
		if item.ExpireAt != nil && time.Now().After(item.ExpireAt.Time) {
			delete(s.items, key)
		}
	}
	return nil
}

// TestConsumeStateRejectsRepeat verifies OAuth state is single-use.
func TestConsumeStateRejectsRepeat(t *testing.T) {
	ctx := context.Background()
	svc := &serviceImpl{
		configSvc:  testAuthProviderConfigReader{cfg: &configsvc.AuthProvidersConfig{StateTTL: time.Minute}},
		kvCacheSvc: newMemoryKVCache(),
	}
	record := authStateRecord{
		Schema:      authProviderStateSchema,
		ProviderKey: "google",
		State:       "state-1",
		Purpose:     authProviderPurposeLogin,
		ExpiresAt:   time.Now().Add(time.Minute),
	}
	if err := svc.storeState(ctx, record); err != nil {
		t.Fatalf("store state: %v", err)
	}
	if _, ok, err := svc.consumeState(ctx, "google", "state-1"); err != nil || !ok {
		t.Fatalf("expected first consume to succeed ok=%v err=%v", ok, err)
	}
	if _, ok, err := svc.consumeState(ctx, "google", "state-1"); err != nil || ok {
		t.Fatalf("expected repeat consume to be rejected ok=%v err=%v", ok, err)
	}
}

// TestConsumeStateRejectsExpired verifies expired OAuth state is rejected.
func TestConsumeStateRejectsExpired(t *testing.T) {
	ctx := context.Background()
	svc := &serviceImpl{
		configSvc:  testAuthProviderConfigReader{cfg: &configsvc.AuthProvidersConfig{StateTTL: time.Minute}},
		kvCacheSvc: newMemoryKVCache(),
	}
	record := authStateRecord{
		Schema:      authProviderStateSchema,
		ProviderKey: "google",
		State:       "state-expired",
		Purpose:     authProviderPurposeLogin,
		ExpiresAt:   time.Now().Add(-time.Minute),
	}
	if err := svc.storeState(ctx, record); err != nil {
		t.Fatalf("store state: %v", err)
	}
	if _, ok, err := svc.consumeState(ctx, "google", "state-expired"); err != nil || ok {
		t.Fatalf("expected expired consume to be rejected ok=%v err=%v", ok, err)
	}
}

// TestProviderRegistrationIndexFallsBackToProviderType verifies multi-instance
// provider rows can reuse a type-level plugin handler.
func TestProviderRegistrationIndexFallsBackToProviderType(t *testing.T) {
	provider := testAuthProvider{key: "oidc", providerType: "oidc"}
	index := &providerRegistrationIndex{
		byKey: map[string]pluginsvc.AuthProviderRegistration{
			"oidc": {PluginID: "auth-oidc", Provider: provider},
		},
		byType: map[string][]pluginsvc.AuthProviderRegistration{
			"oidc": {{PluginID: "auth-oidc", Provider: provider}},
		},
	}
	registration, ok := index.resolve(&entity.SysAuthProvider{
		ProviderKey:  "google",
		ProviderType: "oidc",
		PluginId:     "auth-oidc",
	})
	if !ok || registration.Provider == nil {
		t.Fatalf("expected provider type fallback to resolve")
	}
	if registration.Provider.ProviderKey() != "oidc" {
		t.Fatalf("expected oidc handler, got %q", registration.Provider.ProviderKey())
	}
}

// testKVKey scopes cache keys by owner type for the test double.
func testKVKey(ownerType kvcache.OwnerType, cacheKey string) string {
	return ownerType.String() + ":" + cacheKey
}

var _ kvcache.Service = (*memoryKVCache)(nil)

// testAuthProvider implements pluginhost.AuthProvider for registration tests.
type testAuthProvider struct {
	key          string
	providerType string
}

// ProviderKey returns the test provider key.
func (p testAuthProvider) ProviderKey() string {
	return p.key
}

// ProviderType returns the test provider type.
func (p testAuthProvider) ProviderType() string {
	return p.providerType
}

// DisplayName returns the test display name.
func (p testAuthProvider) DisplayName() string {
	return "Test Provider"
}

// Icon returns the test icon.
func (p testAuthProvider) Icon() string {
	return "test"
}

// BuildAuthorizeURL is unused by registration-index tests.
func (p testAuthProvider) BuildAuthorizeURL(context.Context, pluginhost.AuthProviderAuthorizeInput) (*pluginhost.AuthProviderAuthorizeOutput, error) {
	return nil, nil
}

// ExchangeCallback is unused by registration-index tests.
func (p testAuthProvider) ExchangeCallback(context.Context, pluginhost.AuthProviderCallbackInput) (*pluginhost.ExternalIdentity, error) {
	return nil, nil
}
