@echo off
chcp 65001 >nul
echo ========================================
echo   登录跳转功能测试
echo ========================================
echo.
echo [信息] 开始测试登录功能...
echo.

echo [1/4] 检查后端服务...
powershell -Command "Test-NetConnection -ComputerName localhost -Port 8081 -InformationLevel Quiet" >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 后端服务未运行！请先启动后端服务
    pause
    exit /b 1
)
echo [√] 后端服务运行正常 (端口 8081)
echo.

echo [2/4] 检查前端服务...
powershell -Command "Test-NetConnection -ComputerName localhost -Port 3000 -InformationLevel Quiet" >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 前端服务未运行！请先启动前端服务
    pause
    exit /b 1
)
echo [√] 前端服务运行正常 (端口 3000)
echo.

echo [3/4] 测试登录接口...
powershell -Command "$body = @{username='admin'; password='admin123'} | ConvertTo-Json; $response = Invoke-WebRequest -Uri 'http://localhost:8081/api/admin/login' -Method POST -Body $body -ContentType 'application/json' -UseBasicParsing; $response.Content" > test_login_response.json
type test_login_response.json
echo.
echo [√] 登录接口响应已保存到 test_login_response.json
echo.

echo [4/4] 检查登录响应格式...
powershell -Command "$response = Get-Content test_login_response.json | ConvertFrom-Json; if ($response.code -eq 200 -and $response.data.token) { Write-Host '[√] 登录响应格式正确，包含token' } else { Write-Host '[错误] 登录响应格式异常'; exit 1 }"
if %errorlevel% neq 0 (
    echo [错误] 登录响应格式检查失败！
    pause
    exit /b 1
)
echo.

echo ========================================
echo   测试完成！
echo ========================================
echo.
echo 请在浏览器中访问: http://localhost:3000
echo 使用账号: admin / admin123 进行登录测试
echo.
echo 打开浏览器控制台查看详细日志
echo.
pause
