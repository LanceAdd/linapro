-- 013: Pluggable Authentication Providers
-- 013: 可插拔认证提供方

-- Purpose: Stores host-governed external authentication provider metadata.
-- Secrets are intentionally excluded and must be resolved from static
-- configuration or environment variables by the provider plugin.
CREATE TABLE IF NOT EXISTS sys_auth_provider (
    "id"            INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "provider_key"  VARCHAR(128) NOT NULL,
    "plugin_id"     VARCHAR(128) NOT NULL,
    "provider_type" VARCHAR(64) NOT NULL,
    "name"          VARCHAR(128) NOT NULL,
    "icon"          VARCHAR(256) NOT NULL DEFAULT '',
    "enabled"       SMALLINT NOT NULL DEFAULT 0,
    "sort"          INT NOT NULL DEFAULT 0,
    "config_mode"   VARCHAR(32) NOT NULL DEFAULT 'static_config',
    "config_json"   JSONB NOT NULL DEFAULT '{}'::jsonb,
    "created_at"    TIMESTAMP,
    "updated_at"    TIMESTAMP,
    "deleted_at"    TIMESTAMP
);

COMMENT ON TABLE sys_auth_provider IS 'External authentication provider metadata';
COMMENT ON COLUMN sys_auth_provider."id" IS 'Provider record ID';
COMMENT ON COLUMN sys_auth_provider."provider_key" IS 'Stable provider instance key';
COMMENT ON COLUMN sys_auth_provider."plugin_id" IS 'Source plugin ID that owns this provider';
COMMENT ON COLUMN sys_auth_provider."provider_type" IS 'Provider type, for example oidc or wecom';
COMMENT ON COLUMN sys_auth_provider."name" IS 'Display name';
COMMENT ON COLUMN sys_auth_provider."icon" IS 'Icon name or URL';
COMMENT ON COLUMN sys_auth_provider."enabled" IS 'Enabled flag: 0=disabled, 1=enabled';
COMMENT ON COLUMN sys_auth_provider."sort" IS 'Display sort order';
COMMENT ON COLUMN sys_auth_provider."config_mode" IS 'Configuration mode';
COMMENT ON COLUMN sys_auth_provider."config_json" IS 'Non-sensitive provider configuration or reference';
COMMENT ON COLUMN sys_auth_provider."created_at" IS 'Creation time';
COMMENT ON COLUMN sys_auth_provider."updated_at" IS 'Update time';
COMMENT ON COLUMN sys_auth_provider."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_auth_provider_key
    ON sys_auth_provider ("provider_key")
    WHERE "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_auth_provider_enabled_sort
    ON sys_auth_provider ("enabled", "sort")
    WHERE "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_auth_provider_plugin
    ON sys_auth_provider ("plugin_id")
    WHERE "deleted_at" IS NULL;

-- Purpose: Stores explicit bindings between host users and external
-- identities. The host uses this table to map external callbacks to local
-- users and never auto-creates users from external identities in this
-- iteration.
CREATE TABLE IF NOT EXISTS sys_auth_identity (
    "id"                 INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "tenant_id"          INT NOT NULL DEFAULT 0,
    "user_id"            INT NOT NULL,
    "provider_key"       VARCHAR(128) NOT NULL,
    "provider_type"      VARCHAR(64) NOT NULL,
    "subject"            VARCHAR(256) NOT NULL,
    "union_id"           VARCHAR(256) NOT NULL DEFAULT '',
    "open_id"            VARCHAR(256) NOT NULL DEFAULT '',
    "external_tenant_id" VARCHAR(128) NOT NULL DEFAULT '',
    "external_dept_ids"  JSONB NOT NULL DEFAULT '[]'::jsonb,
    "email"              VARCHAR(128) NOT NULL DEFAULT '',
    "email_verified"     SMALLINT NOT NULL DEFAULT 0,
    "mobile"             VARCHAR(32) NOT NULL DEFAULT '',
    "display_name"       VARCHAR(128) NOT NULL DEFAULT '',
    "avatar"             VARCHAR(512) NOT NULL DEFAULT '',
    "raw_profile"        JSONB NOT NULL DEFAULT '{}'::jsonb,
    "last_login_at"      TIMESTAMP,
    "bound_at"           TIMESTAMP,
    "created_at"         TIMESTAMP,
    "updated_at"         TIMESTAMP,
    "deleted_at"         TIMESTAMP
);

COMMENT ON TABLE sys_auth_identity IS 'External authentication identity bindings';
COMMENT ON COLUMN sys_auth_identity."id" IS 'Identity binding record ID';
COMMENT ON COLUMN sys_auth_identity."tenant_id" IS 'Host tenant scope for the binding';
COMMENT ON COLUMN sys_auth_identity."user_id" IS 'Host user ID';
COMMENT ON COLUMN sys_auth_identity."provider_key" IS 'Stable provider instance key';
COMMENT ON COLUMN sys_auth_identity."provider_type" IS 'Provider type, for example oidc or wecom';
COMMENT ON COLUMN sys_auth_identity."subject" IS 'Provider subject identifier';
COMMENT ON COLUMN sys_auth_identity."union_id" IS 'Provider union ID if available';
COMMENT ON COLUMN sys_auth_identity."open_id" IS 'Provider open ID if available';
COMMENT ON COLUMN sys_auth_identity."external_tenant_id" IS 'External tenant/corp ID if available';
COMMENT ON COLUMN sys_auth_identity."external_dept_ids" IS 'External department identifiers';
COMMENT ON COLUMN sys_auth_identity."email" IS 'External email address';
COMMENT ON COLUMN sys_auth_identity."email_verified" IS 'Email verified flag: 0=false, 1=true';
COMMENT ON COLUMN sys_auth_identity."mobile" IS 'External mobile phone number';
COMMENT ON COLUMN sys_auth_identity."display_name" IS 'External display name';
COMMENT ON COLUMN sys_auth_identity."avatar" IS 'External avatar URL';
COMMENT ON COLUMN sys_auth_identity."raw_profile" IS 'Original non-secret external profile';
COMMENT ON COLUMN sys_auth_identity."last_login_at" IS 'Last successful external login time';
COMMENT ON COLUMN sys_auth_identity."bound_at" IS 'Binding creation time';
COMMENT ON COLUMN sys_auth_identity."created_at" IS 'Creation time';
COMMENT ON COLUMN sys_auth_identity."updated_at" IS 'Update time';
COMMENT ON COLUMN sys_auth_identity."deleted_at" IS 'Deletion time';

CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_auth_identity_provider_subject
    ON sys_auth_identity ("provider_key", "subject")
    WHERE "deleted_at" IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS uk_sys_auth_identity_user_provider
    ON sys_auth_identity ("user_id", "provider_key")
    WHERE "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_auth_identity_tenant_user
    ON sys_auth_identity ("tenant_id", "user_id")
    WHERE "deleted_at" IS NULL;
CREATE INDEX IF NOT EXISTS idx_sys_auth_identity_external_tenant
    ON sys_auth_identity ("provider_key", "external_tenant_id")
    WHERE "deleted_at" IS NULL;
