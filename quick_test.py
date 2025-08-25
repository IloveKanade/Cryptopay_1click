#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
快速测试脚本 - 中间人钱包系统
简化版本，用于快速验证系统功能
"""

import requests
import json
from datetime import datetime

def test_api_endpoint(url, method="GET", data=None, headers=None):
    """测试API端点"""
    try:
        if method.upper() == "GET":
            response = requests.get(url, headers=headers)
        elif method.upper() == "POST":
            response = requests.post(url, json=data, headers=headers)
        elif method.upper() == "PUT":
            response = requests.put(url, json=data, headers=headers)
        elif method.upper() == "DELETE":
            response = requests.delete(url, headers=headers)
        
        return response.status_code, response.json() if response.text else {}
    except Exception as e:
        return 0, {"error": str(e)}

def print_test_result(test_name, success, message=""):
    """打印测试结果"""
    status = "✅ 通过" if success else "❌ 失败"
    print(f"{status} {test_name}")
    if message:
        print(f"    {message}")

def main():
    base_url = "http://localhost:8080"
    print("🚀 快速测试中间人钱包系统")
    print(f"📡 测试服务器: {base_url}")
    print(f"⏰ 开始时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    
    # 1. 测试服务器连接
    print("\n" + "="*50)
    print("测试服务器连接")
    print("="*50)
    
    status, data = test_api_endpoint(f"{base_url}/api/health")
    if status == 200:
        print_test_result("服务器连接", True)
    else:
        print_test_result("服务器连接", False, f"状态码: {status}")
        print("❌ 服务器未启动，请先启动系统")
        return
    
    # 2. 管理员登录
    print("\n" + "="*50)
    print("管理员登录测试")
    print("="*50)
    
    login_data = {
        "email": "admin@example.com",
        "password": "admin123"
    }
    
    status, data = test_api_endpoint(f"{base_url}/api/auth/login", "POST", login_data)
    if status == 200:
        admin_token = data.get("token")
        headers = {"Authorization": f"Bearer {admin_token}"}
        print_test_result("管理员登录", True, "获取到访问令牌")
    else:
        print_test_result("管理员登录", False, f"状态码: {status}")
        print("❌ 管理员登录失败，请检查默认账户")
        return
    
    # 3. 测试中间人钱包管理
    print("\n" + "="*50)
    print("中间人钱包管理测试")
    print("="*50)
    
    # 获取钱包列表
    status, data = test_api_endpoint(f"{base_url}/api/admin/middleman-wallets", headers=headers)
    if status == 200:
        wallets = data.get("data", [])
        print_test_result("获取钱包列表", True, f"当前有 {len(wallets)} 个钱包")
    else:
        print_test_result("获取钱包列表", False, f"状态码: {status}")
    
    # 添加测试钱包
    test_wallet = {
        "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1"
    }
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/middleman-wallets", "POST", test_wallet, headers)
    if status == 201:
        print_test_result("添加测试钱包", True, "钱包添加成功")
    elif status == 400 and "已存在" in str(data):
        print_test_result("添加测试钱包", True, "钱包已存在")
    else:
        print_test_result("添加测试钱包", False, f"状态码: {status}")
    
    # 4. 测试费率配置
    print("\n" + "="*50)
    print("费率配置测试")
    print("="*50)
    
    # 获取费率配置
    status, data = test_api_endpoint(f"{base_url}/api/admin/fee-config", headers=headers)
    if status == 200:
        fee_config = data.get("data")
        if fee_config:
            print_test_result("获取费率配置", True, 
                            f"费率: {fee_config.get('fee_rate', 0)}%, "
                            f"最小: {fee_config.get('min_fee', 0)} USDT, "
                            f"最大: {fee_config.get('max_fee', 0)} USDT, "
                            f"网络手续费: {fee_config.get('network_fee_min', 0)}-{fee_config.get('network_fee_max', 0)} USDT, "
                            f"转账次数: {fee_config.get('network_fee_count', 0)}次")
        else:
            print_test_result("获取费率配置", True, "无配置")
    else:
        print_test_result("获取费率配置", False, f"状态码: {status}")
    
    # 创建费率配置
    fee_data = {
        "fee_rate": 1.0,
        "min_fee": 0.1,
        "max_fee": 10.0,
        "network_fee_min": 0.001,
        "network_fee_max": 0.01,
        "network_fee_count": 2
    }
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/fee-config", "POST", fee_data, headers)
    if status == 201:
        print_test_result("创建费率配置", True, "1%费率配置成功")
    elif status == 400 and "已存在" in str(data):
        print_test_result("创建费率配置", True, "配置已存在")
    else:
        print_test_result("创建费率配置", False, f"状态码: {status}")
    
    # 5. 测试用户管理
    print("\n" + "="*50)
    print("用户管理测试")
    print("="*50)
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/users", headers=headers)
    if status == 200:
        users = data.get("data", [])
        print_test_result("获取用户列表", True, f"共有 {len(users)} 个用户")
        
        # 显示用户信息
        for user in users[:3]:
            print(f"    👤 {user.get('name', 'N/A')} ({user.get('email', 'N/A')}) - {user.get('role', 'N/A')}")
    else:
        print_test_result("获取用户列表", False, f"状态码: {status}")
    
    # 6. 测试订单管理
    print("\n" + "="*50)
    print("订单管理测试")
    print("="*50)
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/orders", headers=headers)
    if status == 200:
        orders = data.get("data", [])
        print_test_result("获取订单列表", True, f"共有 {len(orders)} 个订单")
        
        # 显示订单信息
        for order in orders[:3]:
            print(f"    📦 {order.get('order_id', 'N/A')} - {order.get('amount', 0)} USDT - {order.get('status', 'N/A')}")
    else:
        print_test_result("获取订单列表", False, f"状态码: {status}")
    
    # 7. 测试转账记录
    print("\n" + "="*50)
    print("转账记录测试")
    print("="*50)
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/transfer-records", headers=headers)
    if status == 200:
        records = data.get("data", [])
        print_test_result("获取转账记录", True, f"共有 {len(records)} 条记录")
        
        # 显示转账记录
        for record in records[:3]:
            print(f"    💸 {record.get('order_id', 'N/A')} - {record.get('amount', 0)} USDT - {record.get('status', 'N/A')}")
    else:
        print_test_result("获取转账记录", False, f"状态码: {status}")
    
    # 8. 测试仪表板统计
    print("\n" + "="*50)
    print("仪表板统计测试")
    print("="*50)
    
    status, data = test_api_endpoint(f"{base_url}/api/admin/dashboard/stats", headers=headers)
    if status == 200:
        stats = data.get("data", {})
        print_test_result("获取仪表板统计", True, 
                        f"用户数: {stats.get('total_users', 0)}, "
                        f"订单数: {stats.get('total_orders', 0)}, "
                        f"收入: {stats.get('total_revenue', 0)} USDT")
    else:
        print_test_result("获取仪表板统计", False, f"状态码: {status}")
    
    # 9. 测试用户注册和登录
    print("\n" + "="*50)
    print("用户认证测试")
    print("="*50)
    
    # 用户注册
    user_data = {
        "name": "快速测试用户",
        "email": "quicktest@example.com",
        "password": "test123456"
    }
    
    status, data = test_api_endpoint(f"{base_url}/api/auth/register", "POST", user_data)
    if status == 201:
        print_test_result("用户注册", True, "注册成功")
    elif status == 400 and "已存在" in str(data):
        print_test_result("用户注册", True, "用户已存在")
    else:
        print_test_result("用户注册", False, f"状态码: {status}")
    
    # 用户登录
    login_data = {
        "email": "quicktest@example.com",
        "password": "test123456"
    }
    
    status, data = test_api_endpoint(f"{base_url}/api/auth/login", "POST", login_data)
    if status == 200:
        user_token = data.get("token")
        user_headers = {"Authorization": f"Bearer {user_token}"}
        print_test_result("用户登录", True, "登录成功")
        
        # 测试用户功能
        status, data = test_api_endpoint(f"{base_url}/api/payment-links", headers=user_headers)
        if status == 200:
            links = data.get("data", [])
            print_test_result("获取用户支付链接", True, f"共有 {len(links)} 个链接")
        else:
            print_test_result("获取用户支付链接", False, f"状态码: {status}")
    else:
        print_test_result("用户登录", False, f"状态码: {status}")
    
    print(f"\n{'='*50}")
    print(f"🎉 快速测试完成")
    print(f"⏰ 结束时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"{'='*50}")
    print("\n📋 测试总结:")
    print("✅ 如果大部分测试通过，说明系统运行正常")
    print("❌ 如果测试失败，请检查:")
    print("   1. 系统是否已启动 (start.bat)")
    print("   2. 数据库是否正常")
    print("   3. API地址是否正确")
    print("   4. 网络连接是否正常")

if __name__ == "__main__":
    main()
