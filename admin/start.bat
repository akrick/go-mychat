@echo off
chcp 65001 >nul
echo ========================================
echo    MyChat 管理后台启动脚本
echo ========================================
echo.

echo [1/3] 检查后端依赖...
cd backend
if not exist go.mod (
    echo [错误] go.mod 文件不存在！
    pause
    exit /b 1
)
go mod tidy
if %errorlevel% neq 0 (
    echo [错误] Go 依赖安装失败！
    pause
    exit /b 1
)
echo [√] 后端依赖检查完成
cd ..
echo.

echo [2/3] 检查前端依赖...
cd frontend
if not exist node_modules (
    echo [提示] 正在安装前端依赖，请稍候...
    call npm install
    if %errorlevel% neq 0 (
        echo [错误] 前端依赖安装失败！
        pause
        exit /b 1
    )
)
echo [√] 前端依赖检查完成
cd ..
echo.

echo [3/3] 启动服务...
echo.
echo ========================================
echo    选择启动模式
echo ========================================
echo [1] 仅启动后端服务 (端口 8081)
echo [2] 仅启动前端服务 (端口 3000)
echo [3] 同时启动后端和前端
echo [0] 退出
echo ========================================
set /p choice="请选择 (0-3): "

if "%choice%"=="1" (
    echo.
    echo [启动] 后端服务...
    cd backend
    start "MyChat Admin Backend" go run main.go
    echo [√] 后端服务已启动: http://localhost:8081
) else if "%choice%"=="2" (
    echo.
    echo [启动] 前端服务...
    cd frontend
    start "MyChat Admin Frontend" npm run dev
    echo [√] 前端服务已启动: http://localhost:3000
) else if "%choice%"=="3" (
    echo.
    echo [启动] 后端服务...
    cd backend
    start "MyChat Admin Backend" go run main.go
    echo [√] 后端服务已启动: http://localhost:8081
    cd ..
    timeout /t 2 >nul
    echo.
    echo [启动] 前端服务...
    cd frontend
    start "MyChat Admin Frontend" npm run dev
    echo [√] 前端服务已启动: http://localhost:3000
) else if "%choice%"=="0" (
    echo.
    echo [退出] 已取消启动
    pause
    exit /b 0
) else (
    echo.
    echo [错误] 无效的选择！
    pause
    exit /b 1
)

echo.
echo ========================================
echo    启动完成
echo ========================================
echo 管理后台前端: http://localhost:3000
echo 管理后台后端: http://localhost:8081
echo ========================================
echo.
echo 默认账号: admin
echo 默认密码: admin123
echo.
pause
