// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysAuthIdentityDao is the data access object for the table sys_auth_identity.
type SysAuthIdentityDao struct {
	table    string                 // table is the underlying table name of the DAO.
	group    string                 // group is the database configuration group name of the current DAO.
	columns  SysAuthIdentityColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler     // handlers for customized model modification.
}

// SysAuthIdentityColumns defines and stores column names for the table sys_auth_identity.
type SysAuthIdentityColumns struct {
	Id               string // Identity binding record ID
	TenantId         string // Host tenant scope for the binding
	UserId           string // Host user ID
	ProviderKey      string // Stable provider instance key
	ProviderType     string // Provider type, for example oidc or wecom
	Subject          string // Provider subject identifier
	UnionId          string // Provider union ID if available
	OpenId           string // Provider open ID if available
	ExternalTenantId string // External tenant/corp ID if available
	ExternalDeptIds  string // External department identifiers
	Email            string // External email address
	EmailVerified    string // Email verified flag: 0=false, 1=true
	Mobile           string // External mobile phone number
	DisplayName      string // External display name
	Avatar           string // External avatar URL
	RawProfile       string // Original non-secret external profile
	LastLoginAt      string // Last successful external login time
	BoundAt          string // Binding creation time
	CreatedAt        string // Creation time
	UpdatedAt        string // Update time
	DeletedAt        string // Deletion time
}

// sysAuthIdentityColumns holds the columns for the table sys_auth_identity.
var sysAuthIdentityColumns = SysAuthIdentityColumns{
	Id:               "id",
	TenantId:         "tenant_id",
	UserId:           "user_id",
	ProviderKey:      "provider_key",
	ProviderType:     "provider_type",
	Subject:          "subject",
	UnionId:          "union_id",
	OpenId:           "open_id",
	ExternalTenantId: "external_tenant_id",
	ExternalDeptIds:  "external_dept_ids",
	Email:            "email",
	EmailVerified:    "email_verified",
	Mobile:           "mobile",
	DisplayName:      "display_name",
	Avatar:           "avatar",
	RawProfile:       "raw_profile",
	LastLoginAt:      "last_login_at",
	BoundAt:          "bound_at",
	CreatedAt:        "created_at",
	UpdatedAt:        "updated_at",
	DeletedAt:        "deleted_at",
}

// NewSysAuthIdentityDao creates and returns a new DAO object for table data access.
func NewSysAuthIdentityDao(handlers ...gdb.ModelHandler) *SysAuthIdentityDao {
	return &SysAuthIdentityDao{
		group:    "default",
		table:    "sys_auth_identity",
		columns:  sysAuthIdentityColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysAuthIdentityDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysAuthIdentityDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysAuthIdentityDao) Columns() SysAuthIdentityColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysAuthIdentityDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysAuthIdentityDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysAuthIdentityDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
