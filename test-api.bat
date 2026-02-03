@echo off
cd /d d:\gospace\src\akrick.com\mychat
start /B api.exe
timeout /t 2 /nobreak >nul
echo Testing API...
curl -s -m 5 "http://localhost:3002/api/counselor/list?page=1&page_size=10"
