@echo off
cd /d "d:\gospace\src\akrick.com\mychat\admin\frontend"
echo Checking node_modules...
if not exist node_modules (
    echo Installing dependencies...
    call npm install
)
echo.
echo Starting frontend server...
npm run dev
pause
