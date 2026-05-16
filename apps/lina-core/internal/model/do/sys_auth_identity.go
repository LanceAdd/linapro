// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuthIdentity is the golang structure of table sys_auth_identity for DAO operations like Where/Data.
type SysAuthIdentity struct {
	g.Meta           `orm:"table:sys_auth_identity, do:true"`
	Id               any         // Identity binding record ID
	TenantId         any         // Host tenant scope for the binding
	UserId           any         // Host user ID
	ProviderKey      any         // Stable provider instance key
	ProviderType     any         // Provider type, for example oidc or wecom
	Subject          any         // Provider subject identifier
	UnionId          any         // Provider union ID if available
	OpenId           any         // Provider open ID if available
	ExternalTenantId any         // External tenant/corp ID if available
	ExternalDeptIds  any         // External department identifiers
	Email            any         // External email address
	EmailVerified    any         // Email verified flag: 0=false, 1=true
	Mobile           any         // External mobile phone number
	DisplayName      any         // External display name
	Avatar           any         // External avatar URL
	RawProfile       any         // Original non-secret external profile
	LastLoginAt      *gtime.Time // Last successful external login time
	BoundAt          *gtime.Time // Binding creation time
	CreatedAt        *gtime.Time // Creation time
	UpdatedAt        *gtime.Time // Update time
	DeletedAt        *gtime.Time // Deletion time
}
