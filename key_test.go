package version_cache

import "testing"

//所有key
func TestCache_AllKey(t *testing.T) {
	op := &Option{Key: "test"}
	c := New(op)
	if c.VersionKey() != "test:version" {
		t.Error("version key error")
	}

	if c.DataKey("v1") != "test:data_v1" {
		t.Error("data key error")
	}

	if c.TagKey() != "test:tag" {
		t.Error("tag key error")
	}

	if c.UpdateKey() != "test:update" {
		t.Error("update key error")
	}
}
