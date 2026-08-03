@echo off
chcp 65001 >nul
setlocal

set "DEPLOY_ROOT=%~dp0"
if "%DEPLOY_ROOT:~-1%"=="\" set "DEPLOY_ROOT=%DEPLOY_ROOT:~0,-1%"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "RESET_SCRIPT=%DEPLOY_ROOT%\scripts\reset-competition.ps1"
set "APP_PORT="

if exist "%CONFIG_PATH%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "$path = '%CONFIG_PATH%'; $inServer = $false; foreach ($line in Get-Content -LiteralPath $path) { if ($line -match '^\s*server\s*:') { $inServer = $true; continue }; if ($inServer -and $line -match '^\S') { $inServer = $false }; if ($inServer -and $line -match '^\s*port\s*:\s*([0-9]+)') { Write-Output $matches[1]; break } }"`) do set "APP_PORT=%%i"
)
if not defined APP_PORT set "APP_PORT=8080"

echo.
echo ==============================================
echo   沙盘经营系统 - 赛前重置当前比赛
echo ==============================================
echo.
echo [强提醒] 本脚本会清空当前配置连接的比赛数据库并重新初始化。
echo [强提醒] 如果管理员电脑已经有比赛数据，必须先运行“备份当前比赛数据库.bat”。
echo.
set /p BACKUP_TEXT=如果你确认已经完成备份，请输入 BACKUP 后回车继续：
if /I not "%BACKUP_TEXT%"=="BACKUP" (
  echo 已取消，没有执行任何重置操作。
  pause
  exit /b 0
)
set /p RESET_TEXT=如果你确认现在是赛前重置，请输入 RESET 后回车继续：
if /I not "%RESET_TEXT%"=="RESET" (
  echo 已取消，没有执行任何重置操作。
  pause
  exit /b 0
)

if not exist "%RESET_SCRIPT%" (
  echo [错误] 未找到重置脚本：%RESET_SCRIPT%
  pause
  exit /b 1
)

if not exist "%CONFIG_PATH%" (
  echo [错误] 未找到配置文件：%CONFIG_PATH%
  pause
  exit /b 1
)

echo [1/2] 停止当前系统进程...
powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 1 >nul

echo [2/2] 正在重置当前比赛数据库...
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%RESET_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [错误] 数据库重置失败。
  pause
  exit /b 1
)

echo.
echo [完成] 当前比赛数据库已重置。
echo 下一步：请双击“启动沙盘系统.bat”，然后在网页赛前配置页重新初始化比赛。
echo.
pause
exit /b 0
