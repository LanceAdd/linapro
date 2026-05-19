## ADDED Requirements

### Requirement: 宿主必须提供独立的 C 端插件 API 服务面

系统 SHALL 提供独立于后台管理 API 的 C 端插件 API 服务面，默认前缀为 `/api/c/v1`。源码插件暴露 C 端接口时，公共路径 SHALL 使用 `/api/c/v1/<plugin-id>/...`，其中 `<plugin-id>` 必须与插件清单中的插件 ID 一致。

#### Scenario: 插件注册 C 端 API 路由

- **WHEN** 已启用插件 `mall` 注册 C 端商品列表接口
- **THEN** 该接口的公共路径 SHALL 位于 `/api/c/v1/mall/...`
- **AND** 该接口不得注册到后台管理 API 前缀 `/api/v1` 下作为 C 端入口

#### Scenario: 插件禁用后 C 端 API 不可访问

- **WHEN** 插件 `mall` 已被禁用
- **AND** 调用方请求 `/api/c/v1/mall/products`
- **THEN** 宿主 SHALL 拒绝该请求或按插件不可用语义返回稳定失败
- **AND** 请求不得继续进入插件业务处理器

### Requirement: C 端请求必须注入消费者上下文

系统 SHALL 为 C 端请求提供消费者上下文。该上下文至少 SHALL 能表达当前插件 ID、租户 ID、语言环境、匿名标识、设备标识和渠道标识。后台管理用户上下文不得被当作 C 端身份来源。

#### Scenario: 匿名 C 端请求获得基础上下文

- **WHEN** 匿名调用方请求一个 `public` C 端接口
- **THEN** 宿主 SHALL 注入插件 ID、租户边界和语言环境
- **AND** 宿主 SHALL NOT 注入统一消费者主体

#### Scenario: C 端请求携带插件自有 token

- **WHEN** 调用方携带 `Authorization` 请求 C 端接口
- **THEN** 宿主 SHALL NOT 解析该 token 为统一消费者主体
- **AND** 请求是否登录、token 是否有效、会话是否可用 SHALL 由插件自有逻辑判断

### Requirement: 宿主不得提供 C 端认证 Provider 或默认 C 端 JWT

系统 SHALL NOT 在宿主内提供 Consumer Auth Provider、Consumer Principal、默认 host-signed consumer JWT、`ConsumerAuthOptional` 中间件或 `ConsumerAuthRequired` 中间件。C 端登录、注册、会话、token 签发、token 验证、token 撤销和业务授权 SHALL 由业务插件自理。

#### Scenario: 插件实现自己的登录和会话

- **WHEN** 插件提供 `/api/c/v1/<plugin-id>/sessions` 或等价登录接口
- **THEN** token 签发、存储、验证和撤销 SHALL 由插件完成
- **AND** 宿主 SHALL 只提供 C 端路由、请求上下文、租户边界和资源治理能力

#### Scenario: 插件保护自己的 C 端业务接口

- **WHEN** 插件需要保护会员资料、订单、内容编辑或其他 C 端业务接口
- **THEN** 插件 SHALL 使用自有 middleware、handler 或 service 校验消费者登录态和业务权限
- **AND** 宿主 SHALL NOT 要求该消费者写入后台 `sys_user` 或通过后台 RBAC 授权

### Requirement: C 端身份不得强制复用后台管理员用户模型

系统 SHALL 将后台管理员用户与 C 端消费者作为不同身份域治理。C 端消费者可以由插件自有账号、外部身份提供方或插件自有账号模块产生，但不得要求所有 C 端用户必须写入后台 `sys_user` 并绑定后台角色菜单权限。

#### Scenario: 插件使用自有消费者账号表

- **WHEN** 业务插件使用自有表维护消费者账号
- **THEN** 宿主 SHALL 允许该插件在插件边界内使用自有账号识别 C 端消费者
- **AND** 后台 RBAC 菜单权限不得成为该消费者登录的必需前置条件

### Requirement: 插件消费者服务面不得定义具体 C 端产品模型

系统 SHALL 将插件消费者服务面限定为插件公开消费者服务时使用的宿主接入与治理层。宿主可以提供 C 端 API 服务面、请求上下文、租户边界、前端资源托管、资源索引、治理快照、API 文档和生命周期失效能力；宿主不得内建统一商品、内容、订单、会员、门户登录、CMS 登录、统一 C 端用户中心、统一 C 端认证 Provider、默认 C 端 JWT 或统一 C 端业务后台模型。

#### Scenario: 商城插件自主管理业务模型

- **WHEN** 商城插件提供商品、订单、会员和登录能力
- **THEN** 这些业务表、后台管理页面、业务权限和登录实现 SHALL 由商城插件自理
- **AND** 宿主 SHALL 只按插件边界提供 C 端路由、上下文、资源托管和治理能力

#### Scenario: 宿主不处理插件 C 端 token

- **WHEN** 调用方请求 C 端接口并携带插件自有 token
- **THEN** 宿主 SHALL NOT 调用 Consumer Auth Provider 或默认 JWT 解析逻辑
- **AND** 不得为了具体插件登录方式而在宿主内硬编码插件账号表、插件 session 表或插件登录流程

### Requirement: C 端租户边界必须在认证前可解析

系统 SHALL 支持在 C 端请求认证前解析租户边界。解析来源可以包括域名、路径、Header、查询参数或插件提供的解析器。租户能力启用时，受租户约束的 C 端接口在无法解析租户时 MUST 失败关闭。

#### Scenario: 匿名公开接口仍具有租户边界

- **WHEN** 匿名调用方请求某租户下的公开商品列表
- **THEN** 宿主 SHALL 在进入插件业务处理器前解析租户边界
- **AND** 插件查询 SHALL 能基于该租户边界过滤数据

#### Scenario: 租户启用但无法解析租户

- **WHEN** 租户能力已启用
- **AND** C 端请求需要租户边界但宿主无法解析租户
- **THEN** 宿主 SHALL 拒绝该请求
- **AND** 请求不得退化为跨租户或平台全量访问

### Requirement: C 端缓存和 i18n 必须使用显式作用域

系统 SHALL 要求 C 端相关缓存和运行时翻译资源使用显式作用域。作用域至少 SHALL 包含插件 ID，并在适用时包含租户 ID、语言、渠道或消费者主体。普通 C 端业务路径不得无理由清空所有插件、所有租户或所有语言的缓存。

#### Scenario: 插件刷新 C 端翻译资源

- **WHEN** 插件 `mall` 的 C 端语言资源发生变化
- **THEN** 宿主 SHALL 按插件 ID 和语言范围失效相关运行时消息缓存
- **AND** 不得清空无关插件的语言资源缓存

#### Scenario: C 端业务缓存按租户和插件隔离

- **WHEN** 插件缓存某租户的公开商品列表
- **THEN** 缓存键或失效作用域 SHALL 包含插件 ID 和租户 ID
- **AND** 其他租户请求不得命中该租户缓存内容

### Requirement: C 端业务授权必须由消费者策略或插件业务规则承担

系统 SHALL 不把后台菜单 RBAC 作为 C 端资源访问的默认授权模型。C 端业务资源访问 SHALL 基于插件自有消费者身份、租户边界、资源归属、渠道或插件业务策略进行校验。

#### Scenario: 消费者只能访问自己的订单

- **WHEN** 已登录消费者请求订单详情
- **THEN** 插件 SHALL 在业务层校验订单归属当前消费者或符合业务策略
- **AND** 不得仅依赖后台菜单权限判断该 C 端请求是否可访问

### Requirement: 插件可以选择托管 C 端前端资产

系统 SHALL 允许业务插件提供可选 C 端前端资产。没有 C 端前端资产的插件仍可只提供 C 端 API。宿主托管 C 端前端资产时，资产 URL SHALL 包含插件 ID 和版本或等价发布标识，以支持隔离和升级。

#### Scenario: 插件只提供 C 端 API

- **WHEN** 插件未声明 C 端前端资产
- **THEN** 宿主 SHALL 允许该插件继续注册 C 端 API
- **AND** 插件安装、启用和后台管理页面不得因缺少 C 端前端资产而失败

#### Scenario: 插件提供 C 端前端资产

- **WHEN** 插件声明 C 端前端资产
- **THEN** 宿主 SHALL 以包含插件 ID 和发布标识的稳定路径托管这些资产
- **AND** 插件禁用后对应 C 端入口 SHALL 不再作为有效业务入口暴露

### Requirement: C 端 API 文档必须与后台 API 分组展示

系统 SHALL 在 API 文档中区分后台管理 API 和 C 端插件 API。C 端 API SHALL 能按插件 ID 或 C 端服务面过滤展示，避免与后台管理接口混淆。

#### Scenario: 查看 C 端 API 文档

- **WHEN** 用户查看系统 API 文档
- **THEN** 文档 SHALL 能识别 `/api/c/v1/<plugin-id>/...` 为 C 端插件 API
- **AND** 文档 SHALL 能与 `/api/v1` 后台管理 API 区分展示

### Requirement: 宿主必须提供 C 端服务面治理快照

系统 SHALL 在宿主插件服务内提供按需构建的 C 端服务面治理快照，用于统一观察源码插件的 C 端 API route binding、C 端前端挂载资源、插件版本、启用态和租户治理声明。该快照 SHALL 是宿主治理投影，不得引入 Consumer App Registry，不得把多个插件强制组合成一个 App，也不得要求具体业务插件实现额外后台管理能力。

#### Scenario: 构建包含 API 与前端挂载的 C 端服务面快照

- **WHEN** 源码插件声明 C 端 API route binding 或 C 端前端挂载
- **THEN** 宿主 SHALL 在治理快照中按插件 ID 聚合其 C 端 API 路由数量、路由元数据、前端挂载路径、入口文件、SPA fallback、资源数量、插件版本、启用状态和租户治理声明
- **AND** 后台管理 API 路由不得被计入 C 端服务面

#### Scenario: 快照不引入新的长期缓存

- **WHEN** 宿主构建 C 端服务面治理快照
- **THEN** 宿主 SHALL 从当前 manifest、route binding、C 端前端资源索引和插件启用态按需生成结果
- **AND** 不得新增绕过既有插件生命周期失效与 runtime revision 感知的新长期缓存
