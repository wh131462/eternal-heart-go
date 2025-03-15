package utils

import (
	"strconv"
	"strings"
)

// Map 通用 Map 函数 可以遍历元素数组 返回目标数组
func Map[T any, R any](collection []T, mapper func(T, int) R) []R {
	result := make([]R, len(collection))
	for i, item := range collection {
		result[i] = mapper(item, i)
	}
	return result
}

// NumberEqual 是否和数字相等
func NumberEqual(str string, num int64) bool {
	strNum, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return false
	} else {
		return strNum == num
	}
}

// CreateListText 创建列表文本
func CreateListText(list []string) string {
	list = Map(list, func(t string, i int) string {
		return strconv.Itoa(i+1) + "." + t
	})
	return strings.Join(list, "\n")
}
