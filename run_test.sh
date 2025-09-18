#!/bin/bash

# 中间人钱包系统测试工具

echo "========================================"
echo "    中间人钱包系统测试工具"
echo "========================================"
echo

# 检查Python环境
echo "正在检查Python环境..."
if ! command -v python3 &> /dev/null; then
    echo "❌ 未找到Python3，请先安装Python 3.7+"
    exit 1
fi

echo "✅ Python环境检查通过"
echo

# 检查requests库
echo "正在检查requests库..."
if ! python3 -c "import requests" &> /dev/null; then
    echo "⚠️  未找到requests库，正在安装..."
    pip3 install requests
    if [ $? -ne 0 ]; then
        echo "❌ 安装requests库失败"
        exit 1
    fi
fi

echo "✅ requests库检查通过"
echo

# 显示菜单
echo "请选择测试类型："
echo "1. 快速测试 (推荐)"
echo "2. 完整测试"
echo "3. 钱包管理测试"
echo "4. 费率配置测试"
echo "5. 支付流程测试"
echo "6. 用户管理测试"
echo "7. 订单管理测试"
echo "8. 转账记录测试"
echo "9. 仪表板统计测试"
echo "0. 退出"
echo

read -p "请输入选择 (0-9): " choice

case $choice in
    1)
        echo
        echo "🚀 开始快速测试..."
        python3 quick_test.py
        ;;
    2)
        echo
        echo "🚀 开始完整测试..."
        python3 test_middleman_system.py
        ;;
    3)
        echo
        echo "🚀 开始钱包管理测试..."
        python3 test_middleman_system.py --test wallet
        ;;
    4)
        echo
        echo "🚀 开始费率配置测试..."
        python3 test_middleman_system.py --test fee
        ;;
    5)
        echo
        echo "🚀 开始支付流程测试..."
        python3 test_middleman_system.py --test payment
        ;;
    6)
        echo
        echo "🚀 开始用户管理测试..."
        python3 test_middleman_system.py --test user
        ;;
    7)
        echo
        echo "🚀 开始订单管理测试..."
        python3 test_middleman_system.py --test order
        ;;
    8)
        echo
        echo "🚀 开始转账记录测试..."
        python3 test_middleman_system.py --test transfer
        ;;
    9)
        echo
        echo "🚀 开始仪表板统计测试..."
        python3 test_middleman_system.py --test stats
        ;;
    0)
        echo "退出测试工具"
        exit 0
        ;;
    *)
        echo "❌ 无效选择，请重新运行"
        exit 1
        ;;
esac

echo
echo "========================================"
echo "    测试完成"
echo "========================================"
echo
echo "📋 测试结果说明："
echo "✅ 绿色表示测试通过"
echo "❌ 红色表示测试失败"
echo
echo "🔧 如果测试失败，请检查："
echo "   1. 系统是否已启动 (./start-all.sh)"
echo "   2. 数据库是否正常"
echo "   3. API地址是否正确"
echo "   4. 网络连接是否正常"
echo
