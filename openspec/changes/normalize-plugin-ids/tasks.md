## 1. 插件 ID 治理契约

- [ ] 1.1 新增或抽取插件 ID 解析/校验组件，覆盖 `<author>-<domain>-<capability>`、domain 白名单、`core` 官方保留 capability 和 64 字符长度限制
- [ ] 1.2 将 manifest ID、动态 artifact manifest ID、插件依赖 ID 和源码插件注册 ID 校验统一接入新的插件 ID 治理组件
- [ ] 1.3 为插件 ID 解析、白名单、`core` 保留、旧 ID 拒绝、依赖 ID 拒绝和长度限制补充后端单元测试
- [ ] 1.4 更新错误信息、接口文档和中英文 i18n 资源，确保插件 ID 校验失败返回稳定错误码、messageKey 和英文 fallback

## 2. 官方插件破坏式改名

- [ ] 2.1 按映射重命名官方插件目录、`plugin.yaml` ID、README、manifest README 和示例文本
- [ ] 2.2 同步更新官方插件 Go module 名称、import 路径、`apps/lina-plugins/go.mod` replace、`apps/lina-plugins/lina-plugins.go` 和 GoFrame 生成配置
- [ ] 2.3 更新官方插件源码注册常量和内部常量，包括 `linapro-ops-demo-guard` 的中间件启用态检测
- [ ] 2.4 更新宿主官方插件常量、稳定菜单父级映射、`orgcap.ProviderPluginID`、`tenantcap.ProviderPluginID`、启动一致性检查和 provider 检测逻辑
- [ ] 2.5 更新 `plugin.autoEnable` 默认/开发/镜像配置和配置解析测试中的官方插件 ID

## 3. 插件自有存储与生成代码

- [ ] 3.1 将官方插件自有 SQL 表、索引、约束、mock 数据、uninstall SQL 和 SQL 注释按新插件 ID snake_case 范围重命名
- [ ] 3.2 重新生成或同步官方插件本地 DAO/DO/Entity，确保生成配置、表名、服务代码和测试 fixture 与新表名一致
- [ ] 3.3 更新动态插件 host service data 资源授权表名、示例 SQL、资源声明和相关权限审查展示测试
- [ ] 3.4 补充 SQL/DAO 静态验证，确认官方插件运行时代码不再引用旧插件表名

## 4. 派生运行时身份同步

- [ ] 4.1 更新官方插件 manifest 菜单 key、parent_key、权限字符串、路由 path、cron handlerRef、job i18n key 和菜单 i18n key
- [ ] 4.2 更新动态插件 artifact 构建逻辑和测试 fixture，确保 `.wasm` 文件名、manifest、`/plugin-assets/<id>/...` 和 `/api/v1/extensions/<id>/...` 使用新 ID
- [ ] 4.3 更新插件管理、菜单、用户首页、API 文档、任务管理、运行时路由、host service 授权和 scheduler 相关测试中的插件 ID
- [ ] 4.4 增加旧官方 ID 残留扫描，排除归档历史说明后，运行时代码、配置、测试和活跃 OpenSpec 文档不得再使用旧 ID

## 5. i18n、apidoc 与文档

- [ ] 5.1 更新官方插件 `manifest/i18n/<locale>/*.json` 中的 `plugin.<plugin-id>.`、菜单 key、错误 key、job key 和页面文案引用
- [ ] 5.2 更新官方插件和宿主 `manifest/i18n/<locale>/apidoc/**/*.json` 中的 `plugins.<plugin_id_snake_case>.` namespace
- [ ] 5.3 更新根 README、插件工作区 README、插件开发说明、i18n README、E2E README 和相关 OpenSpec baseline 示例中的插件 ID
- [ ] 5.4 运行 JSON 校验和 runtime i18n 治理扫描，确认无硬编码旧插件 ID、无无效翻译键、无旧 apidoc namespace

## 6. 前端与 E2E

- [ ] 6.1 更新前端插件管理、动态页面、菜单路由、用户首页、测试 page object 和 fixture 中的官方插件 ID
- [ ] 6.2 更新或新增 Playwright E2E，覆盖插件列表显示新官方 ID、`linapro-ops-demo-guard` 启用后写请求阻断、动态插件资产路径和扩展 API 新路径
- [ ] 6.3 按 `lina-e2e` 规范为新增 E2E 分配 TC ID，并运行 `pnpm -C hack/tests test:validate`
- [ ] 6.4 运行前端 typecheck 和受影响 E2E，确认页面无旧插件 ID、无原始 i18n key、无动态资源路径断裂

## 7. 后端编译、缓存与一致性验证

- [ ] 7.1 运行受影响宿主包测试，至少覆盖插件 catalog、runtime、integration、jobhandler、jobmgmt、menu、orgcap、tenantcap 和启动绑定包
- [ ] 7.2 运行官方插件后端包测试，覆盖 `linapro-content-notice`、`linapro-monitor-loginlog`、`linapro-monitor-operlog`、`linapro-monitor-online`、`linapro-monitor-server`、`linapro-tenant-core`、`linapro-org-core`、`linapro-demo-source` 和 `linapro-ops-demo-guard`
- [ ] 7.3 重新构建并验证 `linapro-demo-dynamic` 动态插件 artifact，确认 manifest、host service、frontend、SQL、i18n 和 apidoc 均使用新 ID
- [ ] 7.4 明确记录缓存一致性结论：本变更不新增缓存类型，所有状态、菜单、路由、cron、i18n 和 apidoc 刷新继续使用插件 ID scope 精确失效，集群模式沿用既有广播/共享修订号机制
- [ ] 7.5 明确记录数据权限结论：本变更不新增业务数据访问接口；官方插件改名后的列表、详情、导出、写操作和插件 host service 访问仍复用既有租户与数据权限边界

## 8. OpenSpec 与最终审查

- [ ] 8.1 运行 `openspec validate normalize-plugin-ids --strict`
- [ ] 8.2 运行 `git diff --check` 覆盖本变更涉及的代码、测试、i18n、SQL 和 OpenSpec 文档
- [ ] 8.3 汇总所有 Go 测试、前端 typecheck、E2E、i18n 扫描、旧 ID 残留扫描和 OpenSpec 校验结果到本任务记录
- [ ] 8.4 完成实现后调用 `lina-review`，重点审查官方插件 ID 映射一致性、旧 ID 残留、manifest 校验覆盖、i18n/apidoc namespace、缓存作用域、数据权限边界和后端 Go 编译门禁
