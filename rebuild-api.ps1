# 重新编译并启动 API 服务
Set-Location "d:\gospace\src\akrick.com\mychat\api"
go build -o ..\api.exe .
Start-Process powershell -ArgumentList "-NoExit", "-Command", "..\api.exe"
