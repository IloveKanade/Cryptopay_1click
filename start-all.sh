#!/bin/bash

echo "========================================"
echo "   支付系统一键启动脚本 (Linux/Mac)"
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

echo "[信息] 正在启动所有服务..."
echo ""

# 切换到项目根目录
cd "$(dirname "$0")"

# 启动Epusdt系统
echo "[步骤1/2] 启动Epusdt支付系统..."
docker-compose up -d mysql redis epusdt

if [ $? -ne 0 ]; then
    echo "[错误] Epusdt系统启动失败"
    exit 1
fi

echo "[成功] Epusdt系统启动完成"
echo ""

# 等待数据库完全启动
echo "[信息] 等待数据库服务完全启动..."
sleep 10

# 启动Payment Link MVP系统
echo "[步骤2/2] 启动Payment Link MVP系统..."
cd payment-link-mvp
docker-compose up -d

if [ $? -ne 0 ]; then
    echo "[错误] Payment Link MVP系统启动失败"
    exit 1
fi

echo "[成功] Payment Link MVP系统启动完成"
echo ""

# 等待所有服务启动
echo "[信息] 等待所有服务完全启动..."
sleep 15

echo "========================================"
echo "          服务启动成功！"
echo "========================================"
echo ""
echo "[服务地址]"
echo "Epusdt管理后台: http://localhost:8000"
echo "Payment Link前端: http://localhost:3001"
echo "Payment Link后端API: http://localhost:8080"
echo ""
echo "[数据库]"
echo "MySQL: localhost:3306"
echo "Redis: localhost:6379"
echo ""
echo "[默认账户]"
echo "Epusdt管理后台: admin/admin123"
echo "Payment Link: 请先注册账户"
echo ""

# 检查系统类型并打开浏览器
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    open http://localhost:3001
    open http://localhost:8000
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux
    if command -v xdg-open &> /dev/null; then
        xdg-open http://localhost:3001
        xdg-open http://localhost:8000
    elif command -v gnome-open &> /dev/null; then
        gnome-open http://localhost:3001
        gnome-open http://localhost:8000
    else
        echo "请手动打开浏览器访问:"
        echo "  - http://localhost:3001"
        echo "  - http://localhost:8000"
    fi
else
    echo "请手动打开浏览器访问:"
    echo "  - http://localhost:3001"
    echo "  - http://localhost:8000"
fi

echo ""
echo "[管理命令]"
echo "停止所有服务: ./stop-all.sh"
echo "查看服务日志: ./logs-all.sh"
echo "重启所有服务: ./restart-all.sh"
echo ""
