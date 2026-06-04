@echo off
chcp 65001 >nul
setlocal

echo.
echo ==============================================
echo   沙盘经营系统 - 停止全部服务
echo ==============================================
echo.
echo 正在停止制造业版和服务版服务进程，请稍候...

powershell.exe -NoProfile -Command "$ports = @(18080, 28080); foreach ($port in $ports) { Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue } }"
timeout /t 2 >nul

echo.
echo [完成] 两套游戏服务进程都已停止。
echo 注意：本脚本不会停止 MySQL 服务。
echo.
pause
exit /b 0
