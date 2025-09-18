@echo off
echo ========================================
echo    支付系统状态检查脚本 (Windows)
echo ========================================
echo.

REM 检查Docker是否安装
docker --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未找到Docker
    pause
    exit /b 1
)

echo [信息] 正在检查服务状态...
echo.

REM 切换到项目根目录
cd /d "%~dp0"

echo [Epusdt系统状态]
docker-compose ps mysql redis epusdt
echo.

echo [Payment Link MVP系统状态]
cd payment-link-mvp
docker-compose ps
echo.

REM 返回根目录
cd ..

echo [网络状态]
docker network ls | findstr epusdt
docker network ls | findstr payment
echo.

echo [端口占用情况]
netstat -an | findstr ":3000"
netstat -an | findstr ":8000"
netstat -an | findstr ":8080"
netstat -an | findstr ":3306"
netstat -an | findstr ":6379"
echo.

echo [服务健康检查]
echo 检查Epusdt管理后台...
curl -s -o nul -w "HTTP状态码: %%{http_code}\n" http://localhost:8000 || echo "无法连接到Epusdt管理后台"

echo 检查Payment Link前端...
curl -s -o nul -w "HTTP状态码: %%{http_code}\n" http://localhost:3000 || echo "无法连接到Payment Link前端"

echo 检查Payment Link后端API...
curl -s -o nul -w "HTTP状态码: %%{http_code}\n" http://localhost:8080/api/health || echo "无法连接到Payment Link后端API"
echo.

echo ========================================
echo           状态检查完成！
echo ========================================
echo.
pause
