@echo off
cd /d "d:\gospace\src\akrick.com\mychat\admin\backend"
echo Installing dependencies...
go mod tidy
echo.
echo Starting backend server...
go run main.go
pause
