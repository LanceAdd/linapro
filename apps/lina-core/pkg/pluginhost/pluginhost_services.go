// This file defines the host-published service directory exposed to source
// plugin registration callbacks.

package pluginhost

import (
	"lina-core/pkg/pluginservice/contract"
)

// HostServices exposes host-owned pluginservice adapters to source plugins.
type HostServices interface {
	// APIDoc returns the host API-documentation localization adapter used by
	// source plugins to resolve stable OpenAPI operation keys into localized
	// module titles and operation summaries. The method itself returns no error;
	// implementations may return nil only when the host service directory is not
	// initialized, and callers must treat the returned adapter as read-only host
	// catalog access.
	//
	// APIDoc 返回宿主 API 文档本地化适配器，供源码插件将稳定的 OpenAPI 操作键解析为本地化模块标题和操作摘要。
	// 该方法本身不返回错误；仅当宿主服务目录未初始化时实现才可以返回 nil，调用方必须把返回的适配器视为只读的宿主目录访问能力。
	APIDoc() contract.APIDocService

	// Auth returns the host tenant-auth adapter used by source plugins to
	// delegate tenant selection, tenant switching, and governed impersonation
	// token lifecycle operations back to the host. The returned service owns
	// token signing, session registration, revocation, and permission-cache
	// priming; plugins must still perform their own business authorization before
	// requesting privileged token operations.
	//
	// Auth 返回宿主租户认证适配器，供源码插件把租户选择、租户切换以及受治理的模拟登录令牌生命周期操作委托给宿主。
	// 返回服务负责令牌签发、会话注册、撤销和权限缓存预热；插件在请求高权限令牌操作前仍必须先完成自身业务授权判断。
	Auth() contract.AuthService

	// BizCtx returns the host business-context adapter that exposes a stable,
	// read-only projection of request identity, tenant, platform-bypass, and
	// impersonation metadata. Plugins should use this adapter instead of host
	// internal context types; absent request context fields are represented by
	// zero values in the returned contract model.
	//
	// BizCtx 返回宿主业务上下文适配器，用于暴露稳定、只读的请求身份、租户、平台绕过和模拟登录元数据投影。
	// 插件应使用该适配器而不是宿主内部上下文类型；请求上下文字段不存在时会在契约模型中体现为零值。
	BizCtx() contract.BizCtxService

	// Cache returns the plugin-scoped host cache adapter for transient,
	// tenant-aware plugin runtime data. The adapter binds cache access to the
	// current plugin identity and host cache backend; an unscoped base directory
	// may return nil because cache reads and writes require a plugin-bound service
	// view.
	//
	// Cache 返回插件作用域的宿主缓存适配器，用于存取具备租户感知能力的插件临时运行时数据。
	// 该适配器会把缓存访问绑定到当前插件身份和宿主缓存后端；未绑定插件的基础目录可以返回 nil，
	// 因为缓存读写必须使用插件绑定后的服务视图。
	Cache() contract.CacheService

	// Config returns the plugin-scoped static configuration adapter for reading
	// host-approved plugin configuration values. The adapter must not fall back to
	// unrestricted host-wide configuration when plugin identity is blank; an
	// unscoped base directory may return nil until HostServicesForPlugin binds a
	// plugin ID.
	//
	// Config 返回插件作用域的静态配置适配器，用于读取经过宿主认可的插件配置值。
	// 当插件身份为空时，该适配器不得回退读取不受限制的宿主全局配置；
	// 未绑定插件的基础目录可以在 HostServicesForPlugin 绑定插件 ID 前返回 nil。
	Config() contract.ConfigService

	// HostConfig returns the public host config adapter for whitelisted
	// configuration values that plugins may read without depending on private host
	// configuration models. Missing or unavailable keys are handled by the
	// returned service contract, commonly through default-value reads or errors
	// from typed accessors.
	//
	// HostConfig 返回宿主公开配置适配器，用于读取白名单内、允许插件访问的宿主配置值，避免插件依赖宿主私有配置模型。
	// 键不存在或不可用时由返回服务的契约处理，通常通过类型化读取方法返回默认值或错误。
	HostConfig() contract.HostConfigService

	// I18n returns the host runtime translation adapter for resolving the current
	// request locale, runtime message keys, fallback text, and localized keyword
	// searches. The returned service is read-only and must preserve the host
	// locale resolution rules instead of letting plugins manage translation
	// caches directly.
	//
	// I18n 返回宿主运行时翻译适配器，用于解析当前请求语言、运行时消息键、兜底文本以及本地化关键词搜索。
	// 返回服务是只读能力，并且必须保持宿主语言解析规则，插件不得通过该能力直接管理翻译缓存。
	I18n() contract.I18nService

	// Manifest returns the plugin-scoped manifest resource adapter for read-only
	// access to declaration resources under the plugin manifest root. Paths are
	// resolved relative to the plugin manifest boundary; an unscoped base
	// directory may return nil because manifest reads require a plugin-bound
	// service view.
	//
	// Manifest 返回插件作用域的清单资源适配器，用于只读访问插件 manifest 根目录下的声明资源。
	// 路径会在插件清单边界内按相对路径解析；未绑定插件的基础目录可以返回 nil，因为清单读取必须使用插件绑定后的服务视图。
	Manifest() contract.ManifestService

	// Notify returns the host notification adapter used by source plugins to
	// publish governed messages into the host inbox pipeline or remove messages
	// by declared business source. The adapter owns host delivery records and
	// fan-out behavior; plugins provide source identifiers, message content, and
	// category intent through the contract input models.
	//
	// Notify 返回宿主通知适配器，供源码插件向宿主站内信管道发布受治理的消息，或按声明的业务来源删除消息。
	// 该适配器负责宿主投递记录和分发行为；插件通过契约输入模型提供来源标识、消息内容和分类意图。
	Notify() contract.NotifyService

	// PluginState returns the host plugin enablement adapter for checking whether
	// a plugin is installed, enabled, and visible for the current request scope.
	// The returned service may expose both snapshot-backed fast reads and
	// authoritative persisted-state reads; callers choose the method according to
	// the freshness requirements of their control path.
	//
	// PluginState 返回宿主插件启用状态适配器，用于判断插件在当前请求范围内是否已安装、已启用且可见。
	// 返回服务可以同时提供基于快照的快速读取和基于持久化状态的权威读取；调用方应按控制路径的新鲜度要求选择具体方法。
	PluginState() contract.PluginStateService

	// PluginLifecycle returns the host plugin lifecycle orchestration adapter for
	// tenant-scoped plugin disable and tenant deletion coordination. The returned
	// service runs precondition checks that may return errors before destructive
	// governance actions and also publishes best-effort post-action
	// notifications after those actions complete.
	//
	// PluginLifecycle 返回宿主插件生命周期编排适配器，用于协调租户范围的插件禁用和租户删除流程。
	// 返回服务会在破坏性治理动作前执行可能返回错误的前置检查，并在动作完成后发布尽力而为的后置通知。
	PluginLifecycle() contract.PluginLifecycleService

	// Route returns the host dynamic-route metadata adapter for reading metadata
	// attached to the current dynamic-plugin HTTP request. It exposes the matched
	// plugin ID, method, public path, route text, metadata, and captured bridge
	// response details without exposing the host router internals.
	//
	// Route 返回宿主动态路由元数据适配器，用于读取当前动态插件 HTTP 请求上附加的路由元数据。
	// 它会暴露命中的插件 ID、方法、公开路径、路由文案、元数据以及捕获的桥接响应信息，但不会暴露宿主路由器内部实现。
	Route() contract.RouteService

	// Session returns the host online-session adapter for governed session list
	// and revocation operations. The returned service projects stable session
	// fields for plugins while preserving host authorization, tenant visibility,
	// filtering, pagination, and revocation semantics.
	//
	// Session 返回宿主在线会话适配器，用于受治理的会话列表查询和会话撤销操作。
	// 返回服务向插件投影稳定的会话字段，同时保持宿主授权、租户可见性、过滤、分页和撤销语义。
	Session() contract.SessionService

	// TenantFilter returns the host tenant-filter adapter for applying the
	// conventional tenant predicate to plugin-owned database queries and for
	// reading tenant/audit identity metadata. The adapter centralizes platform
	// bypass and impersonation decisions so plugins do not duplicate host tenant
	// visibility rules.
	//
	// TenantFilter 返回宿主租户过滤适配器，用于向插件自有数据库查询应用约定的租户谓词，并读取租户和审计身份元数据。
	// 该适配器集中处理平台绕过和模拟登录判断，避免插件重复实现宿主租户可见性规则。
	TenantFilter() contract.TenantFilterService
}
