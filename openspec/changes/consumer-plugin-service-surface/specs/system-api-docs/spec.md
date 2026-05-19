## ADDED Requirements

### Requirement: API 文档必须区分后台管理接口和 C 端插件接口

系统 SHALL 在 OpenAPI 文档中区分后台管理接口和 C 端插件接口。路径以 `/api/c/v1/<plugin-id>/...` 开头的接口 SHALL 被标记为 C 端插件接口，并能按 C 端服务面或插件 ID 过滤展示。

#### Scenario: C 端插件接口进入独立分组

- **WHEN** 插件 `mall` 注册 `GET /api/c/v1/mall/products`
- **THEN** OpenAPI 文档 SHALL 将该接口识别为 C 端插件接口
- **AND** 该接口 SHALL 能与后台 `/api/v1` 管理接口区分展示

#### Scenario: 后台接口分组保持不变

- **WHEN** 系统生成现有后台管理 API 文档
- **THEN** `/api/v1` 下的后台接口 SHALL 继续按现有规则展示
- **AND** 新增 C 端接口分组不得改变后台接口的请求地址生成规则
