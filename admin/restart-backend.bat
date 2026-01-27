@echo off
chcp 65001 >nul
echo ========================================
echo   重启后端服务
echo ========================================
echo.

cd /d "%~dp0backend"

echo [1/3] 检查并停止现有进程...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8081') do (
    echo 找到进程 %%a，正在停止...
    taskkill /F /PID %%a 2>nul
)
echo [√] 旧进程已停止
echo.

echo [2/3] 启动后端服务...
start "MyChat Admin Backend" cmd /k "go run main.go"
echo [√] 后端服务启动中...
echo.

echo [3/3] 等待服务就绪...
timeout /t 5 >nul

echo.
echo ========================================
echo   检查服务状态
echo ========================================
echo.

powershell -Command "if (Test-NetConnection -ComputerName localhost -Port 8081 -InformationLevel Quiet) { Write-Host '[✓] 后端服务启动成功！' } else { Write-Host '[✗] 后端服务启动失败' }"
echo.

echo ========================================
echo   服务信息
echo ========================================
echo 后端地址: http://localhost:8081
echo API文档: http://localhost:8081/swagger/index.html
echo ========================================
echo.

echo 按任意键退出...
pause >nul
