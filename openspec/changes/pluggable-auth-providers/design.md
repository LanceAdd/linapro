## Context

当前认证主链路由宿主提供用户名/密码登录，成功后由宿主签发 JWT/Refresh Token、创建 `sys_online_session`、发布认证生命周期事件。插件机制已经可以参与登录后的 Hook、前端 Slot 和宿主服务扩展，但尚未定义“插件参与登录前身份认证”的稳定契约。

外部身份源存在明显差异：OIDC 可覆盖 Google、Microsoft、Azure AD/Entra ID、Keycloak、Okta、Authing 等标准化场景；企业微信、钉钉、飞书、QQ、个人微信、SAML 等需要各自协议适配。宿主不能把每种协议都硬编码进认证服务，否则会让核心认证边界、密钥配置、回调处理和审计逻辑持续膨胀。

本次设计把宿主定位为认证中心，把插件定位为身份源适配器。插件只负责完成外部协议交互并返回标准化外部身份；宿主负责本地用户映射、租户与用户状态校验、Token 签发、会话创建、登录事件和审计。

## Goals / Non-Goals

**Goals:**

- 建立宿主受控、插件可扩展的外部登录提供方模型。
- 首期支持通用 OIDC 和企业微信登录。
- 支持外部身份与既有本地用户绑定，并通过绑定关系完成登录。
- 保持 JWT、Refresh Token、在线会话、用户状态、租户边界、权限加载和登录事件由宿主统一处理。
- 让登录页和个人中心根据启用状态动态展示外部登录和绑定入口。
- 在单机和集群部署下都能可靠处理 OAuth/OIDC 授权临时状态。
- 明确 i18n、缓存一致性、数据权限和测试边界。

**Non-Goals:**

- 首期不支持自动创建本地用户。
- 首期不支持外部身份自动分配租户、部门、岗位或角色。
- 首期不提供后端页面编辑敏感密钥，密钥通过 `config.yaml` 或环境变量管理。
- 首期不实现钉钉、飞书、QQ、个人微信、SAML 专用插件，但保留扩展契约。
- 首期不替代现有用户名/密码登录能力。

## Decisions

### 1. 宿主认证中心 + 插件身份源适配

宿主新增认证提供方注册表和外部身份绑定表。源码插件通过宿主发布的认证提供方注册契约声明自身可处理的 `provider_key`、`provider_type`、展示元数据和回调处理器。

选择该方案是因为宿主必须持续掌握安全边界：本地用户是否启用、租户是否匹配、是否存在绑定、Token 如何签发、会话如何记录、登录日志如何落地。插件直接签发 LinaPro Token 或直接创建会话会绕过宿主安全治理，因此不采用。

备选方案是每种登录方式都直接写入宿主认证服务。该方案短期简单，但后续接入企业微信、钉钉、飞书、Google、Microsoft、SAML 时会不断修改核心认证代码，不符合插件化边界。

### 2. 标准化 ExternalIdentity

插件回调处理成功后只返回标准化 `ExternalIdentity`，至少包含 `provider_key`、`provider_type`、`subject`，并可包含 `union_id`、`open_id`、`external_tenant_id`、邮箱、手机号、显示名、头像、原始资料等字段。

宿主使用 `provider_key + subject` 查找 `sys_auth_identity`，再定位本地 `sys_user`。未绑定时拒绝登录并返回稳定业务错误，前端引导用户先用已登录账号完成绑定。

该选择避免“邮箱相同自动登录”“手机号相同自动合并”等高风险推断，降低账号接管风险。

### 3. 首期不自动创建用户

外部身份无法天然表达 LinaPro 的本地租户、角色、菜单权限、数据权限范围和组织归属。首期要求管理员或用户先完成明确绑定，登录时只消费绑定结果。

未来如需自动开户，应作为独立变更引入可审计的 JIT Provisioning 策略，包括租户解析、默认角色、审批、域名白名单、企业通讯录同步和冲突处理。

### 4. 提供方配置分层

`sys_auth_provider` 保存提供方元数据、启用状态、展示顺序、插件来源、非敏感公开配置和配置引用。客户端密钥、企业微信 Secret、OIDC Client Secret 等敏感值首期从 `config.yaml` 或环境变量读取。

这样可以先避免密钥落库、加密轮换和后台编辑审计的复杂度。后续如需运行时密钥治理，可扩展为专门的密钥管理能力，而不是在认证表中直接保存明文。

### 5. 数据模型采用通用提供方表 + 通用身份绑定表

宿主新增：

- `sys_auth_provider`：`provider_key`、`plugin_id`、`provider_type`、`name`、`icon`、`enabled`、`sort`、`config_mode`、`config_json`、时间戳和软删除字段。
- `sys_auth_identity`：`tenant_id`、`user_id`、`provider_key`、`provider_type`、`subject`、`union_id`、`open_id`、`external_tenant_id`、外部部门、邮箱、手机号、展示资料、`raw_profile`、`last_login_at`、`bound_at`、时间戳和软删除字段。

首期唯一约束：

- `UNIQUE(provider_key, subject)`，防止同一个外部身份绑定到多个本地用户。
- `UNIQUE(user_id, provider_key)`，首期限制同一用户对同一提供方只绑定一个外部账号。

不为每个登录方式新增宿主专用表。协议特有字段尽量进入插件配置或 `raw_profile`；只有宿主需要参与查询、唯一性判断、安全校验或审计的字段才进入通用表。

### 6. REST API 边界

宿主新增认证相关 API：

- `GET /api/v1/auth/providers`：读取当前可展示的外部登录提供方。
- `GET /api/v1/auth/providers/{providerKey}/authorize`：创建授权临时状态并返回跳转地址。
- `GET /api/v1/auth/providers/{providerKey}/callback`：处理 OAuth/OIDC 回调并转换为宿主登录。
- `POST /api/v1/auth/providers/{providerKey}/callback`：保留给未来需要 POST 回调的提供方。
- `GET /api/v1/user/auth-identities`：读取当前用户已绑定外部身份。
- `POST /api/v1/user/auth-identities/{providerKey}/bind`：当前已登录用户发起或完成绑定。
- `DELETE /api/v1/user/auth-identities/{providerKey}`：解绑当前用户的外部身份。

读取使用 GET，绑定是创建关系使用 POST，解绑使用 DELETE，符合项目 REST 规范。

### 7. 授权临时状态使用统一缓存/协调能力

`state`、`nonce`、PKCE verifier、绑定意图、回跳地址和过期时间必须写入宿主统一缓存能力。`cluster.enabled=false` 时可使用本地进程缓存或 SQL 单机分支；`cluster.enabled=true` 时必须使用共享缓存、分布式 KV、协调服务或等价机制，确保任一实例收到回调都能校验状态。

临时状态键必须包含提供方、用途、随机状态值和过期时间，消费后立即失效。失效操作必须幂等，重复回调必须被拒绝。

### 8. 插件启停与缓存失效

提供方列表和插件注册信息可以缓存，但权威数据源是数据库中的提供方元数据与当前插件运行态。插件启用、禁用、卸载、升级或提供方状态变更后，宿主必须按显式 scope 失效认证提供方缓存。

集群模式下，缓存失效必须通过现有集群拓扑、事件广播、共享修订号或分布式缓存传播，不允许只清理当前节点。

### 9. 前端集成方式

登录页保留用户名密码表单，同时调用提供方列表接口展示外部登录按钮。按钮点击后进入宿主授权接口，由宿主返回跳转地址或直接重定向。

个人中心安全设置页新增外部身份绑定列表和绑定/解绑操作。首期不在系统用户管理页提供管理员代绑定，避免管理权限、审计和身份确认流程过早扩大。

外部登录入口、绑定状态、错误提示和插件名称必须接入前端运行时 i18n；插件自身展示名称和图标通过提供方元数据返回。

## Risks / Trade-offs

- 外部身份被错误绑定导致账号接管 -> 绑定前必须要求用户已登录本地账号，并使用 `provider_key + subject` 唯一约束防止重复绑定。
- 邮箱或手机号自动匹配存在误绑定风险 -> 首期不做自动创建或自动匹配，只允许显式绑定。
- 集群回调命中不同实例导致 `state` 丢失 -> 授权临时状态必须使用统一缓存/协调能力，并按 `cluster.enabled` 区分单机和集群策略。
- 插件故障影响登录页 -> 提供方列表只展示宿主确认可用的启用提供方，插件回调失败必须隔离并返回稳定业务错误。
- 敏感配置治理不足 -> 首期密钥不落库，后续运行时密钥管理作为独立能力设计。
- OIDC 提供方差异较大 -> 首期 OIDC 插件支持标准 discovery、authorization code + PKCE 和 claim 映射配置；非标准差异通过配置扩展，无法覆盖时新增专用插件。
- 企业微信多企业、多应用映射复杂 -> 通过 `provider_key` 区分实例，`external_tenant_id` 保存 `corp_id`，避免把企业维度硬编码到宿主枚举。
- 数据权限边界容易被忽略 -> 当前用户自助绑定接口只允许访问自己的绑定记录；未来管理端绑定必须单独接入用户管理数据权限。

## Migration Plan

1. 新增宿主 SQL `013-pluggable-auth-providers.sql`，创建 `sys_auth_provider` 与 `sys_auth_identity`，不修改旧 SQL。
2. 执行 `make init` 或项目初始化流程加载迁移，再执行 `make dao` 生成 DAO/DO/Entity。
3. 实现宿主认证提供方服务、API、缓存和插件注册契约。
4. 新增 `auth-oidc` 与 `auth-wecom` 源码插件，并在插件清单、manifest、i18n 和配置示例中声明能力。
5. 更新前端登录页和个人中心安全设置页。
6. 新增单元测试、集成测试和 E2E 测试。
7. 回滚时可禁用所有外部登录提供方，用户名密码登录保持可用；新增表保留但不参与主链路。

## Open Questions

- 企业微信首期是否只支持网页授权登录，还是同时支持扫码登录，需要在实现任务开始前根据产品入口确认。
- OIDC claim 映射的首期配置范围需要控制在 `subject`、邮箱、手机号、显示名、头像和外部租户标识，复杂组织同步留到后续变更。
- 管理员代绑定、批量导入绑定关系、企业通讯录同步和自动开户是否需要独立排期。
