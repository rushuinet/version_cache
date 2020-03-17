package version_cache

import (
	"fmt"
	"log"
	"time"
)

func (c *Cache) Generate(f func(str string)) {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()

	for {
		<-t.C
		//是否生成数据
		if c.isGenerate() {
			c.lock()
			c.delUpdateKey()
			c.generate(f)
		}
		t.Reset(time.Second * time.Duration(c.CheckTime))
	}
}

//是否生成新版本数据
func (c *Cache) isGenerate() bool {
	b := false

	all := c.Redis.HGetAll(c.TagKey()).Val()
	newVer := c.Redis.LIndex(c.VersionKey(), 0).Val() //新版本
	useVer := all[useVersion]                         //使用版本
	newVerTime := StrToInt64(all[effectTime])         //新版本生效时间
	lockTime := StrToInt64(all[lockTime])             //loadData锁定
	t := time.Now().Unix()                            //当前时间

	//如果使用版本为空测使用最新版本
	if useVer == "" {
		c.Redis.HSet(c.TagKey(), useVersion, newVer)
		useVer = newVer
		log.Println("切换当前版本:" + useVer + "为新版本" + newVer)
	}

	//新老版本不同情况下如果当前时间 > load锁定时间 说明所有应用都加载完成，变更当前版本为最新版本
	if newVer != useVer && t > lockTime {
		newVerTime = t + c.CheckTime*3
		m := make(map[string]interface{})
		m[useVersion] = newVer
		m[effectTime] = newVerTime
		if err := c.Redis.HMSet(c.TagKey(), m).Err(); err != nil {
			fmt.Println(err)
		}
		log.Println("切换当前版本:" + useVer + "为新版本" + newVer)
	}

	//如果使用版本号=新版本号  && 当前时间 > 新版本生效锁时间+版本生成间隔 则生成
	if newVer == useVer && t > newVerTime+c.VersionGenerateTime {
		return true
	}
	return b
}

//生成数据
func (c *Cache) generate(f func(str string)) {
	a := time.Now()
	//产生新版本号
	versionNum := Int64ToStr(a.Unix())

	//产生新版本数据
	f(c.DataKey(versionNum))

	//加入版本列表
	c.Redis.LPush(c.VersionKey(), versionNum)

	//删除多余版本
	keepVersionNum := c.Redis.LLen(c.VersionKey()).Val()
	if keepVersionNum > c.KeepVersionNum {
		for i := 0; i < int(keepVersionNum-c.KeepVersionNum); i++ {
			c.Redis.Del(c.DataKey(c.Redis.RPop(c.VersionKey()).Val()))

		}
	}
	log.Println("generate version:" + versionNum + ",use time:" + Int64ToStr(time.Now().Unix()-a.Unix()))
}

//删除更新字段
func (c *Cache)delUpdateKey()  {
	c.Redis.Del(c.UpdateKey())
}
