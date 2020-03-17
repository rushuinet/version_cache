package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/rushuinet/version_cache"
	"log"
	"time"
)

type Up struct {
	Ids []string `form:"ids" json:"ids" binding:"required"`
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
		CheckTime:           15,
	}

	for i := 1; i < 5; i++ {
		go func(i int) {
			c := version_cache.New(config)
			c.FirstLoad()
			go c.Load()

			r := gin.Default()
			r.GET("/test", func(context *gin.Context) {
				context.JSON(200, c.Get())
			})

			r.GET("/monitor", func(context *gin.Context) {
				context.JSON(200, c.Monitor())
			})

			r.POST("/update", func(context *gin.Context) {
				var form Up
				if err := context.ShouldBindJSON(&form); err != nil {
					context.JSON(200, err.Error())
					return
				}
				c.Update(form.Ids, func(use, new string, ids []string) {
					for _, id := range ids {
						c.Redis.HSet(use, id, "aa"+id+version_cache.Int64ToStr(time.Now().Unix()))
					}
				})
				context.JSON(200, "更新成功")
			})

			r.GET("/clear", func(context *gin.Context) {
				c.Clear()
				context.JSON(200, "删除成功")
			})
			r.Run(":888" + version_cache.Int64ToStr(int64(i)))
		}(i)
	}
	select {}
}
