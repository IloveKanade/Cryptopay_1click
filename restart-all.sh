#!/bin/bash

echo "========================================"
echo "   支付系统一键重启脚本 (Linux/Mac)"
echo "========================================"
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "[错误] 未找到Docker，请先安装Docker"
    exit 1
fi

# 检查Docker是否运行
if ! docker info &> /dev/null; then
    echo "[错误] Docker未运行，请启动Docker服务"
    exit 1
fi

echo "[信息] 正在重启所有服务..."
echo ""

# 切换到项目根目录
cd "$(dirname "$0")"

# 停止所有服务
echo "[步骤1/3] 停止所有服务..."
./stop-all.sh

if [ $? -ne 0 ]; then
    echo "[警告] 停止服务时出现错误，继续执行..."
fi

echo ""
echo "[步骤2/3] 等待服务完全停止..."
sleep 5

# 启动所有服务
echo "[步骤3/3] 重新启动所有服务..."
./start-all.sh

echo ""
echo "[完成] 所有服务重启完成！"
