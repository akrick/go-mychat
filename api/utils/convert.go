package utils

import (
	"os"
	"strconv"
)

// ParseStringToInt 将字符串转换为整数
func ParseStringToInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

// ParseInt 解析字符串为整数（别名）
func ParseInt(s string) int {
	return ParseStringToInt(s)
}

// ParseStringToUint 将字符串转换为无符号整数
func ParseStringToUint(s string) uint {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}

// EnsureDir 确保目录存在，如果不存在则创建
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// 为了向后兼容，提供旧名称
func parseInt(s string) int {
	return ParseStringToInt(s)
}

func parseUint(s string) uint {
	return ParseStringToUint(s)
}
