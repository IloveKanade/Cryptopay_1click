@echo off
echo ========================================
echo    支付系统日志查看脚本 (Windows)
echo ========================================
echo.

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到Docker
    pause
    exit /b 1
)

echo [信息] 正在显示所有服务日志...
echo 按 Ctrl+C 退出日志查看
echo.

REM 切换到项目根目录
cd /d "%~dp0"

REM 显示所有服务的日志
docker-compose logs -f --tail=100

pause
