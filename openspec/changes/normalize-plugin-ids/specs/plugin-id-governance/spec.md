## ADDED Requirements

### Requirement: 插件 ID 必须使用结构化命名

系统 SHALL 要求所有插件 ID 使用 `<author>-<domain>-<capability>` 结构化命名。`author`、`domain` 和 `capability` 均使用小写字母、数字和 hyphen 组成；`author` 与 `domain` 必须各自是单个 slug，`capability` 可以由一个或多个 kebab-case 单词组成。插件 ID 总长度 MUST 不超过运行时 `plugin_id` 字段允许的 64 字符。

#### Scenario: 接受结构化插件 ID
- **WHEN** 插件 manifest 声明 `id: linapro-content-notice`
- **THEN** 系统解析 `author=linapro`
- **AND** 系统解析 `domain=content`
- **AND** 系统解析 `capability=notice`
- **AND** 插件 ID 结构校验通过

#### Scenario: 接受多词 capability
- **WHEN** 插件 manifest 声明 `id: linapro-ops-demo-guard`
- **THEN** 系统解析 `author=linapro`
- **AND** 系统解析 `domain=ops`
- **AND** 系统解析 `capability=demo-guard`
- **AND** 插件 ID 结构校验通过

#### Scenario: 拒绝缺少结构段的插件 ID
- **WHEN** 插件 manifest 声明 `id: demo-control`
- **THEN** 系统拒绝该 manifest
- **AND** 错误说明插件 ID 必须使用 `<author>-<domain>-<capability>` 结构

### Requirement: 插件 domain 必须来自白名单

系统 SHALL 只接受宿主维护的插件 domain 白名单。初始白名单 MUST 包含 `tenant`、`org`、`iam`、`content`、`monitor`、`ops`、`storage`、`workflow`、`integration`、`ai`、`dev` 和 `demo`。新增 domain MUST 通过 OpenSpec 变更说明领域边界后进入白名单，不得由插件作者在 manifest 中自行定义。

#### Scenario: 接受白名单 domain
- **WHEN** 插件 manifest 声明 `id: acme-integration-feishu`
- **THEN** 系统识别 `integration` 为白名单 domain
- **AND** 插件 ID 结构校验通过

#### Scenario: 拒绝非白名单 domain
- **WHEN** 插件 manifest 声明 `id: acme-random-report`
- **THEN** 系统拒绝该 manifest
- **AND** 错误说明 `random` 不是允许的插件 domain

### Requirement: `core` capability 必须仅用于官方基础能力插件

系统 SHALL 将 `core` 作为官方保留 capability。只有 LinaPro 官方基础能力插件可以使用 `linapro-<domain>-core`，第三方插件和非基础能力官方插件 MUST 使用更具体的 capability 名称。

#### Scenario: 接受官方基础能力插件
- **WHEN** 插件 manifest 声明 `id: linapro-tenant-core`
- **THEN** 系统接受该官方基础能力插件 ID

#### Scenario: 拒绝第三方 core capability
- **WHEN** 插件 manifest 声明 `id: acme-org-core`
- **THEN** 系统拒绝该 manifest
- **AND** 错误说明 `core` 是 LinaPro 官方保留 capability

### Requirement: 官方插件 ID 必须使用规范化映射

系统 SHALL 将官方插件 ID 规范化为以下映射，并不得继续在运行时配置、manifest、源码注册、菜单、权限、cron、i18n、apidoc、测试或文档中使用旧官方 ID。

| 旧 ID | 新 ID |
| --- | --- |
| `content-notice` | `linapro-content-notice` |
| `monitor-loginlog` | `linapro-monitor-loginlog` |
| `monitor-operlog` | `linapro-monitor-operlog` |
| `monitor-online` | `linapro-monitor-online` |
| `monitor-server` | `linapro-monitor-server` |
| `multi-tenant` | `linapro-tenant-core` |
| `org-center` | `linapro-org-core` |
| `plugin-demo-dynamic` | `linapro-demo-dynamic` |
| `plugin-demo-source` | `linapro-demo-source` |
| `demo-control` | `linapro-ops-demo-guard` |

#### Scenario: 官方插件清单使用新 ID
- **WHEN** 宿主扫描 `apps/lina-plugins/linapro-org-core/plugin.yaml`
- **THEN** manifest ID 为 `linapro-org-core`
- **AND** 宿主不得发现 `org-center` 作为官方插件 ID

#### Scenario: 官方自动启用配置使用新 ID
- **WHEN** 宿主读取 `plugin.autoEnable`
- **THEN** 官方插件项使用规范化新 ID
- **AND** 配置中不得继续使用 `multi-tenant`、`org-center` 或其他旧官方 ID

### Requirement: 插件运行时身份必须只使用规范化 ID

系统 SHALL 在运行时身份边界只接受规范化插件 ID，不得为旧 ID 提供 alias、重定向或兼容查询。该边界包括插件管理 API、扩展 API、动态前端资产 URL、菜单 key、权限字符串、cron handlerRef、插件状态表、发布表、迁移表、资源引用表、节点状态表、插件 KV 状态表和 host service 授权记录。

#### Scenario: 新扩展 API 路径使用规范化 ID
- **WHEN** 动态插件 `linapro-demo-dynamic` 暴露扩展 API
- **THEN** 宿主公开路径使用 `/api/v1/extensions/linapro-demo-dynamic/...`
- **AND** 宿主不得通过 `/api/v1/extensions/plugin-demo-dynamic/...` 暴露同一插件

#### Scenario: 新动态资产路径使用规范化 ID
- **WHEN** 动态插件 `linapro-demo-dynamic` 提供前端资产
- **THEN** 宿主资产路径使用 `/plugin-assets/linapro-demo-dynamic/<version>/...`
- **AND** 宿主不得通过 `/plugin-assets/plugin-demo-dynamic/<version>/...` 暴露同一资产

#### Scenario: 新 cron handlerRef 使用规范化 ID
- **WHEN** 插件 `linapro-monitor-server` 注册内置定时任务
- **THEN** handlerRef 使用 `plugin:linapro-monitor-server/cron:<name>`
- **AND** 系统不得继续生成 `plugin:monitor-server/cron:<name>`

### Requirement: 仓库治理扫描必须验证插件 ID 一致性

系统 SHALL 提供自动化验证，确保插件目录名、manifest ID、源码插件注册 ID、动态 artifact manifest、依赖声明、菜单 key、运行时 i18n key、apidoc i18n key、配置和测试 fixture 使用同一个规范化插件 ID。验证失败时变更不得通过。

#### Scenario: 目录名与 manifest ID 不一致
- **WHEN** 插件目录为 `apps/lina-plugins/linapro-content-notice`
- **AND** 该目录的 `plugin.yaml` 声明 `id: content-notice`
- **THEN** 治理验证失败
- **AND** 错误指出目录名与 manifest ID 不一致

#### Scenario: i18n namespace 使用旧 ID
- **WHEN** 插件 `linapro-content-notice` 的运行时语言包包含 `plugin.content-notice.name`
- **THEN** 治理验证失败
- **AND** 错误指出运行时 i18n key 必须使用 `plugin.linapro-content-notice.` 前缀

#### Scenario: apidoc namespace 使用旧 ID
- **WHEN** 插件 `linapro-demo-dynamic` 的 apidoc 语言包包含 `plugins.plugin_demo_dynamic`
- **THEN** 治理验证失败
- **AND** 错误指出 apidoc key 必须使用 `plugins.linapro_demo_dynamic` 前缀
