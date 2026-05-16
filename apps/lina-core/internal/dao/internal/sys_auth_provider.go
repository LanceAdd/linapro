// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysAuthProviderDao is the data access object for the table sys_auth_provider.
type SysAuthProviderDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SysAuthProviderColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SysAuthProviderColumns defines and stores column names for the table sys_auth_provider.
type SysAuthProviderColumns struct {
	Id           string // Provider record ID
	ProviderKey  string // Stable provider instance key
	PluginId     string // Source plugin ID that owns this provider
	ProviderType string // Provider type, for example oidc or wecom
	Name         string // Display name
	Icon         string // Icon name or URL
	Enabled      string // Enabled flag: 0=disabled, 1=enabled
	Sort         string // Display sort order
	ConfigMode   string // Configuration mode
	ConfigJson   string // Non-sensitive provider configuration or reference
	CreatedAt    string // Creation time
	UpdatedAt    string // Update time
	DeletedAt    string // Deletion time
}

// sysAuthProviderColumns holds the columns for the table sys_auth_provider.
var sysAuthProviderColumns = SysAuthProviderColumns{
	Id:           "id",
	ProviderKey:  "provider_key",
	PluginId:     "plugin_id",
	ProviderType: "provider_type",
	Name:         "name",
	Icon:         "icon",
	Enabled:      "enabled",
	Sort:         "sort",
	ConfigMode:   "config_mode",
	ConfigJson:   "config_json",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
	DeletedAt:    "deleted_at",
}

// NewSysAuthProviderDao creates and returns a new DAO object for table data access.
func NewSysAuthProviderDao(handlers ...gdb.ModelHandler) *SysAuthProviderDao {
	return &SysAuthProviderDao{
		group:    "default",
		table:    "sys_auth_provider",
		columns:  sysAuthProviderColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysAuthProviderDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysAuthProviderDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysAuthProviderDao) Columns() SysAuthProviderColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysAuthProviderDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysAuthProviderDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysAuthProviderDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
