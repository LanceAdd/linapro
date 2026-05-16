## ADDED Requirements

### Requirement: 外部身份登录
系统 SHALL 支持通过已启用外部登录提供方完成用户认证，并将认证结果转换为宿主本地登录会话。

#### Scenario: 外部身份登录成功
- **WHEN** 用户通过外部登录提供方完成认证且外部身份已绑定到启用的本地用户时
- **THEN** 系统 MUST 返回与用户名密码登录一致的 Token、用户信息、角色、菜单和权限数据

#### Scenario: 外部身份登录发布认证事件
- **WHEN** 外部身份登录成功时
- **THEN** 系统 MUST 发布与宿主登录成功一致的认证生命周期事件，并在事件上下文中包含登录方式和提供方标识

#### Scenario: 外部身份登录失败不创建会话
- **WHEN** 外部身份未绑定、提供方禁用、插件不可用或本地用户状态无效时
- **THEN** 系统 MUST 拒绝登录且不得创建 `sys_online_session` 记录

### Requirement: 宿主统一签发外部登录 Token
系统 SHALL 由宿主统一签发外部登录后的 JWT 和 Refresh Token。

#### Scenario: 插件返回外部身份后由宿主签发 Token
- **WHEN** 认证提供方插件返回有效 `ExternalIdentity` 时
- **THEN** 宿主 MUST 根据本地绑定用户签发 Token，且 Token 有效期继续遵循现有 JWT 运行时配置

#### Scenario: 插件不得改变用户权限结果
- **WHEN** 外部身份登录转换为本地登录时
- **THEN** 用户角色、菜单、权限和首页路径 MUST 继续从本地用户权限体系计算
