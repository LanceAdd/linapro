## ADDED Requirements

### Requirement: 外部登录提供方治理
系统 SHALL 维护宿主级外部登录提供方注册信息，并以 `provider_key` 作为稳定业务标识区分不同提供方实例。

#### Scenario: 管理已启用提供方元数据
- **WHEN** 宿主加载外部登录提供方配置时
- **THEN** 系统 MUST 保存提供方的 `provider_key`、`plugin_id`、`provider_type`、名称、图标、启用状态、排序值和非敏感配置引用

#### Scenario: 禁用提供方后不再展示
- **WHEN** 某外部登录提供方被禁用时
- **THEN** `GET /api/v1/auth/providers` MUST 不返回该提供方

#### Scenario: 同一协议存在多个实例
- **WHEN** 系统配置多个 OIDC 或企业微信实例时
- **THEN** 每个实例 MUST 使用不同 `provider_key` 独立参与展示、授权、绑定和登录

### Requirement: 外部身份绑定
系统 SHALL 使用宿主维护的外部身份绑定关系将外部身份映射到既有本地用户。

#### Scenario: 绑定外部身份到当前用户
- **WHEN** 已登录用户完成某提供方的绑定授权回调时
- **THEN** 系统 MUST 创建 `sys_auth_identity` 记录并关联当前用户

#### Scenario: 外部身份不得绑定多个用户
- **WHEN** 某 `provider_key + subject` 已绑定到一个本地用户时
- **THEN** 系统 MUST 拒绝将同一外部身份绑定到其他本地用户

#### Scenario: 同一用户同一提供方只绑定一个账号
- **WHEN** 当前用户已经绑定某 `provider_key` 时
- **THEN** 系统 MUST 拒绝为该用户再次绑定同一提供方的第二个外部账号

### Requirement: 外部登录回调转换为宿主登录
系统 SHALL 在外部登录回调成功后，由宿主根据绑定关系完成本地登录转换。

#### Scenario: 已绑定外部身份登录成功
- **WHEN** 插件返回有效 `ExternalIdentity` 且宿主找到对应本地用户绑定关系时
- **THEN** 系统 MUST 校验用户状态、租户边界和提供方状态，并由宿主签发 Token、创建在线会话、发布登录成功事件

#### Scenario: 未绑定外部身份登录失败
- **WHEN** 插件返回有效 `ExternalIdentity` 但宿主找不到绑定关系时
- **THEN** 系统 MUST 拒绝登录并返回稳定业务错误，且不得创建本地用户、Token 或在线会话

#### Scenario: 被禁用用户不能通过外部登录进入系统
- **WHEN** 外部身份绑定的本地用户处于禁用状态时
- **THEN** 系统 MUST 拒绝登录并记录登录失败原因

### Requirement: 插件返回标准化外部身份
系统 SHALL 要求认证提供方插件只返回标准化 `ExternalIdentity`，不得直接签发 LinaPro Token 或直接创建宿主会话。

#### Scenario: 插件完成协议交换
- **WHEN** 外部登录回调到达宿主并委托给提供方插件处理时
- **THEN** 插件 MUST 返回包含 `provider_key`、`provider_type` 和 `subject` 的 `ExternalIdentity`

#### Scenario: 插件不得越过宿主认证中心
- **WHEN** 插件完成外部身份认证时
- **THEN** 插件 MUST NOT 直接签发 LinaPro JWT、Refresh Token 或写入 `sys_online_session`

### Requirement: 授权临时状态安全
系统 SHALL 使用宿主统一缓存或协调能力保存 OAuth/OIDC 授权临时状态，并在回调校验后一次性消费。

#### Scenario: 集群模式下任一实例可处理回调
- **WHEN** `cluster.enabled=true` 且授权请求与回调请求落到不同宿主实例时
- **THEN** 回调实例 MUST 能读取并校验 `state`、`nonce` 和 PKCE verifier

#### Scenario: 重复回调被拒绝
- **WHEN** 同一授权 `state` 已经被成功消费后再次提交回调时
- **THEN** 系统 MUST 拒绝该回调并不得签发 Token

#### Scenario: 授权状态过期
- **WHEN** 回调中的 `state` 已超过配置的有效期时
- **THEN** 系统 MUST 拒绝回调并清理过期临时状态

### Requirement: 首期 OIDC 与企业微信提供方
系统 SHALL 首期提供通用 OIDC 插件和企业微信插件作为源码插件接入外部登录。

#### Scenario: OIDC 标准登录
- **WHEN** 用户选择已启用的 OIDC 提供方登录时
- **THEN** 系统 MUST 使用 authorization code 流程、PKCE 和 nonce 校验完成外部身份认证

#### Scenario: OIDC 覆盖 Google 和 Microsoft
- **WHEN** 管理员配置 Google、Microsoft、Azure AD/Entra ID、Keycloak 或 Okta 等标准 OIDC 提供方实例时
- **THEN** 系统 MUST 通过通用 OIDC 插件处理这些实例，而不要求新增宿主认证模型

#### Scenario: 企业微信登录
- **WHEN** 用户选择已启用的企业微信提供方登录时
- **THEN** 系统 MUST 由企业微信插件完成企业微信授权交换，并向宿主返回标准化 `ExternalIdentity`

### Requirement: 外部身份自助绑定管理
系统 SHALL 为当前登录用户提供外部身份绑定列表、绑定和解绑能力。

#### Scenario: 查看当前用户绑定列表
- **WHEN** 已登录用户调用 `GET /api/v1/user/auth-identities`
- **THEN** 系统 MUST 只返回当前用户自己的外部身份绑定记录

#### Scenario: 当前用户解绑外部身份
- **WHEN** 已登录用户调用 `DELETE /api/v1/user/auth-identities/{providerKey}`
- **THEN** 系统 MUST 只删除当前用户对该提供方的绑定关系

#### Scenario: 禁止越权管理他人绑定
- **WHEN** 用户通过自助绑定接口尝试读取或删除其他用户绑定关系时
- **THEN** 系统 MUST 拒绝该操作
