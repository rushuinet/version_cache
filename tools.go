package version_cache

import (
	"errors"
	"time"
)

//更新
func (c *Cache) Update(ids []string, f func(use, new string, ids []string)) error {
	//验证是否支持更新
	if f == nil {
		return errors.New("not defined update function")
	}
	//更新当前版本与最新版本数据
	all := c.Redis.HGetAll(c.TagKey()).Val()
	useVer := all[useVersion]                                //使用版本
	newVersionNum := c.Redis.LIndex(c.VersionKey(), 0).Val() //当前版本
	f(c.DataKey(useVer), newVersionNum, ids)
	_, err := c.Redis.SAdd(c.UpdateKey(), ids).Result()
	if err != nil {
		return err
	}
	c.Redis.HSet(c.TagKey(), updateVersion, time.Now().Unix()) //更新的最新版本
	return nil
}

//清空
func (c *Cache) Clear() {
	//清空数据
	list := c.Redis.LRange(c.VersionKey(), 0, -1).Val()
	for _, val := range list {
		c.Redis.Del(c.DataKey(val))
	}
	//删除版本
	c.Redis.Del(c.VersionKey())
	//删除tag
	c.Redis.Del(c.TagKey())
}

//内存数据
func (c *Cache) Data() map[string]interface{} {
	if c.new.effectTime > 0 && c.new.effectTime < time.Now().Unix() {
		c.switchData()
	}
	m := make(map[string]interface{})
	newVersionNum := c.Redis.LIndex(c.VersionKey(), 0).Val()
	m[c.VersionKey()] = c.Redis.LRange(c.VersionKey(), 0, -1).Val()
	m[c.TagKey()] = c.Redis.HGetAll(c.TagKey()).Val()
	m["newVersionNum"] = newVersionNum
	m["useVer"] = c.use.id
	m["data"] = c.use.data
	m["new_load_data"] = c.new.data
	m["time"] = time.Now().Unix()
	m["effect_time"] = time.Now().Unix() - c.new.effectTime
	return m
}

//监控
func (c *Cache) Monitor() map[string]interface{} {
	m := make(map[string]interface{})
	m["len"] = c.Len()
	m["useVersion"] = c.new
	m["newVersion"] = c.new.id
	m["effect_time"] = time.Now().Unix() - c.new.effectTime
	m[c.TagKey()] = c.Redis.HGetAll(c.TagKey()).Val()
	return m
}
