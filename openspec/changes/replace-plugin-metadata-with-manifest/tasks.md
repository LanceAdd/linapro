## 1. 规则与历史语义核对

- [x] 1.1 实现前重新读取命中的规则文件：`openspec.md`、`architecture.md`、`plugin.md`、`backend-go.md`、`testing.md`、`database.md`、`i18n.md`、`cache-consistency.md`，并在任务记录中写明影响判断。
- [x] 1.2 静态检索历史 `Metadata`、`metadata.get`、`service: metadata`、`metadata.yaml`、`声明型资源`、`HostServices.Manifest()` 相关引用，按设计中的语义分流表确认需要删除、重命名、保留或仅更新注释的项目。
- [x] 1.3 确认宿主框架自身 `apps/lina-core/manifest/config/metadata.yaml`、动态路由元数据、审计元数据、错误元数据、数据库表元数据和 cron 展示 metadata 不属于插件 manifest 资源读取服务，避免误删历史功能。

执行记录：已读取命中规则文件。影响判断：本变更属于`apps/lina-core`插件宿主通用能力和动态插件 host service 授权边界调整；不新增 HTTP API、前端页面、数据库迁移、DAO、运行期依赖或缓存域。`i18n`影响仅为允许读取插件`manifest/i18n/`原文，不改变加载、聚合、缓存失效或翻译治理；SQL影响仅为允许读取插件`manifest/sql/`原文，不执行 SQL 或改变迁移账本；数据权限无业务数据读取/写入影响；开发工具跨平台无影响。静态检索确认无插件可见`Metadata()`、`MetadataService`、`metadata.get`、`service: metadata`入口；保留的 metadata 用法属于宿主系统信息、动态路由、审计/错误/数据库表/cron 展示等非插件 manifest 读取语义。

## 2. 插件 Manifest 契约收敛

- [x] 2.1 更新 `pkg/plugin/capability/contract`、`pkg/plugin/capability`、`pkg/plugin/capability/guest` 中 `Manifest()` 的注释和接口语义，明确其为插件自有 `manifest/` 原始资源只读视图。
- [x] 2.2 删除或迁移任何插件可见的旧 `Metadata()`、`MetadataService`、`metadata.get`、`service: metadata` 或等价读取入口，不保留 deprecated alias。
- [x] 2.3 更新 `Scan` 相关注释和测试命名，明确它是 YAML 便捷扫描能力，不代表 `Manifest()` 只能读取 YAML。

执行记录：已更新`Manifest()`契约和实现注释为插件自有`manifest/`原始资源只读视图；未发现可删除的独立旧`Metadata`公开入口；`Scan`继续只作为 YAML 便捷扫描方法。

## 3. 源码插件读取实现

- [x] 3.1 修改源码插件 manifest resource path 规范化逻辑，移除 `config/`、`sql/`、`i18n/` 排除规则，保留空根、绝对路径、URL、Windows drive path、路径穿越、重复 `manifest/` 前缀和跨插件路径拒绝。
- [x] 3.2 保持源码插件读取来源顺序和作用域：优先读取当前插件嵌入文件系统，再读取当前仓库开发目录中的同一插件 `manifest/`，不得读取宿主或其他插件资源。
- [x] 3.3 补充源码插件单元测试，覆盖读取 `metadata.yaml`、`config/config.example.yaml`、`sql/001-schema.sql`、`i18n/zh-CN/plugin.json`、缺失资源、路径穿越和跨插件拒绝。

执行记录：`manifest`路径规范化已移除`config/sql/i18n`保留目录拒绝，读取来源顺序和`readContainedFile`作用域保护保持不变；单元测试已覆盖新允许路径、缺失资源和越界拒绝。

## 4. 动态插件授权与 host service

- [x] 4.1 修改 `pluginbridge` host service path 校验和 WASM host call 授权逻辑，允许 `manifest` 服务声明和访问 `config/`、`sql/`、`i18n/` 相对路径。
- [x] 4.2 保持动态插件 `service: manifest` 的 `methods: [get]` 和 `resources.paths` 授权快照校验，未授权路径必须拒绝。
- [x] 4.3 补充 `pluginbridge` 和 `wasm` 单元测试，覆盖动态插件授权读取 `config/config.example.yaml`、`sql/001-schema.sql`、`i18n/zh-CN/plugin.json`，以及未授权路径拒绝。

执行记录：`pluginbridge`声明校验和 WASM 授权匹配均允许合法专用目录相对路径；动态插件仍必须具备`service: manifest`、`get`方法和`resources.paths`授权，测试覆盖授权读取和未授权拒绝。

## 5. 动态 artifact 资源视图和打包语义

- [x] 5.1 修改 dynamic active release artifact manifest resource projection，完整投影 artifact 中 `manifest/` 下实际携带的资源，不再只投影 `*.yaml` 或排除 `config/`、`sql/`、`i18n/`。
- [x] 5.2 保持 `HostServices.Config()` 的默认配置读取语义，继续从 active release 中识别 `manifest/config/config.yaml`，但不把 `Manifest()` 读取结果作为运行期有效配置自动生效。
- [x] 5.3 检查动态插件构建和 artifact 解析逻辑，确保 `go:embed` 和目录扫描资源来源都保留 `manifest/config/`、`manifest/sql/`、`manifest/i18n/` 及其他 manifest 资源原始路径。
- [x] 5.4 更新资源计数、示例插件 `plugin.yaml` 或测试 fixture，使动态插件通过 `service: manifest` 显式声明需要读取的 manifest 资源路径。

执行记录：artifact 解析不再限制专用目录或`.yaml`扩展名；active release manifest resource projection 完整输出`manifest/`下资源。`buildArtifactDefaultConfig`仍只读取`manifest/config/config.yaml`作为默认配置来源，`Manifest()`只返回原始字节。未修改动态插件构建器或示例插件目录；通过 runtime/WASM 测试 fixture 更新资源计数与`service: manifest`授权路径覆盖。

## 6. 历史命名和文档清理

- [x] 6.1 清理生产代码、测试、fixture 和注释中把 `Manifest()` 描述为 `Metadata` 或“声明型资源读取器”的历史表述。
- [x] 6.2 更新相关 README 或开发说明时同步维护中英文镜像；若最终无需修改目录级说明文档，在任务记录中明确无文档镜像影响。
- [x] 6.3 运行静态检索确认插件资源读取路径中不存在旧 `Metadata()`、`MetadataService`、`metadata.get`、`service: metadata` 或 deprecated alias。

执行记录：已清理`Manifest()`相关旧“declaration resource”表述；未修改 README 或目录级说明文档，无中英文镜像影响。静态检索确认插件资源读取路径中无旧 Metadata 读取入口残留。

## 7. 验证与审查

- [x] 7.1 运行 `go test ./pkg/plugin/capability/manifest ./pkg/plugin/capability/guest ./pkg/plugin/pluginbridge/internal/hostservice -count=1`。
- [x] 7.2 运行 `go test ./internal/service/plugin/internal/runtime ./internal/service/plugin/internal/wasm -count=1`。
- [x] 7.3 若修改动态插件构建器或官方示例插件，运行对应构建器/示例插件测试，并记录跨平台影响；若未修改开发工具，记录开发工具无影响。
- [x] 7.4 运行 `openspec validate replace-plugin-metadata-with-manifest --strict`。
- [x] 7.5 运行 `git diff --check` 和旧 Metadata 静态检索，确认没有格式问题和旧读取入口残留。
- [x] 7.6 完成任务后调用 `lina-review`，审查结论必须覆盖 `i18n`、缓存一致性、数据权限、SQL、开发工具跨平台、DI 来源和测试策略影响。

执行记录：两组 Go 测试均已通过。未修改动态插件构建器、开发工具、脚本、CI 或官方示例插件目录，开发工具跨平台无影响。`openspec validate replace-plugin-metadata-with-manifest --strict`、`git diff --check`、旧 Metadata 读取入口静态检索和旧“声明型资源/专用管线排除/YAML-only”表述检索均已通过。

审查记录：已按`lina-review`读取`AGENTS.md`以及`openspec.md`、`architecture.md`、`plugin.md`、`backend-go.md`、`testing.md`、`database.md`、`i18n.md`、`cache-consistency.md`、`data-permission.md`、`documentation.md`、`dev-tooling.md`、`api-contract.md`。未发现阻塞问题。`i18n`影响已限定为读取原文且不注册翻译或触发缓存失效；缓存一致性无新增缓存域，仍依赖 active release artifact 快照；数据权限无业务数据接口或租户/组织数据可见性变化；SQL 仅允许读取原文且不执行 SQL、不改账本；开发工具跨平台无影响；DI 来源无新增运行期依赖或构造函数变更；测试策略由 manifest、pluginbridge、runtime、wasm 单元测试和静态治理校验覆盖。

## Feedback

- [x] **FB-1**: `linapro-demo-dynamic`缺少`service: manifest`配置示例

执行记录：根因是本变更已要求动态插件通过`plugin.yaml`中的`service: manifest`、`methods: [get]`和`resources.paths`显式授权读取插件`manifest/`原始资源，但官方动态插件示例只提供了`manifest/config/`、SQL 和 i18n 资源文件，未在`hostServices`中展示`manifest.get`授权写法，也未在中英文说明文档中解释`manifest.get`与专用`config`服务的边界。已更新`apps/lina-plugins/linapro-demo-dynamic/plugin.yaml`，新增`service: manifest`示例并授权`config/config.example.yaml`、`config/config.yaml`、SQL 和 i18n JSON 路径；同步更新该插件根目录和`manifest/`目录的中英文 README，说明`manifest.get`仅读取打包原文，不替代配置、SQL 或 i18n 生命周期管线。

影响分析：已按`AGENTS.md`读取`openspec.md`、`plugin.md`、`documentation.md`、`i18n.md`、`architecture.md`、`testing.md`、`cache-consistency.md`、`data-permission.md`、`dev-tooling.md`、`backend-go.md`、`database.md`和`api-contract.md`。插件根目录无本地`AGENTS.md`。本反馈修改官方动态插件示例的`plugin.yaml`和文档，不新增或修改 Go 生产代码、HTTP API、数据库迁移、SQL 内容、前端页面或运行期依赖；DI 来源无新增依赖或构造函数变更。`i18n`影响为文档和插件清单注释说明，不新增运行时用户可见文案或翻译键；缓存一致性无新增缓存域或失效策略变更；数据权限无业务数据读写接口影响；SQL 影响仅为在授权示例中引用已有 SQL 文件路径，不修改 SQL 内容、账本、幂等性或执行管线；API 契约无影响；开发工具跨平台无脚本、CI、构建入口实现变更；测试策略按治理类反馈使用 YAML 解析、OpenSpec 校验、格式检查、清单解析测试和动态插件构建器示例测试覆盖。

验证记录：已通过`ruby -e 'require "yaml"; ...'`解析`plugin.yaml`和配置 YAML；`openspec validate replace-plugin-metadata-with-manifest --strict`通过；顶层`git diff --check`和`git -C apps/lina-plugins diff --check`通过；`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`通过；`go test ./hack/tools/linactl/internal/wasmbuilder -run 'Test.*Demo|Test.*Official' -count=1`通过；静态检索确认示例中存在`service: manifest`及关键授权路径，且无`service: metadata`或`metadata.get`残留。

审查记录：已按`lina-review`执行反馈级审查，范围包含`apps/lina-plugins/linapro-demo-dynamic/plugin.yaml`、根目录中英文 README、`manifest/`中英文 README 和本反馈任务记录。未发现阻塞问题。规则域结论：插件目录规范通过，动态插件`manifest`服务使用显式`methods`和`resources.paths`；文档规范通过，中英文镜像同步维护；OpenSpec 记录通过，反馈已先记录、再修复并补充根因、影响和验证；`i18n`、缓存一致性、数据权限、SQL、API、开发工具跨平台和 DI 来源均无运行时行为变更风险；本反馈为治理类示例补齐，无需新增单元测试或 E2E。

- [x] **FB-2**: `linapro-demo-dynamic`的`manifest`示例授权过宽且缺少声明到使用的可见闭环

执行记录：根因是 FB-1 为了展示`manifest.get`可读取专用资源目录，把`config.example`、SQL 和 i18n 资源都放入了官方动态插件示例的`hostServices.resources.paths`，示例授权范围过宽，且页面或接口没有展示这些声明路径被真实读取，无法证明“先声明、后使用”的流程已走通。已新增`manifest/config/profile.yaml`作为专用演示资源，并将`plugin.yaml`中的`service: manifest`授权收敛为仅`config/config.yaml`和`config/profile.yaml`。后端新增只读`GET /api/v1/manifest-demo`接口，复用`Manifest()`host service 读取`config/profile.yaml`与`config/config.yaml`，同时在原`host-call-demo`响应中补充 manifest 摘要；内嵌页面新增`Manifest Host Service`面板，加载`manifest-demo`并展示 profile 路径、profile 值、config 路径和配置预览。同步更新插件根目录与`manifest/`目录中英文 README、插件运行时中英文语言包、中文 apidoc 翻译资源，并新增插件自有 E2E `TC005-manifest-host-service-demo.ts`和 POM 定位器。

影响分析：已按`AGENTS.md`读取`openspec.md`、`plugin.md`、`backend-go.md`、`api-contract.md`、`frontend-ui.md`、`testing.md`、`documentation.md`、`i18n.md`、`architecture.md`、`data-permission.md`、`cache-consistency.md`、`dev-tooling.md`和`database.md`，并使用`goframe-v2`、`lina-feedback`、`lina-review`、`lina-e2e`和`frontend-design`技能。插件根目录无本地`AGENTS.md`。本反馈修改动态插件目录、插件 host service 授权、Go 后端 service/controller/API DTO、插件前端页面、插件运行时 i18n、插件 apidoc i18n、E2E 和中英文 README。新增接口为只读`GET /manifest-demo`，不会触发存储、数据或网络副作用；原`host-call-demo`仍为执行型演示接口，仅追加 manifest 响应投影。数据权限影响：`manifest-demo`只读取当前插件打包资源，不读取业务数据、租户数据或组织数据；原`host-call-demo`的数据演示仍沿用既有`data`host service 授权和临时记录清理逻辑。缓存一致性影响：不新增缓存、快照或失效路径；读取来源仍为动态插件 active release artifact/guest host service 当前授权视图。SQL 影响：不修改 SQL 内容、迁移、DAO 或执行账本，且 SQL 资源不再作为示例 manifest 授权路径。开发工具跨平台影响：不修改脚本、Makefile、linactl 或 CI；仅运行既有构建和测试入口。DI 来源：`serviceImpl`新增`manifestSvc`依赖，owner 为插件 guest host service 目录`guest.Default().Manifest()`，在插件 service 构造函数`New()`中与既有`runtime/config/hostConfig`host service client 同源创建，不新增宿主服务图或缓存敏感共享实例。接口性能：页面新增单次只读请求，不产生列表瀑布式调用或`N+1`查询。

验证记录：`GOWORK=off go test ./backend/internal/service/dynamic ./backend/internal/controller/dynamic ./backend/api/dynamic ./backend/api/dynamic/v1 -count=1`通过；`GOWORK=off go test ./backend/api -count=1`通过；`GOWORK=off go test ./backend/internal/service/dynamic -run 'TestRunHostCallDemoManifestReadsAuthorizedResources|TestRunHostCallDemoConfigReadsPluginAndHostConfigValues' -count=1`通过；`cd apps/lina-core && go test ./internal/service/apidoc -count=1`通过；`cd apps/lina-core && go test ./pkg/i18nresource -count=1`通过；`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`通过；`go test ./hack/tools/linactl/internal/wasmbuilder -run 'Test.*Demo|Test.*Official' -count=1`通过；`make wasm p=linapro-demo-dynamic out=../../temp/output`通过；`ruby -e 'require "yaml"; ...' plugin.yaml manifest/config/config.yaml manifest/config/config.example.yaml manifest/config/profile.yaml`通过；`node -e ...`解析插件运行时和 apidoc JSON 通过；`cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`通过；`cd hack/tests && pnpm test:validate`通过；`openspec validate replace-plugin-metadata-with-manifest --strict`通过；顶层`git diff --check`和`git -C apps/lina-plugins diff --check`通过。尝试运行`pnpm exec playwright test ../apps/lina-plugins/linapro-demo-dynamic/hack/tests/e2e/runtime/TC005-manifest-host-service-demo.ts --config=playwright.config.ts --project=chromium`，因本地`http://127.0.0.1:9120/admin/auth/login`未启动返回`ERR_CONNECTION_REFUSED`，未能执行浏览器端断言；已通过 TypeScript 编译、E2E 治理校验和构建验证覆盖测试文件结构与动态路由构建。

审查记录：已按`lina-review`执行反馈级审查，范围包含本反馈新增/修改的插件`plugin.yaml`、`manifest/config/profile.yaml`、Go API/controller/service/test、前端挂载页、插件中英文 i18n、中文 apidoc、E2E、POM、根目录和`manifest/`中英文 README 以及本任务记录。未发现阻塞问题。规则域结论：插件目录规范通过，动态插件`manifest`服务保留显式`methods: [get]`和最小`resources.paths`授权；API 契约通过，新增只读接口使用`GET`且 DTO 文档源文本为英文；后端 Go 通过，未改生成 DAO/DO/Entity，新增 host service 依赖来源清晰且测试覆盖读取逻辑；前端 UI 通过，展示为单页面局部面板且无逐项补查；`i18n`通过，插件启用多语言，运行时语言包和非英文 apidoc 资源已同步；文档通过，中英文 README 镜像事实一致；测试策略通过，后端单测、E2E 文件、TS 编译和治理校验均覆盖，浏览器 E2E 剩余风险仅为当前本地服务未启动；数据权限、缓存一致性、SQL、开发工具跨平台和模块启停均无新增运行时风险。

- [x] **FB-3**: Nightly plugin-full E2E 复用旧动态插件 active release 导致 manifest 配置预览缺失

执行记录：CI 失败任务为`2026-05-29`的`Test verification suite / E2E tests (plugin-full / plugins-2-of-5) / E2E tests`，日志显示`TC005-manifest-host-service-demo.ts`期望`Hello from dynamic plugin`但页面返回`Config PreviewNot returned`，同一分片中`TC004-menu-dynamic-permission-tree.ts`清理阶段还出现`BeforeDisable`生命周期 veto。根因是动态插件 E2E 在构建新 WASM 后，如果`linapro-demo-dynamic`已经处于已安装状态，会复用旧 active release 而跳过重新安装；旧 active release 可能不包含本变更新增的`manifest/config/config.yaml`资源，导致 manifest 配置预览缺失。同时，清理逻辑先调用普通`disable`再卸载，遇到动态插件生命周期保护时会触发 veto。已更新`apps/lina-plugins/linapro-demo-dynamic/hack/tests/e2e/runtime/TC005-manifest-host-service-demo.ts`，在测试前发现动态插件已安装时先执行带`force=true`的卸载并重新`syncPlugins`，随后安装当前构建产物，确保 active release 来自本次构建；恢复阶段在原本未安装或原本已安装但禁用的状态下，也通过强制卸载当前测试产物后按原状态恢复，避免普通`disable`触发生命周期 veto。同步更新`TC004-menu-dynamic-permission-tree.ts`，测试前清理和测试后恢复均绕开动态插件普通禁用路径，避免同一分片前序用例污染后续用例。

影响分析：已按`AGENTS.md`读取`openspec.md`、`architecture.md`、`plugin.md`、`testing.md`、`documentation.md`、`i18n.md`、`cache-consistency.md`、`data-permission.md`、`dev-tooling.md`和`backend-go.md`，并使用`lina-feedback`、`lina-e2e`、`lina-review`和`goframe-v2`技能。插件根目录无本地`AGENTS.md`。本反馈只修改动态插件自有 E2E 测试文件和 OpenSpec 任务记录，不修改 Go 生产代码、HTTP API、DTO、SQL、DAO、前端运行时代码、运行期依赖、构建脚本、CI 或插件清单。`i18n`无运行时文案、语言包、插件清单或 API 文档源文本变更；缓存一致性无缓存实现或失效策略变更，修复点是 E2E 状态准备阶段不再复用陈旧 active release；数据权限无业务数据读写接口或租户/组织可见性影响；开发工具跨平台无脚本或工具入口变更，沿用既有 Playwright/TypeScript 验证入口；SQL 无迁移、seed、mock 或安装卸载 SQL 变更；DI 来源无新增运行期依赖或构造函数变更；接口性能和`N+1`无影响。

验证记录：`cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`通过；`cd hack/tests && pnpm test:validate`通过，验证`239`个 E2E 文件；`openspec validate replace-plugin-metadata-with-manifest --strict`通过；顶层`git diff --check`和`git -C apps/lina-plugins diff --check`通过。本地未重新执行浏览器 E2E，因为当前环境未启动`http://127.0.0.1:9120`服务；修复范围为 E2E 安装/恢复状态隔离，已通过类型检查、E2E 治理校验和 CI 日志根因匹配验证。

审查记录：已按`lina-review`执行反馈级审查，范围包含`TC004-menu-dynamic-permission-tree.ts`、`TC005-manifest-host-service-demo.ts`和本任务记录。未发现阻塞问题。规则域结论：OpenSpec 记录通过，反馈已记录根因、影响和验证后标记完成；插件目录规范通过，变更仅在插件自有 E2E 内部闭环；测试策略通过，两个测试文件保持模块归属和自包含状态准备，清理逻辑不依赖跨文件残留；架构、后端 Go、API、SQL、数据权限、缓存一致性、`i18n`、开发工具跨平台和前端运行时均无生产行为变更风险。剩余风险为本地未跑浏览器端断言，需由 GitHub Actions 的`plugin-full / plugins-2-of-5`分片复验。

- [x] **FB-4**: GitHub Actions `TC005-manifest-host-service-demo` 仍返回 manifest 配置预览缺失

执行记录：`2026-05-31`的`Test verification suite / E2E tests (plugin-full / plugins-2-of-5) / E2E tests`仍在`TC005-manifest-host-service-demo.ts`失败，页面`linapro-demo-dynamic-manifest-config-preview`返回`Config PreviewNot returned`。重新排查确认`plugin.dynamic.storagePath`已锚定仓库根目录，与 E2E 的`make wasm p=linapro-demo-dynamic out=../../temp/output`输出路径一致；`/plugins/sync`会重新扫描该目录下的动态 artifact，且 E2E 已在安装前强制卸载旧 active release 后重新同步和安装。根因是动态 WASM 构建器仍沿用旧的`manifest`资源过滤逻辑，只打包少量 YAML 资源并排除`manifest/sql/`、`manifest/i18n/`及部分`manifest/config/`资源，违反本变更规格中“动态 artifact 完整保留`manifest/`下实际资源”的要求，导致 active release manifest 资源视图不完整。已修改`collectManifestResources`，使`go:embed`和目录扫描回退都完整收集`manifest/`下实际文件，仍跳过`.`和`_`开头的治理隐藏目录或文件；同步更新构建器单测和官方动态插件构建测试，锁定`manifest/config/config.yaml`、`manifest/config/profile.yaml`、SQL、i18n 和非 YAML 资源均进入通用 manifest 原文资源视图。

影响分析：已按`AGENTS.md`读取`openspec.md`、`documentation.md`、`dev-tooling.md`、`backend-go.md`、`plugin.md`、`api-contract.md`、`testing.md`、`i18n.md`、`architecture.md`、`data-permission.md`、`cache-consistency.md`和`database.md`，并使用`lina-feedback`、`goframe-v2`、`lina-e2e`和`lina-review`技能。插件根目录无本地`AGENTS.md`，本反馈未修改插件目录文件。本反馈修改`linactl`动态 WASM 构建器和 Go 测试，属于开发工具、动态插件产物资源视图和后端 Go 测试变更；不新增 HTTP API、DTO、数据库迁移、DAO、前端运行时代码、运行期依赖或构造函数。`i18n`影响：打包产物的通用 manifest 原文视图会携带已有`manifest/i18n/`文件，仍不改变 i18n 专用聚合、缓存失效或翻译治理；缓存一致性影响：active release 资源视图仍绑定 artifact checksum 和既有 runtime 缓存失效，不新增缓存域；数据权限影响：仅改变插件自有打包资源读取视图，不读取业务数据、租户数据或组织数据；SQL 影响：SQL 文件仍由专用 SQL custom section 执行，通用 manifest 视图仅额外保留原文字节，不改变安装、卸载、mock 执行账本或幂等语义；开发工具跨平台影响：变更使用 Go 标准库`filepath.WalkDir`和既有跨平台构建入口，无 shell 平台依赖；DI 来源无新增运行期依赖；接口性能和`N+1`无影响。

验证记录：`go test ./hack/tools/linactl/internal/wasmbuilder -count=1`通过；`go test ./hack/tools/linactl/internal/wasmbuilder -run 'TestBuildRuntimeWasmArtifactFromSourceEmbedsDeclaredAssets|TestCollectManifestResourcesScansDirectoryFallback|TestPluginDemoDynamicRuntimeArtifactEmbedsReviewedAssets' -count=1`通过；`make wasm p=linapro-demo-dynamic out=../../temp/output`通过；`strings temp/output/linapro-demo-dynamic.wasm`静态检查确认产物包含`manifest/config/profile.yaml`、`manifest/sql/001-linapro-demo-dynamic-records.sql`和`manifest/i18n/zh-CN/plugin.json`；`cd apps/lina-core && go test ./internal/service/plugin/internal/runtime ./internal/service/plugin/internal/wasm -count=1`通过；`cd apps/lina-core && go test ./internal/service/plugin/internal/catalog -count=1`通过；`openspec validate replace-plugin-metadata-with-manifest --strict`通过；顶层`git diff --check`和`git -C apps/lina-plugins diff --check`通过。本地未执行浏览器 E2E，因为当前环境未启动`http://127.0.0.1:9120`服务；本次修复用构建器测试和产物静态检查直接覆盖 CI 失败的 active release manifest 资源缺失根因。

审查记录：已按`lina-review`执行反馈级审查，范围包含`wasmbuilder_embed.go`、`wasmbuilder_test.go`、`wasmbuilder_plugin_demo_test.go`和本任务记录。未发现阻塞问题。规则域结论：开发工具跨平台通过，变更使用 Go 标准库目录扫描并沿用既有`make wasm`入口；插件目录规范通过，动态 artifact 通用`manifest/`原文资源视图与源码目录保持一致，仍跳过`.`和`_`开头的治理隐藏项；后端 Go 与测试策略通过，变更包测试、自包含 fixture、官方动态插件构建测试和 artifact 静态检查均覆盖；OpenSpec 记录通过，FB-4 已记录根因、影响和验证后标记完成；`i18n`、SQL、缓存一致性、数据权限、HTTP API、前端运行时和 DI 来源均无新增运行时风险。剩余风险为本地未启动`http://127.0.0.1:9120`，未重跑浏览器 E2E；该路径将由 GitHub Actions 的`plugin-full / plugins-2-of-5`分片复验。

- [x] **FB-5**: GitHub Actions clean checkout 缺少被忽略的动态插件 manifest 配置 fixture

执行记录：`2026-06-01`的 GitHub Actions 任务`Test verification suite / E2E tests (plugin-full / plugins-2-of-5) / E2E tests`执行`pnpm test:module -- plugins -- --shard=2/5`时，`TC005-manifest-host-service-demo.ts`断言`linapro-demo-dynamic-manifest-config-preview`包含`Hello from dynamic plugin`失败，实际页面文本为`Config PreviewNot returned`。根因是`linapro-demo-dynamic/manifest/config/config.yaml`被`apps/lina-plugins/.gitignore`忽略且未纳入 submodule 版本控制，CI clean checkout 只包含`config.example.yaml`和`profile.yaml`；本地工作区存在被忽略的`config.yaml`掩盖了该问题。已更新`TC005-manifest-host-service-demo.ts`，在构建动态 WASM artifact 前，如果`config.yaml`缺失，则从已提交的`config.example.yaml`临时复制生成 fixture，并仅在本测试创建该文件时于失败或结束清理；本地已有用户配置时不覆盖、不删除。

影响分析：已按`AGENTS.md`读取`openspec.md`、`plugin.md`、`testing.md`、`documentation.md`、`i18n.md`、`dev-tooling.md`和`api-contract.md`，并使用`lina-feedback`和`lina-e2e`技能。插件根目录无本地`AGENTS.md`。本反馈只修改动态插件自有 E2E 测试文件和 OpenSpec 任务记录，不修改 Go 生产代码、HTTP API、DTO、SQL、DAO、前端运行时代码、运行期依赖、构建器、CI 或插件清单；DI 来源无新增依赖或构造函数变更。`i18n`无运行时文案、语言包、插件清单或 API 文档源文本变更；缓存一致性无缓存域、快照或失效策略变更；数据权限无业务数据读写接口、租户/组织边界或存在性暴露影响；开发工具跨平台无脚本或工具入口变更，测试 fixture 使用 Node `fs`与`path` API；SQL 无迁移、seed、mock、安装卸载 SQL 或执行账本影响；接口性能和`N+1`无影响。

验证记录：`gh run view 26731053471 --job 78775010061 --log | rg ...`确认失败命令、用例和`Config PreviewNot returned`断言输出；`git -C apps/lina-plugins ls-files ...config.yaml`确认`config.yaml`未被跟踪，`git -C apps/lina-plugins check-ignore -v linapro-demo-dynamic/manifest/config/config.yaml`确认其被`.gitignore`忽略；`cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`通过；`cd hack/tests && pnpm test:validate`通过，验证`242`个 E2E 文件；临时备份并移走本地忽略的`manifest/config/config.yaml`后，运行`cd hack/tests && pnpm exec playwright test ../apps/lina-plugins/linapro-demo-dynamic/hack/tests/e2e/runtime/TC005-manifest-host-service-demo.ts --workers=1`通过，验证 clean checkout 缺失 fixture 场景；`openspec validate replace-plugin-metadata-with-manifest --strict`通过；顶层`git diff --check`和`git -C apps/lina-plugins diff --check`通过。

审查记录：已按`lina-review`执行反馈级审查，范围包含`TC005-manifest-host-service-demo.ts`和本任务记录；`git status --short`另有根目录`README.md`、`README.zh-CN.md`预存改动，不属于本反馈审查范围且未修改。已读取`AGENTS.md`以及`openspec.md`、`plugin.md`、`testing.md`、`documentation.md`、`i18n.md`、`dev-tooling.md`和`api-contract.md`。未发现阻塞问题。规则域结论：OpenSpec 记录通过，FB-5 已先记录、后修复并补充根因、影响和验证；插件目录规范通过，变更保留在动态插件自有测试目录且插件根目录无本地`AGENTS.md`；测试策略通过，E2E 文件归属和编号不变，测试自包含准备缺失 fixture 并只清理自身创建文件，clean checkout 复现场景已由单个 Playwright 用例验证；文档治理通过，仅修改 OpenSpec 任务记录；`i18n`、缓存一致性、数据权限、SQL、HTTP API、前端运行时、开发工具入口、DI 来源和接口性能均无生产行为变更风险。剩余风险为仍需 GitHub Actions 的`plugin-full / plugins-2-of-5`分片在云端 clean checkout 中复验。
