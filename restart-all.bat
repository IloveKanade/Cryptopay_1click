@echo off
echo ========================================
echo    支付系统一键重启脚本 (Windows)
echo ========================================
echo.

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到Docker，请先安装Docker Desktop
    pause
    exit /b 1
)

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo [错误] Docker未运行，请启动Docker Desktop
    pause
    exit /b 1
)

echo [信息] 正在重启所有服务...
echo.

REM 切换到项目根目录
cd /d "%~dp0"

REM 停止所有服务
echo [步骤1/3] 停止所有服务...
call stop-all.bat

if errorlevel 1 (
    echo [警告] 停止服务时出现错误，继续执行...
)

echo.
echo [步骤2/3] 等待服务完全停止...
timeout /t 5 /nobreak >nul

REM 启动所有服务
echo [步骤3/3] 重新启动所有服务...
call start-all.bat

echo.
echo [完成] 所有服务重启完成！
pause
