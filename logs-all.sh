#!/bin/bash

echo "========================================"
echo "   支付系统日志查看脚本 (Linux/Mac)"
echo "========================================"
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "[错误] 未找到Docker"
    exit 1
fi

echo "[信息] 正在显示所有服务日志..."
echo "按 Ctrl+C 退出日志查看"
echo ""

# 切换到项目根目录
cd "$(dirname "$0")"

# 显示所有服务的日志
docker-compose logs -f --tail=100
