@echo off
REM MyChat 服务停止脚本
REM 停止所有三个服务

echo ========================================
echo    MyChat 服务停止脚本
echo ========================================
echo.

REM 停止端口 8080
echo [1/3] 停止用户端 API (端口 8080)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8080 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo 用户端 API 已停止

REM 停止端口 8081
echo [2/3] 停止管理后台 API (端口 8081)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8081 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo 管理后台 API 已停止

REM 停止端口 8082
echo [3/3] 停止 WebSocket 服务 (端口 8082)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :8082 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo WebSocket 服务已停止

echo.
echo ========================================
echo    所有服务已停止
echo ========================================
echo.
pause
