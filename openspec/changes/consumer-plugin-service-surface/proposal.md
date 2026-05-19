## Why

LinaPro 当前已经具备核心宿主、管理工作台和插件运行时，但插件暴露的 HTTP 能力主要服务后台管理语义，尚未为插件承载独立 C 端业务服务提供稳定入口。为保持架构轻量，本变更不引入多插件组装 App，而是让一个业务插件能够完整、隔离、规范地提供一个 C 端服务。

## What Changes

- 新增插件 C 端服务面，约定 C 端 API 使用 `/api/c/v1/<plugin-id>/...` 前缀，与现有后台 `/api/v1` 分离。
- 新增 C 端请求上下文，用于承载 `tenantId`、`pluginId`、渠道、设备、匿名标识等请求元数据。
- 明确 C 端登录、会话、token 和 `public`、`optional`、`login` 等访问语义由插件自理，宿主不提供统一 C 端认证中间件。
- 允许源码插件通过宿主发布的中间件注册 C 端路由，并保持插件启用状态、租户边界、缓存、i18n 和错误响应治理。
- 为动态插件预留 C 端 route surface 声明能力，但第一期不要求完整实现动态插件 C 端执行。
- 允许插件自带可选 C 端前端资产，由宿主托管；插件后台管理页面仍挂载到现有管理工作台。
- API 文档需要区分 Admin API 与 Consumer API，避免 C 端接口混入后台管理接口分组。
- 不引入 Consumer App Registry、不做多插件组装、不内建商城、会员、商品或订单等具体 C 端业务。

## Capabilities

### New Capabilities

- `consumer-plugin-service-surface`: 定义插件提供独立 C 端服务的 API 前缀、请求上下文、租户边界、i18n、缓存、审计和插件资源托管要求。

### Modified Capabilities

- `plugin-http-slot-extension`: 补充源码插件注册 C 端路由时的前缀、中间件和插件启用治理要求。
- `plugin-ui-integration`: 补充插件可选 C 端前端资产的目录、托管和隔离要求。
- `system-api-docs`: 补充 Admin API 与 Consumer API 分组展示和文档过滤要求。

## Impact

- 后端宿主启动和路由装配：`apps/lina-core/internal/cmd`、`apps/lina-core/internal/service/middleware`、`apps/lina-core/pkg/pluginhost`。
- 插件注册契约：源码插件 HTTP registrar、动态插件 route contract 的扩展点和文档。
- 消费者请求上下文：新增宿主请求元数据契约，不复用后台 `sys_user` 作为 C 端身份模型，也不内建统一 C 端认证基础设施。
- API 文档：OpenAPI 聚合、前端 API 文档页面和插件 route projection。
- 插件前端资产：`apps/lina-plugins/<plugin-id>/frontend/consumer/` 的可选资源托管。
- i18n 和缓存治理：C 端接口和资产必须显式纳入 `tenantId + pluginId + locale` 等作用域。
