package database

import (
	"fmt"
	"akrick.com/mychat/admin/backend/models"
	"akrick.com/mychat/admin/backend/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "root:123456@tcp(localhost:3306)/mychat?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移表结构
	err = DB.AutoMigrate(
		&models.User{},
		&models.Administrator{},
		&models.Counselor{},
		&models.CounselorApplication{},
		&models.Order{},
		&models.Payment{},
		&models.PaymentConfig{},
		&models.Review{},
		&models.CounselorStatistics{},
		&models.Notification{},
		&models.ChatSession{},
		&models.ChatMessage{},
		&models.File{},
		&models.ChatBilling{},
		&models.CounselorAccount{},
		&models.WithdrawRecord{},
		&models.Role{},
		&models.Permission{},
		&models.Menu{},
		&models.LowcodeForm{},
		&models.LowcodePage{},
		&models.LowcodeFormData{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.SystemLog{},
		&models.SystemConfig{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 创建默认管理员账号
	createDefaultAdmin()

	return nil
}

// createDefaultAdmin 创建默认管理员账号
func createDefaultAdmin() {
	var count int64
	DB.Model(&models.Administrator{}).Where("username = ?", "admin").Count(&count)
	if count == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := models.Administrator{
			Username:  "admin",
			Password:  hashedPassword,
			RealName:  "超级管理员",
			Email:     "admin@mychat.com",
			Role:      "super_admin",
			Status:    1,
		}
		DB.Create(&admin)
		fmt.Println("默认管理员账号创建成功: admin / admin123")
	}

	// 创建默认角色和权限
	createDefaultRoleAndPermissions()
}

// createDefaultRoleAndPermissions 创建默认角色和权限
func createDefaultRoleAndPermissions() {
	// 检查是否已有默认权限
	var permissionCount int64
	DB.Model(&models.Permission{}).Count(&permissionCount)
	if permissionCount > 0 {
		return // 已有数据，不需要初始化
	}

	fmt.Println("开始初始化权限和角色数据...")

	// 创建系统管理模块权限
	permissions := []models.Permission{
		// 系统管理
		{Name: "系统管理", Code: "system", Type: "menu", Path: "/system", Icon: "Setting", Sort: 1, Status: 1},
		{Name: "用户管理", Code: "system:user", Type: "menu", Path: "/user", Icon: "User", ParentID: 1, Sort: 1, Status: 1},
		{Name: "查看用户", Code: "system:user:list", Type: "button", ParentID: 2, Sort: 1, Status: 1},
		{Name: "创建用户", Code: "system:user:create", Type: "button", ParentID: 2, Sort: 2, Status: 1},
		{Name: "编辑用户", Code: "system:user:edit", Type: "button", ParentID: 2, Sort: 3, Status: 1},
		{Name: "删除用户", Code: "system:user:delete", Type: "button", ParentID: 2, Sort: 4, Status: 1},
		{Name: "重置密码", Code: "system:user:reset", Type: "button", ParentID: 2, Sort: 5, Status: 1},

		{Name: "角色管理", Code: "system:role", Type: "menu", Path: "/roles", Icon: "Avatar", ParentID: 1, Sort: 2, Status: 1},
		{Name: "查看角色", Code: "system:role:list", Type: "button", ParentID: 8, Sort: 1, Status: 1},
		{Name: "创建角色", Code: "system:role:create", Type: "button", ParentID: 8, Sort: 2, Status: 1},
		{Name: "编辑角色", Code: "system:role:edit", Type: "button", ParentID: 8, Sort: 3, Status: 1},
		{Name: "删除角色", Code: "system:role:delete", Type: "button", ParentID: 8, Sort: 4, Status: 1},

		{Name: "权限管理", Code: "system:permission", Type: "menu", Path: "/permissions", Icon: "Key", ParentID: 1, Sort: 3, Status: 1},
		{Name: "查看权限", Code: "system:permission:list", Type: "button", ParentID: 13, Sort: 1, Status: 1},
		{Name: "创建权限", Code: "system:permission:create", Type: "button", ParentID: 13, Sort: 2, Status: 1},
		{Name: "编辑权限", Code: "system:permission:edit", Type: "button", ParentID: 13, Sort: 3, Status: 1},
		{Name: "删除权限", Code: "system:permission:delete", Type: "button", ParentID: 13, Sort: 4, Status: 1},

		{Name: "菜单管理", Code: "system:menu", Type: "menu", Path: "/menus", Icon: "Menu", ParentID: 1, Sort: 4, Status: 1},
		{Name: "查看菜单", Code: "system:menu:list", Type: "button", ParentID: 18, Sort: 1, Status: 1},
		{Name: "创建菜单", Code: "system:menu:create", Type: "button", ParentID: 18, Sort: 2, Status: 1},
		{Name: "编辑菜单", Code: "system:menu:edit", Type: "button", ParentID: 18, Sort: 3, Status: 1},
		{Name: "删除菜单", Code: "system:menu:delete", Type: "button", ParentID: 18, Sort: 4, Status: 1},

		{Name: "在线用户", Code: "system:online", Type: "menu", Path: "/online", Icon: "User", ParentID: 1, Sort: 5, Status: 1},
		{Name: "查看在线用户", Code: "system:online:list", Type: "button", ParentID: 23, Sort: 1, Status: 1},

		// 业务管理
		{Name: "业务管理", Code: "business", Type: "menu", Path: "/business", Icon: "Document", Sort: 2, Status: 1},
		{Name: "订单管理", Code: "business:order", Type: "menu", Path: "/order", Icon: "Document", ParentID: 25, Sort: 1, Status: 1},
		{Name: "查看订单", Code: "business:order:list", Type: "button", ParentID: 26, Sort: 1, Status: 1},
		{Name: "编辑订单", Code: "business:order:edit", Type: "button", ParentID: 26, Sort: 2, Status: 1},
		{Name: "聊天记录", Code: "business:chat", Type: "menu", Path: "/chat", Icon: "ChatDotRound", ParentID: 25, Sort: 2, Status: 1},
		{Name: "查看聊天", Code: "business:chat:list", Type: "button", ParentID: 28, Sort: 1, Status: 1},

		// 财务管理
		{Name: "财务管理", Code: "finance", Type: "menu", Path: "/finance", Icon: "Wallet", Sort: 3, Status: 1},
		{Name: "提现审核", Code: "finance:withdraw", Type: "menu", Path: "/withdraw", Icon: "Wallet", ParentID: 30, Sort: 1, Status: 1},
		{Name: "查看提现", Code: "finance:withdraw:list", Type: "button", ParentID: 31, Sort: 1, Status: 1},
		{Name: "审核提现", Code: "finance:withdraw:approve", Type: "button", ParentID: 31, Sort: 2, Status: 1},
		{Name: "财务统计", Code: "finance:statistics", Type: "menu", Path: "/statistics", Icon: "DataLine", ParentID: 30, Sort: 2, Status: 1},
		{Name: "查看统计", Code: "finance:statistics:view", Type: "button", ParentID: 33, Sort: 1, Status: 1},

		// 低代码平台
		{Name: "低代码平台", Code: "lowcode", Type: "menu", Path: "/lowcode", Icon: "Grid", Sort: 4, Status: 1},
		{Name: "表单设计", Code: "lowcode:form", Type: "menu", Path: "/lowcode/forms", Icon: "Edit", ParentID: 35, Sort: 1, Status: 1},
		{Name: "设计表单", Code: "lowcode:form:design", Type: "button", ParentID: 36, Sort: 1, Status: 1},
		{Name: "页面设计", Code: "lowcode:page", Type: "menu", Path: "/lowcode/pages", Icon: "Grid", ParentID: 35, Sort: 2, Status: 1},
		{Name: "设计页面", Code: "lowcode:page:design", Type: "button", ParentID: 37, Sort: 1, Status: 1},
		{Name: "数据管理", Code: "lowcode:data", Type: "menu", Path: "/lowcode/data", Icon: "Database", ParentID: 35, Sort: 3, Status: 1},
		{Name: "管理数据", Code: "lowcode:data:manage", Type: "button", ParentID: 38, Sort: 1, Status: 1},
	}

	// 批量创建权限
	for i := range permissions {
		DB.Create(&permissions[i])
	}
	fmt.Printf("成功创建 %d 条权限数据\n", len(permissions))

	// 创建超级管理员角色
	var adminRole models.Role
	DB.Where("code = ?", "admin").FirstOrCreate(&adminRole, models.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Description: "拥有所有权限的超级管理员",
		Sort:        1,
		Status:      1,
	})
	fmt.Println("超级管理员角色创建成功")

	// 为超级管理员角色分配所有权限
	var allPermissions []models.Permission
	DB.Find(&allPermissions)
	for _, perm := range allPermissions {
		DB.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE role_id=role_id", adminRole.ID, perm.ID)
	}
	fmt.Println("已为超级管理员分配所有权限")

	// 创建默认菜单数据
	menus := []models.Menu{
		{Name: "系统管理", Type: 1, Path: "/system", Icon: "Setting", Sort: 1, Status: 1},
		{Name: "用户管理", Type: 2, Path: "/user", Component: "/views/system/user/index", Icon: "User", ParentID: intPtr(1), Sort: 1, Status: 1},
		{Name: "咨询师管理", Type: 2, Path: "/counselor", Component: "/views/system/counselor/index", Icon: "UserFilled", ParentID: intPtr(1), Sort: 2, Status: 1},
		{Name: "角色管理", Type: 2, Path: "/roles", Component: "/views/system/roles/index", Icon: "Avatar", ParentID: intPtr(1), Sort: 3, Status: 1},
		{Name: "权限管理", Type: 2, Path: "/permissions", Component: "/views/system/permissions/index", Icon: "Key", ParentID: intPtr(1), Sort: 4, Status: 1},
		{Name: "菜单管理", Type: 2, Path: "/menus", Component: "/views/system/menus/index", Icon: "Menu", ParentID: intPtr(1), Sort: 5, Status: 1},
		{Name: "在线用户", Type: 2, Path: "/online", Component: "/views/system/online/index", Icon: "User", ParentID: intPtr(1), Sort: 6, Status: 1},

		{Name: "业务管理", Type: 1, Path: "/business", Icon: "Document", Sort: 2, Status: 1},
		{Name: "订单管理", Type: 2, Path: "/order", Component: "/views/business/order/index", Icon: "Document", ParentID: intPtr(7), Sort: 1, Status: 1},
		{Name: "聊天记录", Type: 2, Path: "/chat", Component: "/views/business/chat/index", Icon: "ChatDotRound", ParentID: intPtr(7), Sort: 2, Status: 1},

		{Name: "财务管理", Type: 1, Path: "/finance", Icon: "Wallet", Sort: 3, Status: 1},
		{Name: "提现审核", Type: 2, Path: "/withdraw", Component: "/views/finance/withdraw/index", Icon: "Wallet", ParentID: intPtr(9), Sort: 1, Status: 1},
		{Name: "财务统计", Type: 2, Path: "/statistics", Component: "/views/finance/statistics/index", Icon: "DataLine", ParentID: intPtr(9), Sort: 2, Status: 1},

		{Name: "低代码平台", Type: 1, Path: "/lowcode", Icon: "Grid", Sort: 4, Status: 1},
		{Name: "表单设计", Type: 2, Path: "/lowcode/forms", Component: "/views/lowcode/forms/index", Icon: "Edit", ParentID: intPtr(11), Sort: 1, Status: 1},
		{Name: "页面设计", Type: 2, Path: "/lowcode/pages", Component: "/views/lowcode/pages/index", Icon: "Grid", ParentID: intPtr(11), Sort: 2, Status: 1},
		{Name: "数据管理", Type: 2, Path: "/lowcode/data", Component: "/views/lowcode/data/index", Icon: "Database", ParentID: intPtr(11), Sort: 3, Status: 1},
	}

	// 先检查是否已有菜单数据
	var menuCount int64
	DB.Model(&models.Menu{}).Count(&menuCount)
	if menuCount > 0 {
		fmt.Printf("数据库中已有 %d 条菜单数据，跳过初始化\n", menuCount)
		return
	}

	// 批量创建菜单
	for i := range menus {
		if err := DB.Create(&menus[i]).Error; err != nil {
			fmt.Printf("创建菜单失败: %s, 错误: %v\n", menus[i].Name, err)
		} else {
			fmt.Printf("创建菜单成功: %s (ID: %d)\n", menus[i].Name, menus[i].ID)
		}
	}
	fmt.Printf("成功创建 %d 条菜单数据\n", len(menus))
}

// intPtr 返回int指针的辅助函数
func intPtr(i int) *int {
	return &i
}
