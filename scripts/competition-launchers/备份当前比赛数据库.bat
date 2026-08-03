@echo off
chcp 65001 >nul
setlocal

set "DEPLOY_ROOT=%~dp0"
if "%DEPLOY_ROOT:~-1%"=="\" set "DEPLOY_ROOT=%DEPLOY_ROOT:~0,-1%"
set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "BACKUP_SCRIPT=%DEPLOY_ROOT%\scripts\backup-competition-database.ps1"

echo.
echo ==============================================
echo   沙盘经营系统 - 备份当前比赛数据库
echo ==============================================
echo.
echo 本脚本不会清空数据库，也不会修改比赛数据。
echo 备份文件会生成到当前系统目录下的 backup 文件夹。
echo.

if not exist "%DEPLOY_ROOT%" (
  echo [错误] 未找到系统目录：%DEPLOY_ROOT%
  pause
  exit /b 1
)

if not exist "%CONFIG_PATH%" (
  echo [错误] 未找到配置文件：%CONFIG_PATH%
  pause
  exit /b 1
)

if not exist "%BACKUP_SCRIPT%" (
  echo [错误] 未找到备份脚本：%BACKUP_SCRIPT%
  pause
  exit /b 1
)

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%BACKUP_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo.
  echo [错误] 数据库备份失败，请联系技术支持查看上方错误信息。
  pause
  exit /b 1
)

echo.
echo [完成] 数据库备份已完成。
echo.
pause
exit /b 0
