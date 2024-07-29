package util

import (
	tourMemCache "github.com/go-programming-tour-book/tourMemCache/cache"
)

func CalcLen(value interface{}) int {
	switch value.(type) {
	case bool, uint8, int8:
		return 1
	case int16, uint16:
		return 2
	case int32, uint32:
		return 4
	case int64, uint64:
		return 8
	case float32, float64:
		return 16
	case string:
		return len(value.(string))
	case tourMemCache.Cache:
		return value.(tourMemCache.Cache).Len()
	case int, uint:
		return 16
	default:
		panic("无法接受的变量类型")
	}
}
