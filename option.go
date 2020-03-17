package version_cache

import (
	"github.com/go-redis/redis/v7"
)

type Option struct {
	Redis               *redis.Client
	Key                 string
	CheckTime           int64 `json:"check_time"`            //检查版本时间间隔（秒）
	VersionGenerateTime int64 `json:"version_generate_time"` //版本生成间隔(完成一次版本生效后)
	KeepVersionNum      int64 `json:"keep_version_num"`      //保留的版本数
}
