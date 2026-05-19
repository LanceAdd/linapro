## ADDED Requirements

### Requirement: 插件可以提供可选 C 端前端资产

系统 SHALL 允许插件在后台管理页面之外提供可选 C 端前端资产。源码插件的 C 端前端资源 SHOULD 放置在 `frontend/consumer/` 或等价清晰目录下，后台管理页面继续使用现有后台页面集成路径。缺少 C 端前端资产不得阻止插件提供 C 端 API。

#### Scenario: 源码插件提供 C 端前端资源

- **WHEN** 源码插件包含 `frontend/consumer/` 资源
- **THEN** 宿主 SHALL 能按插件 ID 和发布标识托管或构建这些资源
- **AND** 这些资源不得被当作后台管理工作台页面自动挂入后台菜单

#### Scenario: 源码插件不提供 C 端前端资源

- **WHEN** 源码插件只提供 C 端 API 而没有 `frontend/consumer/` 资源
- **THEN** 插件 SHALL 仍可正常安装、启用和注册 C 端 API
- **AND** 使用者可以通过外部前端消费该插件 C 端 API

### Requirement: C 端前端入口必须与插件启用状态一致

系统 SHALL 将插件 C 端前端入口与插件启用状态绑定。插件禁用、卸载或租户未开通时，宿主不得继续将该插件的 C 端前端入口作为有效业务入口暴露。

#### Scenario: 插件禁用后隐藏 C 端入口

- **WHEN** 插件 `mall` 被禁用
- **THEN** 宿主 SHALL 停止暴露该插件的有效 C 端前端入口
- **AND** 直接访问已知入口时 SHALL 返回稳定的不可用响应或降级页面

### Requirement: 插件 C 端前端挂载契约必须稳定且不遮蔽宿主路径

系统 SHALL 允许源码插件通过 `consumer.frontend` 声明稳定的 C 端前端挂载入口。`mount_path` 必须是非根绝对路径，宿主 SHALL 对前后斜杠进行归一化，并拒绝 `/`、包含路径穿越或重复分隔符的路径，以及会遮蔽 `/api`、`/plugin-assets`、`/consumer-plugin-assets`、`/swagger`、`/api.json`、`/openapi` 等宿主保留前缀的路径。`index` 缺省为 `index.html`，并且必须是安全的插件相对路径；`spa_fallback` 缺省为关闭，只有插件显式声明 `spa_fallback: true` 时，clean route 才回退到入口文件。

#### Scenario: 插件声明稳定 C 端挂载入口

- **WHEN** 源码插件在 `plugin.yaml` 中声明 `consumer.frontend.mount_path: /portal`
- **THEN** 宿主 SHALL 将 `/portal` 作为该插件的稳定 C 端访问入口
- **AND** 宿主 SHALL 将 `/portal/...` 下的请求解析到该插件 `frontend/consumer/` 资产
- **AND** 该稳定入口不得替代 `/consumer-plugin-assets/<plugin-id>/<version>/...` 这种插件和版本作用域的调试资产路径

#### Scenario: 插件声明非法挂载路径

- **WHEN** 插件声明 `consumer.frontend.mount_path: /api/mall` 或 `/`
- **THEN** 宿主 SHALL 拒绝该 C 端前端挂载声明
- **AND** 该插件不得通过非法路径遮蔽宿主 API、静态资产、API 文档或宿主前端路由

#### Scenario: 显式开启 SPA clean route 回退到入口文件

- **WHEN** 插件声明 `consumer.frontend.spa_fallback: true`
- **AND** 调用方访问 `/portal/login` 且该路径不是已声明的静态资源
- **THEN** 宿主 SHALL 返回该插件声明的入口文件
- **AND** 宿主 SHALL 确保 HTML entry 中的相对资源按挂载根路径解析

#### Scenario: 未开启 SPA fallback 时 clean route 返回不存在

- **WHEN** 插件省略或关闭 `consumer.frontend.spa_fallback`
- **AND** 调用方访问 `/portal/login` 且该路径不是已声明的静态资源
- **THEN** 宿主 SHALL 返回稳定的资源不存在响应
- **AND** 宿主不得将该请求回退到 HTML entry

#### Scenario: 缺失静态资源不回退到入口文件

- **WHEN** 调用方访问 `/portal/assets/missing.js`
- **AND** 该静态资源不存在
- **THEN** 宿主 SHALL 返回稳定的资源不存在响应
- **AND** 宿主不得将该请求回退到 HTML entry

### Requirement: C 端前端资源索引必须有明确权威来源和缓存边界

系统 SHALL 为源码插件 C 端前端挂载构建显式资源索引。索引的权威来源 SHALL 是当前进程可见的源码插件嵌入 manifest 与 `frontend/consumer/` 资产清单；索引条目 SHALL 至少包含插件 ID、插件版本、稳定挂载路径、入口文件、SPA fallback 策略和该插件声明的 C 端前端资产集合。稳定挂载路径不得重复，也不得互为父子路径，以避免一个插件的 C 端服务嵌入另一个插件的 URL 空间。索引缓存 MAY 使用进程内缓存，但必须在插件安装、卸载、启停、源码插件升级或运行时插件缓存修订号变化后失效或刷新；集群模式下 SHALL 复用宿主既有插件 runtime cache revision 或等价协调机制感知其他节点变更。

#### Scenario: 宿主构建 C 端前端资源索引

- **WHEN** 宿主首次解析源码插件 C 端前端挂载请求
- **THEN** 宿主 SHALL 从当前源码插件 manifest 与资产清单构建资源索引
- **AND** 索引 SHALL 记录每个挂载入口对应的插件 ID、版本、入口文件、SPA fallback 策略和资产集合
- **AND** 重复或互为父子路径的稳定挂载路径 SHALL 使索引构建失败并返回内部错误

#### Scenario: 资源索引缓存命中

- **WHEN** C 端前端资源索引已经构建且未失效
- **THEN** 后续挂载解析 MAY 复用该进程内索引
- **AND** 调用方不得能修改缓存中的索引状态

#### Scenario: 插件生命周期变更后索引失效

- **WHEN** 源码插件安装、卸载、启用状态变化或源码插件升级成功
- **THEN** 宿主 SHALL 失效本节点 C 端前端资源索引
- **AND** 下一次解析 SHALL 基于当前权威来源重建索引
