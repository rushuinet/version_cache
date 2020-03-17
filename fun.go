package version_cache

import "strconv"

//int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

//str to int64
func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}
