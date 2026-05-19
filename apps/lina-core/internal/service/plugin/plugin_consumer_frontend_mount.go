// This file resolves source-plugin consumer frontend assets through stable
// manifest-declared mount paths and applies mount-specific response policies.

package plugin

import (
	"context"
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/internal/service/plugin/internal/catalog"
	"lina-core/internal/service/plugin/internal/frontend"
)

const (
	// sourceConsumerFrontendHTMLContentType is the MIME type used when serving rewritten HTML entries.
	sourceConsumerFrontendHTMLContentType = "text/html; charset=utf-8"
)

var (
	// errSourceConsumerFrontendMountNotFound identifies requests outside declared mounts.
	errSourceConsumerFrontendMountNotFound = errors.New("source consumer frontend mount path does not exist")
	// errSourceConsumerFrontendMountDisabled identifies a matched mount whose plugin is not enabled.
	errSourceConsumerFrontendMountDisabled = errors.New("source consumer frontend mount plugin is not enabled")
	// errSourceConsumerFrontendMountAssetNotFound identifies a matched mount whose requested asset is absent.
	errSourceConsumerFrontendMountAssetNotFound = errors.New("source consumer frontend mount asset does not exist")
	// sourceConsumerFrontendHeadPattern locates the document head start tag used for base injection.
	sourceConsumerFrontendHeadPattern = regexp.MustCompile(`(?i)<head(?:\s[^>]*)?>`)
)

// ResolveSourceConsumerFrontendMountAsset resolves one enabled source-plugin
// consumer frontend asset by a manifest-declared user-facing mount path.
func (s *serviceImpl) ResolveSourceConsumerFrontendMountAsset(
	ctx context.Context,
	requestPath string,
) (*SourceConsumerFrontendMountAssetOutput, error) {
	normalizedRequestPath := normalizeSourceConsumerFrontendMountRequestPath(requestPath)
	if normalizedRequestPath == "" {
		return nil, errSourceConsumerFrontendMountNotFound
	}

	index, err := s.loadSourceConsumerFrontendResourceIndex(ctx)
	if err != nil {
		return nil, err
	}

	matchedMount, matchedRelativePath := index.match(normalizedRequestPath)
	if matchedMount == nil {
		return nil, errSourceConsumerFrontendMountNotFound
	}
	if !s.IsEnabled(ctx, matchedMount.pluginID) {
		return nil, errSourceConsumerFrontendMountDisabled
	}

	if matchedRelativePath == "" || matchedRelativePath == "/" {
		matchedRelativePath = matchedMount.index
	}
	if !matchedMount.assetDeclared(matchedRelativePath) {
		if !matchedMount.spaFallback || looksLikeSourceConsumerStaticAsset(matchedRelativePath) {
			return nil, gerror.Wrap(
				errSourceConsumerFrontendMountAssetNotFound,
				"current source plugin consumer frontend asset does not exist",
			)
		}
		matchedRelativePath = matchedMount.index
	}
	if matchedRelativePath == matchedMount.index && matchedMount.indexAsset != nil {
		return cloneFrontendAssetOutput(matchedMount.indexAsset), nil
	}
	out, err := s.ResolveSourceConsumerFrontendAsset(
		ctx,
		matchedMount.pluginID,
		matchedMount.version,
		matchedRelativePath,
	)
	if err == nil {
		return applySourceConsumerMountAssetPolicy(
			rewriteSourceConsumerHTMLBase(out, matchedMount.mountPath),
			matchedMount,
			matchedRelativePath,
		), nil
	}
	if !matchedMount.spaFallback || looksLikeSourceConsumerStaticAsset(matchedRelativePath) {
		return nil, err
	}

	if matchedMount.indexAsset != nil {
		return cloneFrontendAssetOutput(matchedMount.indexAsset), nil
	}
	out, err = s.ResolveSourceConsumerFrontendAsset(
		ctx,
		matchedMount.pluginID,
		matchedMount.version,
		matchedMount.index,
	)
	if err != nil {
		return nil, err
	}
	return applySourceConsumerMountAssetPolicy(
		rewriteSourceConsumerHTMLBase(out, matchedMount.mountPath),
		matchedMount,
		matchedMount.index,
	), nil
}

// normalizeSourceConsumerFrontendMountRequestPath normalizes browser request paths for mount matching.
func normalizeSourceConsumerFrontendMountRequestPath(requestPath string) string {
	trimmed := strings.TrimSpace(requestPath)
	if trimmed == "" {
		return ""
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	return strings.TrimRight(trimmed, "/")
}

// matchSourceConsumerFrontendMountPath returns the asset-relative path when a request is under a mount.
func matchSourceConsumerFrontendMountPath(requestPath string, mountPath string) (string, bool) {
	normalizedMount := strings.TrimRight(strings.TrimSpace(mountPath), "/")
	if normalizedMount == "" {
		return "", false
	}
	if requestPath == normalizedMount {
		return "", true
	}
	prefix := normalizedMount + "/"
	if !strings.HasPrefix(requestPath, prefix) {
		return "", false
	}
	return strings.TrimPrefix(requestPath, prefix), true
}

// sourceConsumerSPAFallbackEnabled reports whether missing clean routes should
// serve the frontend entry. Plugins must opt in explicitly so stable mounts do
// not swallow arbitrary extensionless paths by default.
func sourceConsumerSPAFallbackEnabled(frontendSpec *catalog.ConsumerFrontendSpec) bool {
	if frontendSpec == nil || frontendSpec.SPAFallback == nil {
		return false
	}
	return *frontendSpec.SPAFallback
}

// looksLikeSourceConsumerStaticAsset avoids serving the SPA entry for missing concrete assets.
func looksLikeSourceConsumerStaticAsset(relativePath string) bool {
	return filepath.Ext(strings.TrimSpace(relativePath)) != ""
}

// buildSourceConsumerFrontendMountIndexAsset precomputes the stable mounted
// entry HTML so repeated clean-route requests do not rewrite it on every hit.
func (s *serviceImpl) buildSourceConsumerFrontendMountIndexAsset(
	_ context.Context,
	manifest *catalog.Manifest,
	mount *sourceConsumerFrontendMountEntry,
) *frontend.RuntimeFrontendAssetOutput {
	if manifest == nil || mount == nil || !mount.assetDeclared(mount.index) {
		return nil
	}
	assetPath, err := normalizeSourceConsumerFrontendAssetPath(mount.index)
	if err != nil {
		return nil
	}
	contentBytes, err := s.catalogSvc.ReadSourcePluginAssetBytes(manifest, assetPath)
	if err != nil {
		return nil
	}
	contentType := mime.TypeByExtension(filepath.Ext(assetPath))
	if contentType == "" {
		contentType = http.DetectContentType(contentBytes)
	}
	out := &frontend.RuntimeFrontendAssetOutput{
		Content:      contentBytes,
		ContentType:  contentType,
		ETag:         frontend.BuildAssetETag(contentBytes, "source-consumer-mount", mount.pluginID, mount.version, mount.mountPath, assetPath),
		CacheControl: frontend.CacheControlRevalidate,
	}
	return applySourceConsumerMountAssetPolicy(rewriteSourceConsumerHTMLBase(out, mount.mountPath), mount, mount.index)
}

// applySourceConsumerMountAssetPolicy keeps stable mount paths validator-based
// even when the underlying asset came from a versioned debug namespace.
func applySourceConsumerMountAssetPolicy(
	out *SourceConsumerFrontendMountAssetOutput,
	mount *sourceConsumerFrontendMountEntry,
	relativePath string,
) *SourceConsumerFrontendMountAssetOutput {
	if out == nil {
		return nil
	}
	out.CacheControl = frontend.CacheControlRevalidate
	if mount != nil {
		out.ETag = frontend.BuildAssetETag(
			out.Content,
			"source-consumer-mount",
			mount.pluginID,
			mount.version,
			mount.mountPath,
			relativePath,
		)
		return out
	}
	out.ETag = frontend.BuildAssetETag(out.Content, "source-consumer-mount", out.ContentType)
	return out
}

// rewriteSourceConsumerHTMLBase injects a base tag so clean mounted routes load relative assets.
func rewriteSourceConsumerHTMLBase(
	out *SourceConsumerFrontendMountAssetOutput,
	mountPath string,
) *SourceConsumerFrontendMountAssetOutput {
	if out == nil || !strings.Contains(strings.ToLower(out.ContentType), "text/html") {
		return out
	}
	content := string(out.Content)
	if strings.Contains(strings.ToLower(content), "<base ") {
		return out
	}
	baseHref := strings.TrimRight(mountPath, "/") + "/"
	match := sourceConsumerFrontendHeadPattern.FindStringIndex(content)
	var rewritten string
	if match == nil {
		rewritten = "<base href=\"" + baseHref + "\" />\n" + content
	} else {
		insertionPoint := match[1]
		rewritten = content[:insertionPoint] + "\n    <base href=\"" + baseHref + "\" />" + content[insertionPoint:]
	}
	return &frontend.RuntimeFrontendAssetOutput{
		Content:      []byte(rewritten),
		ContentType:  sourceConsumerFrontendHTMLContentType,
		ETag:         frontend.BuildAssetETag([]byte(rewritten), "source-consumer-mount-html", mountPath),
		CacheControl: frontend.CacheControlRevalidate,
	}
}

// isSourceConsumerFrontendMountNotFound reports whether an error is a no-match result.
func isSourceConsumerFrontendMountNotFound(err error) bool {
	return errors.Is(err, errSourceConsumerFrontendMountNotFound)
}

// isSourceConsumerFrontendMountDisabled reports whether a matched mount belongs to a disabled plugin.
func isSourceConsumerFrontendMountDisabled(err error) bool {
	return errors.Is(err, errSourceConsumerFrontendMountDisabled)
}

// isSourceConsumerFrontendMountAssetNotFound reports whether a matched mount missed one concrete asset.
func isSourceConsumerFrontendMountAssetNotFound(err error) bool {
	return errors.Is(err, errSourceConsumerFrontendMountAssetNotFound)
}
