@echo off
chcp 65001 >nul
title USDT支付系统启动

echo.
echo ========================================
echo    USDT支付系统启动
echo ========================================
echo.

echo 当前工作目录: %cd%
echo.

:: 检查Docker
echo 检查Docker环境...
docker --version
if %errorlevel% neq 0 (
    echo [错误] Docker未安装或未启动
    pause
    exit /b 1
)

:: 检查Docker Compose
echo 检查Docker Compose...
docker-compose --version
if %errorlevel% neq 0 (
    echo [错误] Docker Compose未安装
    pause
    exit /b 1
)

:: 创建目录
echo 创建必要目录...
if not exist "epusdt\data\mysql" mkdir epusdt\data\mysql
if not exist "epusdt\data\redis" mkdir epusdt\data\redis
if not exist "epusdt\logs" mkdir epusdt\logs

:: 清理旧容器
echo 清理旧容器...
docker-compose down
docker stop epusdt-mysql epusdt-redis epusdt-app payment-backend payment-frontend 2>nul
docker rm epusdt-mysql epusdt-redis epusdt-app payment-backend payment-frontend 2>nul

:: 创建Dockerfile
echo 创建Dockerfile文件...

:: Epusdt Dockerfile
cd epusdt\src
echo FROM golang:1.21-alpine > Dockerfile
echo RUN apk add --no-cache git >> Dockerfile
echo WORKDIR /app >> Dockerfile
echo COPY . . >> Dockerfile
echo RUN go mod tidy >> Dockerfile
echo RUN go build -o epusdt main.go >> Dockerfile
echo EXPOSE 8000 >> Dockerfile
echo CMD ["./epusdt", "http", "start"] >> Dockerfile
cd ..\..

:: 后端 Dockerfile
cd payment-link-mvp\backend
echo FROM golang:1.21-alpine > Dockerfile
echo RUN apk add --no-cache git >> Dockerfile
echo WORKDIR /app >> Dockerfile
echo COPY . . >> Dockerfile
echo RUN go mod tidy >> Dockerfile
echo RUN go build -o main main.go >> Dockerfile
echo EXPOSE 8080 >> Dockerfile
echo CMD ["./main"] >> Dockerfile
cd ..\..

:: 前端 Dockerfile
cd payment-link-mvp\frontend
echo FROM node:18-alpine AS builder > Dockerfile
echo WORKDIR /app >> Dockerfile
echo COPY package*.json ./ >> Dockerfile
echo RUN npm ci >> Dockerfile
echo COPY . . >> Dockerfile
echo RUN npm run build >> Dockerfile
echo FROM nginx:alpine >> Dockerfile
echo COPY --from=builder /app/build /usr/share/nginx/html >> Dockerfile
echo COPY nginx.conf /etc/nginx/nginx.conf >> Dockerfile
echo EXPOSE 3001 >> Dockerfile
echo CMD ["nginx", "-g", "daemon off;"] >> Dockerfile
cd ..\..

:: 启动服务
echo 启动所有服务...
docker-compose up -d --build

if %errorlevel% equ 0 (
    echo.
    echo ========================================
    echo    启动成功！
    echo ========================================
    echo.
    echo 服务访问地址:
    echo    - 前端界面: http://localhost:3001
    echo    - 后端API:  http://localhost:8080
    echo    - Epusdt管理: http://localhost:8000
    echo.
    echo 容器状态:
    docker-compose ps
    echo.
    
    set /p open_browser="是否打开浏览器访问系统? (y/n): "
    if /i "%open_browser%"=="y" (
        start http://localhost:3001
        timeout /t 2 /nobreak >nul
        start http://localhost:8000
    )
) else (
    echo.
    echo ========================================
    echo    启动失败！
    echo ========================================
    echo.
    echo 请检查错误信息并重试
    echo 查看详细日志: docker-compose logs
    echo.
)

echo.
echo 按任意键退出...
pause
