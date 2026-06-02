## Feedback

- [x] **FB-1**: 部门管理页面在部门树为空时无法创建首个顶级部门

### FB-1 根因

后端 `POST /dept` 和服务层已经支持 `parentId=0` 创建顶级部门，且现有部门管理规范包含“创建根部门”场景。问题出在前端抽屉：`parentId` 使用 `TreeSelect` 且配置了必选校验，候选项只来自 `dept/tree` 或 `dept/exclude/{id}` 的返回结果。部门表为空时树数据为空，页面没有“顶级部门”候选项，也不会提交 `parentId=0`。

### FB-1 影响分析

- 修改范围：`apps/lina-plugins/linapro-org-core` 源码插件前端页面、插件自有 E2E 和本 OpenSpec 变更文档。
- 宿主边界：仅修复工作台前端适配层，不修改 `apps/lina-core` 核心领域契约、通用 service 语义或存储模型。
- API 契约影响：无 HTTP 路由、DTO、OpenAPI 元数据或响应结构变更；沿用现有 `parentId=0` 语义。
- 后端 Go 影响：无 Go 生产代码、Controller、Service、DAO 或运行期依赖变更。
- 数据权限影响：无新增或修改数据访问路径；创建后的租户边界继续由现有后端 `tenantFilter` 处理。
- 缓存一致性影响：无缓存、快照、失效或集群一致性变更。
- 数据库影响：无 SQL、Seed、Mock、DAO 生成或索引变更。
- i18n 影响：新增插件运行时 UI 文案“顶级部门”，需维护 `linapro-org-core` 插件自身 `zh-CN` 与 `en-US` 语言资源。
- 开发工具跨平台影响：无 Makefile、脚本、CI、代码生成或 `linactl` 变更。
- 测试策略：新增插件自有 E2E `apps/lina-plugins/linapro-org-core/hack/tests/e2e/dept/TC004-dept-empty-tree-root-create.ts` 覆盖空部门树下创建首个顶级部门。
- 已读取规则：`openspec.md`、`documentation.md`、`architecture.md`、`plugin.md`、`frontend-ui.md`、`testing.md`、`i18n.md`、`api-contract.md`、`backend-go.md`、`data-permission.md`。

## Tasks

- [x] 修复部门新增抽屉：为空树和普通新增场景提供“顶级部门”选项，默认提交 `parentId=0`
- [x] 补齐插件 `zh-CN` 与 `en-US` 运行时语言资源
- [x] 创建 E2E — TC004 部门空树创建顶级部门
- [x] 运行相关前端类型检查或静态验证、E2E 定向验证、`openspec validate fix-dept-root-create-empty-tree --strict`

### FB-1 验证记录

- `pnpm -C apps/lina-vben --filter @lina/web-antd typecheck`：通过。
- `npm exec -- tsc --noEmit --project tsconfig.json`（`hack/tests`）：通过。
- `npm run test:validate`（`hack/tests`）：通过。
- `npm exec -- playwright test ../apps/lina-plugins/linapro-org-core/hack/tests/e2e/dept/TC004-dept-empty-tree-root-create.ts --project=chromium`（`hack/tests`）：通过。
- `openspec validate fix-dept-root-create-empty-tree --strict`：通过。
