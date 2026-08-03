@echo off
chcp 65001 >nul
setlocal

set "MYSQL_SERVICE="
set "MYSQL_SERVICE_CANDIDATES=SandboxGameMySQL,MySQL,MySQL80,MySQL57,mysql,MariaDB,mariadb"
set "DEPLOY_ROOT=%~dp0"
if "%DEPLOY_ROOT:~-1%"=="\" set "DEPLOY_ROOT=%DEPLOY_ROOT:~0,-1%"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "INIT_SCRIPT=%DEPLOY_ROOT%\scripts\init-competition.ps1"
set "APP_PORT="

if exist "%CONFIG_PATH%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "$path = '%CONFIG_PATH%'; $inServer = $false; foreach ($line in Get-Content -LiteralPath $path) { if ($line -match '^\s*server\s*:') { $inServer = $true; continue }; if ($inServer -and $line -match '^\S') { $inServer = $false }; if ($inServer -and $line -match '^\s*port\s*:\s*([0-9]+)') { Write-Output $matches[1]; break } }"`) do set "APP_PORT=%%i"
)
if not defined APP_PORT set "APP_PORT=8080"

echo.
echo ==============================================
echo   沙盘经营系统 - 首次初始化数据库
echo ==============================================
echo.
echo 本脚本用于新部署目录第一次初始化数据库。
echo 如果比赛已经开始，请不要使用这个脚本。
echo.
set /p CONFIRM_TEXT=如果你确认现在是首次部署或赛前首次准备，请输入 INIT 后回车继续：
if /I not "%CONFIRM_TEXT%"=="INIT" (
  echo 已取消，没有执行任何初始化操作。
  pause
  exit /b 0
)

if not exist "%INIT_SCRIPT%" (
  echo [错误] 未找到初始化脚本：%INIT_SCRIPT%
  pause
  exit /b 1
)

if not exist "%CONFIG_PATH%" (
  echo [错误] 未找到配置文件：%CONFIG_PATH%
  pause
  exit /b 1
)

powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 1 >nul

echo [1/2] 检查 MySQL 服务...
if not defined MYSQL_SERVICE (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "$candidates = '%MYSQL_SERVICE_CANDIDATES%'.Split(','); foreach ($name in $candidates) { try { $svc = Get-Service -Name $name -ErrorAction Stop; if ($svc) { Write-Output $svc.Name; break } } catch {} }"`) do set "MYSQL_SERVICE=%%i"
)
if not defined MYSQL_SERVICE (
  echo [错误] 未找到可用的 MySQL 服务。
  pause
  exit /b 1
)
powershell.exe -NoProfile -Command "$svc = Get-Service -Name '%MYSQL_SERVICE%' -ErrorAction SilentlyContinue; if (-not $svc) { exit 2 }; if ($svc.Status -eq 'Running') { exit 0 } else { exit 1 }"
if errorlevel 2 (
  echo [错误] 未找到 MySQL 服务：%MYSQL_SERVICE%
  pause
  exit /b 1
)
if errorlevel 1 (
  echo [提示] MySQL 未启动，正在尝试启动...
  powershell.exe -NoProfile -Command "Start-Service -Name '%MYSQL_SERVICE%'"
  if errorlevel 1 (
    echo [错误] MySQL 服务启动失败。
    pause
    exit /b 1
  )
)

echo [2/2] 正在初始化数据库...
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%INIT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [错误] 数据库初始化失败。
  pause
  exit /b 1
)

echo.
echo [完成] 数据库已初始化。
echo 下一步：请双击“启动沙盘系统.bat”，然后在网页赛前配置页选择本场比赛版本包。
echo.
pause
exit /b 0
