@echo off
chcp 65001 >nul
setlocal

for %%A in ("%~dp0.") do set "DEPLOY_ROOT=%%~fA"

set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "BACKUP_SCRIPT=%DEPLOY_ROOT%\scripts\backup-competition-database.ps1"

echo.
echo ==============================================
echo   Sandbox Game - Backup Current Database
echo ==============================================
echo.
echo Deploy root: %DEPLOY_ROOT%
echo Config file: %CONFIG_PATH%
echo.
echo This script only exports data. It will not reset or modify the database.
echo Backup zip files will be created under the backup folder.
echo.

if not exist "%CONFIG_PATH%" (
  echo [ERROR] Config file not found: %CONFIG_PATH%
  pause
  exit /b 1
)

if not exist "%BACKUP_SCRIPT%" (
  echo [ERROR] Backup script not found: %BACKUP_SCRIPT%
  pause
  exit /b 1
)

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%BACKUP_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo.
  echo [ERROR] Database backup failed. Please check the error above.
  pause
  exit /b 1
)

echo.
echo [DONE] Database backup finished.
echo.
pause
exit /b 0
