// This file defines public authentication provider list DTOs.

package v1

import "github.com/gogf/gf/v2/frame/g"

// ProviderListReq defines the request for public auth provider listing.
type ProviderListReq struct {
	g.Meta `path:"/auth/providers" method:"get" tags:"Authentication" summary:"List authentication providers" dc:"Lists enabled external authentication providers that can be displayed on the login page."`
}

// ProviderListRes defines the public auth provider list response.
type ProviderListRes struct {
	List []*AuthProviderEntity `json:"list" dc:"Authentication provider list"`
}

// AuthProviderEntity is a public-safe authentication provider projection.
type AuthProviderEntity struct {
	ProviderKey  string `json:"providerKey" dc:"Provider instance key" eg:"google"`
	ProviderType string `json:"providerType" dc:"Provider type" eg:"oidc"`
	Name         string `json:"name" dc:"Provider display name" eg:"Google"`
	Icon         string `json:"icon" dc:"Provider icon name or URL" eg:"brand-google"`
	Sort         int    `json:"sort" dc:"Display sort order" eg:"10"`
}

// ProviderAuthorizeReq defines the external provider authorization request.
type ProviderAuthorizeReq struct {
	g.Meta      `path:"/auth/providers/{providerKey}/authorize" method:"get" tags:"Authentication" summary:"Authorize authentication provider" dc:"Creates authorization state and returns an external provider redirect URL."`
	ProviderKey string `json:"providerKey" v:"required#validation.auth.provider.providerKey.required" dc:"Provider instance key" eg:"google"`
	Purpose     string `json:"purpose" dc:"Authorization purpose: login or bind" eg:"login"`
	RedirectUri string `json:"redirectUri" dc:"Frontend return URI after host callback" eg:"/auth/login"`
}

// ProviderAuthorizeRes defines the external provider authorization response.
type ProviderAuthorizeRes struct {
	RedirectUrl string `json:"redirectUrl" dc:"External provider redirect URL" eg:"https://accounts.google.com/o/oauth2/v2/auth?..."`
	State       string `json:"state" dc:"Host authorization state" eg:"state_123"`
}

// ProviderCallbackReq defines the GET callback request for external providers.
type ProviderCallbackReq struct {
	g.Meta      `path:"/auth/providers/{providerKey}/callback" method:"get" tags:"Authentication" summary:"Authentication provider callback" dc:"Handles an external provider authorization callback."`
	ProviderKey string `json:"providerKey" v:"required#validation.auth.provider.providerKey.required" dc:"Provider instance key" eg:"google"`
	State       string `json:"state" v:"required#validation.auth.provider.state.required" dc:"Authorization state" eg:"state_123"`
	Code        string `json:"code" dc:"Authorization code" eg:"code_123"`
	Error       string `json:"error" dc:"Provider error" eg:"access_denied"`
}

// ProviderCallbackRes defines the external provider callback login response.
type ProviderCallbackRes struct {
	AccessToken  string               `json:"accessToken" dc:"JWT token. Empty when tenant selection is required." eg:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string               `json:"refreshToken" dc:"JWT refresh token. Empty when tenant selection is required." eg:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	PreToken     string               `json:"preToken" dc:"Short-lived pre-login token when tenant selection is required." eg:"pre_8f4f..."`
	Tenants      []*LoginTenantEntity `json:"tenants" dc:"Tenant candidates when tenant selection is required." eg:"[]"`
	RedirectUri  string               `json:"redirectUri" dc:"Frontend return URI after host callback" eg:"/auth/login"`
}

// ProviderCallbackPostReq defines the POST callback request for external providers.
type ProviderCallbackPostReq struct {
	g.Meta      `path:"/auth/providers/{providerKey}/callback" method:"post" tags:"Authentication" summary:"Authentication provider POST callback" dc:"Handles an external provider POST authorization callback."`
	ProviderKey string `json:"providerKey" v:"required#validation.auth.provider.providerKey.required" dc:"Provider instance key" eg:"google"`
	State       string `json:"state" v:"required#validation.auth.provider.state.required" dc:"Authorization state" eg:"state_123"`
	Code        string `json:"code" dc:"Authorization code" eg:"code_123"`
	Error       string `json:"error" dc:"Provider error" eg:"access_denied"`
}

// ProviderCallbackPostRes defines the POST callback login response.
type ProviderCallbackPostRes = ProviderCallbackRes
