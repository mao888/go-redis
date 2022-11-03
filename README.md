# gredis

#### 通过名称获取Redis连接（Client、哨兵、集群）
* 创建Redis连接
```go
    gredis.MustInit(constants.RedisGiaSentinelName,
        gredis.WithSentinel(),
        gredis.WithHosts([]string{"localhost:26379","localhost:26380","localhost:26381"}),
        gredis.WithPoolSize(10),
        gredis.WithPassWord(""),
        gredis.WithMinIdleCons(1),
        gredis.WithMasterName("master"),
    )
```

* 获取Redis连接
```go
    gredis.Redis(constants.RedisGiaSentinelName)
```