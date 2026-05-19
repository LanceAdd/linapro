// This file defines the public consumer-service contracts shared by the host
// middleware, source plugins, route binding snapshots, and OpenAPI projection.
// The contracts deliberately stop at request metadata so consumer account,
// login, session, token, and authorization models remain plugin-owned.

package pluginhost

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// AdminAPIPrefix is the stable host administration API prefix.
	AdminAPIPrefix = "/api/v1"
	// ConsumerAPIPrefix is the stable consumer plugin API prefix.
	ConsumerAPIPrefix = "/api/c/v1"
)

const (
	// ConsumerHeaderAnonymousID carries one caller-provided anonymous visitor identifier.
	ConsumerHeaderAnonymousID = "X-Consumer-Anonymous-Id"
	// ConsumerHeaderDeviceID carries one caller-provided consumer device identifier.
	ConsumerHeaderDeviceID = "X-Consumer-Device-Id"
	// ConsumerHeaderChannel carries one caller-provided consumer channel identifier.
	ConsumerHeaderChannel = "X-Consumer-Channel"
	// ConsumerHeaderTenantID carries the requested tenant boundary for public consumer APIs.
	ConsumerHeaderTenantID = "X-Tenant-Id"
)

const (
	// consumerContextKey stores the consumer request context in GoFrame request context values.
	consumerContextKey gctx.StrKey = "ConsumerCtx"
)

// ServiceSurface identifies the host-facing surface one HTTP route belongs to.
type ServiceSurface string

const (
	// SurfaceAdmin marks administrator and management routes.
	SurfaceAdmin ServiceSurface = "admin"
	// SurfaceConsumer marks consumer plugin routes.
	SurfaceConsumer ServiceSurface = "consumer"
)

// ConsumerContext stores request-scoped consumer metadata before and after
// optional authentication. It is in-memory only for the active request and must
// not be treated as durable session state.
type ConsumerContext struct {
	// PluginID is the source plugin that owns the matched consumer route.
	PluginID string `json:"pluginId"`
	// TenantID is the resolved tenant boundary for this request.
	TenantID int `json:"tenantId"`
	// TenantResolved reports whether a tenant value came from a request source.
	TenantResolved bool `json:"tenantResolved"`
	// Locale is the resolved runtime language.
	Locale string `json:"locale"`
	// AnonymousID is an optional caller-provided anonymous visitor identifier.
	AnonymousID string `json:"anonymousId"`
	// DeviceID is an optional caller-provided device identifier.
	DeviceID string `json:"deviceId"`
	// Channel is an optional caller-provided consumer channel.
	Channel string `json:"channel"`
}

// RouteSurfaceFromPath returns the stable service surface for a public route path.
func RouteSurfaceFromPath(path string) ServiceSurface {
	normalized := normalizeRoutePath(path)
	if normalized == ConsumerAPIPrefix || strings.HasPrefix(normalized, ConsumerAPIPrefix+"/") {
		return SurfaceConsumer
	}
	return SurfaceAdmin
}

// ConsumerPluginIDFromPath extracts the plugin id from /api/c/v1/<plugin-id>/...
// paths. It returns an empty string for non-consumer or malformed paths.
func ConsumerPluginIDFromPath(path string) string {
	normalized := normalizeRoutePath(path)
	if normalized == ConsumerAPIPrefix {
		return ""
	}
	remaining := strings.TrimPrefix(normalized, ConsumerAPIPrefix+"/")
	if remaining == normalized || remaining == "" {
		return ""
	}
	parts := strings.SplitN(remaining, "/", 2)
	return strings.TrimSpace(parts[0])
}

// ResolveConsumerTenantID reads the lightweight consumer tenant declaration
// supported by the host surface. Header wins over query string.
func ResolveConsumerTenantID(request *ghttp.Request) (int, bool, error) {
	if request == nil {
		return 0, false, nil
	}
	rawValue := strings.TrimSpace(request.GetHeader(ConsumerHeaderTenantID))
	if rawValue == "" {
		rawValue = strings.TrimSpace(request.GetQuery("tenantId").String())
	}
	if rawValue == "" {
		return 0, false, nil
	}
	tenantID, err := strconv.Atoi(rawValue)
	if err != nil {
		return 0, false, err
	}
	if tenantID < 0 {
		return 0, false, fmt.Errorf("consumer tenant id must be non-negative: %d", tenantID)
	}
	return tenantID, true, nil
}

// SetConsumerContext stores the consumer request context on the active request.
func SetConsumerContext(request *ghttp.Request, consumerCtx *ConsumerContext) {
	if request == nil || consumerCtx == nil {
		return
	}
	request.SetCtxVar(consumerContextKey, consumerCtx)
}

// ConsumerContextFromRequest returns the consumer context from a GoFrame request.
func ConsumerContextFromRequest(request *ghttp.Request) (*ConsumerContext, bool) {
	if request == nil {
		return nil, false
	}
	return ConsumerContextFromContext(request.Context())
}

// ConsumerContextFromContext returns the consumer context from a standard context.
func ConsumerContextFromContext(ctx context.Context) (*ConsumerContext, bool) {
	if ctx == nil {
		return nil, false
	}
	consumerCtx, ok := ctx.Value(consumerContextKey).(*ConsumerContext)
	return consumerCtx, ok && consumerCtx != nil
}

// normalizeRoutePath canonicalizes URL paths before prefix checks.
func normalizeRoutePath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "/"
	}
	if parsed, err := url.Parse(trimmed); err == nil && parsed.Path != "" {
		trimmed = parsed.Path
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	if trimmed != "/" {
		trimmed = strings.TrimRight(trimmed, "/")
	}
	return trimmed
}
