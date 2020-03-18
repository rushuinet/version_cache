package version_cache

import (
	"log"
	"time"
)

//加载使用数据
func (c *Cache) firstLoad(verKey string) {
	var (
		cursor uint64 = 0
		values []string
	)
	for {
		values, cursor = c.Redis.HScan(c.DataKey(verKey), cursor, "", loadCount).Val()
		for i := 0; i < len(values); i += 2 {
			if c.LoadFun != nil {
				c.use.Set(values[i], c.LoadFun(values[i+1]))

			} else {
				c.use.Set(values[i], values[i+1])
			}
		}
		if cursor < 1 {
			break
		}
	}
	c.use.id = StrToInt64(verKey)
	c.use.effectTime = time.Now().Unix()
	log.Println("首次数据加载成功")
}

func (c *Cache) Load() {
	t := time.NewTimer(time.Second * 0)
	defer t.Stop()
	for {
		<-t.C
		//各种标记
		all := c.Redis.HGetAll(c.TagKey()).Val()
		newVer := StrToInt64(c.Redis.LIndex(c.VersionKey(), 0).Val()) //新版本
		newVerTime := StrToInt64(all[effectTime])                     //新版本生效时间
		useVer := StrToInt64(all[useVersion])                         //生效版本
		updateVer := StrToInt64(all[updateVersion])                   //即时更新版本

		//更新紧急数据
		if updateVer > 0 && updateVer != c.updateVersion {
			//加载中标记
			c.isLoading = true
			//先锁定
			c.lock()
			c.loadUpdate()
			c.updateVersion = updateVer
		}

		//如果设置的生效时间大于新版本时间，说明设置了新生效时间load
		if c.new.id > 0 && newVerTime > c.new.id {
			c.new.effectTime = newVerTime
		}

		//存在load新版本生效时间 && 当前时间大于新版本生效时间则切换（心跳切换）
		if c.new.effectTime > 0 && c.new.effectTime < time.Now().Unix() {
			c.switchData()
		}
		//新版本与使用版本不同 && 新版本也不是load的新版本
		if newVer != useVer && newVer != c.new.id {
			//加载中标记
			c.isLoading = true
			//先锁定
			c.lock()
			c.load(Int64ToStr(newVer))
			c.new.id = newVer
		}

		c.isLoading = false

		// 需要重置Reset 使 t 重新开始计时
		t.Reset(time.Second * time.Duration(c.CheckTime))
	}
}

//加载新数据
func (c *Cache) load(verKey string) {
	var (
		cursor uint64 = 0
		values []string
	)
	for {
		values, cursor = c.Redis.HScan(c.DataKey(verKey), cursor, "", loadCount).Val()
		for i := 0; i < len(values); i += 2 {
			if c.LoadFun != nil {
				c.new.Set(values[i], c.LoadFun(values[i+1]))

			} else {
				c.new.Set(values[i], values[i+1])
			}
		}
		if cursor < 1 {
			break
		}
	}
	log.Println(c.Key + " load new data version:" + verKey)
}

//load紧急更新的数据
func (c *Cache) loadUpdate() {
	//取要更新的key
	keys := c.Redis.SMembers(c.UpdateKey()).Val()
	dataKey := c.DataKey(Int64ToStr(c.use.id))
	log.Println("search key :" + dataKey)
	data := c.Redis.HMGet(dataKey, keys...).Val()
	//取最新数据
	for i, key := range keys {
		if data[i] == nil {
			c.use.Del(key)
			if c.new.id > 0 {
				c.new.Del(key)
			}
			log.Println(key + " delete")
		} else {
			c.use.Set(key, data[i])
			if c.new.id > 0 {
				c.new.Set(key, data[i])
			}
			str, _ := data[i].(string)
			log.Println(key + " update :" + str)
		}
	}
}
