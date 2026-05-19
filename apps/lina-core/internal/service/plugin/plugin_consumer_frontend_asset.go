// This file resolves direct source-plugin consumer frontend assets from the
// debug namespace and shared asset metadata helpers.

package plugin

import (
	"context"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/internal/service/plugin/internal/frontend"
	"lina-core/pkg/pluginfs"
)

// ResolveSourceConsumerFrontendAsset resolves one enabled source-plugin
// consumer frontend asset for public serving.
func (s *serviceImpl) ResolveSourceConsumerFrontendAsset(
	ctx context.Context,
	pluginID string,
	version string,
	relativePath string,
) (*SourceConsumerFrontendAssetOutput, error) {
	normalizedPluginID := strings.TrimSpace(pluginID)
	normalizedVersion := strings.TrimSpace(version)
	if normalizedPluginID == "" || normalizedVersion == "" {
		return nil, gerror.New("source consumer frontend asset identity is incomplete")
	}
	if !s.IsEnabled(ctx, normalizedPluginID) {
		return nil, gerror.New("current source plugin is not enabled")
	}

	manifest, err := s.catalogSvc.GetActiveManifest(ctx, normalizedPluginID)
	if err != nil {
		return nil, err
	}
	if manifest == nil || catalog.NormalizeType(manifest.Type) != catalog.TypeSource {
		return nil, gerror.New("current plugin is not a source plugin")
	}
	if strings.TrimSpace(manifest.Version) != normalizedVersion {
		return nil, gerror.New("current source plugin version does not exist or has switched")
	}

	assetPath, err := normalizeSourceConsumerFrontendAssetPath(relativePath)
	if err != nil {
		return nil, err
	}
	if !sourceConsumerFrontendAssetDeclared(s.catalogSvc.ListConsumerFrontendPaths(manifest), assetPath) {
		return nil, gerror.Wrap(
			errSourceConsumerFrontendMountAssetNotFound,
			"current source plugin consumer frontend asset does not exist",
		)
	}

	contentBytes, err := s.catalogSvc.ReadSourcePluginAssetBytes(manifest, assetPath)
	if err != nil {
		return nil, err
	}
	contentType := mime.TypeByExtension(filepath.Ext(assetPath))
	if contentType == "" {
		contentType = http.DetectContentType(contentBytes)
	}
	return &frontend.RuntimeFrontendAssetOutput{
		Content:      contentBytes,
		ContentType:  contentType,
		ETag:         frontend.BuildAssetETag(contentBytes, "source-consumer", normalizedPluginID, normalizedVersion, assetPath),
		CacheControl: frontend.CacheControlForContentType(contentType),
	}, nil
}

// BuildSourceConsumerFrontendPublicBaseURL returns the stable public base URL
// for source-plugin consumer frontend assets.
func (s *serviceImpl) BuildSourceConsumerFrontendPublicBaseURL(pluginID string, version string) string {
	return "/consumer-plugin-assets/" + strings.TrimSpace(pluginID) + "/" + strings.TrimSpace(version) + "/"
}

// normalizeSourceConsumerFrontendAssetPath converts browser-facing consumer
// asset requests into plugin-relative frontend/consumer paths.
func normalizeSourceConsumerFrontendAssetPath(relativePath string) (string, error) {
	trimmedPath := strings.TrimSpace(relativePath)
	if trimmedPath == "" || trimmedPath == "/" {
		trimmedPath = "index.html"
	}
	normalizedPath, err := pluginfs.NormalizeRelativePath(trimmedPath)
	if err != nil {
		return "", err
	}
	return "frontend/consumer/" + strings.TrimPrefix(normalizedPath, "frontend/consumer/"), nil
}

// sourceConsumerFrontendAssetDeclared reports whether one normalized asset path
// is part of the source plugin's declared consumer frontend directory.
func sourceConsumerFrontendAssetDeclared(paths []string, assetPath string) bool {
	return buildSourceConsumerFrontendAssetSet(paths).has(assetPath)
}

// sourceConsumerFrontendAssetSet stores normalized source-plugin consumer frontend asset paths.
type sourceConsumerFrontendAssetSet map[string]struct{}

// buildSourceConsumerFrontendAssetSet builds a lookup set from catalog asset listings.
func buildSourceConsumerFrontendAssetSet(paths []string) sourceConsumerFrontendAssetSet {
	assets := make(sourceConsumerFrontendAssetSet, len(paths))
	for _, item := range paths {
		normalizedPath := strings.TrimSpace(item)
		if normalizedPath == "" {
			continue
		}
		assets[normalizedPath] = struct{}{}
	}
	return assets
}

// has reports whether the normalized asset path is declared in the set.
func (assets sourceConsumerFrontendAssetSet) has(assetPath string) bool {
	if len(assets) == 0 {
		return false
	}
	_, ok := assets[strings.TrimSpace(assetPath)]
	return ok
}

// cloneFrontendAssetOutput copies one frontend output so cached bytes cannot be
// mutated by callers.
func cloneFrontendAssetOutput(out *frontend.RuntimeFrontendAssetOutput) *frontend.RuntimeFrontendAssetOutput {
	if out == nil {
		return nil
	}
	cloned := *out
	if len(out.Content) > 0 {
		cloned.Content = append([]byte(nil), out.Content...)
	}
	return &cloned
}
