package version_cache

//版本列表key类型为list
func (c *Cache) VersionKey() string {
	return c.Key + ":version"
}

//版本数据key类型为set
func (c *Cache) DataKey(versionKey string) string {
	return c.Key + ":data_" + versionKey
}

//各种标记key类型为set
func (c *Cache) TagKey() string {
	return c.Key + ":tag"
}

//更新数据key类型为
func (c *Cache) UpdateKey() string {
	return c.Key + ":update"
}