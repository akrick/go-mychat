package handlers

import "fmt"

// parseInt 字符串转整数
func parseInt(s string) int {
	result := 0
	fmt.Sscanf(s, "%d", &result)
	return result
}

// parseUint 字符串转uint
func parseUint(s string) uint {
	result := uint(0)
	fmt.Sscanf(s, "%d", &result)
	return result
}
