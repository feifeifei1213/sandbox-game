@echo off
chcp 65001 >nul
setlocal

for %%A in ("%~dp0.") do set "DEPLOY_ROOT=%%~fA"

set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "PORT_SCRIPT=%DEPLOY_ROOT%\scripts\get-competition-port.ps1"
set "STOP_PORT_SCRIPT=%DEPLOY_ROOT%\scripts\stop-competition-port.ps1"
set "MYSQL_SCRIPT=%DEPLOY_ROOT%\scripts\ensure-mysql-service.ps1"
set "START_BACKGROUND_SCRIPT=%DEPLOY_ROOT%\scripts\start-competition-background.ps1"
set "APP_PORT=8080"

if exist "%PORT_SCRIPT%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%PORT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"`) do set "APP_PORT=%%i"
)

set "LOGIN_URL=http://127.0.0.1:%APP_PORT%/sandbox-game/login"

echo.
echo ==============================================
echo   Sandbox Game - Start System
echo ==============================================
echo.
echo Deploy root: %DEPLOY_ROOT%
echo Config file: %CONFIG_PATH%
echo App port: %APP_PORT%
echo.

if not exist "%CONFIG_PATH%" (
  echo [ERROR] Config file not found: %CONFIG_PATH%
  pause
  exit /b 1
)

if not exist "%DEPLOY_ROOT%\frontend\dist\index.html" (
  echo [ERROR] Frontend file not found: %DEPLOY_ROOT%\frontend\dist\index.html
  pause
  exit /b 1
)

if not exist "%START_BACKGROUND_SCRIPT%" (
  echo [ERROR] Start helper not found: %START_BACKGROUND_SCRIPT%
  pause
  exit /b 1
)

if exist "%MYSQL_SCRIPT%" (
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%MYSQL_SCRIPT%"
  if errorlevel 1 (
    echo [ERROR] MySQL service check failed.
    pause
    exit /b 1
  )
)

if exist "%STOP_PORT_SCRIPT%" (
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%STOP_PORT_SCRIPT%" -Port %APP_PORT%
)

powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%START_BACKGROUND_SCRIPT%" -ConfigPath "%CONFIG_PATH%" -Port %APP_PORT% -LogDir "%DEPLOY_ROOT%\logs"
if errorlevel 1 (
  echo [ERROR] Failed to start Sandbox Game.
  pause
  exit /b 1
)

timeout /t 2 >nul
start "" "%LOGIN_URL%"

echo.
echo [DONE] Sandbox Game started.
echo Local login URL: %LOGIN_URL%
echo LAN users should replace 127.0.0.1 with this computer's LAN IP.
echo.
pause
exit /b 0
