package version_cache

import (
	"log"
	"sync"
	"time"
)

const (
	useVersion    string = "use_version"    //使用版本
	updateVersion string = "update_version" //更新版本（即时更新时用）
	effectTime    string = "effect_time"    //新版本生效时间
	lockTime      string = "lock_time"      //加载数据锁
	loadCount     int64  = 1000             //每次从redis里HScan的数据条数
)

type IData interface {
	Set(key string, val interface{})
	Get(key string) interface{}
	Len() int
	Exists(key string) bool
}

type Cache struct {
	*Option
	mu            sync.Mutex
	isLoading     bool     //是否加载中
	updateVersion int64    //更新版本号
	use           *dataMap //使用的
	new           *dataMap //最新的
}

func New(op *Option) *Cache {
	cache := &Cache{Option: op}
	cache.isLoading = false
	cache.use = newData()
	cache.new = newData()
	go cache.check() //心跳锁
	return cache
}

//第一次启动时加载
func (c *Cache) FirstLoad() {
	useVer := c.Redis.HGet(c.TagKey(), useVersion).Val() //新版本
	//加载中标记
	c.isLoading = true
	//先锁定
	c.lock()
	//首次启动需加载使用版本数据
	c.firstLoad(useVer)
	c.isLoading = false
}

//切换
func (c *Cache) switchData() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.new.effectTime > 0 && c.new.effectTime < time.Now().Unix() {
		//加写锁
		c.use = c.new
		c.new = newData()
		log.Println(c.Key + " 切换为新版本" + Int64ToStr(c.use.id))
	}
	return
}

//设置数据
func (c *Cache) Set(key string, val interface{}) {
	c.use.Set(key, val)
}

//获取数据
func (c *Cache) Get(key string) interface{} {
	if c.new.effectTime > 0 && c.new.effectTime < time.Now().Unix() {
		c.switchData()
	}
	return c.use.Get(key)
}

//长度
func (c *Cache) Len() int {
	return c.use.Len()
}

//Key是否存在
func (c *Cache) Exists(key string) bool {
	return c.use.Exists(key)
}

//是否需要锁定
func (c *Cache) check() {
	t := time.NewTimer(time.Second * time.Duration(c.CheckTime))
	defer t.Stop()
	for {
		<-t.C
		//load数据时加锁
		if c.isLoading {
			c.lock()
		}
		t.Reset(time.Second * time.Duration(c.CheckTime))
	}
}

//加锁
func (c *Cache) lock() {
	t := time.Now().Unix() + c.CheckTime*2
	c.Redis.HSet(c.TagKey(), lockTime, t)
	log.Println("locking " + Int64ToStr(t))
}
