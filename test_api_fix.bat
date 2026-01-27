@echo off
echo ========================================
echo 管理后台 API 数据访问修复测试
echo ========================================
echo.
echo 修复内容:
echo 1. 修复前端数据访问问题 (data.data -> res)
echo 2. 添加通用辅助函数 (utils.go)
echo 3. 修复所有页面的数据加载
echo.
echo 修复的文件列表:
echo - views/system/user/index.vue
echo - views/system/users/index.vue
echo - views/system/roles/index.vue
echo - views/system/permissions/index.vue
echo - views/system/menus/index.vue
echo - views/system/counselor/index.vue
echo - views/business/order/index.vue
echo - views/business/chat/index.vue
echo - views/finance/withdraw/index.vue
echo - views/finance/statistics/index.vue
echo - views/finance/reports/index.vue
echo - views/lowcode/data/index.vue
echo - views/dashboard/index.vue
echo - handlers/utils.go (新增)
echo.
echo ========================================
echo 启动服务进行测试...
echo ========================================
echo.
cd ..\api
echo 启动后端服务 (端口 8080)...
start "MyChat API" go run main.go
timeout /t 3 /nobreak >nul
cd ..\admin\frontend
echo 启动前端服务 (端口 3000)...
start "MyChat Admin" npm run dev
echo.
echo ========================================
echo 服务已启动，请在浏览器中访问:
echo - 前端: http://localhost:3000
echo - 后端: http://localhost:8080
echo - API文档: http://localhost:8080/swagger/index.html
echo ========================================
echo.
echo 默认登录账号:
echo - 用户名: admin
echo - 密码: admin123
echo.
echo 测试步骤:
echo 1. 打开浏览器访问 http://localhost:3000
echo 2. 使用 admin/admin123 登录
echo 3. 点击各个菜单项，检查是否正常加载数据
echo 4. 检查浏览器控制台是否还有错误
echo.
pause
