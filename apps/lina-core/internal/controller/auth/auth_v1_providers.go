// This file handles pluggable authentication provider endpoints.

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"

	v1 "lina-core/api/auth/v1"
	authsvc "lina-core/internal/service/auth"
	authprovidersvc "lina-core/internal/service/authprovider"
)

const (
	authProviderCallbackPurposeBind = "bind"
	authProviderFrontendCallback    = "/auth/providers/callback"
	authProviderDefaultLoginReturn  = "/auth/login"
	authProviderDefaultBindReturn   = "/profile"
)

// ProviderList lists enabled external authentication providers.
func (c *ControllerV1) ProviderList(ctx context.Context, _ *v1.ProviderListReq) (res *v1.ProviderListRes, err error) {
	items, err := c.authProviderSvc.ListPublicProviders(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.ProviderListRes{List: toAuthProviderEntities(items)}, nil
}

// ProviderAuthorize starts an external authentication authorization flow.
func (c *ControllerV1) ProviderAuthorize(ctx context.Context, req *v1.ProviderAuthorizeReq) (res *v1.ProviderAuthorizeRes, err error) {
	output, err := c.authProviderSvc.Authorize(ctx, authprovidersvc.AuthorizeInput{
		ProviderKey: req.ProviderKey,
		Purpose:     req.Purpose,
		RedirectURI: req.RedirectUri,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ProviderAuthorizeRes{
		RedirectUrl: output.RedirectURL,
		State:       output.State,
	}, nil
}

// ProviderCallback handles a GET external authentication provider callback.
func (c *ControllerV1) ProviderCallback(ctx context.Context, req *v1.ProviderCallbackReq) (res *v1.ProviderCallbackRes, err error) {
	output, err := c.handleProviderCallback(ctx, req.ProviderKey, req.State, map[string]string{
		"code":  req.Code,
		"error": req.Error,
	}, nil)
	if err != nil {
		return nil, err
	}
	redirectProviderCallback(ctx, req.ProviderKey, output)
	return output, nil
}

// ProviderCallbackPost handles a POST external authentication provider callback.
func (c *ControllerV1) ProviderCallbackPost(ctx context.Context, req *v1.ProviderCallbackPostReq) (res *v1.ProviderCallbackPostRes, err error) {
	output, err := c.handleProviderCallback(ctx, req.ProviderKey, req.State, nil, map[string]string{
		"code":  req.Code,
		"error": req.Error,
	})
	if err != nil {
		return nil, err
	}
	redirectProviderCallback(ctx, req.ProviderKey, output)
	return output, nil
}

// handleProviderCallback consumes state and completes either login or binding.
func (c *ControllerV1) handleProviderCallback(
	ctx context.Context,
	providerKey string,
	state string,
	query map[string]string,
	form map[string]string,
) (*v1.ProviderCallbackRes, error) {
	callback, err := c.authProviderSvc.Callback(ctx, authprovidersvc.CallbackInput{
		ProviderKey: providerKey,
		State:       state,
		Query:       query,
		Form:        form,
	})
	if err != nil {
		return nil, err
	}
	if callback.Purpose == authProviderCallbackPurposeBind {
		if err = c.authProviderSvc.BindCurrentUserIdentity(ctx, callback.UserID, callback.ProviderKey, callback.Identity); err != nil {
			return nil, err
		}
		return &v1.ProviderCallbackRes{RedirectUri: callback.RedirectURI}, nil
	}
	loginOutput, err := c.authSvc.LoginWithExternalIdentity(ctx, authsvc.ExternalLoginInput{
		ProviderKey:  callback.ProviderKey,
		ProviderType: callback.ProviderType,
		Identity:     callback.Identity,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ProviderCallbackRes{
		AccessToken:  loginOutput.AccessToken,
		RefreshToken: loginOutput.RefreshToken,
		PreToken:     loginOutput.PreToken,
		Tenants:      toLoginTenants(loginOutput.Tenants),
		RedirectUri:  callback.RedirectURI,
	}, nil
}

// redirectProviderCallback sends browser callbacks back to the frontend bridge.
func redirectProviderCallback(ctx context.Context, providerKey string, output *v1.ProviderCallbackRes) {
	request := g.RequestFromCtx(ctx)
	if request == nil || output == nil {
		return
	}
	target := buildProviderFrontendCallbackURL(providerKey, output)
	request.Response.RedirectTo(target)
	request.ExitAll()
}

// buildProviderFrontendCallbackURL builds a relative frontend callback URL.
func buildProviderFrontendCallbackURL(providerKey string, output *v1.ProviderCallbackRes) string {
	values := url.Values{}
	values.Set("providerKey", providerKey)
	values.Set("payload", encodeProviderCallbackPayload(output))
	if strings.TrimSpace(output.RedirectUri) != "" {
		values.Set("redirectUri", output.RedirectUri)
	}
	target := authProviderDefaultLoginReturn
	if output.AccessToken == "" && output.PreToken == "" {
		target = authProviderDefaultBindReturn
	}
	if strings.TrimSpace(output.RedirectUri) != "" {
		target = output.RedirectUri
	}
	values.Set("target", target)
	return authProviderFrontendCallback + "#" + authProviderFrontendCallback + "?" + values.Encode()
}

// encodeProviderCallbackPayload serializes callback output for frontend pickup.
func encodeProviderCallbackPayload(output *v1.ProviderCallbackRes) string {
	data, err := json.Marshal(output)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(data)
}

// toAuthProviderEntities converts service provider items into API DTOs.
func toAuthProviderEntities(items []authprovidersvc.ProviderItem) []*v1.AuthProviderEntity {
	if len(items) == 0 {
		return []*v1.AuthProviderEntity{}
	}
	out := make([]*v1.AuthProviderEntity, 0, len(items))
	for _, item := range items {
		out = append(out, &v1.AuthProviderEntity{
			ProviderKey:  item.ProviderKey,
			ProviderType: item.ProviderType,
			Name:         item.Name,
			Icon:         item.Icon,
			Sort:         item.Sort,
		})
	}
	return out
}
