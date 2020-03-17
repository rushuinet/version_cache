package version_cache

import (
	"github.com/go-redis/redis/v7"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	//依赖的redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	if pong != "PONG" || err != nil {
		t.Fatal("redis conn error")
	}
	//配置
	config := &Option{
		Redis:               client,
		Key:                 "test",
		VersionGenerateTime: 60,
		CheckTime:           5,
		KeepVersionNum:      6,
	}
	//生成
	cache := New(config)
	go cache.Generate(func(key string) {
		for i := 0; i < 10; i++ {
			cache.Redis.HSet(key, strconv.Itoa(i), "aa"+strconv.Itoa(i)+Int64ToStr(time.Now().Unix()))
		}
	})
	time.Sleep(time.Second) //延迟1秒保证生成数据
	//load
	cache.FirstLoad()
	go cache.Load()

	if cache.Len() < 1 {
		t.Error("cache data is nil")
	}

}


