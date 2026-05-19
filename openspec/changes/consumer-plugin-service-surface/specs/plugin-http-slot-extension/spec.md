## ADDED Requirements

### Requirement: 源码插件必须能够注册 C 端路由组

系统 SHALL 允许源码插件通过宿主 HTTP 注册入口注册 C 端路由组。C 端路由组 SHALL 复用宿主发布的响应、CORS、请求大小、C 端上下文和 C 端租户中间件，而不是要求插件直接持有 `*ghttp.Server` 或自行重建宿主服务图。插件自有认证中间件可以在插件路由组内自行组合。

#### Scenario: 插件注册 C 端公开路由组

- **WHEN** 源码插件通过 HTTP registrar 注册 `/api/c/v1/<plugin-id>` 路由组
- **THEN** 宿主 SHALL 在启动期装配该路由组
- **AND** 路由请求 SHALL 接受插件启用状态治理

#### Scenario: 插件组合自有 C 端认证中间件

- **WHEN** 同一插件同时声明公开 C 端接口和登录后 C 端接口
- **THEN** 插件 SHALL 能在不同路由子组中组合插件自有认证中间件
- **AND** 中间件组合不得影响该插件已有后台 `/api/v1` 路由

### Requirement: C 端路由注册必须暴露插件归属

系统 SHALL 在源码插件 C 端路由请求上下文中暴露当前插件 ID，使日志、i18n、缓存失效、API 文档和错误响应能够按插件归属治理。

#### Scenario: C 端路由请求带有插件 ID

- **WHEN** 请求命中源码插件注册的 `/api/c/v1/mall/products`
- **THEN** 宿主 SHALL 在请求上下文中记录插件 ID `mall`
- **AND** 插件业务处理器和宿主中间件 SHALL 能读取该插件 ID
