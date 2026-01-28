@echo off
REM MyChat 服务启动脚本
REM 启动所有三个服务: 用户端API(8080)、管理后台API(8081)、WebSocket(8082)

echo ========================================
echo    MyChat 服务启动脚本
echo ========================================
echo.

REM 检查是否已安装 Go
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo [错误] 未检测到 Go 环境，请先安装 Go 1.20+
    pause
    exit /b 1
)

echo [1/3] 启动用户端 API (端口 8080)...
start "MyChat User API - 8080" cmd /k "cd /d %~dp0api && go run main.go"
timeout /t 2 >nul

echo [2/3] 启动管理后台 API (端口 8081)...
start "MyChat Admin API - 8081" cmd /k "cd /d %~dp0admin\backend && go run main.go"
timeout /t 2 >nul

echo [3/3] 启动 WebSocket 服务 (端口 8082)...
start "MyChat WebSocket - 8082" cmd /k "cd /d %~dp0websocket && go run main.go"
timeout /t 2 >nul

echo.
echo ========================================
echo    所有服务已启动！
echo ========================================
echo.
echo 服务地址:
echo   - 用户端 API:    http://localhost:8080
echo   - 管理后台 API:  http://localhost:8081
echo   - WebSocket:     ws://localhost:8082
echo.
echo 文档:
echo   - Swagger文档:   http://localhost:8080/swagger/index.html
echo   - 项目文档:      docs/README.md
echo.
echo 提示: 每个服务在独立的窗口中运行，关闭窗口即停止服务
echo.
pause
