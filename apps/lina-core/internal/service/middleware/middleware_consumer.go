// This file implements the lightweight consumer service-surface middleware.
// It creates a consumer request context and resolves the tenant boundary while
// leaving login, session, token, and authorization behavior to each plugin.

package middleware

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"

	"lina-core/internal/model"
	"lina-core/pkg/bizerr"
	"lina-core/pkg/pluginhost"
	pkgtenantcap "lina-core/pkg/tenantcap"
)

// ConsumerCtx injects the base consumer context and request locale.
func (s *serviceImpl) ConsumerCtx(r *ghttp.Request) {
	if r == nil {
		return
	}
	customCtx := &model.Context{}
	s.bizCtxSvc.Init(r, customCtx)

	locale := s.i18nSvc.ResolveRequestLocale(r)
	r.SetCtx(gi18n.WithLanguage(r.Context(), locale))
	s.bizCtxSvc.SetLocale(r.Context(), locale)
	r.Response.Header().Set("Content-Language", locale)

	pluginID := pluginhost.SourcePluginIDFromRequest(r)
	if pluginID == "" {
		pluginID = pluginhost.ConsumerPluginIDFromPath(r.URL.Path)
	}
	if pluginID == "" {
		r.SetError(bizerr.NewCode(CodeMiddlewareConsumerPluginRequired))
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}

	consumerCtx := &pluginhost.ConsumerContext{
		PluginID:    pluginID,
		Locale:      locale,
		AnonymousID: strings.TrimSpace(r.GetHeader(pluginhost.ConsumerHeaderAnonymousID)),
		DeviceID:    strings.TrimSpace(r.GetHeader(pluginhost.ConsumerHeaderDeviceID)),
		Channel:     strings.TrimSpace(r.GetHeader(pluginhost.ConsumerHeaderChannel)),
	}
	pluginhost.SetConsumerContext(r, consumerCtx)
	r.Middleware.Next()
}

// ConsumerTenant resolves the consumer tenant boundary before authentication.
func (s *serviceImpl) ConsumerTenant(r *ghttp.Request) {
	if r == nil {
		return
	}
	consumerCtx, _ := pluginhost.ConsumerContextFromRequest(r)
	if consumerCtx == nil {
		r.SetError(bizerr.NewCode(CodeMiddlewareConsumerPluginRequired))
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	tenantID, matched, err := pluginhost.ResolveConsumerTenantID(r)
	if err != nil {
		r.SetError(bizerr.WrapCode(err, CodeMiddlewareConsumerTenantRequired))
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	if matched {
		consumerCtx.TenantID = tenantID
		consumerCtx.TenantResolved = true
		s.bizCtxSvc.SetTenant(r.Context(), tenantID)
		r.Middleware.Next()
		return
	}
	if s == nil || s.bizCtxSvc == nil {
		r.SetError(bizerr.NewCode(CodeMiddlewareConsumerTenantRequired))
		r.Response.WriteStatus(http.StatusUnauthorized)
		return
	}
	if s.tenantSvc == nil || !s.tenantSvc.Enabled(r.Context()) {
		consumerCtx.TenantID = int(pkgtenantcap.PLATFORM)
		consumerCtx.TenantResolved = true
		s.bizCtxSvc.SetTenant(r.Context(), int(pkgtenantcap.PLATFORM))
		r.Middleware.Next()
		return
	}

	r.SetError(bizerr.NewCode(CodeMiddlewareConsumerTenantRequired))
	r.Response.WriteStatus(http.StatusUnauthorized)
	return
}
