// This file contains shared host-owned login completion helpers.

package auth

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/mssola/useragent"

	"lina-core/internal/dao"
	"lina-core/internal/model/do"
	"lina-core/internal/model/entity"
	pluginsvc "lina-core/internal/service/plugin"
	"lina-core/pkg/bizerr"
	"lina-core/pkg/logger"
	"lina-core/pkg/pluginhost"
)

// loginAuditContext carries client and provider metadata for login hooks.
type loginAuditContext struct {
	Username string
	IP       string
	Browser  string
	OS       string
	Method   string
	Provider string
}

// completeLogin signs tokens, persists session state, and dispatches success hooks.
func (s *serviceImpl) completeLogin(
	ctx context.Context,
	user *entity.SysUser,
	tenantID int,
	audit loginAuditContext,
) (*LoginOutput, error) {
	accessToken, refreshToken, tokenID, err := s.generateTokenPair(ctx, user, tenantID)
	if err != nil {
		return nil, err
	}
	if _, err = dao.SysUser.Ctx(ctx).
		Where(do.SysUser{Id: user.Id}).
		Data(do.SysUser{LoginDate: gtime.Now()}).
		Update(); err != nil {
		return nil, bizerr.WrapCode(err, CodeAuthLoginStateUpdateFailed)
	}
	if err = s.createSession(ctx, user, tenantID, tokenID); err != nil {
		logger.Warningf(ctx, "create online session failed tokenId=%s err=%v", tokenID, err)
	}
	s.dispatchLoginSucceeded(ctx, audit)
	return &LoginOutput{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

// dispatchLoginSucceeded sends a best-effort login success hook.
func (s *serviceImpl) dispatchLoginSucceeded(ctx context.Context, audit loginAuditContext) {
	if s == nil || s.pluginSvc == nil {
		return
	}
	if audit.IP == "" && audit.Browser == "" && audit.OS == "" {
		audit.IP, audit.Browser, audit.OS = requestClientInfo(ctx)
	}
	if err := s.pluginSvc.HandleAuthLoginSucceeded(ctx, pluginsvc.AuthLoginSucceededInput{
		UserName:   audit.Username,
		Status:     authLoginStatusSuccess,
		Ip:         audit.IP,
		ClientType: "web",
		Browser:    audit.Browser,
		Os:         audit.OS,
		Message:    pluginsvc.AuthEventMessageLoginSuccessful,
		Reason:     pluginhost.AuthHookReasonLoginSuccessful,
		Method:     audit.Method,
		Provider:   audit.Provider,
	}); err != nil {
		logger.Warningf(ctx, "plugin login succeeded hook failed: %v", err)
	}
}

// dispatchLoginFailed sends a best-effort login failure hook.
func (s *serviceImpl) dispatchLoginFailed(ctx context.Context, username string, msg string, reason string) {
	if s == nil || s.pluginSvc == nil {
		return
	}
	ip, browser, osName := requestClientInfo(ctx)
	if hookErr := s.pluginSvc.HandleAuthLoginFailed(ctx, pluginsvc.AuthLoginSucceededInput{
		UserName:   username,
		Status:     authLoginStatusFail,
		Ip:         ip,
		ClientType: "web",
		Browser:    browser,
		Os:         osName,
		Message:    msg,
		Reason:     reason,
	}); hookErr != nil {
		logger.Warningf(ctx, "plugin login failed hook failed: %v", hookErr)
	}
}

// requestClientInfo extracts IP, browser, and OS from the current request.
func requestClientInfo(ctx context.Context) (string, string, string) {
	var ip, browser, osName string
	if r := g.RequestFromCtx(ctx); r != nil {
		ip = r.GetClientIp()
		ua := useragent.New(r.GetHeader("User-Agent"))
		browserName, browserVersion := ua.Browser()
		browser = browserName + " " + browserVersion
		osName = ua.OS()
	}
	return ip, browser, osName
}
