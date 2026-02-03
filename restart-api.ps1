# 重启 API 服务
cd d:\gospace\src\akrick.com\mychat

# 停止所有 api.exe 进程
Get-Process -Name "api" -ErrorAction SilentlyContinue | Stop-Process -Force
Start-Sleep -Seconds 2

# 启动新服务
Start-Process -FilePath ".\api.exe" -NoNewWindow -RedirectStandardOutput "api.log" -RedirectStandardError "api-error.log"

# 等待服务启动
Start-Sleep -Seconds 3

# 测试 API
$Response = Invoke-WebRequest -Uri "http://localhost:3002/api/counselor/list?page=1&page_size=10" -TimeoutSec 5 -UseBasicParsing
Write-Host "Status: $($Response.StatusCode)"
Write-Host "Content: $($Response.Content)"
