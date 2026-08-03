@echo off
chcp 65001 >nul
setlocal

set "DEPLOY_ROOT=%~dp0"
if "%DEPLOY_ROOT:~-1%"=="\" set "DEPLOY_ROOT=%DEPLOY_ROOT:~0,-1%"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "APP_PORT="

if exist "%CONFIG_PATH%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "$path = '%CONFIG_PATH%'; $inServer = $false; foreach ($line in Get-Content -LiteralPath $path) { if ($line -match '^\s*server\s*:') { $inServer = $true; continue }; if ($inServer -and $line -match '^\S') { $inServer = $false }; if ($inServer -and $line -match '^\s*port\s*:\s*([0-9]+)') { Write-Output $matches[1]; break } }"`) do set "APP_PORT=%%i"
)
if not defined APP_PORT set "APP_PORT=8080"

echo.
echo ==============================================
echo   沙盘经营系统 - 停止服务
echo ==============================================
echo.
echo 正在停止端口 %APP_PORT% 上的沙盘系统进程，请稍候...

powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 2 >nul

echo.
echo [完成] 沙盘系统服务已停止。
echo 注意：本脚本不会停止 MySQL 服务。
echo.
if not defined SANDBOX_GAME_SKIP_PAUSE pause
exit /b 0
