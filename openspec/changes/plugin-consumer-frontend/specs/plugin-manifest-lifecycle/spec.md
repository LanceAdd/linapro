## ADDED Requirements

### Requirement: 插件清单可以声明消费端前端挂载

系统 SHALL 支持源码插件在`plugin.yaml`中通过`consumer.frontend`声明消费端前端挂载能力。

#### Scenario: 校验合法的消费端前端挂载声明
- **WHEN** 源码插件声明`consumer.frontend.mount_path`为`/portal`
- **AND** 声明`index`为`index.html`
- **THEN** 宿主接受该清单
- **AND** 将入口文件解析为`frontend/consumer/index.html`

#### Scenario: 拒绝根路径挂载
- **WHEN** 源码插件声明`consumer.frontend.mount_path`为`/`
- **THEN** 宿主拒绝该清单

#### Scenario: 拒绝覆盖宿主保留前缀
- **WHEN** 源码插件声明`consumer.frontend.mount_path`为`/api/plugin`
- **THEN** 宿主拒绝该清单

### Requirement: 消费端前端挂载索引必须随插件生命周期失效

系统 SHALL 将消费端前端挂载索引视为由插件清单、资产列表、插件版本和启用状态派生的缓存，并在插件生命周期或运行时修订号变化后失效。

#### Scenario: 源码插件状态变化后本地索引失效
- **WHEN** 源码插件安装、卸载、启用、禁用或升级成功
- **THEN** 宿主清空当前进程内的消费端前端挂载索引
- **AND** 下一次访问消费端前端挂载时重新从权威来源构建索引

#### Scenario: 集群运行时修订号变化后索引失效
- **WHEN** `cluster.enabled=true`
- **AND** 当前实例在读路径观察到插件运行时共享修订号变化
- **THEN** 宿主刷新插件运行时派生缓存
- **AND** 清空当前实例内的消费端前端挂载索引

### Requirement: 消费端治理快照只投影前端能力

系统 SHALL 提供只读的消费端前端治理快照，用于展示源码插件消费端前端挂载和插件治理元数据。

#### Scenario: 快照包含消费端前端挂载
- **WHEN** 源码插件声明了有效消费端前端挂载
- **THEN** 治理快照包含插件 ID、版本、启用状态、租户治理声明、挂载路径、入口文件、`SPA`回退状态和资产数量

#### Scenario: 未声明消费端前端能力的插件不出现在快照中
- **WHEN** 源码插件没有声明有效消费端前端挂载
- **THEN** 治理快照不为该插件生成消费端前端能力项
