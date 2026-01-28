package utils

import "strconv"

// ParseInt 字符串转整数
func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// ParseUint 字符串转无符号整数
func ParseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
