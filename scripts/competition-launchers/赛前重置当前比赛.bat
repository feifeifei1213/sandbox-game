@echo off
chcp 65001 >nul
setlocal

for %%A in ("%~dp0.") do set "DEPLOY_ROOT=%%~fA"

set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "RESET_SCRIPT=%DEPLOY_ROOT%\scripts\reset-competition.ps1"
set "PORT_SCRIPT=%DEPLOY_ROOT%\scripts\get-competition-port.ps1"
set "STOP_PORT_SCRIPT=%DEPLOY_ROOT%\scripts\stop-competition-port.ps1"
set "MYSQL_SCRIPT=%DEPLOY_ROOT%\scripts\ensure-mysql-service.ps1"
set "APP_PORT=8080"

if exist "%PORT_SCRIPT%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%PORT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"`) do set "APP_PORT=%%i"
)

echo.
echo ==============================================
echo   Sandbox Game - Reset Current Competition
echo ==============================================
echo.
echo Deploy root: %DEPLOY_ROOT%
echo Config file: %CONFIG_PATH%
echo App port: %APP_PORT%
echo.
echo WARNING: This will reset the database configured in competition.yaml.
echo Please run "备份当前比赛数据库.bat" before continuing.
echo.

set /p BACKUP_CONFIRM=Type BACKUP if you have already backed up data:
if /I not "%BACKUP_CONFIRM%"=="BACKUP" (
  echo Canceled. Nothing was changed.
  pause
  exit /b 0
)

set /p RESET_CONFIRM=Type RESET to confirm database reset:
if /I not "%RESET_CONFIRM%"=="RESET" (
  echo Canceled. Nothing was changed.
  pause
  exit /b 0
)

if not exist "%CONFIG_PATH%" (
  echo [ERROR] Config file not found: %CONFIG_PATH%
  pause
  exit /b 1
)

if not exist "%RESET_SCRIPT%" (
  echo [ERROR] Reset script not found: %RESET_SCRIPT%
  pause
  exit /b 1
)

if exist "%STOP_PORT_SCRIPT%" (
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%STOP_PORT_SCRIPT%" -Port %APP_PORT%
)

if exist "%MYSQL_SCRIPT%" (
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%MYSQL_SCRIPT%"
  if errorlevel 1 (
    echo [ERROR] MySQL service check failed.
    pause
    exit /b 1
  )
)

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%RESET_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [ERROR] Database reset failed.
  pause
  exit /b 1
)

echo.
echo [DONE] Database has been reset to pre-setup state.
echo Next:
echo   1. Double-click "启动沙盘系统.bat".
echo   2. Open the admin page.
echo   3. Configure group count and edition.
echo   4. Confirm game initialization.
echo.
pause
exit /b 0
