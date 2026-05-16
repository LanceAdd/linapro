// This file tests external-identity login conversion.

package auth

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	_ "lina-core/pkg/dbdriver"

	"lina-core/internal/dao"
	"lina-core/internal/model/do"
	"lina-core/internal/service/kvcache"
	"lina-core/internal/service/orgcap"
	"lina-core/internal/service/session"
	"lina-core/pkg/bizerr"
	"lina-core/pkg/pluginhost"
)

const externalLoginSQLiteChildEnv = "LINA_EXTERNAL_LOGIN_SQLITE_CHILD"

// TestLoginWithExternalIdentityRejectsUnboundIdentity verifies external login
// never auto-creates users when no binding exists.
func TestLoginWithExternalIdentityRejectsUnboundIdentity(t *testing.T) {
	if runExternalLoginSQLiteChild(t) {
		return
	}
	ctx := context.Background()
	setupExternalLoginSQLite(t, ctx)
	svc := New(configTestService{}, nil, orgcap.New(nil), roleTestService{}, disabledTenantAuthTestService{}, session.NewDBStore(), kvcache.New())
	_, err := svc.LoginWithExternalIdentity(ctx, ExternalLoginInput{
		ProviderKey:  "google",
		ProviderType: "oidc",
		Identity:     &pluginhost.ExternalIdentity{Subject: "missing-subject"},
	})
	assertBizCode(t, err, CodeAuthExternalIdentityNotBound.RuntimeCode())
}

// TestLoginWithExternalIdentityRejectsDisabledUser verifies disabled bound users
// cannot log in through external providers.
func TestLoginWithExternalIdentityRejectsDisabledUser(t *testing.T) {
	if runExternalLoginSQLiteChild(t) {
		return
	}
	ctx := context.Background()
	setupExternalLoginSQLite(t, ctx)
	insertExternalLoginUserAndBinding(t, ctx, 0)
	svc := New(configTestService{}, nil, orgcap.New(nil), roleTestService{}, disabledTenantAuthTestService{}, session.NewDBStore(), kvcache.New())
	_, err := svc.LoginWithExternalIdentity(ctx, ExternalLoginInput{
		ProviderKey:  "google",
		ProviderType: "oidc",
		Identity:     &pluginhost.ExternalIdentity{Subject: "subject-1"},
	})
	assertBizCode(t, err, CodeAuthUserDisabled.RuntimeCode())
}

// TestLoginWithExternalIdentityIssuesHostToken verifies successful external
// login reuses host token and session issuance.
func TestLoginWithExternalIdentityIssuesHostToken(t *testing.T) {
	if runExternalLoginSQLiteChild(t) {
		return
	}
	ctx := context.Background()
	setupExternalLoginSQLite(t, ctx)
	insertExternalLoginUserAndBinding(t, ctx, 1)
	svc := New(configTestService{}, nil, orgcap.New(nil), roleTestService{}, disabledTenantAuthTestService{}, session.NewDBStore(), kvcache.New())
	output, err := svc.LoginWithExternalIdentity(ctx, ExternalLoginInput{
		ProviderKey:  "google",
		ProviderType: "oidc",
		Identity:     &pluginhost.ExternalIdentity{Subject: "subject-1"},
	})
	if err != nil {
		t.Fatalf("external login: %v", err)
	}
	if output == nil || output.AccessToken == "" || output.RefreshToken == "" {
		t.Fatalf("expected host tokens, got %+v", output)
	}
	count, err := dao.SysOnlineSession.Ctx(ctx).Count()
	if err != nil {
		t.Fatalf("count sessions: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected one online session, got %d", count)
	}
}

// setupExternalLoginSQLite prepares minimal auth tables for external login tests.
func setupExternalLoginSQLite(t *testing.T, ctx context.Context) {
	t.Helper()

	link := "sqlite::@file(" + filepath.ToSlash(t.TempDir()) + "/external-login.db)"
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
		`CREATE TABLE sys_online_session (
			tenant_id INTEGER NOT NULL DEFAULT 0,
			token_id TEXT NOT NULL PRIMARY KEY,
			user_id INTEGER NOT NULL DEFAULT 0,
			username TEXT NOT NULL DEFAULT '',
			dept_name TEXT NOT NULL DEFAULT '',
			ip TEXT NOT NULL DEFAULT '',
			browser TEXT NOT NULL DEFAULT '',
			os TEXT NOT NULL DEFAULT '',
			login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_active_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE sys_kv_cache (
			owner_type TEXT NOT NULL,
			cache_key TEXT NOT NULL,
			value_kind TEXT NOT NULL DEFAULT 'string',
			value TEXT NOT NULL DEFAULT '',
			int_value INTEGER NOT NULL DEFAULT 0,
			expire_at TIMESTAMP NULL,
			created_at TIMESTAMP NULL,
			updated_at TIMESTAMP NULL,
			PRIMARY KEY (owner_type, cache_key)
		);`,
	}
	for _, statement := range statements {
		if _, err := g.DB().Exec(ctx, statement); err != nil {
			t.Fatalf("create table: %v", err)
		}
	}
}

// runExternalLoginSQLiteChild isolates GoFrame's process-global DB config from
// the rest of the auth package tests.
func runExternalLoginSQLiteChild(t *testing.T) bool {
	t.Helper()
	if os.Getenv(externalLoginSQLiteChildEnv) == t.Name() {
		return false
	}
	cmd := exec.Command(os.Args[0], "-test.run", "^"+regexp.QuoteMeta(t.Name())+"$", "-test.v")
	cmd.Env = append(os.Environ(), externalLoginSQLiteChildEnv+"="+t.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("external login sqlite child failed: %v\n%s", err, string(output))
	}
	return true
}

// insertExternalLoginUserAndBinding inserts one user and its external identity binding.
func insertExternalLoginUserAndBinding(t *testing.T, ctx context.Context, status int) {
	t.Helper()
	_, err := dao.SysUser.Ctx(ctx).Data(do.SysUser{
		TenantId: 0,
		Username: "external-user",
		Status:   status,
	}).Insert()
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	_, err = dao.SysAuthIdentity.Ctx(ctx).Data(do.SysAuthIdentity{
		TenantId:     0,
		UserId:       1,
		ProviderKey:  "google",
		ProviderType: "oidc",
		Subject:      "subject-1",
	}).Insert()
	if err != nil {
		t.Fatalf("insert auth identity: %v", err)
	}
}

// assertBizCode verifies a business error code.
func assertBizCode(t *testing.T, err error, code string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected biz error %s", code)
	}
	bizCode, ok := bizerr.As(err)
	if !ok || bizCode.RuntimeCode() != code {
		t.Fatalf("expected biz code %s, got %v", code, err)
	}
}
