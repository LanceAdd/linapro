## 1. 数据模型与配置

- [x] 1.1 新增宿主 SQL `apps/lina-core/manifest/sql/013-pluggable-auth-providers.sql`，创建 `sys_auth_provider` 与 `sys_auth_identity`，包含软删除、时间戳、必要索引和幂等约束。
- [x] 1.2 执行项目初始化和 DAO 生成流程，确保新增表对应 DAO/DO/Entity 由生成工具维护，且不手写生成文件。
- [x] 1.3 新增外部登录提供方配置读取结构，支持 OIDC 与企业微信敏感密钥从 `config.yaml` 或环境变量读取，数据库仅保存非敏感配置和引用。
- [x] 1.4 补充默认配置示例和插件 manifest 配置说明，明确 `provider_key`、回调地址、密钥来源和多实例配置方式。

## 2. 宿主认证提供方服务

- [x] 2.1 实现认证提供方注册表服务，支持插件注册、启停过滤、提供方列表查询、缓存读取和显式 scope 失效。
- [x] 2.2 实现外部身份绑定服务，支持当前用户绑定列表、绑定、解绑、唯一约束冲突处理和用户状态校验。
- [x] 2.3 实现外部登录转换服务，根据 `ExternalIdentity` 查找绑定用户，并复用宿主现有 JWT/Refresh Token、在线会话、用户信息和权限加载流程。
- [x] 2.4 为外部登录失败场景补充 `bizerr` 错误码，包括提供方禁用、插件不可用、授权状态无效、身份未绑定、重复绑定和用户禁用。
- [x] 2.5 登录成功与失败事件中补充登录方式、提供方标识和外部身份摘要，确保登录日志插件可区分外部登录来源。

## 3. 授权流程与缓存一致性

- [x] 3.1 实现授权临时状态存储，包含 `state`、`nonce`、PKCE verifier、用途、回跳地址和过期时间。
- [x] 3.2 按 `cluster.enabled` 区分单机与集群策略：单机可使用本地或 SQL 分支，集群必须使用共享缓存、协调服务、分布式 KV 或等价机制。
- [x] 3.3 实现授权临时状态一次性消费、过期清理、重复回调拒绝和幂等失效。
- [x] 3.4 实现提供方状态、插件状态、配置变更后的跨实例缓存失效，并记录权威数据源、一致性模型和最大可接受陈旧时间。

## 4. 后端 API 与路由

- [x] 4.1 新增 `GET /api/v1/auth/providers`，返回当前可展示外部登录提供方列表。
- [x] 4.2 新增 `GET /api/v1/auth/providers/{providerKey}/authorize`，创建授权临时状态并返回外部授权跳转信息。
- [x] 4.3 新增 `GET /api/v1/auth/providers/{providerKey}/callback` 和预留 `POST /api/v1/auth/providers/{providerKey}/callback`，完成外部登录回调处理。
- [x] 4.4 新增 `GET /api/v1/user/auth-identities`、`POST /api/v1/user/auth-identities/{providerKey}/bind`、`DELETE /api/v1/user/auth-identities/{providerKey}`。
- [x] 4.5 路由、Controller 构造和服务依赖必须使用显式依赖注入，不新增聚合依赖结构或业务路径隐式 `New()`。

## 5. 插件注册契约

- [x] 5.1 在宿主插件扩展契约中新增认证提供方注册类型、`ExternalIdentity` 类型和回调处理接口。
- [x] 5.2 更新源码插件 registrar，使插件可注册一个或多个认证提供方实例，并在插件禁用、卸载、升级时失效对应提供方。
- [x] 5.3 明确认证提供方扩展点与 `auth.login.succeeded`、`auth.logout.succeeded` Hook 的边界，确保插件不能通过 Hook 获得登录前回调处理能力。
- [x] 5.4 更新插件开发文档的英文 `README.md` 和中文 `README.zh-CN.md`，说明认证提供方契约和示例。

## 6. 首期认证插件

- [x] 6.1 新增源码插件 `apps/lina-plugins/auth-oidc/`，包含 `plugin.yaml`、`plugin_embed.go`、`backend/`、`frontend/`、`manifest/` 和插件配置示例。
- [x] 6.2 实现 OIDC discovery、authorization code、PKCE、nonce 校验、token exchange、userinfo/ID Token claim 解析和标准化 `ExternalIdentity` 返回。
- [x] 6.3 新增源码插件 `apps/lina-plugins/auth-wecom/`，包含标准源码插件目录结构、manifest、i18n 和配置示例。
- [x] 6.4 实现企业微信授权交换、企业用户身份解析、`corp_id` 映射到 `external_tenant_id`，并返回标准化 `ExternalIdentity`。
- [x] 6.5 两个插件均不得直接签发宿主 Token、写入 `sys_online_session` 或绕过宿主绑定关系。

## 7. 前端体验

- [x] 7.1 更新登录页，在用户名密码表单之外动态加载并展示已启用外部登录提供方入口。
- [x] 7.2 实现外部登录按钮点击流程，调用宿主授权入口并处理跳转、回调结果和稳定错误提示。
- [x] 7.3 更新个人中心安全设置页，展示当前用户外部身份绑定状态、可绑定提供方、绑定和解绑操作。
- [x] 7.4 更新前端 API 客户端、类型定义和错误处理，确保提供方为空或接口失败时登录页仍可使用用户名密码登录。
- [x] 7.5 外部登录入口、绑定状态、按钮、错误提示和插件展示文案接入前端运行时 i18n。

## 8. 安全、权限与治理

- [x] 8.1 确认外部登录回调校验 `state`、`nonce`、PKCE、回调用途和过期时间，拒绝重复回调和跨用途复用。
- [x] 8.2 明确首期不自动创建用户、不自动分配租户、不自动分配角色，并在错误提示和文档中体现未绑定处理路径。
- [x] 8.3 自助绑定接口只允许访问当前用户自己的绑定记录；若后续新增管理端绑定，必须另行接入用户管理数据权限。
- [x] 8.4 运行依赖治理扫描或等价静态验证，确认认证、会话、插件状态、缓存敏感服务未在非启动边界隐式构造新实例。
- [x] 8.5 检查新增或修改的枚举语义值、错误码、日志输出、时间长度配置和 API 文档标签符合项目后端规范。

## 9. i18n 与资源

- [x] 9.1 更新宿主前端运行时语言包，覆盖登录页外部入口、个人中心绑定区、按钮、状态和错误提示。
- [x] 9.2 为 `auth-oidc` 与 `auth-wecom` 新增 `manifest/i18n/<locale>/*.json` 和必要的 apidoc i18n 资源。
- [x] 9.3 插件启用、禁用、升级后按语言、插件和业务 scope 精细失效运行时翻译包缓存，不清空无关语言或无关 sector。
- [x] 9.4 确认 `i18n.enabled=false` 时前端仍按默认语言展示外部登录和绑定文案，并隐藏语言切换入口。

## 10. 测试与验证

- [x] 10.1 为外部身份绑定服务、唯一约束冲突、未绑定登录拒绝、禁用用户拒绝登录新增后端单元测试。
- [x] 10.2 为授权临时状态新增单机与集群策略测试，覆盖过期、重复消费、跨实例读取和幂等失效。
- [x] 10.3 为 OIDC 插件新增协议回调、nonce/PKCE 校验、claim 映射和错误路径测试。
- [x] 10.4 为企业微信插件新增身份解析、`corp_id` 映射和错误路径测试。
- [x] 10.5 新增 E2E 测试用例，覆盖登录页外部入口展示、空提供方降级、个人中心绑定列表、绑定入口和解绑操作，TC ID 按 `lina-e2e` 规范分配。
- [x] 10.6 运行变更包 Go 测试；涉及 Controller、路由绑定或启动编排时，运行 `cd apps/lina-core && go test ./internal/cmd -count=1` 或等价覆盖测试。
- [x] 10.7 运行前端类型检查、单元测试和新增 E2E 测试，确认登录页和个人中心在桌面与移动视口无文本重叠。
- [x] 10.8 运行 `openspec validate pluggable-auth-providers --strict`，并在实现完成后触发 `lina-review` 审查。
