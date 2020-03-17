package version_cache

import (
	"testing"
)

//test set
func TestDataMap_Get(t *testing.T) {
	data := newData()
	data.Set("testKey", "testValue")
	if data.Get("testKey") != "testValue" {
		t.Error("get error")
	}
}

//test del
func TestDataMap_Del(t *testing.T) {
	data := newData()
	data.Set("testKey", "testValue")
	data.Del("testKey")
	if data.Get("testKey") != nil {
		t.Error("del error")
	}
}

//test exists
func TestDataMap_Exists(t *testing.T) {
	data := newData()
	data.Set("testKey", "testValue")
	if !data.Exists("testKey") {
		t.Error("exists error")
	}
}

//test len
func TestDataMap_Len(t *testing.T) {
	data := newData()
	data.Set("testKey", "testValue")
	if data.Len() < 1 {
		t.Error("exists error")
	}
}