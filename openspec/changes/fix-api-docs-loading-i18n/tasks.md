## Feedback

- [x] **FB-1**: 接口文档页内容较多时缺少加载 Loading 状态
- [x] **FB-2**: 中文语言环境下部分接口标题仍显示英文
- [x] **FB-3**: 中文接口标题本地化修复缺少完整 E2E 覆盖

### FB-1 根因

接口文档页面由 `apps/lina-vben/apps/web-antd/public/stoplight/apidocs.html` 中的独立 iframe 页面承载。原实现会立即追加 `<elements-api>` 并等待 `web-components.min.js` 加载和 `/api.json` 解析渲染，但在 Stoplight Elements 完成侧边栏内容渲染前，iframe 内没有任何可见占位或状态提示。接口数量较多时，脚本加载和 OpenAPI 文档解析时间变长，用户看到的是空白区域，容易误判为页面失败。

### FB-1 影响分析

- 修改范围：`apps/lina-vben/apps/web-antd/public/stoplight/apidocs.html`、`hack/tests/e2e/about/TC001-api-docs-page.ts`、`openspec/changes/fix-api-docs-loading-i18n/specs/system-api-docs/spec.md`。
- 用户可观察影响：iframe 内新增 Loading、长耗时提示和脚本加载失败提示；Stoplight 内容渲染完成后自动隐藏 Loading。
- i18n 影响：新增接口文档静态页运行时可见文案，已在页面内按 `lang` 参数维护 `zh-CN` 与 `en-US` 文案；不复用前端运行时语言包，符合接口文档资源隔离。
- API 契约影响：无 HTTP 路由、DTO、OpenAPI 响应结构或调用频次变更；仅改变 iframe 静态页面加载反馈。
- 后端 Go 影响：无 Go 生产代码影响。
- 数据权限影响：无数据读取、写入、导出、聚合或权限边界变化。
- 缓存一致性影响：无缓存、快照、失效或集群一致性变化。
- 开发工具跨平台影响：无 Makefile、脚本、CI、代码生成或 `linactl` 变更。
- 测试策略：更新宿主 E2E `TC001-api-docs-page.ts`，新增 `TC001j` 使用延迟 `/api.json` 的最小 OpenAPI 文档覆盖 Loading 显示与隐藏。
- 已读取规则：`openspec.md`、`documentation.md`、`frontend-ui.md`、`testing.md`、`i18n.md`、`api-contract.md`、`backend-go.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`cache-consistency.md`。

### FB-1 验证

- `cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`：通过。
- `cd hack/tests && pnpm test:validate`：通过。
- `cd hack/tests && pnpm exec playwright test e2e/about/TC001-api-docs-page.ts --config=playwright.config.ts --project=chromium`：10 项通过。

### FB-2 根因

`/api.json` 的请求语言上下文已正确接入 `lang` 参数和 `Content-Language`，基础 `apidoc` 翻译资源覆盖校验也存在。实际回退风险在 OpenAPI 操作级本地化键推断：GET/DELETE 等静态接口没有 request body 时，原逻辑优先通过请求 DTO schema 或 `g.Meta dc` 描述反推 operation key。多个接口使用相同 `dc` 描述时会被判定为歧义并删除索引，导致这些接口只能回退到路径派生键；生产翻译资源维护的是 DTO 稳定键，所以部分中文接口标题仍显示英文。

### FB-2 影响分析

- 修改范围：`apps/lina-core/internal/service/apidoc/apidoc_builder.go`、`apps/lina-core/internal/service/apidoc/apidoc_i18n.go`、`apps/lina-core/internal/service/apidoc/apidoc_builder_test.go`、`openspec/changes/fix-api-docs-loading-i18n/specs/system-api-docs/spec.md`。
- 宿主边界：修复 `lina-core` 接口文档生成这一框架通用能力，不绑定具体工作台页面结构。
- API 契约影响：不新增或修改业务 API、DTO、HTTP 方法、权限标签或 OpenAPI 公开字段；内部 `x-lina-apidoc-operation-key` 仅用于本地化并在返回文档前删除。
- 接口性能影响：无数据库访问、聚合、列表装配或额外 HTTP 调用；只在已有 OpenAPI 构建内对每个静态路由记录一次 DTO key，复杂度随路由数线性且不产生外部 I/O。
- 后端 Go 影响：修改 `apidoc` service 内部构建与本地化逻辑，无新增运行期依赖、无 DI 来源变化。
- 数据权限影响：无业务数据读取、写入或可见性变化；`/api.json` 仍是接口元数据文档。
- 缓存一致性影响：未新增或修改 apidoc 翻译缓存、失效机制或集群同步策略。
- 插件影响：源码插件静态路由同样受益于 DTO 稳定 key；未修改 `apps/lina-plugins/<plugin-id>/` 文件，动态插件路径派生 key 逻辑保持不变。
- i18n 影响：修复接口文档 `zh-CN` operation `tags` 与 `summary` 的稳定 key 解析；未新增生产翻译资源。
- 开发工具跨平台影响：无 Makefile、脚本、CI、代码生成或 `linactl` 变更。
- 测试策略：更新 `TestBuildLocalizesOpenAPIForRequestLocale`，构造两个共享 `dc` 的 GET 静态接口，验证中文摘要仍按 DTO key 翻译，且内部 marker 不泄露。
- 已读取规则：`openspec.md`、`documentation.md`、`architecture.md`、`api-contract.md`、`backend-go.md`、`testing.md`、`i18n.md`、`data-permission.md`、`plugin.md`、`cache-consistency.md`。

### FB-2 验证

- `cd apps/lina-core && go test ./internal/service/apidoc -run 'TestBuildLocalizesOpenAPIForRequestLocale|TestBuildProjectsHostAndEnabledPluginRoutes' -count=1`：通过。
- `cd apps/lina-core && go test ./internal/service/apidoc -count=1`：通过。
- `cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`：通过。
- `cd hack/tests && pnpm exec playwright test e2e/about/TC001-api-docs-page.ts --config=playwright.config.ts --project=chromium`：10 项通过。

### FB-3 根因

此前 `TC001i` 只抽样断言了 `用户登录`、`身份认证` 和公共响应字段等少量中文文案，无法系统发现 `/api.json?lang=zh-CN` 中仍有部分 operation `tags` 或 `summary` 回退英文。补充运行态扫描后进一步确认，宿主 `/api/v1/...` 接口已被 DTO 稳定 key 覆盖，但启用 `i18n` 的源码插件接口挂载在 `/x/<plugin-id>/...` 下时，本地化逻辑先把该路径识别为动态插件路径，导致源码插件 DTO key 和插件自有 `apidoc` 翻译资源未被优先使用。

### FB-3 影响分析

- 修改范围：`apps/lina-core/internal/service/apidoc/apidoc_i18n.go`、`apps/lina-core/internal/service/apidoc/apidoc_route_text_test.go`、`hack/tests/e2e/about/TC001-api-docs-page.ts`、`apps/lina-plugins/linapro-org-core/hack/tests/e2e/host-integration/TC008-api-docs-localization.ts`、`openspec/changes/fix-api-docs-loading-i18n/specs/system-api-docs/spec.md`。
- 宿主边界：修复 `lina-core` 接口文档本地化 key 选择顺序，保持框架通用能力；插件特定 E2E 按所有权放在 `linapro-org-core` 插件目录内。
- API 契约影响：不新增或修改业务 HTTP API、DTO、权限标签或公开 OpenAPI 字段；仅调整 OpenAPI 文档构建后的本地化文本选择。
- 接口性能影响：无数据库访问、无额外 HTTP 调用、无前端瀑布式调用；只改变已有 operation key 的优先级判断。
- 后端 Go 影响：修改 `apidoc` service 内部逻辑，无新增运行期依赖、无 DI 来源变化。
- 数据权限影响：无业务数据读取、写入、导出、聚合或可见性边界变化；`/api.json` 仍只返回接口元数据。
- 缓存一致性影响：未修改 apidoc 翻译缓存结构、失效触发、跨实例同步或动态插件资源加载；仅使用已加载 catalog。
- 插件影响：新增 `linapro-org-core` 插件自有 E2E；该插件根目录不存在本地 `AGENTS.md`，已按项目顶层插件规范处理。未修改插件生产代码或插件资源。
- i18n 影响：补强宿主与启用 `i18n` 的源码插件 API 文档标题本地化验证；未新增生产翻译资源。
- 开发工具跨平台影响：未修改 Makefile、脚本、CI、代码生成或 `linactl`。
- 测试策略：宿主 `TC001i` 扫描所有宿主 `/api/v1/...` operation 标题是否含中文，并精确断言 `/api/v1/user`；插件新增 `TC008` 扫描 `linapro-org-core` 的 `/x/linapro-org-core/...` operation 标题是否含中文，并精确断言部门、岗位接口。
- 已读取规则：`openspec.md`、`documentation.md`、`frontend-ui.md`、`testing.md`、`i18n.md`、`api-contract.md`、`backend-go.md`、`architecture.md`、`data-permission.md`、`plugin.md`、`cache-consistency.md`、`dev-tooling.md`。

### FB-3 验证

- `cd apps/lina-core && go test ./internal/service/apidoc -run 'TestOperationBaseKeyPrefersStaticMarkerForPluginNamespace|TestOperationBaseKeyIgnoresDynamicOperationID|TestBuildLocalizesOpenAPIForRequestLocale' -count=1`：通过。
- `cd apps/lina-core && go test ./internal/service/apidoc -count=1`：通过。
- `cd hack/tests && ./node_modules/.bin/tsc -p tsconfig.json --noEmit --pretty false`：通过。
- `cd hack/tests && pnpm test:validate`：通过，校验 241 个 E2E 文件。
- `cd hack/tests && pnpm exec playwright test apps/lina-plugins/linapro-org-core/hack/tests/e2e/host-integration/TC008-api-docs-localization.ts --config=playwright.config.ts --project=chromium`：1 项通过。
- `cd hack/tests && pnpm exec playwright test e2e/about/TC001-api-docs-page.ts --config=playwright.config.ts --project=chromium`：10 项通过。

### 综合验证

- `openspec validate fix-api-docs-loading-i18n --strict`：通过。
- `git diff --check`：通过。
