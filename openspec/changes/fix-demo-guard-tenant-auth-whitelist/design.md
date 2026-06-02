## 决策：在演示保护插件内补齐租户会话路径白名单

本次问题来自`linapro-ops-demo-guard`的请求分类逻辑，而不是宿主认证服务或`linapro-tenant-core`插件契约变化。前端通过`pluginApiPath("linapro-tenant-core", "auth/select-tenant")`和`pluginApiPath("linapro-tenant-core", "auth/switch-tenant")`生成`/x/linapro-tenant-core/api/v1/auth/*`路径，这是源码插件能力挂载路径的正常形态。

修复应在`linapro-ops-demo-guard`内收敛最小会话白名单，使其识别租户核心插件挂载路径下的`select-tenant`和`switch-tenant`动作，并删除已无真实 HTTP 入口的宿主`/api/v1/auth/select-tenant`与`/api/v1/auth/switch-tenant`历史白名单。该白名单只覆盖`POST`方法和当前真实会话路径，不放行其他`/x/{plugin}/api/v1/**`写请求，避免削弱演示只读保护。

## 边界判断

- 核心宿主边界：不修改`apps/lina-core`核心认证、插件路由或通用 service 契约；问题来源于官方演示保护插件的适配逻辑。
- 插件边界：变更只发生在`apps/lina-plugins/linapro-ops-demo-guard`内，不修改`linapro-tenant-core`业务实现或接口定义。
- 接口性能：请求分类为常量级字符串规范化和精确匹配，不引入数据库访问、远程调用或随数据量增长的装配路径。
- 数据权限：放行后仍由原认证/租户切换链路校验`preToken`、membership、token 和租户边界；中间件本身不读取或暴露业务数据。
- 缓存一致性：不新增缓存或失效路径，插件启用状态读取保持现状。
- `i18n`：不新增用户可见文案或语言资源。

## 测试策略

- 单元测试：扩展`linapro-ops-demo-guard`中间件测试，证明启用演示保护时租户核心插件挂载路径的`select-tenant`和`switch-tenant`可以通过，同时历史宿主租户会话路径和普通插件挂载写请求仍被拒绝。
- E2E：更新插件自有生命周期 E2E 的会话白名单用例，使用前端实际插件挂载路径验证租户用户选择租户和 admin 切换租户请求不会被只读保护拦截。
- 验证：运行插件 Go 测试、E2E TypeScript/治理校验和`openspec validate fix-demo-guard-tenant-auth-whitelist --strict`。
