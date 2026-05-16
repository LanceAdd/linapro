// This file defines current-user external authentication identity DTOs.

package v1

import "github.com/gogf/gf/v2/frame/g"

// AuthIdentityListReq defines the request for current-user auth identity listing.
type AuthIdentityListReq struct {
	g.Meta `path:"/user/auth-identities" method:"get" tags:"User" summary:"List current user auth identities" dc:"Lists external authentication identities bound to the current user."`
}

// AuthIdentityListRes defines the current-user auth identity list response.
type AuthIdentityListRes struct {
	List []*AuthIdentityEntity `json:"list" dc:"External authentication identity list"`
}

// AuthIdentityUnbindReq defines the request for current-user auth identity unbind.
type AuthIdentityUnbindReq struct {
	g.Meta      `path:"/user/auth-identities/{providerKey}" method:"delete" tags:"User" summary:"Unbind current user auth identity" dc:"Deletes the current user's binding for one external authentication provider."`
	ProviderKey string `json:"providerKey" v:"required#validation.auth.identity.providerKey.required" dc:"Provider instance key" eg:"google"`
}

// AuthIdentityUnbindRes defines the current-user auth identity unbind response.
type AuthIdentityUnbindRes struct{}

// AuthIdentityBindReq defines the request for starting current-user auth identity binding.
type AuthIdentityBindReq struct {
	g.Meta      `path:"/user/auth-identities/{providerKey}/bind" method:"post" tags:"User" summary:"Bind current user auth identity" dc:"Creates authorization state for binding one external authentication provider to the current user."`
	ProviderKey string `json:"providerKey" v:"required#validation.auth.identity.providerKey.required" dc:"Provider instance key" eg:"google"`
	RedirectUri string `json:"redirectUri" dc:"Frontend return URI after host callback" eg:"/profile/security"`
}

// AuthIdentityBindRes defines the current-user auth identity bind response.
type AuthIdentityBindRes struct {
	RedirectUrl string `json:"redirectUrl" dc:"External provider redirect URL" eg:"https://accounts.google.com/o/oauth2/v2/auth?..."`
	State       string `json:"state" dc:"Host authorization state" eg:"state_123"`
}

// AuthIdentityEntity is a current-user auth identity projection.
type AuthIdentityEntity struct {
	ProviderKey      string `json:"providerKey" dc:"Provider instance key" eg:"google"`
	ProviderType     string `json:"providerType" dc:"Provider type" eg:"oidc"`
	Subject          string `json:"subject" dc:"External subject identifier" eg:"123456"`
	ExternalTenantId string `json:"externalTenantId" dc:"External tenant or corp identifier" eg:"contoso"`
	Email            string `json:"email" dc:"External email" eg:"user@example.com"`
	Mobile           string `json:"mobile" dc:"External mobile phone number" eg:"13800000000"`
	DisplayName      string `json:"displayName" dc:"External display name" eg:"Alex"`
	Avatar           string `json:"avatar" dc:"External avatar URL" eg:"https://example.com/avatar.png"`
	LastLoginAt      string `json:"lastLoginAt" dc:"Last successful external login time" eg:"2026-05-16 12:00:00"`
	BoundAt          string `json:"boundAt" dc:"Binding creation time" eg:"2026-05-16 12:00:00"`
}
