# 中间人钱包支付系统

## 功能概述

本系统实现了首付款流程使用随机收付款钱包地址，并通过该地址向用户登记的钱包地址转账的功能。管理员可以配置中间人扣费费率。

## 核心功能

### 1. 中间人钱包管理
- **随机钱包选择**: 系统从多个中间人钱包中随机选择一个用于收款
- **钱包地址管理**: 管理员可以添加、删除、启用/禁用中间人钱包
- **余额跟踪**: 自动跟踪每个中间人钱包的手续费余额

### 2. 自动转账流程
- **支付成功触发**: 当用户支付成功后，系统自动触发转账流程
- **手续费扣除**: 根据配置的费率自动计算并扣除手续费
- **转账执行**: 从中间人钱包向用户登记的钱包地址转账

### 3. 费率配置
- **灵活费率**: 支持百分比费率配置（0-100%）
- **最小/最大手续费**: 设置手续费的最小值和最大值
- **实时计算**: 提供实时手续费计算功能

## 数据库表结构

### MiddlemanWallet (中间人钱包表)
```sql
CREATE TABLE middleman_wallets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    wallet_address VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    balance DECIMAL(20,8) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### TransferRecord (转账记录表)
```sql
CREATE TABLE transfer_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    payment_order_id BIGINT NOT NULL,
    middleman_wallet_id BIGINT NOT NULL,
    from_address VARCHAR(255) NOT NULL,
    to_address VARCHAR(255) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    fee_amount DECIMAL(20,8) DEFAULT 0,
    net_amount DECIMAL(20,8) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    transaction_hash VARCHAR(255),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

### FeeConfig (费率配置表)
```sql
CREATE TABLE fee_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    fee_rate DECIMAL(10,8) NOT NULL DEFAULT 0.01,
    min_fee DECIMAL(20,8) NOT NULL DEFAULT 0.1,
    max_fee DECIMAL(20,8) NOT NULL DEFAULT 10,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## API接口

### 中间人钱包管理

#### 获取中间人钱包列表
```
GET /api/admin/middleman-wallets
```

#### 添加中间人钱包
```
POST /api/admin/middleman-wallets
Content-Type: application/json

{
    "wallet_address": "TQn9Y2khDD95J42FQtQTdwVVRj2qE8qKqK"
}
```

#### 更新中间人钱包状态
```
PUT /api/admin/middleman-wallets/:id/status
Content-Type: application/json

{
    "status": "inactive"
}
```

#### 删除中间人钱包
```
DELETE /api/admin/middleman-wallets/:id
```

### 转账记录管理

#### 获取转账记录列表
```
GET /api/admin/transfer-records?page=1&page_size=10
```

#### 获取单个转账记录
```
GET /api/admin/transfer-records/:id
```

#### 手动处理转账
```
POST /api/admin/transfer-records/:id/process
```

### 费率配置管理

#### 获取当前费率配置
```
GET /api/admin/fee-config
```

#### 更新费率配置
```
PUT /api/admin/fee-config
Content-Type: application/json

{
    "fee_rate": 0.02,
    "min_fee": 0.5,
    "max_fee": 20.0
}
```

#### 创建新的费率配置
```
POST /api/admin/fee-config
Content-Type: application/json

{
    "fee_rate": 0.015,
    "min_fee": 0.3,
    "max_fee": 15.0
}
```

#### 获取费率配置历史
```
GET /api/admin/fee-config/history?page=1&page_size=10
```

#### 计算手续费
```
GET /api/admin/fee-config/calculate?amount=100
```

## 支付流程

### 1. 创建支付订单
1. 用户创建支付链接，指定收款钱包地址
2. 系统随机选择一个中间人钱包地址
3. 使用中间人钱包地址创建Epusdt支付订单

### 2. 用户支付
1. 用户向中间人钱包地址支付USDT
2. Epusdt系统检测到支付并回调通知

### 3. 自动转账
1. 系统收到支付成功回调
2. 自动创建转账记录
3. 计算手续费（根据配置的费率）
4. 从中间人钱包向用户登记的钱包地址转账
5. 更新转账记录状态和中间人钱包余额

## 配置说明

### 默认费率配置
- **费率**: 1% (0.01)
- **最小手续费**: 0.1 USDT
- **最大手续费**: 10 USDT

### 费率计算示例
- 支付100 USDT: 手续费 = 1 USDT (1%)
- 支付5 USDT: 手续费 = 0.1 USDT (最小手续费)
- 支付2000 USDT: 手续费 = 10 USDT (最大手续费)

## 安全考虑

1. **钱包地址验证**: 所有钱包地址都经过TRC20格式验证
2. **权限控制**: 只有管理员可以管理中间人钱包和费率配置
3. **转账记录**: 所有转账操作都有完整的审计记录
4. **错误处理**: 转账失败时记录详细错误信息

## 部署说明

1. 确保数据库已正确配置
2. 运行数据库迁移: `go run main.go`
3. 通过管理员接口添加中间人钱包地址
4. 配置合适的费率参数
5. 测试支付和转账流程

## 注意事项

1. 中间人钱包需要有足够的USDT余额用于转账
2. 建议定期检查转账记录和钱包余额
3. 费率调整会影响所有新创建的转账记录
4. 转账失败时需要手动处理或重试
