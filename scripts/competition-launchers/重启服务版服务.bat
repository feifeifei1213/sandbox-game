@echo off
chcp 65001 >nul
setlocal

set "DEPLOY_ROOT=D:\deploy\sandbox-game-guibing"
set "APP_PORT=28080"

echo.
echo ==============================================
echo   沙盘经营系统 - 重启服务版服务
echo ==============================================
echo.
echo 正在停止服务版服务，请稍候...

powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 2 >nul

call "%~dp0启动服务版.bat"
exit /b %errorlevel%
