@echo off
chcp 65001 >nul
title 中间人钱包系统测试

echo ========================================
echo    中间人钱包系统测试工具
echo ========================================
echo.

echo 正在检查Python环境...
python --version >nul 2>&1
if errorlevel 1 (
    echo ❌ 未找到Python，请先安装Python 3.7+
    pause
    exit /b 1
)

echo ✅ Python环境检查通过
echo.

echo 正在检查requests库...
python -c "import requests" >nul 2>&1
if errorlevel 1 (
    echo ⚠️  未找到requests库，正在安装...
    pip install requests
    if errorlevel 1 (
        echo ❌ 安装requests库失败
        pause
        exit /b 1
    )
)

echo ✅ requests库检查通过
echo.

echo 请选择测试类型：
echo 1. 快速测试 (推荐)
echo 2. 完整测试
echo 3. 钱包管理测试
echo 4. 费率配置测试
echo 5. 支付流程测试
echo 6. 用户管理测试
echo 7. 订单管理测试
echo 8. 转账记录测试
echo 9. 仪表板统计测试
echo 0. 退出
echo.

set /p choice="请输入选择 (0-9): "

if "%choice%"=="1" (
    echo.
    echo 🚀 开始快速测试...
    python quick_test.py
) else if "%choice%"=="2" (
    echo.
    echo 🚀 开始完整测试...
    python test_middleman_system.py
) else if "%choice%"=="3" (
    echo.
    echo 🚀 开始钱包管理测试...
    python test_middleman_system.py --test wallet
) else if "%choice%"=="4" (
    echo.
    echo 🚀 开始费率配置测试...
    python test_middleman_system.py --test fee
) else if "%choice%"=="5" (
    echo.
    echo 🚀 开始支付流程测试...
    python test_middleman_system.py --test payment
) else if "%choice%"=="6" (
    echo.
    echo 🚀 开始用户管理测试...
    python test_middleman_system.py --test user
) else if "%choice%"=="7" (
    echo.
    echo 🚀 开始订单管理测试...
    python test_middleman_system.py --test order
) else if "%choice%"=="8" (
    echo.
    echo 🚀 开始转账记录测试...
    python test_middleman_system.py --test transfer
) else if "%choice%"=="9" (
    echo.
    echo 🚀 开始仪表板统计测试...
    python test_middleman_system.py --test stats
) else if "%choice%"=="0" (
    echo 退出测试工具
    exit /b 0
) else (
    echo ❌ 无效选择，请重新运行
    pause
    exit /b 1
)

echo.
echo ========================================
echo    测试完成
echo ========================================
echo.
echo 📋 测试结果说明：
echo ✅ 绿色表示测试通过
echo ❌ 红色表示测试失败
echo.
echo 🔧 如果测试失败，请检查：
echo    1. 系统是否已启动 (start.bat)
echo    2. 数据库是否正常
echo    3. API地址是否正确
echo    4. 网络连接是否正常
echo.
pause
