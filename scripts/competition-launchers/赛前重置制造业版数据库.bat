@echo off
chcp 65001 >nul
setlocal

set "MYSQL_SERVICE="
set "MYSQL_SERVICE_CANDIDATES=SandboxGameMySQL,MySQL,MySQL80,MySQL57,mysql,MariaDB,mariadb"
set "DEPLOY_ROOT=D:\deploy\sandbox-game-shengchan"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "RESET_SCRIPT=%DEPLOY_ROOT%\scripts\reset-competition.ps1"
set "INIT_SCRIPT=%DEPLOY_ROOT%\scripts\init-competition.ps1"
set "APP_PORT=18080"

echo.
echo ==========================================================
echo   警告：这是“赛前重置制造业版数据库”脚本
echo ==========================================================
echo.
echo 这个操作会清空制造业版当前比赛数据。
echo 比赛正式开始后，不要使用这个脚本。
echo.
set /p CONFIRM_TEXT=如果你确认现在仍处于赛前，请输入 RESET 后回车继续：
if /I not "%CONFIRM_TEXT%"=="RESET" (
  echo 已取消，没有执行任何重置操作。
  pause
  exit /b 0
)

powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 2 >nul

echo [1/3] 检查 MySQL 服务...
if not defined MYSQL_SERVICE (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "$candidates = '%MYSQL_SERVICE_CANDIDATES%'.Split(','); foreach ($name in $candidates) { try { $svc = Get-Service -Name $name -ErrorAction Stop; if ($svc) { Write-Output $svc.Name; break } } catch {} }"`) do set "MYSQL_SERVICE=%%i"
)
if not defined MYSQL_SERVICE (
  echo [错误] 未找到可用的 MySQL 服务。
  echo 已尝试自动识别：SandboxGameMySQL、MySQL、MySQL80、MySQL57、mysql、MariaDB、mariadb
  echo 如果你的 MySQL 服务名是自定义的，请编辑本脚本顶部的 MYSQL_SERVICE 变量。
  pause
  exit /b 1
)
echo [通过] 已识别 MySQL 服务：%MYSQL_SERVICE%
powershell.exe -NoProfile -Command "$svc = Get-Service -Name '%MYSQL_SERVICE%' -ErrorAction SilentlyContinue; if (-not $svc) { exit 2 }; if ($svc.Status -eq 'Running') { exit 0 } else { exit 1 }"
if errorlevel 2 (
  echo [错误] 未找到 MySQL 服务：%MYSQL_SERVICE%
  echo 请先确认本机 MySQL 已安装，或把脚本顶部的 MYSQL_SERVICE 改成真实服务名。
  pause
  exit /b 1
)
if errorlevel 1 (
  powershell.exe -NoProfile -Command "Start-Service -Name '%MYSQL_SERVICE%'"
  if errorlevel 1 (
    echo [错误] MySQL 服务启动失败。
    pause
    exit /b 1
  )
)

if not exist "%RESET_SCRIPT%" (
  echo [错误] 未找到重置脚本：%RESET_SCRIPT%
  pause
  exit /b 1
)

if not exist "%INIT_SCRIPT%" (
  echo [错误] 未找到初始化脚本：%INIT_SCRIPT%
  pause
  exit /b 1
)

echo [2/3] 正在重置制造业版数据库...
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%RESET_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [错误] 制造业版数据库重置失败。
  pause
  exit /b 1
)

echo [3/3] 正在重新初始化制造业版数据库...
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%INIT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [错误] 制造业版数据库初始化失败。
  pause
  exit /b 1
)

echo.
echo [完成] 制造业版数据库已重置并重新初始化。
echo 下一步：请双击“启动制造业版.bat”重新开赛前准备。
echo.
pause
exit /b 0
