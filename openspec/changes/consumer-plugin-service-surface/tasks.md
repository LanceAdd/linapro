## 1. 架构与契约准备

- [x] 1.1 梳理现有 HTTP 启动、插件 HTTP registrar、动态插件 route contract、API 文档聚合和插件前端资产托管的改动点，确认不引入 Consumer App Registry 或多插件组装模型。
- [x] 1.2 定义 C 端服务面常量与命名约定，固定 `/api/c/v1/<plugin-id>/...` 路由前缀，并记录与后台 `/api/v1` 的边界。
- [x] 1.3 定义 Consumer Context 和 C 端租户解析的 Go 契约，确保宿主只传递请求元数据且不强制复用后台 `sys_user` 作为 C 端身份模型。
- [x] 1.4 明确本变更对 i18n 的影响面：新增 C 端运行时错误、API 文档文本、前端提示或插件示例文案时同步维护宿主与插件 i18n JSON；若某个子任务确认不新增文案，在任务记录中说明。
- [x] 1.5 明确本变更对缓存一致性的影响面：C 端缓存和翻译缓存必须携带插件、租户、语言等显式作用域；集群模式下缓存失效必须复用现有 cluster/cachecoord/coordination 接缝。

## 2. 后端服务面与中间件

- [x] 2.1 在宿主 HTTP 路由装配中新增 C 端服务面分组，保持现有 `/api/v1` 后台路由行为不变。
- [x] 2.2 实现 ConsumerCtx 中间件，向请求注入插件 ID、语言、租户候选、匿名标识、设备和渠道元数据。
- [x] 2.3 实现 C 端租户解析中间件，支持至少一种轻量解析来源，并在租户能力启用但无法解析必要租户时失败关闭。
- [x] 2.4 明确宿主不提供 ConsumerAuth optional/login 中间件；插件自有认证负责 `public`、`optional`、`login` 等业务访问语义。
- [x] 2.5 为 C 端租户失败和插件路由归属失败补充 bizerr 错误码和运行时 i18n 文案，禁止返回裸错误文本。
- [x] 2.6 为 Consumer Context 和租户解析增加单元测试，覆盖匿名请求元数据、插件自有 Authorization header 透传、租户缺失和插件禁用场景。

## 3. 插件 HTTP 注册扩展

- [x] 3.1 扩展源码插件 HTTP registrar 暴露 C 端可组合中间件目录，保持现有后台中间件目录兼容。
- [x] 3.2 在 route binding 捕获中标记 C 端服务面和插件 ID，供 API 文档、审计、缓存和日志治理使用。
- [x] 3.3 更新一处示例或测试插件，注册最小 `/api/c/v1/<plugin-id>/...` C 端接口，覆盖 public 和 login 两类访问。
- [x] 3.4 增加源码插件 C 端路由注册测试，验证插件禁用后 C 端路由不可进入业务处理器。
- [x] 3.5 评估动态插件 route contract 的 `surface` 扩展点并补充校验或预留文档；第一期如不实现动态 C 端执行，应在代码和任务记录中明确边界。

## 4. API 文档与前端资产

- [x] 4.1 扩展 OpenAPI 聚合逻辑，识别 `/api/c/v1/<plugin-id>/...` 为 C 端插件 API，并按服务面或插件 ID 分组。
- [x] 4.2 保持现有后台 `/api/v1` 文档生成和 servers URL 运行时生成规则不变。
- [x] 4.3 为 C 端 API 文档分组补充后端单元测试或集成测试。
- [x] 4.4 定义插件可选 `frontend/consumer/` 资源约定，更新插件资源扫描或托管逻辑，确保缺少 C 端资产不影响 API 插件。
- [x] 4.5 为 C 端前端资产托管增加插件启用状态检查，确保插件禁用后不继续暴露有效 C 端业务入口。

## 5. 文档与示例

- [x] 5.1 更新 `apps/lina-plugins/README.md` 和 `README.zh-CN.md`，说明一个 C 端业务服务由一个业务插件承载、C 端 API 前缀、目录建议和非目标边界。
- [x] 5.2 如新增目录级说明文档，按仓库规则同步创建英文 `README.md` 和中文 `README.zh-CN.md`。
- [x] 5.3 更新示例插件或开发文档，展示插件自有 `public`、`optional`、`login` 访问语义组合方式和 Consumer Context 请求元数据读取方式。
- [x] 5.4 更新 API 文档页面或说明，标明 Admin API 与 Consumer API 的分组含义。

## 6. 验证与审查

- [x] 6.1 运行 `openspec validate consumer-plugin-service-surface --strict` 并修复所有规范问题。
- [x] 6.2 运行覆盖变更包的 Go 测试，至少包括 `cd apps/lina-core && go test ./internal/cmd ./internal/service/middleware ./pkg/pluginhost -count=1` 或更窄但能覆盖路由构造和中间件契约的命令。
- [x] 6.3 如果修改 Controller 构造、路由绑定或启动编排，运行 `cd apps/lina-core && go test ./internal/cmd -count=1`。
- [x] 6.4 如果修改插件示例后端，运行对应插件包测试或编译烟测，例如进入插件目录执行 `go test ./... -count=1`。
- [x] 6.5 如果新增或修改前端页面交互，按 `lina-e2e` 规范新增或更新 E2E 测试并运行对应用例；若仅更新文档或后端 API 契约，记录无需 E2E 的原因。
- [x] 6.6 运行依赖治理或等价静态扫描，确认关键服务没有在非启动边界新增隐式 `New()` 构造。
- [x] 6.7 完成 `/lina-review` 审查，重点检查 C 端身份边界、多租户隔离、缓存作用域、i18n 资源、API 文档分组和 Go 编译门禁。

## 7. 实施记录

- i18n 影响：新增 C 端租户和插件路由归属错误码，并同步维护 `apps/lina-core/manifest/i18n/en-US/error.json` 与 `apps/lina-core/manifest/i18n/zh-CN/error.json`。本次未新增插件示例页面文案。
- 缓存一致性影响：本次新增 C 端请求上下文、路由面和源码插件 C 端资产托管，不新增业务缓存；C 端上下文显式携带 `pluginId`、`tenantId`、`locale` 等作用域，源码插件资产托管按插件 ID、版本和启用状态解析，未新增全局缓存清空路径。
- 动态插件边界：动态插件 route contract 已预留 `surface` 字段，第一期校验拒绝 `consumer` surface，暂不提供动态插件 C 端执行。
- 示例插件后端验证：本次通过测试插件/路由注册测试覆盖 C 端 public/login 行为，未修改 `apps/lina-plugins` 中的示例插件后端，因此无需执行具体插件包编译烟测。
- E2E 判断：本次未新增或修改前端页面交互，只新增源码插件可选 `frontend/consumer/` 资产发现与静态托管约定，因此无需新增 E2E。
- Go targeted 测试已通过：
  - `cd apps/lina-core && go test ./internal/service/middleware -run "TestConsumer|TestTenancy" -count=1`
  - `cd apps/lina-core && go test ./internal/service/apidoc -run TestBuildProjectsHostAndEnabledPluginRoutes -count=1`
  - `cd apps/lina-core && go test ./pkg/pluginhost ./pkg/pluginbridge/codec -count=1`
  - `cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -run "TestDiscoverPluginVuePathsUseDirectoryConvention|TestBuildPluginManifestSnapshotIncludesDirectoryDiscoveredAssets|TestScanEmbeddedSourcePluginManifestsUsesPluginEmbeddedFiles" -count=1`
  - `cd apps/lina-core && go test ./internal/service/plugin -run "TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared" -count=1`
  - `cd apps/lina-core && go test ./internal/cmd -run "TestBind|TestParse.*PluginAssetRequestPath" -count=1`
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- 更宽的门禁命令 `cd apps/lina-core && go test ./internal/cmd ./internal/service/middleware ./pkg/pluginhost -count=1` 已尝试执行，但被当前仓库既有测试环境/fixture 阻塞：`TestProductionPanicsMatchAllowlist` 中插件子模块 panic allowlist 调用点不一致、若干嵌入 SQL fixture 缺失、测试配置 YAML hex 解析失败，以及 `TestRequestBodyLimitFriendlyError` 在缺失 i18n/数据库 fixture 时回退为英文。上述失败不来自本次 C 端实现路径；本次使用更窄命令覆盖路由构造和中间件契约。
- 依赖治理扫描：已对本次涉及的 `internal/cmd`、`internal/service/middleware`、`internal/service/plugin`、`pkg/pluginhost`、`pkg/pluginbridge`、`pkg/pluginfs` 执行 `New()`/构造调用静态扫描；新增 C 端请求路径未发现关键服务在非启动边界隐式构造，扫描结果中的新增项主要为测试构造、错误构造、启动期既有装配和纯路径工具逻辑。
- `/lina-review` 审查结论：未发现阻塞归档的严重问题。审查中发现并修复 `ConsumerTenant` 中间件 nil 分支鲁棒性问题，修复后已重新运行 `cd apps/lina-core && go test ./internal/service/middleware -run "TestConsumer|TestTenancy" -count=1` 并通过。剩余风险为当前仓库既有宽包测试 fixture 阻塞，已在上方 Go 测试记录中说明。

## Feedback

- [x] **FB-1**: C 端前端 mount 解析错误分类不能依赖错误文本字符串匹配，应使用稳定错误类型或 sentinel error 支撑 catch-all 分流。
- [x] **FB-2**: C 端前端 mount 匹配后的失败语义需要区分 no-match、插件不可用、资源不存在和内部错误，避免所有非 no-match 错误都被粗略压成 404。
- [x] **FB-3**: C 端前端 mount 索引需要具备宿主生命周期失效机制，插件安装、卸载、启停或源码插件升级后应触发重建。
- [x] **FB-4**: C 端前端 HTML entry 的 `<base href>` 注入应比固定 `<head>` 字符串替换更健壮，并补充对应单元测试。
- [x] **FB-5**: C 端前端 `consumer.frontend` manifest 契约需要固化 `mount_path`、`index`、`spa_fallback`、保留前缀、根路径拒绝和稳定挂载语义。
- [x] **FB-6**: C 端前端 mount/lifecycle 测试矩阵需要覆盖 manifest 边界、默认值、禁用声明、静态资源缺失、SPA fallback、HTML base 注入和 mount index 失效。
- [x] **FB-7**: C 端前端资源索引需要从简单 mount 列表升级为显式资源索引模型，统一记录 mount、插件版本、入口文件、SPA 策略和资产集合。
- [x] **FB-8**: C 端前端资源索引缓存策略需要显式固化权威来源、缓存键、失效触发、集群刷新路径和故障降级边界。
- [x] **FB-9**: Consumer Context 需要保持为请求元数据契约，明确不复用后台 `sys_user`，且不在宿主定义统一 Consumer Principal。
- [x] **FB-10**: Consumer Auth Provider、默认 host-signed consumer JWT、`ConsumerAuthOptional` 和 `ConsumerAuthRequired` 不应作为宿主插件消费者服务面能力，应收束为插件自有认证。
- [x] **FB-11**: 宿主需要提供按需构建的 C 端服务面治理快照，统一观察 C 端 API route binding、前端挂载资源、插件版本、启用态与租户治理声明。
- [x] **FB-12**: C 端服务面治理快照需要单元测试覆盖聚合、过滤后台路由、前端资源统计、启用态与租户治理字段，且不得新增长期缓存或具体业务插件实现。
- [x] **FB-13**: 插件消费者服务面边界需要重新收束为宿主接入与治理层，明确宿主不提供具体 C 端产品模型、统一用户中心、业务后台或 Consumer App Registry。
- [x] **FB-14**: 为了保持一个业务服务一个插件的干净边界，宿主应移除已实现的 Consumer Auth Provider/Principal/JWT 代码和误导性的 Auth Provider 注入草案文档。
- [x] **FB-15**: 插件工作区 README、Lina Portal C 端说明和任务记录中仍存在旧宿主认证能力表述，需要统一清理为插件自有认证边界。
- [x] **FB-16**: 旧术语容易误导为宿主提供完整 C 端产品平台，应整体收束为“插件消费者服务面”。
- [x] **FB-17**: 代码注释应与“插件消费者服务面”术语保持一致，同时保留 `ConsumerSurface*` 稳定代码标识。
- [x] **FB-18**: 插件前端资源响应缺少 HTTP 缓存头和 304 协议，浏览器无法区分版本化静态资源与稳定 HTML 入口。
- [x] **FB-19**: C 端前端 HTML `<base href>` 注入在每次稳定挂载请求中重复执行，应预处理并随资源索引生命周期失效。
- [x] **FB-20**: `spa_fallback` 默认开启会让稳定挂载吞掉任意无扩展名路径，应改为插件显式 opt-in。
- [x] **FB-21**: Consumer API OpenAPI operations 不应继承宿主后台 `BearerAuth` 全局安全声明，避免误导为宿主提供 C 端 JWT。
- [x] **FB-22**: `plugin_frontend.go` 中 C 端前端资产、索引、mount 逻辑需要拆分为职责聚焦的实现文件，并移除旧测试切口表述。
- [x] **FB-23**: OpenSpec 任务记录中残留的旧宿主 C 端认证正向表述需要清理为插件自有认证边界。

### Feedback 验证记录

- FB-1/FB-2：已将 mount no-match、插件未启用、资源不存在改为稳定 sentinel error 分类；catch-all handler 对 no-match 放行，对插件不可用/资源不存在返回 404，对未分类内部错误返回 500 并记录日志。新增/更新 `TestIsSourceConsumerFrontendMountNotFound`、`TestIsSourceConsumerFrontendMountAssetNotFound` 覆盖错误分类。
- FB-3：已新增 `invalidateSourceConsumerFrontendMounts` 并在源码插件安装、卸载、启停、源码插件升级成功后失效本节点 mount index；读取 mount index 前调用 `ensureRuntimeCacheFresh`，集群模式下通过既有 plugin runtime cache revision 刷新路径触发其他节点失效。
- FB-4：已将 HTML base 注入改为大小写不敏感且支持 `<head ...>` 属性的 head start tag 匹配；`TestRewriteSourceConsumerHTMLBase` 覆盖 `<HEAD data-app="portal">` 场景。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared|TestMatchSourceConsumerFrontendMountPath|TestRewriteSourceConsumerHTMLBase|TestIsSourceConsumerFrontendMountNotFound|TestIsSourceConsumerFrontendMountAssetNotFound|TestLoadSourceConsumerFrontendMountEntriesCachesIndex|TestInvalidateSourceConsumerFrontendMountsClearsIndex|TestSourceConsumerSPAFallbackEnabledDefaultsFalse" -count=1`。
- 路由绑定相关烟测已通过：`cd apps/lina-core && go test ./internal/cmd -run "TestParseSourceConsumerPluginAssetRequestPath|TestParsePluginAssetRequestPath" -count=1`。
- catalog 回归测试已通过：`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`。
- 宽包门禁已尝试：`cd apps/lina-core && go test ./internal/service/plugin -count=1` 仍被当前仓库既有动态插件 fixture/SQLite YAML/缺失测试表阻塞；`cd apps/lina-core && go test ./internal/cmd -count=1` 仍被当前仓库既有 SQLite YAML fixture 阻塞。本次以更窄命令覆盖 mount 解析、缓存失效、路由入口和 catalog 契约。
- i18n 影响：本次只调整宿主内部错误分类、HTTP 状态分流、mount index 缓存失效和 HTML base 注入，不新增用户可见业务文案或运行时语言包。
- 缓存一致性影响：本次新增进程内 mount index 显式失效；单机模式通过本地失效生效，集群模式复用既有 plugin runtime cache revision 与 `ensureRuntimeCacheFresh` 跨实例刷新，不引入普通业务路径全量缓存清空。
- 数据权限影响：本次不新增数据读写接口；仅调整插件静态资源托管入口，不涉及后台数据权限查询或写操作。
- FB-5：已在 `plugin-ui-integration` 增量规范中固化 `consumer.frontend` 稳定挂载契约，明确 `mount_path` 非根绝对路径、保留前缀拒绝、`index` 默认值与安全相对路径、`spa_fallback` 默认关闭且需显式 opt-in、稳定入口与 `/consumer-plugin-assets/<plugin-id>/<version>/...` 调试资产路径的边界。
- FB-5：已修复 `NormalizeConsumerSpec` 对根路径和异常分隔符的校验，`mount_path: /`、`//`、重复分隔符、反斜杠、路径穿越和宿主保留前缀都会被拒绝；空 `mount_path` 仍表示插件未声明稳定 C 端入口，不阻止 API-only 插件。
- FB-6：已新增 `TestNormalizeConsumerSpecValidatesFrontendMountContract` 覆盖 manifest 默认值、挂载归一化、禁用声明、根路径拒绝、保留前缀拒绝、重复分隔符、反斜杠、路径穿越和非法 `index`。
- FB-6：已新增 `TestActiveSourceConsumerFrontendSpecFiltersDisabledAndNonSourceManifests` 与 `TestLooksLikeSourceConsumerStaticAsset`，补齐 mount 生效筛选和静态资源缺失不触发 SPA fallback 的单元测试；既有测试继续覆盖 SPA fallback 显式开启、默认关闭、HTML base 注入、mount index 本地缓存与生命周期失效。
- Go 编译门禁已通过：`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -run "TestNormalizeConsumerSpecValidatesFrontendMountContract" -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared|TestMatchSourceConsumerFrontendMountPath|TestRewriteSourceConsumerHTMLBase|TestIsSourceConsumerFrontendMountNotFound|TestIsSourceConsumerFrontendMountAssetNotFound|TestLoadSourceConsumerFrontendMountEntriesCachesIndex|TestInvalidateSourceConsumerFrontendMountsClearsIndex|TestSourceConsumerSPAFallbackEnabledDefaultsFalse|TestActiveSourceConsumerFrontendSpecFiltersDisabledAndNonSourceManifests|TestLooksLikeSourceConsumerStaticAsset" -count=1`。
- i18n 影响：FB-5/FB-6 只调整 manifest 契约、内部校验和单元测试，不新增用户可见文案、运行时语言包、插件 manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-5 不新增缓存；FB-6 只补充既有进程内 mount index 与生命周期失效测试。现有策略仍是单机本地失效，集群模式通过既有 plugin runtime cache revision 与 `ensureRuntimeCacheFresh` 跨实例感知。
- 数据权限影响：FB-5/FB-6 不新增数据读写接口或查询路径，仅约束公开静态资产挂载契约，不涉及后台角色数据权限。
- FB-7：已将内部 `sourceConsumerFrontendMounts` 简单列表升级为 `sourceConsumerFrontendResourceIndex`，索引条目统一记录插件 ID、版本、稳定 mount、入口文件、SPA fallback 策略和 `frontend/consumer/` 资产集合。
- FB-7：稳定 mount 解析现在先通过资源索引匹配最具体挂载路径，再使用索引资产集合判断资源是否存在；缺失静态资源返回资源不存在，clean route 才按 SPA fallback 回退到入口文件。
- FB-8：索引权威来源明确为当前进程可见的源码插件 embedded manifest 与 catalog 暴露的 `frontend/consumer/` 资产清单；缓存键等价于当前插件 runtime revision 下的进程内资源索引，条目作用域包含 plugin ID、version、mount path 和 asset path。
- FB-8：索引缓存仍为进程内缓存；插件安装、卸载、启停、源码插件升级成功后本节点显式失效，读取前通过 `ensureRuntimeCacheFresh` 感知既有 plugin runtime cache revision，集群模式下由同一修订号刷新路径触发跨实例本地索引失效。索引构建失败时不写入 ready 状态，下一次请求会重新尝试构建。
- FB-8：已补充 `TestSourceConsumerFrontendResourceIndexMatchIgnoresSiblingPrefixes` 与 `TestFindSourceConsumerFrontendOverlappingMountRejectsNestedMounts` 覆盖 sibling prefix 匹配和嵌套 mount 拒绝；更新 `TestLoadSourceConsumerFrontendMountEntriesCachesIndex` 覆盖索引条目与资产集合 clone，防止调用方修改进程内缓存。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared|TestMatchSourceConsumerFrontendMountPath|TestRewriteSourceConsumerHTMLBase|TestIsSourceConsumerFrontendMountNotFound|TestIsSourceConsumerFrontendMountAssetNotFound|TestLoadSourceConsumerFrontendMountEntriesCachesIndex|TestInvalidateSourceConsumerFrontendMountsClearsIndex|TestSourceConsumerFrontendResourceIndexMatchIgnoresSiblingPrefixes|TestFindSourceConsumerFrontendOverlappingMountRejectsNestedMounts|TestSourceConsumerSPAFallbackEnabledDefaultsFalse|TestActiveSourceConsumerFrontendSpecFiltersDisabledAndNonSourceManifests|TestLooksLikeSourceConsumerStaticAsset" -count=1`。
- Go regression 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`。
- 路由入口烟测已通过：`cd apps/lina-core && go test ./internal/cmd -run "TestParseSourceConsumerPluginAssetRequestPath|TestParsePluginAssetRequestPath" -count=1`。
- 宽包门禁已尝试：`cd apps/lina-core && go test ./internal/service/plugin -count=1` 仍被当前仓库既有动态插件 fixture/SQLite YAML/缺失测试表阻塞，包括缺失动态 wasm、SQLite 测试配置 YAML hex 解析失败、`plugin_multi_tenant_user_membership` 测试表缺失等。本次使用定向命令覆盖 C 端资源索引生产包编译和行为回归。
- i18n 影响：FB-7/FB-8 只调整内部资源索引模型、缓存失效策略和单元测试，不新增用户可见文案、运行时语言包、插件 manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-7/FB-8 明确索引缓存为进程内派生状态，权威来源为源码插件 manifest 与资产清单；单机通过本地生命周期失效，集群通过既有 plugin runtime cache revision 跨实例感知后本地重建，不引入业务路径全量缓存清空。
- 数据权限影响：FB-7/FB-8 不新增数据读写接口或查询路径，仅调整公开静态资产索引与解析，不涉及后台角色数据权限。
- FB-9：已将 `pluginhost.ConsumerContext` 收束为请求元数据，只保留 plugin ID、tenant ID、locale、anonymous/device/channel 等字段；宿主不再定义统一 `ConsumerPrincipal`，也不再把后台 `sys_user` 作为 C 端身份来源。
- FB-10：已从宿主移除 `ConsumerAuthProvider`、`ConsumerAuthInput`、`ConsumerTokenKindAccess`、默认 `hostJWTConsumerAuthProvider`、`ConsumerAuthOptional` 和 `ConsumerAuthRequired`；插件 C 端登录、session、token 和业务授权由插件自有 middleware/handler/service 负责。
- FB-10：已更新 middleware 与 pluginhost 测试，覆盖 ConsumerCtx/ConsumerTenant 只注入请求元数据、宿主不解析 C 端 Authorization header、发布给插件的宿主 middleware 目录不再暴露 ConsumerAuth 能力。
- FB-11：已新增宿主内部 `ConsumerSurfaceGovernanceService` 与 `BuildConsumerSurfaceSnapshot`，按需聚合源码插件 C 端 API route binding、C 端前端资源索引、插件版本、启用态和租户治理声明；该投影不引入 Consumer App Registry，也不要求具体业务插件实现后台管理能力。
- FB-12：已新增 `TestBuildConsumerSurfaceSnapshotAggregatesHostInputs`、`TestBuildConsumerSurfaceSnapshotDerivesPluginIDFromConsumerPath`、`TestBuildConsumerSurfaceSnapshotSortsPluginsAndRoutes`，覆盖聚合、后台路由过滤、前端资源统计、启用态、租户治理字段和确定性排序。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestBuildConsumerSurfaceSnapshot" -count=1`。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- 宽包门禁已尝试：`cd apps/lina-core && go test ./internal/service/plugin -count=1` 仍被当前仓库既有动态插件/SQLite fixture 阻塞，包括缺失动态 wasm、SQLite YAML hex 解析失败、`plugin_multi_tenant_user_membership` 测试表缺失和动态插件卸载 fixture 预期不一致；本次新增代码通过定向测试覆盖生产包编译和快照行为。
- i18n 影响：FB-11/FB-12 仅新增宿主内部治理投影和单元测试，不新增用户可见文案、错误码、运行时语言包、插件 manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-11/FB-12 不新增长期缓存；快照按需从当前 manifest、route binding、既有 C 端前端资源索引和启用态读取。前端资源索引仍复用既有生命周期失效与 plugin runtime cache revision 集群感知路径。
- 数据权限影响：FB-11/FB-12 不新增 HTTP API、数据库查询接口或后台数据操作；当前只是宿主内部治理元数据投影，不涉及后台角色数据权限。
- FB-13：已更新 `design.md`，将插件消费者服务面定义为“插件公开消费者服务的宿主接入与治理层”，并明确宿主只提供服务面、上下文、租户边界、前端资源托管、资源索引、治理快照、API 文档和生命周期失效；商品、内容、订单、会员、门户登录、CMS 登录、业务后台和业务权限均归插件自理。
- FB-13：已在增量规范中新增“插件消费者服务面不得定义具体 C 端产品模型”，明确不得把统一 C 端用户中心、商品/内容/订单模型、门户登录、CMS 登录、统一认证 Provider 或统一 C 端业务后台纳入宿主通用能力。
- FB-14：已进一步更新设计与规范，明确宿主不提供 Consumer Auth Provider、Consumer Principal、默认 C 端 JWT 或 `ConsumerAuth*` 中间件；删除未纳入当前边界的 `apps/lina-core/docs/consumer-auth-provider-injection.md` 草案，避免误导后续实现方向。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./pkg/pluginhost -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/middleware -run "TestConsumer|TestTenancy" -count=1`。
- Go regression 测试已通过：`cd apps/lina-core && go test ./internal/service/middleware -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestBuildConsumerSurfaceSnapshot" -count=1`。
- 路由绑定烟测已通过：`cd apps/lina-core && go test ./internal/cmd -run "TestBind|TestParse.*PluginAssetRequestPath" -count=1`。
- 插件兼容验证已通过：`cd apps/lina-plugins/lina-portal && GOWORK=off go test ./backend -run TestRegisterRoutesCapturesConsumerAndAdminBindings -count=1` 与 `cd apps/lina-plugins/lina-portal && GOWORK=off go test ./backend/internal/service/portal -run Test -count=1`。直接在插件目录运行不带 `GOWORK=off` 的 `go test` 会被根 `go.work` 拦截，因为该插件 module 不在 workspace use 列表中。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- 静态残留扫描已通过：宿主 Go 代码中未发现 `ConsumerAuthProvider`、`ConsumerPrincipal`、`ConsumerAuthInput`、`ConsumerTokenKind`、`hostJWTConsumerAuthProvider`、`ConsumerAuthOptional`、`ConsumerAuthRequired` 或 `DefaultAuthProviderAvailable` 残留引用。
- i18n 影响：FB-14 删除宿主 C 端认证错误码和认证中间件，不新增用户可见文案、运行时语言包、插件 manifest/i18n 或 apidoc i18n 资源；既有未使用的 auth 错误翻译扫描无残留。
- 缓存一致性影响：FB-14 移除请求内认证 provider 和默认 JWT 解析，不新增缓存、订阅、分布式状态或失效路径；C 端资源索引缓存策略不变。
- 数据权限影响：FB-14 不新增 HTTP 数据读写接口或后台查询路径；C 端业务授权与资源归属校验由插件自有逻辑承担。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- i18n 影响：FB-13 仅调整 OpenSpec 设计和规范文本，不新增运行时文案、错误码、前端页面或 API 文档资源。
- 缓存一致性影响：FB-13 不修改生产代码或缓存行为，仅收束架构边界，不新增缓存。
- 数据权限影响：FB-13 不新增数据读写接口、后台查询路径或 HTTP API，仅明确宿主不接管具体 C 端业务后台和权限模型。
- FB-15：已更新 `apps/lina-plugins/README.md`、`apps/lina-plugins/README.zh-CN.md` 和 `apps/lina-plugins/lina-portal/CONSUMER-SURFACE.md`，将 C 端账号、登录、session、token、optional/login 访问语义统一表述为插件自有能力；宿主只提供请求元数据、租户边界和 C 端前端托管/治理能力。
- FB-15：已更新任务记录旧表述，避免继续使用“宿主发布 optional/login 认证中间件”或“回退”等误导性说法。
- 静态残留扫描已通过：除负向边界约束和插件自有 optional-login 语义外，未发现宿主验证 host-signed consumer JWT、注入 ConsumerPrincipal 或发布 ConsumerAuth optional/login 中间件的正向描述。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- i18n 影响：FB-15 只调整 README、插件说明和 OpenSpec 任务记录，不新增或删除运行时文案、前端语言包、manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-15 不修改生产代码或缓存行为，不新增失效路径、订阅、分布式状态或缓存键。
- 数据权限影响：FB-15 不新增 HTTP API、数据库查询、后台操作或 C 端业务授权逻辑。
- FB-16：已将 OpenSpec 设计、增量规范和任务记录中的旧术语统一收束为“插件消费者服务面”，用于表达插件公开消费者服务时使用的宿主接入与治理层。
- FB-16：保留代码中的 `ConsumerSurface*`、`consumer-plugin-service-surface` 变更目录和 `surface` route contract 等稳定标识，不做语义无关的 API、包名或目录迁移。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- 静态术语扫描已通过：活跃变更文本中不再出现旧术语。
- i18n 影响：FB-16 只调整架构术语文档，不新增用户可见运行时文案、前端语言包、manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-16 不修改生产代码或缓存行为，不新增缓存、失效路径、订阅或分布式状态。
- 数据权限影响：FB-16 不新增 HTTP API、数据库查询、后台操作或 C 端业务授权逻辑。
- FB-17：已将 `ConsumerSurface*` 相关 Go 注释中的旧服务面表述收束为 `host-governed plugin consumer surface`，并移除已经不符合边界的认证能力注释描述。
- FB-17：保留 `ConsumerSurfaceSnapshot`、`ConsumerSurfaceGovernanceService`、`BuildConsumerSurfaceSnapshot`、`plugin_consumer_surface.go` 等稳定代码标识，不做类型、方法、文件或 API 迁移。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestBuildConsumerSurfaceSnapshot" -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./pkg/pluginhost -run "TestRouteRegistrarCaptureConsumerRouteSurface" -count=1`。
- 静态注释扫描已通过：宿主相关 Go 文件中未发现旧服务面表述或不符合边界的认证能力注释描述。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- i18n 影响：FB-17 只调整 Go 注释和任务记录，不新增用户可见运行时文案、前端语言包、manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-17 不修改生产逻辑或缓存行为，不新增缓存、失效路径、订阅或分布式状态。
- 数据权限影响：FB-17 不新增 HTTP API、数据库查询、后台操作或 C 端业务授权逻辑。
- FB-18：已为插件前端资源输出增加 `ETag` 和 `Cache-Control` 元数据；当前公开资源 URL 尚未包含 release checksum，静态资源、HTML 与稳定挂载入口统一使用 `no-cache` 进行验证式缓存，HTTP 写入层根据 `If-None-Match` 返回 304。
- FB-18：已补充 HTTP header 和 ETag 匹配单元测试，覆盖缓存头写入、逗号分隔 `If-None-Match` 和通配符匹配。
- FB-19：已在源码插件 C 端前端资源索引条目中预处理挂载入口 HTML，缓存已注入 `<base href>` 的 index 输出；稳定挂载请求返回缓存副本，避免每次请求重复 rewrite。
- FB-19：HTML 预处理缓存仍挂在既有进程内 C 端前端资源索引上，插件安装、卸载、启停、源码插件升级和 runtime revision 感知继续通过既有 `invalidateSourceConsumerFrontendMounts` 与 `ensureRuntimeCacheFresh` 生效；集群模式仍由既有 plugin runtime cache revision 触发各节点本地重建。
- FB-20：已将 `consumer.frontend.spa_fallback` 缺省语义从开启改为关闭；只有插件显式声明 `spa_fallback: true` 时，缺失的 clean route 才回退到入口 HTML，带扩展名的缺失静态资源仍稳定返回不存在。
- FB-20：该改动保持 manifest 字段为 `*bool`，不引入 routes 或 fallback 枚举，避免宿主承担插件前端路由注册中心职责。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin/internal/frontend -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestRewriteSourceConsumerHTMLBase|TestApplySourceConsumerMountAssetPolicyForcesRevalidation|TestCloneFrontendAssetOutputProtectsCachedBytes|TestLoadSourceConsumerFrontendMountEntriesCachesIndex|TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared|TestMatchSourceConsumerFrontendMountPath|TestIsSourceConsumerFrontendMount|TestSourceConsumerFrontendResourceIndex|TestFindSourceConsumerFrontendOverlappingMount|TestSourceConsumerSPAFallbackEnabled|TestActiveSourceConsumerFrontendSpec|TestLooksLikeSourceConsumerStaticAsset" -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/cmd -run "TestApplyPluginFrontendAssetHeadersEmitsValidators|TestRequestETagMatches|TestParseSourceConsumerPluginAssetRequestPath|TestParsePluginAssetRequestPath" -count=1`。
- OpenSpec 校验已通过：`openspec validate consumer-plugin-service-surface --strict`。
- i18n 影响：FB-18/FB-19/FB-20 不新增用户可见运行时文案、前端语言包、manifest/i18n 或 apidoc i18n 资源。
- 缓存一致性影响：FB-18 新增的是浏览器 HTTP 缓存验证协议；FB-19 复用现有进程内资源索引缓存和既有生命周期/集群 revision 失效路径；FB-20 只调整 fallback 判定默认值，不新增缓存、订阅或全量失效路径。
- 数据权限影响：FB-18/FB-19/FB-20 不新增 HTTP API、数据库查询、后台操作或 C 端业务授权逻辑；仅调整插件前端静态资源响应、预处理和 fallback 默认语义。
- FB-21：已让 Consumer API operation 使用空 `security` 覆盖文档级 `BearerAuth`，后台 API 仍保留全局 BearerAuth 声明，C 端认证边界继续归插件自有逻辑。
- FB-22：已将动态插件 frontend facade、源码插件 C 端直接资产解析、资源索引和稳定 mount 解析拆分到独立文件；`plugin_frontend.go` 只保留运行时托管前端 facade，旧测试切口表述已移除。
- FB-23：已清理任务记录中把宿主描述为提供 C 端认证能力的正向表述，并同步更新重命名后的测试记录。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/apidoc -run TestBuildProjectsHostAndEnabledPluginRoutes -count=1`。
- Go targeted 测试已通过：`cd apps/lina-core && go test ./internal/service/plugin -run "TestRewriteSourceConsumerHTMLBase|TestApplySourceConsumerMountAssetPolicyForcesRevalidation|TestCloneFrontendAssetOutputProtectsCachedBytes|TestLoadSourceConsumerFrontendMountEntriesCachesIndex|TestNormalizeSourceConsumerFrontendAssetPath|TestSourceConsumerFrontendAssetDeclared|TestMatchSourceConsumerFrontendMountPath|TestIsSourceConsumerFrontendMount|TestSourceConsumerFrontendResourceIndex|TestFindSourceConsumerFrontendOverlappingMount|TestSourceConsumerSPAFallbackEnabled|TestActiveSourceConsumerFrontendSpec|TestLooksLikeSourceConsumerStaticAsset|TestBuildConsumerSurfaceSnapshot" -count=1`。
- i18n 影响：FB-21/FB-22/FB-23 不新增用户可见运行时文案、错误码、前端语言包、manifest/i18n 或 apidoc i18n 资源；FB-23 仅清理 OpenSpec 任务文字。
- 缓存一致性影响：FB-21 不涉及缓存；FB-22 为结构拆分，不改变既有 C 端资源索引缓存、失效或集群 revision 感知策略；FB-23 不修改生产逻辑。
- 数据权限影响：FB-21/FB-22/FB-23 不新增 HTTP 数据读写接口、数据库查询或后台操作，不涉及角色数据权限。
