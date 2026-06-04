@echo off
chcp 65001 >nul
setlocal

set "MYSQL_SERVICE="
set "MYSQL_SERVICE_CANDIDATES=SandboxGameMySQL,MySQL,MySQL80,MySQL57,mysql,MariaDB,mariadb"
set "GAME_NAME=服务版"
set "DEPLOY_ROOT=D:\deploy\sandbox-game-guibing"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "START_SCRIPT=%DEPLOY_ROOT%\scripts\start-competition.ps1"
set "LOG_DIR=%DEPLOY_ROOT%\logs"
set "APP_PORT=28080"
set "LOGIN_URL=http://127.0.0.1:%APP_PORT%/sandbox-game/login"

echo.
echo ==============================================
echo   沙盘经营系统 - %GAME_NAME% 启动脚本
echo ==============================================
echo.

if not exist "%DEPLOY_ROOT%" (
  echo [错误] 未找到正式包目录：%DEPLOY_ROOT%
  echo 请先让技术支持确认服务版正式包已放到约定目录。
  pause
  exit /b 1
)

if not exist "%CONFIG_PATH%" (
  echo [错误] 未找到配置文件：%CONFIG_PATH%
  echo 请先让技术支持确认正式包配置文件已准备好。
  pause
  exit /b 1
)

if not exist "%START_SCRIPT%" (
  echo [错误] 未找到启动脚本：%START_SCRIPT%
  pause
  exit /b 1
)

if not exist "%DEPLOY_ROOT%\frontend\dist\index.html" (
  echo [错误] 未找到前端静态文件：%DEPLOY_ROOT%\frontend\dist\index.html
  echo 请先让技术支持确认正式包中的 frontend\dist 已完整复制。
  pause
  exit /b 1
)

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
  echo [提示] MySQL 未启动，正在尝试启动...
  powershell.exe -NoProfile -Command "Start-Service -Name '%MYSQL_SERVICE%'"
  if errorlevel 1 (
    echo [错误] MySQL 服务启动失败。
    echo 请联系技术支持处理数据库服务。
    pause
    exit /b 1
  )
) else (
  echo [通过] MySQL 已运行。
)

echo [2/3] 停掉旧的服务版进程...
powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 1 >nul

echo [3/3] 启动服务版服务...
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "$ErrorActionPreference='Stop'; New-Item -ItemType Directory -Force -Path '%LOG_DIR%' | Out-Null; Start-Process -FilePath 'powershell.exe' -ArgumentList @('-NoProfile','-ExecutionPolicy','Bypass','-File','%START_SCRIPT%','-ConfigPath','%CONFIG_PATH%') -WorkingDirectory '%DEPLOY_ROOT%' -WindowStyle Hidden -RedirectStandardOutput '%LOG_DIR%\server-%APP_PORT%.log' -RedirectStandardError '%LOG_DIR%\server-%APP_PORT%.err.log' | Out-Null"
if errorlevel 1 (
  echo [错误] 服务版服务启动失败。
  pause
  exit /b 1
)
timeout /t 2 >nul

echo.
echo [完成] 服务版已启动。
echo 本机登录页：%LOGIN_URL%
echo 局域网玩家地址：请把 127.0.0.1 换成管理员电脑的局域网 IP。
echo 注意：competition.yaml 中的 server.port 需要与 %APP_PORT% 保持一致。
echo.
start "" "%LOGIN_URL%"
echo 已自动打开登录页。
echo 如果网页打不开，请先尝试使用“重启服务版服务.bat”。
echo.
pause
exit /b 0
