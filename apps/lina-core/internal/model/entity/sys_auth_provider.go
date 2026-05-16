// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuthProvider is the golang structure for table sys_auth_provider.
type SysAuthProvider struct {
	Id           int         `json:"id"           orm:"id"            description:"Provider record ID"`
	ProviderKey  string      `json:"providerKey"  orm:"provider_key"  description:"Stable provider instance key"`
	PluginId     string      `json:"pluginId"     orm:"plugin_id"     description:"Source plugin ID that owns this provider"`
	ProviderType string      `json:"providerType" orm:"provider_type" description:"Provider type, for example oidc or wecom"`
	Name         string      `json:"name"         orm:"name"          description:"Display name"`
	Icon         string      `json:"icon"         orm:"icon"          description:"Icon name or URL"`
	Enabled      int         `json:"enabled"      orm:"enabled"       description:"Enabled flag: 0=disabled, 1=enabled"`
	Sort         int         `json:"sort"         orm:"sort"          description:"Display sort order"`
	ConfigMode   string      `json:"configMode"   orm:"config_mode"   description:"Configuration mode"`
	ConfigJson   string      `json:"configJson"   orm:"config_json"   description:"Non-sensitive provider configuration or reference"`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"Creation time"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"Update time"`
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"Deletion time"`
}
