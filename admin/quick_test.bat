@echo off
chcp 65001 >nul
echo ========================================
echo   MyChat 管理后台快速测试
echo ========================================
echo.

echo [检查 1/3] 检查后端服务...
powershell -Command "Test-NetConnection -ComputerName localhost -Port 8081 -InformationLevel Quiet" >nul 2>&1
if %errorlevel% neq 0 (
    echo [✗] 后端服务未运行
    echo.
    echo 正在启动后端服务...
    start "MyChat Admin Backend" cmd /c "cd /d %~dp0backend && go run main.go"
    timeout /t 3 >nul
) else (
    echo [✓] 后端服务运行正常
)
echo.

echo [检查 2/3] 检查前端服务...
powershell -Command "Test-NetConnection -ComputerName localhost -Port 3000 -InformationLevel Quiet" >nul 2>&1
if %errorlevel% neq 0 (
    echo [✗] 前端服务未运行
    echo.
    echo 正在启动前端服务...
    start "MyChat Admin Frontend" cmd /c "cd /d %~dp0frontend && npm run dev"
    timeout /t 5 >nul
) else (
    echo [✓] 前端服务运行正常
)
echo.

echo [检查 3/3] 测试API连接...
powershell -Command "try { $response = Invoke-WebRequest -Uri 'http://localhost:8081/api/test' -UseBasicParsing -TimeoutSec 5; Write-Host '[✓] API连接正常'; Write-Host '响应:' $response.Content } catch { Write-Host '[✗] API连接失败' }"
echo.

echo ========================================
echo   测试完成
echo ========================================
echo.
echo 访问地址:
echo   前端: http://localhost:3000
echo   后端: http://localhost:8081
echo.
echo 默认账号:
echo   用户名: admin
echo   密码: admin123
echo.
echo 按任意键打开浏览器...
pause >nul
start http://localhost:3000
