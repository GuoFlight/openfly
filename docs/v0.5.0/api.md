# 健康检查

```shell
curl http://127.0.0.1:1216/v1/health
```

# 鉴权

```shell
# 获取Token
token=$(curl -s -XPOST http://127.0.0.1:1216/v1/login -d "{\"username\":\"admin\",\"password\":\"admin\"}" -H "Content-Type: application/json" | jq -r .data)
echo ${token}
# 访问需要鉴权的接口
curl -H "Authorization: ${token}" http://127.0.0.1:1216/xxx
```

# Nginx相关API

```shell
# 删除配置
curl -XDELETE -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/delete?listen=30001
# 禁用指定端口的配置（启用为on，禁用为off）
curl -XPOST -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/switch -d "listen=30001&switch=off"
# 获取指定端口的配置
curl -s -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/get?listen=30001
# 获取所有配置
curl -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/getAll
# 新增L4配置
curl -i -H "Content-Type: application/json" -XPOST -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/add -d '
{
    "listen":30001,
    "upstream":{
        "hosts":[
            {
                "ip":"1.1.1.1",
                "port":53
            }
        ]
    }
}'
# 更新L4配置
curl -i -H "Content-Type: application/json" -XPOST -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/set -d '
{
    "listen":30001,
    "upstream":{
        "hosts":[
            {
                "ip":"1.1.1.1",
                "port":53
            }
        ]
    }
}'
```

# 所有支持的参数

```shell
curl -i -H "Content-Type: application/json" -XPOST -H "Authorization: ${token}" http://127.0.0.1:1216/v1/admin/nginx/set -d '
{
    "disable": false,
    "listen": 30001,
    "comments": ["我是注释","I am the comment."],
    "includeFiles": ["/etc/nginx/my.conf"],
    "proxyUploadRate": "10M",
    "proxyDownloadRate": "20M",
    "proxyConnectTimeout": "10s",
    "proxyTimeout": "10m",
    "log": {              // 默认值：继承全局配置
        "mod": "local",   // off：不打印日志；local：单独打印日志；global/其他：继承全局配置。
        "path": "/var/log/nginx/test.stream.log",
        "formatName": "my_log_format",
        "buffer": "5k",
        "flush": "5s"
    },
    "upstream":{
        "hosts":[
            {
                "ip": "1.1.1.1",
                "port": 53,
                "isBackup": false,
                "weight": 100,
                "maxFails": 10,
                "failTimeoutSecond":2
            }
        ],
        "isHash":    true,
        "hashField": "test_field",
        "interval":  10,
        "rise":      2,
        "fall":      2,
        "timeout":   1000
    },
    "whiteList": [
        {
            "type": "allow",
            "target": "1.1.1.1"
        },
         {
            "type": "deny",
            "target": "all"
        }
    ]
}'
```