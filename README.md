# version_cache
version_cache是一个分布式一致性缓存解决方案。

原理：job 将数据打包成版本到redis，实例将存在redis的版本load到本地内存并计算最新版本的生效时间，使所有实例的缓存在同一时间生效来达到所有实例数据的一致。

实用场景：数据量少、非及时生效数据、高并发强一致的场景。如：配置服务，门店服务等

优点：
1. 轻松实现水平扩展，实现千万并发的服务不是梦
2. 数据强一至性，不论启动多少实例，同一时间的数据绝对是一致的（服务器时间一致情况下）
3. 使用简单，实现数据生成接口后就可以像使用缓存一样方便，轻松实现高性能服务

缺点：
1. 数据按版本生效，变更的数据会延迟生效（原则上数据量越小处理越快）
2. 不适合大数据缓存

# 架构：

https://www.processon.com/view/5e4b7f2ee4b00aefb7e5d054

# 安装
```
go get github.com/rushuinet/version_cache
```

# 用法:

job 数据生成服务启动一个实例
```
package main

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"github.com/rushuinet/version_cache"
	"log"
	"strconv"
	"time"
)

type Info struct {
	Id         int
	Key        string
	CreateTime int64
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil {
		log.Fatal("redis conn error")
	}
	config := &version_cache.Option{
		Redis:               client,
		Key:                 "test",
		VersionGenerateTime: 60,
		CheckTime:           5,
		KeepVersionNum:      6,
	}
	cache := version_cache.New(config)
	cache.Generate(func(key string) {
		for i := 0; i < 5; i++ {
			info := &Info{
				Id:         i,
				Key:        key,
				CreateTime: time.Now().Unix(),
			}
			data, _ := json.Marshal(info)
			cache.Redis.HSet(key, strconv.Itoa(i), data)
			time.Sleep(1 * time.Second)
		}
	})

	select {}
}


```


服务中使用方式：

```
# 服务启动前加
c := version_cache.New(config)
c.FirstLoad()
go c.Load()

# 业务中使用
c.Get(key)
```