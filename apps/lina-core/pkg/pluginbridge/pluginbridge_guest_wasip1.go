//go:build wasip1

// pluginbridge_guest_wasip1.go provides the WASI guest facade for dynamic-plugin
// host-service helpers.
//
// This file is compiled only for wasip1 guest builds. It keeps the historical
// root-package pluginbridge API available to dynamic-plugin WASM guests while
// delegating contracts and helpers to the guest subcomponent. Real host calls,
// request encoding, response decoding, authorization resource names, and WASI
// import/export details stay below this facade.
//
// Development constraints:
//   - Keep this file as a thin, stable root-package compatibility layer.
//   - Continue reaching host capabilities through governed host services
//     declared by pluginbridge and hostServices contracts.
//   - Add new guest capabilities in the guest subcomponent first, then expose
//     aliases here only when root-package compatibility is required.
//
// 本文件为动态插件 host service 辅助能力提供 WASI guest facade。
//
// 本文件只在 wasip1 guest 构建目标下编译，用于让动态插件 WASM guest 继续通过
// pluginbridge 根包使用历史 API。文件中的契约和辅助入口都委托给 guest 子组件；
// 真实 host call、请求编码、响应解码、授权资源名以及 WASI import/export 细节
// 都保留在该 facade 下层实现中。
//
// 开发约束：
//   - 保持本文件是稳定且轻量的根包兼容层。
//   - 访问宿主能力时，必须继续通过 pluginbridge 与 hostServices 契约声明的
//     受治理 host service。
//   - 新增 guest 能力时，应先在 guest 子组件中实现；只有需要根包兼容入口时，
//     才在这里暴露别名。

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
	// Runtime returns the WASI guest runtime host-service client. It writes
	// structured logs, accesses plugin-scoped runtime state, reads host time,
	// generates UUIDs, and obtains host node identity.
	//
	// Runtime 返回 wasip1 guest 运行时宿主服务客户端，用于写入结构化日志、访问
	// 插件作用域状态、读取宿主时间、生成 UUID 和获取宿主节点信息。
	Runtime = guest.Runtime

	// Storage returns the WASI guest storage host-service client. Object paths,
	// overwrite policy, and accessible resources must be constrained by host
	// hostServices authorization declarations.
	//
	// Storage 返回 wasip1 guest 受治理对象存储客户端。对象路径、覆盖策略和
	// 可访问资源必须经宿主 hostServices 授权声明约束。
	Storage = guest.Storage

	// HTTP returns the WASI guest outbound HTTP host-service client. Real network
	// access must pass through host authorization, auditing, and resource bounds.
	//
	// HTTP 返回 wasip1 guest 受治理外部 HTTP 请求客户端。真实网络访问必须通过
	// 宿主授权、审计和资源边界处理。
	HTTP = guest.HTTP

	// Data returns the WASI guest compatibility structured-data host-service client.
	// New dynamic-plugin code should prefer plugindb for a clearer data contract.
	//
	// Data 返回 wasip1 guest 兼容型结构化数据 host service 客户端。新动态插件
	// 代码应优先使用 plugindb，以获得更清晰的数据访问契约。
	Data = guest.Data

	// Cache returns the WASI guest cache host-service client. Cache namespaces,
	// keys, and expiration policy must be governed by host authorization and
	// consistency rules.
	//
	// Cache 返回 wasip1 guest 受治理缓存客户端，缓存命名空间、键和过期策略必须
	// 由宿主授权模型和一致性策略管理。
	Cache = guest.Cache

	// Lock returns the WASI guest distributed-lock host-service client. Lock
	// lifecycle operations must use host-issued tickets for acquire, renew, and release.
	//
	// Lock 返回 wasip1 guest 受治理分布式锁客户端，锁生命周期必须通过宿主签发
	// 的 ticket 完成获取、续租和释放。
	Lock = guest.Lock

	// Config returns the WASI guest plugin-config host-service client. It reads
	// only read-only configuration values published by the host for the current plugin.
	//
	// Config 返回 wasip1 guest 插件作用域配置客户端，仅读取宿主发布给当前插件
	// 的只读配置值。
	Config = guest.Config

	// Notify returns the WASI guest notification host-service client. Notification
	// channels, templates, and resource boundaries are authorized and audited by the host.
	//
	// Notify 返回 wasip1 guest 受治理通知客户端，通知渠道、模板和资源边界由宿主
	// 统一授权与审计。
	Notify = guest.Notify

	// Cron returns the WASI guest cron declaration host-service client. It submits
	// dynamic-plugin cron contracts to host discovery and scheduling boundaries.
	//
	// Cron 返回 wasip1 guest 定时任务声明客户端，用于把动态插件 cron 契约提交
	// 给宿主发现和调度边界。
	Cron = guest.Cron

	// HostConfig returns the WASI guest client for whitelisted public host configuration.
	// Only public configuration items explicitly published to dynamic plugins are reachable.
	//
	// HostConfig 返回 wasip1 guest 宿主公开配置白名单客户端，只能读取宿主明确
	// 发布给动态插件的公共配置项。
	HostConfig = guest.HostConfig

	// Manifest returns the WASI guest client for plugin-scoped manifest resources.
	// It reads declared resources from the plugin package and is not a general
	// host filesystem access entrypoint.
	//
	// Manifest 返回 wasip1 guest 插件作用域 manifest 资源客户端，用于读取插件
	// 包内声明资源，不提供任意宿主文件系统访问能力。
	Manifest = guest.Manifest

	// HostLog writes a runtime log entry through Runtime for compatibility callers.
	//
	// HostLog 通过 Runtime 写入运行时日志，保留旧版函数式入口。
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
