@echo off
chcp 936 >nul
REM MyChat 服务停止脚本
REM 用于停止所有服务：用户端API、管理后台API、WebSocket

echo ========================================
echo    MyChat 服务停止脚本
echo ========================================
echo.

REM 停止端口 3002
echo [1/3] 正在停止用户端 API (端口 3002)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3002 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo 用户端 API 已停止

REM 停止端口 3003
echo [2/3] 正在停止管理后台 API (端口 3003)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3003 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo 管理后台 API 已停止

REM 停止端口 3004
echo [3/3] 正在停止 WebSocket 服务 (端口 3004)...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3004 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo WebSocket 服务已停止

REM 停止前端开发服务器
echo [4/4] 正在停止前端开发服务器...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr :3000 ^| findstr LISTENING') do (
    taskkill /F /PID %%a 2>nul
)
echo 前端开发服务器已停止

echo.
echo ========================================
echo    所有服务已停止
echo ========================================
echo.
pause
