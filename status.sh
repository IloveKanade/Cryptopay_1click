#!/bin/bash

echo "========================================"
echo "   支付系统状态检查脚本 (Linux/Mac)"
echo "========================================"
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "[错误] 未找到Docker"
    exit 1
fi

echo "[信息] 正在检查服务状态..."
echo ""

# 切换到项目根目录
cd "$(dirname "$0")"

echo "[Epusdt系统状态]"
docker-compose ps mysql redis epusdt
echo ""

echo "[Payment Link MVP系统状态]"
cd payment-link-mvp
docker-compose ps
echo ""

# 返回根目录
cd ..

echo "[网络状态]"
docker network ls | grep epusdt
docker network ls | grep payment
echo ""

echo "[端口占用情况]"
netstat -tuln | grep ":3001" || echo "端口3001未占用"
netstat -tuln | grep ":8000" || echo "端口8000未占用"
netstat -tuln | grep ":8080" || echo "端口8080未占用"
netstat -tuln | grep ":3306" || echo "端口3306未占用"
netstat -tuln | grep ":6379" || echo "端口6379未占用"
echo ""

echo "[服务健康检查]"
echo "检查Epusdt管理后台..."
curl -s -o /dev/null -w "HTTP状态码: %{http_code}\n" http://localhost:8000 || echo "无法连接到Epusdt管理后台"

echo "检查Payment Link前端..."
curl -s -o /dev/null -w "HTTP状态码: %{http_code}\n" http://localhost:3001 || echo "无法连接到Payment Link前端"

echo "检查Payment Link后端API..."
curl -s -o /dev/null -w "HTTP状态码: %{http_code}\n" http://localhost:8080/api/health || echo "无法连接到Payment Link后端API"
echo ""

echo "========================================"
echo "          状态检查完成！"
echo "========================================"
echo ""
