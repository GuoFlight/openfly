# v0.5.0

新特性（Features）

* 支持配置nginx日志(access_log)

优化

* 健康检查接口：检查nginx配置是否正确

# v0.4.0

新特性（Features）

* 支持优雅退出

# v0.3.0

新特性（Features）

* 支持限速(上传与下载)
* 支持超时控制
* 支持启用/禁用某个配置

# v0.2.0

新特性（Features）

* 新增API，支持获取指定监听端口的配置。(GET /v1/admin/nginx/get?listen=xxxx)

# v0.1.0

新特性（Features）

* 支持添加注释

优化（Refactored）

* nginx配置有问题，不再阻塞openfly启动
* 参数校验：端口校验
* 参数校验：白名单(allow和deny语句)
* 参数校验：ip和网段
* 优化【API监听端口失败，openfly依然会启动成功】的问题，改为Fatal退出。

修复（Fixed）

* 修复 nginx配置有问题时，重启openfly会中断配置生成，导致配置缺失的问题。改为忽略有问题的配置。

# v0.0.0

* 初始版本

