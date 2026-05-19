// Package frontend manages in-memory frontend asset bundles built from runtime WASM
// artifacts. The WASM artifact is the single source of truth; bundles are cached in
// memory and rebuilt on demand after a host restart, eliminating extraction to disk.
package frontend

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"lina-core/internal/model/entity"
	"lina-core/internal/service/plugin/internal/catalog"
)

const (
	// CacheControlRevalidate allows browser storage but requires ETag validation before reuse.
	CacheControlRevalidate = "no-cache"
	// CacheControlStaticRevalidate keeps plugin static assets correct across same-version refreshes.
	CacheControlStaticRevalidate = CacheControlRevalidate
)

// RuntimeFrontendAssetOutput contains one resolved frontend asset ready to be served.
type RuntimeFrontendAssetOutput struct {
	// Content is the raw asset body.
	Content []byte
	// ContentType is the HTTP Content-Type header value returned to browsers.
	ContentType string
	// ETag is a strong HTTP validator derived from asset identity and content.
	ETag string
	// CacheControl is the HTTP Cache-Control policy for this asset.
	CacheControl string
}

// Service defines the frontend service contract.
type Service interface {
	// EnsureBundleReader returns a BundleReader for the manifest, building and caching the bundle if needed.
	EnsureBundleReader(ctx context.Context, manifest *catalog.Manifest) (BundleReader, error)
	// ValidateRuntimeFrontendMenuBindings verifies that dynamic plugin menus only reference
	// hosted assets that exist in the plugin's in-memory bundle.
	ValidateRuntimeFrontendMenuBindings(ctx context.Context, manifest *catalog.Manifest) error
	// ValidateHostedMenuBindings is the exported form of validateHostedMenuBindings for cross-package access.
	ValidateHostedMenuBindings(ctx context.Context, manifest *catalog.Manifest, menus []*entity.SysMenu) error
	// PrewarmRuntimeFrontendBundles rebuilds in-memory frontend bundles for all enabled
	// dynamic plugins during host startup. A single failed preload does not stop the host;
	// errors are collected and returned as one joined error.
	PrewarmRuntimeFrontendBundles(ctx context.Context) error
	// ResolveRuntimeFrontendAsset resolves one enabled dynamic plugin frontend asset for public serving.
	ResolveRuntimeFrontendAsset(
		ctx context.Context,
		pluginID string,
		version string,
		relativePath string,
	) (*RuntimeFrontendAssetOutput, error)
	// BuildRuntimeFrontendPublicBaseURL returns the stable public base URL for runtime plugin assets.
	BuildRuntimeFrontendPublicBaseURL(pluginID string, version string) string
	// InvalidateBundle removes all cached bundle entries for the given plugin ID.
	InvalidateBundle(ctx context.Context, pluginID string, reason string)
	// InvalidateAllBundles removes every cached runtime frontend bundle.
	InvalidateAllBundles(ctx context.Context, reason string)
	// EnsureBundle guarantees an in-memory frontend bundle exists for the given manifest,
	// building and caching it if necessary. Returns the bundle for immediate use.
	// This is called by the runtime reconciler to pre-warm bundles after reconciliation.
	EnsureBundle(ctx context.Context, manifest *catalog.Manifest) error
}

// Ensure serviceImpl satisfies the runtime frontend asset contract.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	// catalogSvc provides registry and release lookups for enabled runtime plugins.
	catalogSvc catalog.Service
}

// New creates a frontend Service backed by the shared plugin catalog.
func New(catalogSvc catalog.Service) Service {
	return &serviceImpl{catalogSvc: catalogSvc}
}

// BuildAssetETag returns a strong ETag derived from asset identity parts and
// the current content bytes.
func BuildAssetETag(content []byte, identityParts ...string) string {
	hasher := sha256.New()
	for _, part := range identityParts {
		hasher.Write([]byte(strings.TrimSpace(part)))
		hasher.Write([]byte{0})
	}
	hasher.Write(content)
	return `"` + hex.EncodeToString(hasher.Sum(nil)) + `"`
}

// CacheControlForContentType returns the default cache policy for a plugin
// frontend asset. Plugin URLs do not include release checksums yet, so even
// static assets revalidate to stay correct after same-version refreshes.
func CacheControlForContentType(contentType string) string {
	if strings.Contains(strings.ToLower(strings.TrimSpace(contentType)), "text/html") {
		return CacheControlRevalidate
	}
	return CacheControlStaticRevalidate
}
