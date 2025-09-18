# 🔌 API接口文档

## 📋 基础信息

- **基础URL**: `http://localhost:8080`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

## 🔐 认证接口

### 用户登录
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**响应示例**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "用户名",
    "role": "user",
    "created_at": "2025-08-25T10:00:00Z"
  }
}
```

### 用户注册
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "新用户"
}
```

## 💳 支付链接接口

### 获取支付链接列表
```http
GET /api/payment-links
Authorization: Bearer <token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10)

**响应示例**:
```json
{
  "data": [
    {
      "id": 1,
      "title": "商品支付",
      "amount": 100.00,
      "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1",
      "callback_url": "https://example.com/callback",
      "status": "active",
      "created_at": "2025-08-25T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10
}
```

### 创建支付链接
```http
POST /api/payment-links
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "商品支付",
  "amount": 100.00,
  "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1",
  "callback_url": "https://example.com/callback"
}
```

### 获取支付链接详情
```http
GET /api/payment-links/{id}
Authorization: Bearer <token>
```

### 删除支付链接
```http
DELETE /api/payment-links/{id}
Authorization: Bearer <token>
```

## 💰 支付订单接口

### 创建支付订单
```http
POST /api/payments
Authorization: Bearer <token>
Content-Type: application/json

{
  "payment_link_id": 1,
  "amount": 100.00
}
```

**响应示例**:
```json
{
  "data": {
    "id": 1,
    "order_id": "ORD202508251000001",
    "trade_id": "TRADE202508251000001",
    "amount": 100.00,
    "status": "pending",
    "payment_url": "http://localhost:8000/pay/TRADE202508251000001",
    "created_at": "2025-08-25T10:00:00Z"
  }
}
```

### 获取订单详情
```http
GET /api/payments/{id}
Authorization: Bearer <token>
```

### 查询支付状态
```http
GET /api/payments/{id}/status
Authorization: Bearer <token>
```

## 🏦 管理员接口

### 中间人钱包管理

#### 获取钱包列表
```http
GET /api/admin/middleman-wallets
Authorization: Bearer <admin_token>
```

#### 添加钱包
```http
POST /api/admin/middleman-wallets
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "钱包1",
  "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1",
  "private_key": "your_private_key_here"
}
```

#### 更新钱包状态
```http
PUT /api/admin/middleman-wallets/{id}/status
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "status": "active"
}
```

#### 删除钱包
```http
DELETE /api/admin/middleman-wallets/{id}
Authorization: Bearer <admin_token>
```

### 转账记录管理

#### 获取转账记录
```http
GET /api/admin/transfer-records
Authorization: Bearer <admin_token>
```

**查询参数**:
- `page`: 页码 (默认: 1)
- `page_size`: 每页数量 (默认: 10)
- `status`: 状态筛选 (pending/success/failed)

**响应示例**:
```json
{
  "data": [
    {
      "id": 1,
      "payment_order_id": 1,
      "payment_order": {
        "order_id": "ORD202508251000001",
        "trade_id": "TRADE202508251000001"
      },
      "middleman_wallet_id": 1,
      "middleman_wallet": {
        "name": "钱包1",
        "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1"
      },
      "from_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1",
      "to_address": "TUserWalletAddress",
      "amount": 100.00,
      "fixed_fee_amount": 1.00,
      "network_fee_amount": 0.005,
      "total_fee_amount": 1.005,
      "net_amount": 98.995,
      "status": "success",
      "transaction_hash": "0x1234567890abcdef...",
      "created_at": "2025-08-25T10:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "page_size": 10
}
```

#### 获取单个转账记录
```http
GET /api/admin/transfer-records/{id}
Authorization: Bearer <admin_token>
```

#### 手动处理转账
```http
POST /api/admin/transfer-records/{id}/process
Authorization: Bearer <admin_token>
```

### 费率配置管理

#### 获取费率配置
```http
GET /api/admin/fee-config
Authorization: Bearer <admin_token>
```

**响应示例**:
```json
{
  "data": {
    "id": 1,
    "fee_rate": 0.01,
    "min_fee": 0.1,
    "max_fee": 10.0,
    "network_fee_min": 0.001,
    "network_fee_max": 0.01,
    "created_at": "2025-08-25T10:00:00Z",
    "updated_at": "2025-08-25T10:00:00Z"
  }
}
```

#### 更新费率配置
```http
PUT /api/admin/fee-config
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "fee_rate": 0.01,
  "min_fee": 0.1,
  "max_fee": 10.0,
  "network_fee_min": 0.001,
  "network_fee_max": 0.01
}
```

### 系统设置管理

#### Epusdt配置

##### 获取Epusdt配置
```http
GET /api/admin/epusdt/config
Authorization: Bearer <admin_token>
```

##### 更新Epusdt配置
```http
PUT /api/admin/epusdt/config
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "base_url": "http://localhost:8000"
}
```

##### 测试Epusdt连接
```http
POST /api/admin/epusdt/test-connection
Authorization: Bearer <admin_token>
```

#### 回调配置

##### 获取回调配置
```http
GET /api/admin/callback/config
Authorization: Bearer <admin_token>
```

##### 更新回调配置
```http
PUT /api/admin/callback/config
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "transfer_success_url": "https://example.com/callback/success",
  "transfer_fail_url": "https://example.com/callback/fail",
  "notify_secret": "your_secret_key"
}
```

##### 测试回调
```http
POST /api/admin/callback/test
Authorization: Bearer <admin_token>
```

## 📞 回调接口

### Epusdt支付回调
```http
POST /api/payments/notify
Content-Type: application/x-www-form-urlencoded

trade_id=TRADE202508251000001&amount=100.00&actual_amount=100.00&status=1&signature=md5_signature
```

**回调参数说明**:
- `trade_id`: 交易ID
- `amount`: 订单金额
- `actual_amount`: 实际支付金额
- `status`: 支付状态 (1=成功, 0=失败)
- `signature`: MD5签名

**响应示例**:
```json
{
  "code": 200,
  "message": "success"
}
```

## 🔧 工具接口

### 健康检查
```http
GET /health
```

**响应示例**:
```json
{
  "status": "ok",
  "time": "2025-08-25T10:00:00Z"
}
```

### 计算手续费
```http
GET /api/admin/fee-config/calculate?amount=100.00
Authorization: Bearer <admin_token>
```

**响应示例**:
```json
{
  "amount": 100.00,
  "fixed_fee": 1.00,
  "network_fees": [0.002, 0.003],
  "total_fee": 1.005,
  "net_amount": 98.995
}
```

## 📊 状态码说明

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 201 | 创建成功 |
| 400 | 请求参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## 🔒 错误响应格式

```json
{
  "error": "错误描述信息"
}
```

## 📝 注意事项

1. **认证**: 除健康检查和支付回调外，所有接口都需要JWT认证
2. **权限**: 管理员接口需要管理员权限
3. **签名**: 支付回调需要验证MD5签名
4. **回调**: 转账回调使用HMAC-SHA256签名
5. **分页**: 列表接口支持分页查询
6. **状态**: 转账状态包括 pending/success/failed

## 🧪 测试示例

### 使用curl测试登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'
```

### 使用curl测试创建支付链接
```bash
curl -X POST http://localhost:8080/api/payment-links \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试支付","amount":100.00,"wallet_address":"TQn9Y2khDD95J42FQtQTdwVVRj2qE8qXK1"}'
```

---

*API接口文档 v1.0 - 2025-08-25*
