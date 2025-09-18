# 中间人钱包支付系统使用示例

## 1. 管理员配置中间人钱包

### 添加中间人钱包
```bash
curl -X POST http://localhost:8080/api/admin/middleman-wallets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqK"
  }'
```

### 查看中间人钱包列表
```bash
curl -X GET http://localhost:8080/api/admin/middleman-wallets \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 2. 配置费率

### 设置费率配置
```bash
curl -X PUT http://localhost:8080/api/admin/fee-config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "fee_rate": 0.02,
    "min_fee": 0.5,
    "max_fee": 20.0
  }'
```

### 查看当前费率配置
```bash
curl -X GET http://localhost:8080/api/admin/fee-config \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 计算手续费
```bash
curl -X GET "http://localhost:8080/api/admin/fee-config/calculate?amount=100" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 3. 用户创建支付链接

### 创建支付链接
```bash
curl -X POST http://localhost:8080/api/payment-links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_USER_TOKEN" \
  -d '{
    "title": "测试支付",
    "amount": 100.0,
    "currency": "USDT",
    "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqZ",
    "description": "这是一个测试支付链接"
  }'
```

## 4. 创建支付订单

### 创建支付订单
```bash
curl -X POST http://localhost:8080/api/payment/create \
  -H "Content-Type: application/json" \
  -d '{
    "payment_link_id": 1,
    "amount": 100.0,
    "client_ip": "127.0.0.1"
  }'
```

响应示例：
```json
{
  "order_id": "ORDER_1_1703123456",
  "trade_id": "EPUSDT_TRADE_123456",
  "amount": 100.0,
  "actual_amount": 100.0,
  "payment_url": "https://epusdt.example.com/pay/xxx",
  "expires_at": "2023-12-21T10:30:56Z",
  "status": "pending"
}
```

## 5. 支付流程

1. 用户访问支付URL进行支付
2. 用户向中间人钱包地址支付USDT
3. Epusdt系统检测到支付并回调通知
4. 系统自动创建转账记录
5. 系统从中间人钱包向用户钱包转账（扣除手续费）

## 6. 查看转账记录

### 获取转账记录列表
```bash
curl -X GET "http://localhost:8080/api/admin/transfer-records?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 获取单个转账记录
```bash
curl -X GET http://localhost:8080/api/admin/transfer-records/1 \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 7. 手动处理转账

### 手动处理转账（如果需要）
```bash
curl -X POST http://localhost:8080/api/admin/transfer-records/1/process \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

## 完整流程示例

### 步骤1: 管理员添加中间人钱包
```bash
# 添加多个中间人钱包
curl -X POST http://localhost:8080/api/admin/middleman-wallets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{"wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqK"}'

curl -X POST http://localhost:8080/api/admin/middleman-wallets \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{"wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqL"}'
```

### 步骤2: 配置费率
```bash
curl -X PUT http://localhost:8080/api/admin/fee-config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{
    "fee_rate": 0.015,
    "min_fee": 0.3,
    "max_fee": 15.0
  }'
```

### 步骤3: 用户创建支付链接
```bash
curl -X POST http://localhost:8080/api/payment-links \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer USER_TOKEN" \
  -d '{
    "title": "商品购买",
    "amount": 500.0,
    "currency": "USDT",
    "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqZ",
    "description": "购买商品A"
  }'
```

### 步骤4: 创建支付订单
```bash
curl -X POST http://localhost:8080/api/payment/create \
  -H "Content-Type: application/json" \
  -d '{
    "payment_link_id": 1,
    "amount": 500.0,
    "client_ip": "127.0.0.1"
  }'
```

### 步骤5: 查看转账记录
```bash
curl -X GET http://localhost:8080/api/admin/transfer-records \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

## 注意事项

1. 确保中间人钱包有足够的USDT余额
2. 费率配置会影响所有新创建的转账
3. 转账失败时可以手动重试
4. 所有操作都需要相应的权限
5. 钱包地址必须是有效的TRC20地址格式
