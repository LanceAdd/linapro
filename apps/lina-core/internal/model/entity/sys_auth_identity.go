// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuthIdentity is the golang structure for table sys_auth_identity.
type SysAuthIdentity struct {
	Id               int         `json:"id"               orm:"id"                 description:"Identity binding record ID"`
	TenantId         int         `json:"tenantId"         orm:"tenant_id"          description:"Host tenant scope for the binding"`
	UserId           int         `json:"userId"           orm:"user_id"            description:"Host user ID"`
	ProviderKey      string      `json:"providerKey"      orm:"provider_key"       description:"Stable provider instance key"`
	ProviderType     string      `json:"providerType"     orm:"provider_type"      description:"Provider type, for example oidc or wecom"`
	Subject          string      `json:"subject"          orm:"subject"            description:"Provider subject identifier"`
	UnionId          string      `json:"unionId"          orm:"union_id"           description:"Provider union ID if available"`
	OpenId           string      `json:"openId"           orm:"open_id"            description:"Provider open ID if available"`
	ExternalTenantId string      `json:"externalTenantId" orm:"external_tenant_id" description:"External tenant/corp ID if available"`
	ExternalDeptIds  string      `json:"externalDeptIds"  orm:"external_dept_ids"  description:"External department identifiers"`
	Email            string      `json:"email"            orm:"email"              description:"External email address"`
	EmailVerified    int         `json:"emailVerified"    orm:"email_verified"     description:"Email verified flag: 0=false, 1=true"`
	Mobile           string      `json:"mobile"           orm:"mobile"             description:"External mobile phone number"`
	DisplayName      string      `json:"displayName"      orm:"display_name"       description:"External display name"`
	Avatar           string      `json:"avatar"           orm:"avatar"             description:"External avatar URL"`
	RawProfile       string      `json:"rawProfile"       orm:"raw_profile"        description:"Original non-secret external profile"`
	LastLoginAt      *gtime.Time `json:"lastLoginAt"      orm:"last_login_at"      description:"Last successful external login time"`
	BoundAt          *gtime.Time `json:"boundAt"          orm:"bound_at"           description:"Binding creation time"`
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"         description:"Creation time"`
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"         description:"Update time"`
	DeletedAt        *gtime.Time `json:"deletedAt"        orm:"deleted_at"         description:"Deletion time"`
}
