# 一些测试case

```shell
# 写入有问题的配置
etcdctl put /openfly/l4/30001 '{"listen":30001,"upstream":{"hosts":[{"ip":"","port":53}]}}'
```