package version_cache

import (
	"reflect"
	"testing"
)

func TestInt64ToStr(t *testing.T) {
	var i int64
	str := Int64ToStr(i)
	if reflect.TypeOf(str).Name() != "string" {
		t.Error("data type error")
	}
}

func TestStrToInt64(t *testing.T) {
	var str string
	i := StrToInt64(str)
	if reflect.TypeOf(i).Name() != "int64" {
		t.Error("data type error")
	}
}
