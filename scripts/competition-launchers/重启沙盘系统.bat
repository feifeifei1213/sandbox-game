@echo off
chcp 65001 >nul
setlocal

for %%A in ("%~dp0.") do set "DEPLOY_ROOT=%%~fA"

set "CONFIG_PATH=%DEPLOY_ROOT%\configs\competition.yaml"
set "PORT_SCRIPT=%DEPLOY_ROOT%\scripts\get-competition-port.ps1"
set "STOP_PORT_SCRIPT=%DEPLOY_ROOT%\scripts\stop-competition-port.ps1"
set "APP_PORT=8080"

if exist "%PORT_SCRIPT%" (
  for /f "usebackq delims=" %%i in (`powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%PORT_SCRIPT%" -ConfigPath "%CONFIG_PATH%"`) do set "APP_PORT=%%i"
)

echo.
echo ==============================================
echo   Sandbox Game - Restart System
echo ==============================================
echo.
echo App port: %APP_PORT%
echo.

if exist "%STOP_PORT_SCRIPT%" (
  powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%STOP_PORT_SCRIPT%" -Port %APP_PORT%
)

call "%DEPLOY_ROOT%\启动沙盘系统.bat"
exit /b %errorlevel%
