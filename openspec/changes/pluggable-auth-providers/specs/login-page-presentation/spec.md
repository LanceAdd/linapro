## MODIFIED Requirements

### 需求：当前阶段仅暴露用户名/密码登录入口
系统在当前阶段 SHALL 保留用户名/密码登录能力，并 SHALL 在存在已启用、可展示的外部登录提供方时展示对应外部登录入口。系统 MUST NOT 展示未完成、未启用或当前前端无法渲染的认证入口。

#### Scenario: 标准登录页显示用户名密码表单
- **WHEN** 未认证用户访问 `/auth/login` 时
- **THEN** 页面 MUST 显示用户名、密码、记住我和登录控件

#### Scenario: 标准登录页显示已启用外部登录入口
- **WHEN** `GET /api/v1/auth/providers` 返回一个或多个可展示外部登录提供方时
- **THEN** 页面 MUST 在用户名密码表单之外展示对应外部登录入口

#### Scenario: 无外部登录提供方时保持简化登录页
- **WHEN** `GET /api/v1/auth/providers` 返回空列表或请求失败时
- **THEN** 页面 MUST 继续可用并只展示用户名密码登录能力

#### Scenario: 用户访问未完成的认证子路由
- **WHEN** 用户访问 `/auth/code-login`、`/auth/qrcode-login`、`/auth/forget-password` 或 `/auth/register` 时
- **THEN** 系统 MUST 重定向回标准登录页 `/auth/login`

## ADDED Requirements

### Requirement: 外部登录入口交互
登录页 SHALL 通过宿主认证提供方接口动态渲染外部登录入口，并通过宿主授权接口进入外部授权流程。

#### Scenario: 点击外部登录入口
- **WHEN** 用户点击某外部登录提供方入口时
- **THEN** 前端 MUST 调用该提供方授权入口并按照宿主返回结果跳转到外部授权页

#### Scenario: 提供方展示文案国际化
- **WHEN** 当前语言发生变化时
- **THEN** 外部登录入口的固定文案 MUST 使用当前语言资源刷新，提供方名称 MUST 使用宿主返回的可展示名称
