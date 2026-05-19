## ADDED Requirements

### Requirement: 源码插件可以托管消费端前端资产

系统 SHALL 允许源码插件通过`frontend/consumer/`目录提供消费端浏览器前端资产，并由宿主按插件启用状态、版本和声明的前端挂载路径进行托管。

#### Scenario: 访问启用插件的消费端前端挂载入口
- **WHEN** 源码插件在`plugin.yaml`中声明`consumer.frontend.mount_path`为`/portal`
- **AND** 插件已安装并启用
- **AND** 插件提供`frontend/consumer/index.html`
- **THEN** 浏览器访问`/portal`时，宿主返回该插件的消费端前端入口
- **AND** 响应包含正确的`Content-Type`、`Cache-Control`和`ETag`

#### Scenario: 禁用插件的消费端前端挂载不可访问
- **WHEN** 源码插件声明了消费端前端挂载路径
- **AND** 插件当前未启用
- **THEN** 浏览器访问该挂载路径时，宿主返回`404`

#### Scenario: 消费端前端资产命名空间按版本读取资源
- **WHEN** 源码插件已启用并提供`frontend/consumer/assets/app.js`
- **THEN** 浏览器可通过`/consumer-plugin-assets/<plugin-id>/<version>/assets/app.js`读取该资源
- **AND** 当请求版本与当前有效插件版本不一致时，宿主拒绝返回该资源

### Requirement: 消费端前端挂载必须支持受控的 SPA 回退

系统 SHALL 只在插件显式声明`consumer.frontend.spa_fallback=true`且请求路径不像具体静态资源时，将挂载路径下不存在的子路由回退到入口文件。

#### Scenario: 开启 SPA 回退后访问干净子路由
- **WHEN** 插件声明`spa_fallback=true`
- **AND** 浏览器访问`/portal/orders/123`
- **AND** `orders/123`不是声明过的静态资源
- **THEN** 宿主返回该插件的消费端前端入口文件

#### Scenario: 缺失静态资源不回退到入口文件
- **WHEN** 插件声明`spa_fallback=true`
- **AND** 浏览器访问`/portal/assets/missing.js`
- **THEN** 宿主返回`404`

### Requirement: 消费端前端响应必须治理浏览器缓存

系统 SHALL 为插件消费端前端资产生成内容相关的强`ETag`，并对稳定挂载路径下的入口和静态资源使用重新验证策略，避免同版本资产刷新后浏览器继续使用过期内容。

#### Scenario: ETag 命中时返回 304
- **WHEN** 浏览器请求消费端前端资产并携带匹配当前内容的`If-None-Match`
- **THEN** 宿主返回`304 Not Modified`
- **AND** 不重复写出响应内容

#### Scenario: 挂载入口 HTML 注入 base
- **WHEN** 宿主通过稳定挂载路径返回消费端`HTML`入口
- **THEN** 响应内容包含指向该挂载路径的`base`声明
- **AND** 插件页面中的相对资源可以在该挂载路径下解析
