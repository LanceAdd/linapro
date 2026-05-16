## ADDED Requirements

### Requirement: 认证提供方后端注册契约
系统 SHALL 在插件后端扩展点中发布认证提供方注册契约，允许源码插件注册外部登录提供方处理器。

#### Scenario: 插件注册认证提供方
- **WHEN** 源码插件启动并声明认证提供方能力时
- **THEN** 宿主 MUST 接收提供方元数据和回调处理器，并将其纳入认证提供方治理

#### Scenario: 未启用插件不参与登录
- **WHEN** 某认证提供方所属插件被禁用或卸载时
- **THEN** 宿主 MUST 停止展示该提供方并拒绝其授权和回调请求

#### Scenario: 认证提供方注册失败
- **WHEN** 插件声明的认证提供方缺少 `provider_key`、类型或必要处理器时
- **THEN** 宿主 MUST 拒绝该注册并记录可诊断错误，不得影响其他插件注册

### Requirement: 认证提供方扩展点与认证 Hook 边界
系统 SHALL 区分登录前的认证提供方处理器和登录后的认证生命周期 Hook。

#### Scenario: 外部登录成功后继续触发登录 Hook
- **WHEN** 外部身份登录被宿主转换为本地登录成功时
- **THEN** 宿主 MUST 继续分发 `auth.login.succeeded` Hook

#### Scenario: 认证提供方不得替代登录 Hook
- **WHEN** 插件仅订阅 `auth.login.succeeded` Hook 时
- **THEN** 该插件 MUST NOT 获得处理外部授权回调或返回 `ExternalIdentity` 的能力
