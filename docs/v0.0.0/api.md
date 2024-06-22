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