// This file verifies consumer service-surface context and tenant middleware
// behavior without depending on administrator users or host-owned C-side auth.

package middleware

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"

	"lina-core/internal/service/bizctx"
	"lina-core/internal/service/cachecoord"
	hostconfig "lina-core/internal/service/config"
	i18nsvc "lina-core/internal/service/i18n"
	"lina-core/pkg/pluginhost"
)

// TestConsumerPublicContextAllowsAnonymous verifies public consumer requests
// receive plugin, tenant, locale, and anonymous metadata without authentication.
func TestConsumerPublicContextAllowsAnonymous(t *testing.T) {
	status, body := runConsumerMiddlewareRequest(t, consumerMiddlewareCase{
		path:      "/api/c/v1/plugin-demo/public",
		tenantID:  "42",
		anonymous: "anon-1",
		deviceID:  "device-1",
		channel:   "web",
		register: func(group *ghttp.RouterGroup, svc *serviceImpl) {
			group.Middleware(svc.ConsumerCtx, svc.ConsumerTenant)
		},
		tenantSvc:  &tenancyTestTenantService{enabled: true},
		handlerURL: "/public",
	})

	if status != http.StatusOK {
		t.Fatalf("expected status 200, got %d body=%s", status, body)
	}
	if !strings.Contains(body, "plugin=plugin-demo tenant=42 anonymous=anon-1") {
		t.Fatalf("expected anonymous consumer context body, got %q", body)
	}
	if !strings.Contains(body, "device=device-1 channel=web") {
		t.Fatalf("expected device and channel metadata, got %q", body)
	}
}

// TestConsumerContextIgnoresAuthorizationHeader verifies the host consumer
// middleware does not interpret C-side Authorization headers.
func TestConsumerContextIgnoresAuthorizationHeader(t *testing.T) {
	status, body := runConsumerMiddlewareRequest(t, consumerMiddlewareCase{
		path:     "/api/c/v1/plugin-demo/login",
		tenantID: "42",
		token:    "plugin-owned-token",
		register: func(group *ghttp.RouterGroup, svc *serviceImpl) {
			group.Middleware(svc.ConsumerCtx, svc.ConsumerTenant)
		},
		tenantSvc:  &tenancyTestTenantService{enabled: true},
		handlerURL: "/login",
	})

	if status != http.StatusOK {
		t.Fatalf("expected host middleware to pass through plugin-owned auth header, got %d body=%s", status, body)
	}
	if !strings.Contains(body, "plugin=plugin-demo tenant=42") {
		t.Fatalf("expected consumer context to remain request metadata only, got %q", body)
	}
}

// TestConsumerTenantEnabledRequiresTenant verifies enabled tenancy fails closed
// when no lightweight tenant boundary is supplied.
func TestConsumerTenantEnabledRequiresTenant(t *testing.T) {
	status, body := runConsumerMiddlewareRequest(t, consumerMiddlewareCase{
		path: "/api/c/v1/plugin-demo/public",
		register: func(group *ghttp.RouterGroup, svc *serviceImpl) {
			group.Middleware(svc.ConsumerCtx, svc.ConsumerTenant)
		},
		tenantSvc:  &tenancyTestTenantService{enabled: true},
		handlerURL: "/public",
	})

	if status != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d body=%s", status, body)
	}
	if strings.Contains(body, "plugin=") {
		t.Fatalf("expected tenant failure before handler, got body %q", body)
	}
}

// TestConsumerDisabledPluginDoesNotReachHandler verifies the route registrar
// guard prevents disabled plugins from entering consumer business handlers.
func TestConsumerDisabledPluginDoesNotReachHandler(t *testing.T) {
	svc := newConsumerMiddlewareService(&tenancyTestTenantService{enabled: false})
	server := g.Server("middleware-consumer-disabled-" + guid.S())
	server.SetPort(0)
	server.SetDumpRouterMap(false)

	var rootGroup *ghttp.RouterGroup
	server.Group("/", func(group *ghttp.RouterGroup) {
		rootGroup = group
	})
	called := false
	registrar := pluginhost.NewRouteRegistrar(rootGroup, "plugin-demo", func(context.Context, string) bool {
		return false
	}, svc.PublishedRouteMiddlewares())
	registrar.Group(pluginhost.ConsumerAPIPrefix+"/plugin-demo", func(group pluginhost.RouteGroup) {
		group.Middleware(svc.ConsumerCtx, svc.ConsumerTenant)
		group.GET("/public", func(r *ghttp.Request) {
			called = true
			r.Response.Write("unexpected")
		})
	})

	if err := server.Start(); err != nil {
		t.Fatalf("start disabled consumer middleware server: %v", err)
	}
	t.Cleanup(func() {
		if err := server.Shutdown(); err != nil {
			t.Fatalf("shutdown disabled consumer middleware server: %v", err)
		}
	})
	time.Sleep(100 * time.Millisecond)

	response, err := http.Get("http://" + server.GetListenedAddress() + "/api/c/v1/plugin-demo/public")
	if err != nil {
		t.Fatalf("send disabled consumer request: %v", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			t.Fatalf("close disabled consumer response: %v", err)
		}
	}()
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.StatusCode)
	}
	if called {
		t.Fatalf("expected disabled plugin route not to reach handler")
	}
}

// consumerMiddlewareCase describes one test request through consumer middleware.
type consumerMiddlewareCase struct {
	path       string
	tenantID   string
	anonymous  string
	deviceID   string
	channel    string
	token      string
	register   func(group *ghttp.RouterGroup, svc *serviceImpl)
	tenantSvc  *tenancyTestTenantService
	handlerURL string
}

// runConsumerMiddlewareRequest serves one consumer request through configured middleware.
func runConsumerMiddlewareRequest(t *testing.T, testCase consumerMiddlewareCase) (int, string) {
	t.Helper()

	svc := newConsumerMiddlewareService(testCase.tenantSvc)
	server := g.Server("middleware-consumer-" + guid.S())
	server.SetPort(0)
	server.SetDumpRouterMap(false)
	server.Group(pluginhost.ConsumerAPIPrefix+"/plugin-demo", func(group *ghttp.RouterGroup) {
		testCase.register(group, svc)
		group.GET(testCase.handlerURL, func(r *ghttp.Request) {
			consumerCtx, ok := pluginhost.ConsumerContextFromRequest(r)
			if !ok {
				r.Response.Write("missing consumer context")
				return
			}
			r.Response.Writef(
				"plugin=%s tenant=%d anonymous=%s device=%s channel=%s",
				consumerCtx.PluginID,
				consumerCtx.TenantID,
				consumerCtx.AnonymousID,
				consumerCtx.DeviceID,
				consumerCtx.Channel,
			)
		})
	})
	if err := server.Start(); err != nil {
		t.Fatalf("start consumer middleware server: %v", err)
	}
	t.Cleanup(func() {
		if err := server.Shutdown(); err != nil {
			t.Fatalf("shutdown consumer middleware server: %v", err)
		}
	})
	time.Sleep(100 * time.Millisecond)

	request, err := http.NewRequest(http.MethodGet, "http://"+server.GetListenedAddress()+testCase.path, nil)
	if err != nil {
		t.Fatalf("create consumer middleware request: %v", err)
	}
	if testCase.tenantID != "" {
		request.Header.Set(pluginhost.ConsumerHeaderTenantID, testCase.tenantID)
	}
	if testCase.anonymous != "" {
		request.Header.Set(pluginhost.ConsumerHeaderAnonymousID, testCase.anonymous)
	}
	if testCase.deviceID != "" {
		request.Header.Set(pluginhost.ConsumerHeaderDeviceID, testCase.deviceID)
	}
	if testCase.channel != "" {
		request.Header.Set(pluginhost.ConsumerHeaderChannel, testCase.channel)
	}
	if testCase.token != "" {
		request.Header.Set("Authorization", "Bearer "+testCase.token)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("send consumer middleware request: %v", err)
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			t.Fatalf("close consumer middleware response: %v", err)
		}
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read consumer middleware response: %v", err)
	}
	return response.StatusCode, string(body)
}

// newConsumerMiddlewareService creates a middleware service for consumer tests.
func newConsumerMiddlewareService(tenantSvc *tenancyTestTenantService) *serviceImpl {
	svc := &serviceImpl{
		bizCtxSvc: bizctx.New(),
		configSvc: hostconfig.New(),
		i18nSvc:   i18nsvc.New(bizctx.New(), hostconfig.New(), cachecoord.Default(nil)),
		tenantSvc: tenantSvc,
	}
	return svc
}
