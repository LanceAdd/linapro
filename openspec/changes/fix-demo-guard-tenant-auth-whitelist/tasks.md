## Tasks

- [x] 1. 修复`linapro-ops-demo-guard`最小会话白名单，覆盖`linapro-tenant-core`插件挂载路径下的租户选择与租户切换。
- [x] 2. 补充插件中间件单元测试，覆盖插件挂载路径放行和其他插件挂载写请求继续拒绝。
- [x] 3. 更新插件自有 E2E 覆盖，验证启用演示保护后租户用户可选择租户登录，admin 可通过插件挂载路径切换租户。
- [x] 4. 运行 OpenSpec、Go、E2E 治理和格式验证，并完成反馈审查记录。

## Feedback

- [x] **FB-1**: `linapro-ops-demo-guard`启用后误拦截租户核心插件挂载路径的租户选择和租户切换
- [x] **FB-2**: 删除`linapro-ops-demo-guard`中已无真实 HTTP 入口的宿主租户选择和租户切换白名单

### FB-1 执行记录

- 根因：`apps/lina-vben/apps/web-antd/src/api/tenant/index.ts`通过`pluginApiPath("linapro-tenant-core", "auth/select-tenant")`和`pluginApiPath("linapro-tenant-core", "auth/switch-tenant")`调用租户核心插件挂载路径`/x/linapro-tenant-core/api/v1/auth/select-tenant`和`/x/linapro-tenant-core/api/v1/auth/switch-tenant`。`linapro-ops-demo-guard`中间件只对白名单中的宿主路径`/api/v1/auth/select-tenant`和`/api/v1/auth/switch-tenant`放行，启用后会把插件挂载路径下的两个`POST`误判为普通写请求并返回只读拒绝，导致租户用户无法选择租户登录，admin 也无法切换租户。
- 修复：在`linapro-ops-demo-guard`中间件中新增两个精确会话白名单路径：`/x/linapro-tenant-core/api/v1/auth/select-tenant`和`/x/linapro-tenant-core/api/v1/auth/switch-tenant`；白名单仍只接受`POST`且只匹配这两个路径，不放行其他`/x/linapro-tenant-core/api/v1/**`写请求。同步更新中间件测试服务和断言，证明租户核心插件会话路径放行、`/x/linapro-tenant-core/api/v1/platform/tenants`仍被拒绝；更新插件自有 E2E `TC001`，改用前端实际插件挂载路径验证租户选择和租户切换。
- 修改文件：`apps/lina-plugins/linapro-ops-demo-guard/backend/internal/service/middleware/middleware_guard.go`、`apps/lina-plugins/linapro-ops-demo-guard/backend/internal/service/middleware/middleware_guard_test.go`、`apps/lina-plugins/linapro-ops-demo-guard/hack/tests/e2e/lifecycle/TC001-linapro-ops-demo-guard.ts`、`openspec/changes/fix-demo-guard-tenant-auth-whitelist/**`。
- 影响分析：已按`AGENTS.md`读取`openspec.md`、`documentation.md`、`architecture.md`、`plugin.md`、`backend-go.md`、`api-contract.md`、`testing.md`、`i18n.md`、`cache-consistency.md`和`data-permission.md`，并使用`lina-feedback`、`lina-e2e`、`lina-review`和`goframe-v2`技能。插件根目录无本地`AGENTS.md`。本反馈修改官方源码插件中间件、插件自有 E2E 和 OpenSpec 记录，不修改宿主核心 API、租户核心插件 API、数据库、SQL、DAO、运行期依赖、构建脚本、CI、前端运行时代码或用户可见文案。
- `i18n`影响：无运行时用户可见文案、菜单、按钮、API 文档源文本、错误消息、插件清单或语言包资源变更。
- 缓存一致性影响：不新增缓存、快照、失效、刷新或跨节点同步逻辑；仍复用既有插件启用状态读取。
- 数据权限影响：不新增数据库读写、列表、详情、导出、下载、聚合或租户/组织数据可见性逻辑；放行后仍由原认证链路校验`preToken`、membership、JWT 和租户切换边界。
- 开发工具跨平台影响：不修改脚本、`Makefile`、`make.cmd`、`linactl`、CI 或构建入口；验证沿用现有 Go、OpenSpec 和 Playwright 治理命令。
- DI 来源检查：不新增运行期依赖、构造函数参数或服务实例；中间件仍复用既有`EnablementReader`读取插件启用状态。
- 接口性能：请求分类为常量数量的字符串规范化和精确匹配，不访问数据库、不产生批量装配或`N+1`风险。
- 验证记录：`openspec validate fix-demo-guard-tenant-auth-whitelist --strict`通过；`pnpm -C hack/tests test:validate`通过，验证`239`个 E2E 文件；`pnpm -C hack/tests exec tsc -p tsconfig.json --noEmit --pretty false`通过；`git diff --check && git -C apps/lina-plugins diff --check`通过；使用临时`go.work`包含`apps/lina-core`和`apps/lina-plugins/linapro-ops-demo-guard`后，`go test ./backend/internal/service/middleware -count=1`通过；同样使用临时`go.work`运行`go test ./backend/... -count=1`通过。目标 E2E 命令`pnpm -C hack/tests test:module -- plugin:linapro-ops-demo-guard --grep "TC-1g"`已执行，但当前本地服务中`linapro-ops-demo-guard`未启用，既有前置条件将`TC-1g`标记为 skipped；该文件的类型检查和 E2E 治理校验已通过，启用态路径行为由中间件单元测试直接覆盖。
- 审查记录：已按`lina-review`执行反馈级审查，范围包含本反馈修改的插件中间件、测试、E2E 和 OpenSpec 记录。未发现阻塞问题。规则域结论：OpenSpec 记录通过，反馈已先记录根因再修复；插件目录规范通过，变更局限于官方源码插件目录且插件根目录无本地`AGENTS.md`；后端 Go 通过，未修改生成代码、构造函数或服务依赖；测试策略通过，单元测试复现并验证原问题，E2E 已更新为真实插件挂载路径且通过静态治理；API 契约无新增或修改 HTTP API；数据权限、缓存一致性、`i18n`、SQL、开发工具跨平台和前端 UI 运行时代码均无新增风险。剩余风险：当前本地服务未启用`linapro-ops-demo-guard`，浏览器端启用态断言未实际命中，需要在启用演示保护的 E2E 环境或 CI 配置中复验。

### FB-2 执行记录

- 根因：用户明确指出项目不需要考虑兼容性后，重新检索确认`apps/lina-core/api`和`apps/lina-core/internal/controller`中没有`select-tenant`或`switch-tenant`HTTP 路由定义；当前 HTTP 入口由`linapro-tenant-core`源码插件在`/x/linapro-tenant-core/api/v1/auth/select-tenant`和`/x/linapro-tenant-core/api/v1/auth/switch-tenant`提供。`/api/v1/auth/select-tenant`和`/api/v1/auth/switch-tenant`只剩历史白名单意义，继续保留会扩大演示保护绕过面。宿主`auth.Service`中的`SelectTenant`和`SwitchTenant`仍是插件能力适配使用的内部契约，不属于应删除的历史 HTTP 路由。
- 修复：删除`linapro-ops-demo-guard`中`demoControlAuthSelectTenantPath`和`demoControlAuthSwitchTenantPath`两个宿主历史路径常量，并从最小会话白名单中移除；更新单元测试，移除旧路径放行断言，新增`TestGuardRejectsLegacyHostTenantSessionPaths`验证`POST /api/v1/auth/select-tenant`和`POST /api/v1/auth/switch-tenant`在演示保护启用时会被拒绝；同步更新增量规范、设计和提案，明确当前只支持租户核心插件挂载路径。
- 修改文件：`apps/lina-plugins/linapro-ops-demo-guard/backend/internal/service/middleware/middleware_guard.go`、`apps/lina-plugins/linapro-ops-demo-guard/backend/internal/service/middleware/middleware_guard_test.go`、`openspec/changes/fix-demo-guard-tenant-auth-whitelist/proposal.md`、`openspec/changes/fix-demo-guard-tenant-auth-whitelist/design.md`、`openspec/changes/fix-demo-guard-tenant-auth-whitelist/specs/demo-control-guard/spec.md`、`openspec/changes/fix-demo-guard-tenant-auth-whitelist/tasks.md`。
- 影响分析：已按`AGENTS.md`读取`openspec.md`、`documentation.md`、`architecture.md`、`plugin.md`、`backend-go.md`、`testing.md`、`i18n.md`、`cache-consistency.md`和`data-permission.md`，并使用`lina-feedback`、`lina-review`和`goframe-v2`技能。插件根目录无本地`AGENTS.md`。本反馈删除历史白名单，不修改宿主 API、租户核心插件 API、数据库、SQL、DAO、运行期依赖、构建脚本、CI、前端运行时代码或用户可见文案。
- `i18n`影响：无运行时用户可见文案、菜单、按钮、API 文档源文本、错误消息、插件清单或语言包资源变更。
- 缓存一致性影响：不新增缓存、快照、失效、刷新或跨节点同步逻辑；仍复用既有插件启用状态读取。
- 数据权限影响：不新增数据库读写、列表、详情、导出、下载、聚合或租户/组织数据可见性逻辑；删除历史白名单只减少中间件绕过面。
- 开发工具跨平台影响：不修改脚本、`Makefile`、`make.cmd`、`linactl`、CI 或构建入口。
- DI 来源检查：不新增运行期依赖、构造函数参数或服务实例。
- 接口性能：请求分类仍为常量数量字符串匹配，不访问数据库、不产生批量装配或`N+1`风险。
- 验证记录：`openspec validate fix-demo-guard-tenant-auth-whitelist --strict`通过；`pnpm -C hack/tests test:validate`通过，验证`239`个 E2E 文件；`pnpm -C hack/tests exec tsc -p tsconfig.json --noEmit --pretty false`通过；使用临时`go.work`包含`apps/lina-core`和`apps/lina-plugins/linapro-ops-demo-guard`后，`go test ./backend/internal/service/middleware -count=1`通过。
- 审查记录：已按`lina-review`执行反馈级审查，范围包含本反馈修改的插件中间件、单元测试和 OpenSpec 记录。未发现阻塞问题。规则域结论：OpenSpec 记录通过，反馈已先记录根因再修复；插件目录规范通过；后端 Go 通过，未修改生成代码、构造函数或服务依赖；测试策略通过，单元测试覆盖历史路径拒绝和当前插件路径放行；数据权限、缓存一致性、`i18n`、SQL、开发工具跨平台和前端 UI 运行时代码均无新增风险。
