#!/bin/bash

echo "========================================"
echo "   支付系统一键停止脚本 (Linux/Mac)"
echo "========================================"
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "[错误] 未找到Docker"
    exit 1
fi

echo "[信息] 正在停止所有服务..."
echo ""

# 切换到项目根目录
cd "$(dirname "$0")"

# 停止Payment Link MVP系统
echo "[步骤1/2] 停止Payment Link MVP系统..."
cd payment-link-mvp
docker-compose down

if [ $? -ne 0 ]; then
    echo "[警告] Payment Link MVP系统停止时出现错误"
fi

echo "[成功] Payment Link MVP系统已停止"
echo ""

# 返回根目录
cd ..

# 停止Epusdt系统
echo "[步骤2/2] 停止Epusdt支付系统..."
docker-compose down

if [ $? -ne 0 ]; then
    echo "[警告] Epusdt系统停止时出现错误"
fi

echo "[成功] Epusdt系统已停止"
echo ""

echo "========================================"
echo "          所有服务已停止！"
echo "========================================"
echo ""
echo "[清理命令]"
echo "清理所有容器: docker system prune -f"
echo "清理所有镜像: docker system prune -a -f"
echo "清理所有数据: docker volume prune -f"
echo ""
