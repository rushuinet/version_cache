package version_cache

import (
	"sync"
)

type dataMap struct {
	id         int64 //版本ID
	effectTime int64 //生效时间
	mu         sync.RWMutex
	data       map[string]interface{}
}

func newData() *dataMap {
	return &dataMap{
		id:         0,
		effectTime: 0,
		data:       make(map[string]interface{}),
	}
}

//设置数据
func (d *dataMap) Set(key string, val interface{}) {
	d.mu.Lock()
	d.data[key] = val
	d.mu.Unlock()
}

//获取数据
func (d *dataMap) Get(key string) interface{} {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.data[key]
}

//设置数据
func (d *dataMap) Del(key string) {
	d.mu.Lock()
	delete(d.data, key)
	d.mu.Unlock()
}

//长度
func (d *dataMap) Len() int {
	return len(d.data)
}

//Key是否存在
func (d *dataMap) Exists(key string) bool {
	_, ok := d.data[key]
	return ok
}
