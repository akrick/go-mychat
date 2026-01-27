package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 数据库中的哈希值
	dbHash := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"

	// 测试密码
	testPassword := "admin123"

	fmt.Println("数据库中的哈希值:")
	fmt.Println(dbHash)
	fmt.Println("\n测试密码:", testPassword)

	// 验证密码
	err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(testPassword))
	if err == nil {
		fmt.Println("✅ 密码验证成功!")
	} else {
		fmt.Println("❌ 密码验证失败!")
		fmt.Println("错误信息:", err)
	}

	// 生成新的 admin123 哈希值
	fmt.Println("\n生成新的 admin123 哈希值:")
	newHash, err := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("生成哈希失败:", err)
		return
	}
	fmt.Println(string(newHash))

	// 验证新哈希
	err = bcrypt.CompareHashAndPassword(newHash, []byte(testPassword))
	if err == nil {
		fmt.Println("✅ 新哈希验证成功!")
	} else {
		fmt.Println("❌ 新哈希验证失败!")
	}
}
