// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuthProvider is the golang structure of table sys_auth_provider for DAO operations like Where/Data.
type SysAuthProvider struct {
	g.Meta       `orm:"table:sys_auth_provider, do:true"`
	Id           any         // Provider record ID
	ProviderKey  any         // Stable provider instance key
	PluginId     any         // Source plugin ID that owns this provider
	ProviderType any         // Provider type, for example oidc or wecom
	Name         any         // Display name
	Icon         any         // Icon name or URL
	Enabled      any         // Enabled flag: 0=disabled, 1=enabled
	Sort         any         // Display sort order
	ConfigMode   any         // Configuration mode
	ConfigJson   any         // Non-sensitive provider configuration or reference
	CreatedAt    *gtime.Time // Creation time
	UpdatedAt    *gtime.Time // Update time
	DeletedAt    *gtime.Time // Deletion time
}
