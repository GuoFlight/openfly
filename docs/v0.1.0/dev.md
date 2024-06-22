# 一些测试case

```shell
# 写入有问题的配置
etcdctl put /openfly/l4/30001 '{"listen":30001,"upstream":{"hosts":[{"ip":"","port":53}]}}'
```

# FAQ

添加重复接口会报错吗？
* 会

若配置有问题，会导致线上问题吗？
* 不会。
* nginx -t失败后配置不生效。