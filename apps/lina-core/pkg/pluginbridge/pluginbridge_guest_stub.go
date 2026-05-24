//go:build !wasip1

// pluginbridge_guest_stub.go provides the non-WASI facade for dynamic-plugin
// guest host-service helpers.
//
// This file is compiled only for non-wasip1 builds. It keeps the historical
// root-package pluginbridge API available to ordinary Go builds, host-side
// tests, and compatibility callers while delegating every contract to the guest
// subcomponent. The delegated guest helpers resolve to unsupported stubs in
// this build target, so callers receive an explicit error instead of
// accidentally assuming that real WASI host calls are available.
//
// Development constraints:
//   - Keep this file as a thin facade of aliases and forwarded entrypoints.
//   - Add new guest host-service capabilities in the guest subcomponent first,
//     then expose root-package aliases here only when backward-compatible
//     pluginbridge entrypoints are required.
//   - Keep non-WASI behavior explicit and deterministic. Do not bypass the
//     unsupported stubs to access host resources in ordinary Go builds.
//
// 本文件为动态插件 guest 宿主服务辅助能力提供非 WASI 构建下的 facade。
//
// 本文件只在非 wasip1 构建目标下编译，用于在普通 Go 构建、宿主侧测试和
// 兼容调用场景中保留 pluginbridge 根包的历史 API。文件中的所有契约都委托给
// guest 子组件；在当前构建目标下，这些 guest 辅助入口会解析为 unsupported
// stub，因此调用方会收到明确错误，而不会误以为真实 WASI host call 可用。
//
// 开发约束：
//   - 保持本文件只承担别名和入口转发职责。
//   - 新增 guest 宿主服务能力时，应先在 guest 子组件中定义；只有需要兼容
//     pluginbridge 根包入口时，才在这里暴露别名。
//   - 非 WASI 行为必须保持显式且确定，不得绕过 unsupported stub 在普通 Go
//     构建中访问宿主资源。

package pluginbridge

import "lina-core/pkg/pluginbridge/guest"

type (
	// RuntimeHostService is the guest-side contract for runtime host capabilities.
	//
	// RuntimeHostService 是 guest 侧访问运行时宿主能力的契约，覆盖结构化日志、
	// 插件作用域状态、宿主时间、宿主生成 UUID 和宿主节点信息等基础运行时能力。
	RuntimeHostService = guest.RuntimeHostService

	// StorageHostService is the guest-side contract for governed object storage.
	//
	// StorageHostService 是 guest 侧访问受治理对象存储的契约，所有路径、覆盖
	// 策略和资源边界均由宿主 host service 授权模型约束。
	StorageHostService = guest.StorageHostService

	// HTTPHostService is the guest-side contract for governed outbound HTTP calls.
	//
	// HTTPHostService 是 guest 侧发起受治理外部 HTTP 请求的契约，避免动态插件
	// 绕过宿主网络授权、审计和资源边界直接访问外部网络。
	HTTPHostService = guest.HTTPHostService

	// DataHostService is the guest-side compatibility contract for governed data access.
	//
	// DataHostService 是 guest 侧兼容型结构化数据访问契约；新代码应优先使用
	// plugindb，但该别名仍保留给历史动态插件和桥接调用方。
	DataHostService = guest.DataHostService

	// CacheHostService is the guest-side contract for governed cache namespaces.
	//
	// CacheHostService 是 guest 侧访问受治理缓存命名空间的契约，缓存键、过期
	// 策略和可访问命名空间必须由宿主授权边界控制。
	CacheHostService = guest.CacheHostService

	// LockHostService is the guest-side contract for governed distributed locks.
	//
	// LockHostService 是 guest 侧访问受治理分布式锁的契约，插件必须通过宿主
	// 签发的 ticket 续租和释放锁，不能自行构造宿主内部锁实现。
	LockHostService = guest.LockHostService

	// ConfigHostService is the guest-side contract for plugin-scoped configuration.
	//
	// ConfigHostService 是 guest 侧读取插件作用域配置的契约，配置为只读宿主
	// 能力，调用方应把缺失值和类型转换失败作为显式状态处理。
	ConfigHostService = guest.ConfigHostService

	// NotifyHostService is the guest-side contract for governed notifications.
	//
	// NotifyHostService 是 guest 侧发送受治理通知的契约，通知渠道、模板和资源
	// 范围必须由宿主统一授权与审计。
	NotifyHostService = guest.NotifyHostService

	// CronHostService is the guest-side contract for dynamic-plugin cron declarations.
	//
	// CronHostService 是 guest 侧声明动态插件定时任务的契约，仅用于把 cron
	// 注册信息提交给宿主发现与调度边界，不承载插件业务逻辑。
	CronHostService = guest.CronHostService

	// HostConfigHostService is the guest-side contract for whitelisted host config.
	//
	// HostConfigHostService 是 guest 侧读取宿主公开配置白名单的契约，仅暴露
	// 被宿主明确允许的公共配置值，不能泄漏宿主私有运行时状态。
	HostConfigHostService = guest.HostConfigHostService

	// ManifestHostService is the guest-side contract for plugin manifest resources.
	//
	// ManifestHostService 是 guest 侧读取插件 manifest 资源的契约，用于访问
	// 插件包内受治理资源文件，不应用作任意宿主文件系统访问入口。
	ManifestHostService = guest.ManifestHostService

	// HostDBQueryResult preserves the legacy guest-side tabular query result shape.
	//
	// HostDBQueryResult 保留旧版 guest 数据库查询结果形态，仅用于尚未迁移到
	// 结构化 data host service 或 plugindb 的兼容调用。
	HostDBQueryResult = guest.HostDBQueryResult

	// DataListResult is the decoded compatibility result for governed list queries.
	//
	// DataListResult 是兼容型受治理列表查询的 guest 侧解码结果。
	DataListResult = guest.DataListResult

	// DataGetResult is the decoded compatibility result for governed get queries.
	//
	// DataGetResult 是兼容型受治理单条查询的 guest 侧解码结果。
	DataGetResult = guest.DataGetResult

	// DataMutationResult is the decoded compatibility result for governed mutations.
	//
	// DataMutationResult 是兼容型受治理数据变更的 guest 侧解码结果。
	DataMutationResult = guest.DataMutationResult

	// DataTransactionInput describes one compatibility transaction operation.
	//
	// DataTransactionInput 描述兼容型受治理数据事务中的单个操作。
	DataTransactionInput = guest.DataTransactionInput

	// DataTransactionResult is the decoded compatibility result for transactions.
	//
	// DataTransactionResult 是兼容型受治理数据事务的 guest 侧解码结果。
	DataTransactionResult = guest.DataTransactionResult
)

var (
	// Runtime returns the guest runtime host-service client for the current build target.
	// In non-WASI builds it returns an unsupported stub that reports real host
	// calls as unavailable; the wasip1 facade forwards the same entrypoint to
	// the real WASM host-service client.
	//
	// Runtime 返回当前构建目标下的 guest 运行时宿主服务客户端。非 WASI 构建中
	// 该入口返回 unsupported stub，用于明确提示真实 host call 不可用；wasip1
	// 构建中的同名入口会转发到真实 WASM host service 调用。
	Runtime = guest.Runtime

	// Storage returns the guest storage host-service client for governed object access.
	// In non-WASI builds it exists for compile-time compatibility and tests, and
	// it never accesses real host storage.
	//
	// Storage 返回受治理对象存储的 guest 宿主服务客户端。非 WASI 构建中该入口
	// 仅用于编译兼容和测试兜底，不会访问真实宿主存储。
	Storage = guest.Storage

	// HTTP returns the guest outbound HTTP host-service client. Real network
	// access must be executed in a wasip1 guest through host authorization.
	//
	// HTTP 返回受治理外部 HTTP 请求的 guest 宿主服务客户端。所有真实网络访问
	// 必须在 wasip1 guest 中经宿主授权执行，非 WASI 构建仅返回不可用错误。
	HTTP = guest.HTTP

	// Data returns the compatibility guest structured-data host-service client.
	// New dynamic-plugin code should prefer plugindb.
	//
	// Data 返回兼容型结构化数据 host service 客户端。新动态插件代码应优先使用
	// plugindb，保留该入口是为了兼容已有 pluginbridge 调用方。
	Data = guest.Data

	// Cache returns the guest cache host-service client for authorized namespaces.
	// Real cache reads and writes must remain governed by host authorization and
	// consistency policy, not by a local fallback cache in non-WASI builds.
	//
	// Cache 返回受治理缓存 host service 客户端。真实读写必须受宿主命名空间授权
	// 和一致性策略约束，非 WASI 构建不会落到本地临时缓存。
	Cache = guest.Cache

	// Lock returns the guest distributed-lock host-service client. Lock acquire,
	// renew, and release operations are mediated by the host ticket protocol.
	//
	// Lock 返回受治理分布式锁 host service 客户端。锁获取、续租和释放必须经
	// 宿主 ticket 协议完成，非 WASI 构建仅保留调用形态。
	Lock = guest.Lock

	// Config returns the guest plugin-config host-service client. It reads only
	// configuration values published by the host for the current plugin.
	//
	// Config 返回插件作用域配置 host service 客户端。该入口只读取宿主发布给
	// 当前插件的配置值，不应被用作宿主全局配置或私有状态访问方式。
	Config = guest.Config

	// Notify returns the guest notification host-service client. Notification
	// dispatch must pass through host channel authorization, resource boundaries,
	// and audit handling.
	//
	// Notify 返回受治理通知 host service 客户端。通知发送必须经宿主渠道授权、
	// 资源边界和审计链路处理。
	Notify = guest.Notify

	// Cron returns the guest cron declaration host-service client. It submits
	// plugin-side cron contracts to host discovery and scheduling boundaries.
	//
	// Cron 返回动态插件定时任务声明 host service 客户端，用于把插件侧 cron
	// 契约交给宿主发现和调度，不在该 facade 中执行任务。
	Cron = guest.Cron

	// HostConfig returns the guest client for whitelisted public host configuration.
	// Only host-published public configuration items are reachable through it.
	//
	// HostConfig 返回宿主公开配置白名单 host service 客户端，只能读取宿主明确
	// 发布给动态插件的公共配置项。
	HostConfig = guest.HostConfig

	// Manifest returns the guest client for plugin-scoped manifest resources.
	// It reads declared resources from the plugin package and is not a general
	// host filesystem access entrypoint.
	//
	// Manifest 返回插件作用域 manifest 资源 host service 客户端，用于读取插件
	// 包内声明资源，不提供任意路径文件访问能力。
	Manifest = guest.Manifest

	// HostLog writes a runtime log entry through Runtime for compatibility callers.
	// In non-WASI builds it returns the host-call-unavailable error.
	//
	// HostLog 通过 Runtime 写入运行时日志，保留旧版函数式入口；非 WASI 构建中
	// 返回 host call 不可用错误。
	HostLog = guest.HostLog

	// HostStateGet reads one plugin-scoped runtime state value through Runtime.
	//
	// HostStateGet 通过 Runtime 读取一个插件作用域运行时状态值。
	HostStateGet = guest.HostStateGet

	// HostStateSet writes one plugin-scoped runtime state value through Runtime.
	//
	// HostStateSet 通过 Runtime 写入一个插件作用域运行时状态值。
	HostStateSet = guest.HostStateSet

	// HostStateDelete deletes one plugin-scoped runtime state value through Runtime.
	//
	// HostStateDelete 通过 Runtime 删除一个插件作用域运行时状态值。
	HostStateDelete = guest.HostStateDelete

	// HostStateGetInt reads one integer plugin-scoped runtime state value.
	//
	// HostStateGetInt 通过 Runtime 读取一个整数型插件作用域运行时状态值。
	HostStateGetInt = guest.HostStateGetInt

	// HostStateSetInt writes one integer plugin-scoped runtime state value.
	//
	// HostStateSetInt 通过 Runtime 写入一个整数型插件作用域运行时状态值。
	HostStateSetInt = guest.HostStateSetInt

	// HostDBQuery preserves the legacy host database query helper entrypoint.
	// New code should use governed data host service contracts or plugindb.
	//
	// HostDBQuery 保留旧版宿主数据库查询辅助入口，仅用于兼容历史调用。新代码
	// 应使用受治理的 data host service 或 plugindb 契约。
	HostDBQuery = guest.HostDBQuery

	// HostDBExecute preserves the legacy host database execute helper entrypoint.
	// New code should use governed data host service contracts or plugindb.
	//
	// HostDBExecute 保留旧版宿主数据库执行辅助入口，仅用于兼容历史调用。新代码
	// 应使用受治理的 data host service 或 plugindb 契约。
	HostDBExecute = guest.HostDBExecute
)
