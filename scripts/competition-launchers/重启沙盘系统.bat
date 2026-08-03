@echo off
chcp 65001 >nul
setlocal

set "SANDBOX_GAME_SKIP_PAUSE=1"
call "%~dp0停止沙盘系统.bat"
if errorlevel 1 (
  echo [错误] 停止服务失败，已取消重启。
  pause
  exit /b 1
)

call "%~dp0启动沙盘系统.bat"
exit /b %errorlevel%
