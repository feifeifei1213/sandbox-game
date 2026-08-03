@echo off
chcp 65001 >nul
setlocal

for %%A in ("%~dp0.") do set "DEPLOY_ROOT=%%~fA"

set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "INIT_SCRIPT=%DEPLOY_ROOT%\scripts\init-competition.ps1"
set "PORT_SCRIPT=%DEPLOY_ROOT%\scripts\get-competition-port.ps1"
set "STOP_PORT_SCRIPT=%DEPLOY_ROOT%\scripts\stop-competition-port.ps1"
set "MYSQL_SCRIPT=%DEPLOY_ROOT%\scripts\ensure-mysql-service.ps1"
set "APP_PORT=8080"

if exist "%PORT_SCRIPT%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%PORT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"`) do set "APP_PORT=%%i"
)

echo.
echo ==============================================
echo   Sandbox Game - Init Competition Database
echo ==============================================
echo.
echo Deploy root: %DEPLOY_ROOT%
echo Config file: %CONFIG_PATH%
echo App port: %APP_PORT%
echo.
echo This script initializes the database for a new deployment.
echo If the competition has already started, do not continue.
echo.

set /p CONFIRM_TEXT=Type INIT and press Enter to continue:
if /I not "%CONFIRM_TEXT%"=="INIT" (
  echo Canceled. Nothing was changed.
  pause
  exit /b 0
)

if not exist "%CONFIG_PATH%" (
  echo [ERROR] Config file not found: %CONFIG_PATH%
  pause
  exit /b 1
)

if not exist "%INIT_SCRIPT%" (
  echo [ERROR] Init script not found: %INIT_SCRIPT%
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

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%INIT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"
if errorlevel 1 (
  echo [ERROR] Database initialization failed.
  pause
  exit /b 1
)

echo.
echo [DONE] Database initialized.
echo Next: double-click "启动沙盘系统.bat".
echo.
pause
exit /b 0
