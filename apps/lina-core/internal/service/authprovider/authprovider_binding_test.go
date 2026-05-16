// This file tests external identity binding behavior.

package authprovider

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	_ "lina-core/pkg/dbdriver"

	"lina-core/internal/dao"
	"lina-core/internal/model/do"
	"lina-core/internal/service/plugin"
	"lina-core/pkg/pluginhost"
)

// TestBindCurrentUserIdentityRejectsDuplicateSubject verifies one external
// subject cannot be bound to multiple host users.
func TestBindCurrentUserIdentityRejectsDuplicateSubject(t *testing.T) {
	ctx := context.Background()
	setupAuthProviderSQLite(t, ctx)
	insertAuthProviderTestUsers(t, ctx)
	_, err := dao.SysAuthProvider.Ctx(ctx).Data(do.SysAuthProvider{
		ProviderKey:  "google",
		PluginId:     "auth-oidc",
		ProviderType: "oidc",
		Name:         "Google",
		Enabled:      ProviderStatusEnabled,
	}).Insert()
	if err != nil {
		t.Fatalf("insert provider: %v", err)
	}
	svc := New(nil, testAuthProviderPluginService{
		enabled: true,
		providers: []plugin.AuthProviderRegistration{{
			PluginID: "auth-oidc",
			Provider: testAuthProvider{
				key:          "oidc",
				providerType: "oidc",
			},
		}},
	}, nil)
	identity := &pluginhost.ExternalIdentity{Subject: "subject-1", Email: "subject@example.com"}
	if err = svc.BindCurrentUserIdentity(ctx, 1, "google", identity); err != nil {
		t.Fatalf("bind first user: %v", err)
	}
	err = svc.BindCurrentUserIdentity(ctx, 2, "google", identity)
	if err == nil {
		t.Fatal("expected duplicate subject bind to fail")
	}
}

// setupAuthProviderSQLite prepares minimal auth-provider tables.
func setupAuthProviderSQLite(t *testing.T, ctx context.Context) {
	t.Helper()

	link := "sqlite::@file(" + t.TempDir() + "/authprovider-binding.db)"
	originalConfig := gdb.GetAllConfig()
	if err := gdb.SetConfig(gdb.Config{gdb.DefaultGroupName: gdb.ConfigGroup{{Link: link}}}); err != nil {
		t.Fatalf("configure sqlite database: %v", err)
	}
	db := g.DB()
	t.Cleanup(func() {
		if closeErr := db.Close(ctx); closeErr != nil {
			t.Errorf("close sqlite database: %v", closeErr)
		}
		if err := gdb.SetConfig(originalConfig); err != nil {
			t.Errorf("restore database config: %v", err)
		}
	})
	statements := []string{
		`CREATE TABLE sys_user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL DEFAULT 0,
			username TEXT NOT NULL,
			password TEXT NOT NULL DEFAULT '',
			nickname TEXT NOT NULL DEFAULT '',
			email TEXT NOT NULL DEFAULT '',
			phone TEXT NOT NULL DEFAULT '',
			sex INTEGER NOT NULL DEFAULT 0,
			avatar TEXT NOT NULL DEFAULT '',
			status INTEGER NOT NULL DEFAULT 1,
			remark TEXT NOT NULL DEFAULT '',
			login_date TIMESTAMP NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`,
		`CREATE TABLE sys_auth_provider (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			provider_key TEXT NOT NULL,
			plugin_id TEXT NOT NULL,
			provider_type TEXT NOT NULL,
			name TEXT NOT NULL,
			icon TEXT NOT NULL DEFAULT '',
			enabled INTEGER NOT NULL DEFAULT 0,
			sort INTEGER NOT NULL DEFAULT 0,
			config_mode TEXT NOT NULL DEFAULT 'static_config',
			config_json TEXT NOT NULL DEFAULT '{}',
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`,
		`CREATE TABLE sys_auth_identity (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tenant_id INTEGER NOT NULL DEFAULT 0,
			user_id INTEGER NOT NULL,
			provider_key TEXT NOT NULL,
			provider_type TEXT NOT NULL,
			subject TEXT NOT NULL,
			union_id TEXT NOT NULL DEFAULT '',
			open_id TEXT NOT NULL DEFAULT '',
			external_tenant_id TEXT NOT NULL DEFAULT '',
			external_dept_ids TEXT NOT NULL DEFAULT '[]',
			email TEXT NOT NULL DEFAULT '',
			email_verified INTEGER NOT NULL DEFAULT 0,
			mobile TEXT NOT NULL DEFAULT '',
			display_name TEXT NOT NULL DEFAULT '',
			avatar TEXT NOT NULL DEFAULT '',
			raw_profile TEXT NOT NULL DEFAULT '{}',
			last_login_at TIMESTAMP NULL,
			bound_at TIMESTAMP NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			deleted_at TIMESTAMP NULL
		);`,
	}
	for _, statement := range statements {
		if _, err := g.DB().Exec(ctx, statement); err != nil {
			t.Fatalf("create table: %v", err)
		}
	}
}

// insertAuthProviderTestUsers inserts two enabled test users.
func insertAuthProviderTestUsers(t *testing.T, ctx context.Context) {
	t.Helper()
	for _, username := range []string{"user-1", "user-2"} {
		_, err := dao.SysUser.Ctx(ctx).Data(do.SysUser{
			TenantId: 0,
			Username: username,
			Status:   1,
		}).Insert()
		if err != nil {
			t.Fatalf("insert user %s: %v", username, err)
		}
	}
}

// testAuthProviderPluginService supplies auth provider registrations to tests.
type testAuthProviderPluginService struct {
	enabled   bool
	providers []plugin.AuthProviderRegistration
}

// IsEnabled reports whether the plugin is enabled.
func (s testAuthProviderPluginService) IsEnabled(context.Context, string) bool {
	return s.enabled
}

// ListAuthProviders returns test provider registrations.
func (s testAuthProviderPluginService) ListAuthProviders(context.Context) ([]plugin.AuthProviderRegistration, error) {
	return append([]plugin.AuthProviderRegistration(nil), s.providers...), nil
}
