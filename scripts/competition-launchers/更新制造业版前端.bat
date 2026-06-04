@echo off
chcp 65001 >nul
setlocal

set "GAME_NAME=制造业版"
set "DEPLOY_ROOT=D:\deploy\sandbox-game-shengchan"
set "TARGET_FRONTEND_ROOT=%DEPLOY_ROOT%\frontend"
set "TARGET_DIST=%TARGET_FRONTEND_ROOT%\dist"
set "SOURCE_ROOT=%~dp0"
set "SOURCE_DIST=%SOURCE_ROOT%frontend\dist"
set "BACKUP_ROOT=D:\deploy\_backup\shengchan"
set "APP_PORT=18080"
set "LAUNCHER_ROOT=D:\deploy\launcher"
set "RESTART_SCRIPT=%LAUNCHER_ROOT%\重启制造业版服务.bat"

echo.
echo ==============================================
echo   沙盘经营系统 - 制造业版前端更新
echo ==============================================
echo.
echo 本脚本只适用于“前端页面文字 / 样式 / 提醒项”更新。
echo 请先把整个更新包解压到任意本地临时目录，再双击本脚本。
echo 不要把更新包直接解压到 D:\deploy\sandbox-game-shengchan。
echo.

if not exist "%SOURCE_DIST%\index.html" (
  echo [错误] 未找到更新包中的前端文件：%SOURCE_DIST%\index.html
  echo 请确认当前目录下是否带有 frontend\dist\ 完整内容。
  pause
  exit /b 1
)

if not exist "%DEPLOY_ROOT%" (
  echo [错误] 未找到正式版目录：%DEPLOY_ROOT%
  echo 请先确认管理员电脑上的制造业正式包已经部署完成。
  pause
  exit /b 1
)

if not exist "%TARGET_DIST%\index.html" (
  echo [错误] 未找到当前正式版前端文件：%TARGET_DIST%\index.html
  echo 请先确认制造业正式包目录完整，或联系交付方处理。
  pause
  exit /b 1
)

for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -Command "Get-Date -Format 'yyyyMMdd-HHmmss'"`) do set "TIMESTAMP=%%i"
if not defined TIMESTAMP (
  set "TIMESTAMP=manual"
)
set "BACKUP_PATH=%BACKUP_ROOT%\frontend-dist-%TIMESTAMP%"

echo [1/4] 停止当前制造业版服务...
powershell.exe -NoProfile -Command "$port = %APP_PORT%; Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue }"
timeout /t 1 >nul

echo [2/4] 备份当前前端到：%BACKUP_PATH%
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "$ErrorActionPreference='Stop'; New-Item -ItemType Directory -Force -Path '%BACKUP_ROOT%' | Out-Null; Copy-Item -Path '%TARGET_DIST%' -Destination '%BACKUP_PATH%' -Recurse -Force"
if errorlevel 1 (
  echo [错误] 备份旧前端失败，已停止更新。
  pause
  exit /b 1
)

echo [3/4] 覆盖新的前端文件...
powershell.exe -NoProfile -ExecutionPolicy Bypass -Command "$ErrorActionPreference='Stop'; if (Test-Path '%TARGET_DIST%') { Remove-Item -Path '%TARGET_DIST%' -Recurse -Force }; New-Item -ItemType Directory -Force -Path '%TARGET_FRONTEND_ROOT%' | Out-Null; Copy-Item -Path '%SOURCE_DIST%' -Destination '%TARGET_DIST%' -Recurse -Force"
if errorlevel 1 (
  echo [错误] 覆盖新前端失败。
  echo 旧前端备份仍在：%BACKUP_PATH%
  pause
  exit /b 1
)

echo [4/4] 尝试重启制造业版服务...
if exist "%RESTART_SCRIPT%" (
  call "%RESTART_SCRIPT%"
  exit /b %errorlevel%
)

echo [提示] 已完成前端更新，但未找到重启脚本：%RESTART_SCRIPT%
echo 请手工双击 D:\deploy\launcher\重启制造业版服务.bat
echo 浏览器重新打开后，请按 Ctrl+F5 强制刷新页面。
echo.
pause
exit /b 0
