# 数字货币支付系统集成项目

## 项目背景

本项目是一个完整的数字货币支付解决方案，专为需要集成USDT支付功能的开发者和企业设计。项目整合了多个核心组件，提供从支付网关到用户界面的完整支付生态系统。

### 核心价值
- **私有化部署**: 完全自主控制，无需依赖第三方支付平台
- **零手续费**: 数字货币直接进入您的钱包，无中间商抽成
- **高安全性**: 开源透明，支持私有化部署，避免资金风险
- **易于集成**: 提供完整的API接口和SDK，快速集成到现有系统
- **多场景支持**: 适用于电商、数字商品销售、服务付费等多种业务场景

## 📚 相关文档

- **[系统说明文档](./系统说明文档.md)** - 详细的技术实现文档
  - [技术架构详解](./系统说明文档.md#技术架构详解)
  - [API接口文档](./系统说明文档.md#api接口文档)
  - [部署配置详解](./系统说明文档.md#部署配置详解)
  - [测试与质量保证](./系统说明文档.md#测试与质量保证)
  - [故障排除与运维](./系统说明文档.md#故障排除与运维)

## 快速部署指南

### 环境要求
- **操作系统**: Linux/Windows (支持x86/ARM架构)
- **数据库**: MySQL 5.7+ 或 MariaDB 10.3+
- **缓存**: Redis 5.0+ (可选，用于Payment Link MVP)
- **Go版本**: 1.19+ (仅开发环境需要)

### 一键部署 (推荐)

#### 使用宝塔面板部署
1. 安装宝塔面板并配置MySQL
2. 下载Epusdt二进制文件
3. 按照 [宝塔部署教程](epusdt/wiki/BT_RUN.md) 进行配置
4. 启动服务并测试API连通性

#### 使用Docker部署
```bash
# 克隆项目
git clone [项目地址]
cd Test_Tasks

# 启动Payment Link MVP
cd payment-link-mvp
docker-compose up -d

# 启动Epusdt (需要先配置数据库)
cd ../epusdt
./epusdt
```

### 配置说明

#### Epusdt配置
```yaml
# config.yml 主要配置项
database:
  host: localhost
  port: 3306
  username: epusdt
  password: epusdt123456
  database: epusdt

tron:
  api_key: "your_tron_api_key"
  
telegram:
  bot_token: "your_bot_token"  # 可选
```

#### Payment Link MVP配置
```bash
# 后端环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export REDIS_URL=redis://localhost:6379
```

### 验证部署

1. **Epusdt API测试**:
   ```bash
   curl http://localhost:8080/api/v1/health
   ```

2. **Payment Link访问**:
   - 前端: http://localhost:3000
   - 后端API: http://localhost:8081

3. **数据库连接测试**:
   - 检查数据表是否正确创建
   - 验证示例数据是否导入成功

> 🔧 **详细配置说明**: 更多配置选项和故障排除，请参考 [系统说明文档.md](系统说明文档.md) 中的部署章节

## 核心组件概览

### 1. Epusdt支付网关

<p align="center">
<img src="epusdt/wiki/img/usdtlogo.png" width="200">
</p>

**Epusdt**（Easy Payment Usdt）是一个Go语言编写的私有化USDT支付中间件，支持TRC20网络，提供完整的数字货币支付解决方案。

**核心特性**:
- 🏠 **私有化部署** - 完全自主控制，资金安全
- 🚀 **高性能处理** - 支持高并发订单和异步队列
- 🔌 **零依赖部署** - 仅需一个二进制文件即可运行
- 📡 **完整API** - RESTful接口，支持任何系统集成
- 📱 **Telegram集成** - 实时支付消息推送

**快速集成**:
- 宝塔部署教程: [epusdt/wiki/BT_RUN.md](epusdt/wiki/BT_RUN.md)
- API文档: [epusdt/wiki/API.md](epusdt/wiki/API.md)
- 插件开发: [epusdt/plugins/](epusdt/plugins/)

### 2. Payment Link MVP系统

**Payment Link MVP** 是基于React + Gin的现代化支付链接管理平台，提供直观的支付链接创建和管理功能。

**主要功能**:
- ✅ 访客落地页和用户系统
- ✅ 支付链接管理和公开支付页面
- ✅ 响应式设计，完美适配移动设备
- ✅ 实时订单跟踪和数据统计

**技术栈**: React 18 + TypeScript + Material-UI + Gin + MySQL + Redis

> 📖 **详细技术文档**: 如需了解深度技术实现、API接口详情、中间人钱包机制等，请参考 [系统说明文档.md](系统说明文档.md)

### 3. 插件生态系统

#### 独角数卡集成插件

**适用版本**: 独角数卡 2.0.4以下版本（2.0.4+已内置支持）

**安装步骤**:
1. 将插件文件复制到独角数卡根目录
2. 在后台添加Epusdt支付方式
3. 配置API认证和回调地址

**配置参数**:
| 参数 | 值 | 说明 |
|------|-----|------|
| 商户ID | API认证Token | 从Epusdt后台获取 |
| 商户Key | 空 | 留空即可 |
| 回调地址 | http://127.0.0.1:8000/api/v1/order/create-transaction | 本地部署地址 |

## 管理和运维

### 服务管理脚本

项目提供了完整的服务管理脚本，支持一键操作：

```bash
# 启动所有服务
./start-all.sh

# 停止所有服务
./stop-all.sh

# 重启所有服务
./restart-all.sh

# 查看服务状态
./status.sh

# 查看服务日志
./logs-all.sh
```

### 监控和日志

#### 服务状态检查
```bash
# 检查Docker容器状态
docker-compose ps

# 检查端口占用
netstat -tulpn | grep -E ':(3001|8000|8080|8081|3306|6379)'

# 检查服务健康状态
curl http://localhost:8080/api/health
```

#### 日志管理
```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs epusdt
docker-compose logs backend
docker-compose logs frontend

# 实时跟踪日志
docker-compose logs -f
```

## 故障排除

### 常见问题解决

#### 1. 端口冲突
```bash
# 检查端口占用
netstat -tulpn | grep :3001
netstat -tulpn | grep :8000

# 停止占用端口的进程
sudo lsof -ti:3001 | xargs kill -9
```

#### 2. Docker服务启动失败
```bash
# 检查Docker状态
docker info
systemctl status docker

# 重启Docker服务
sudo systemctl restart docker

# 清理Docker缓存
docker system prune -a
```

#### 3. 数据库连接问题
```bash
# 检查MySQL容器状态
docker-compose ps mysql

# 查看MySQL日志
docker-compose logs mysql

# 重启MySQL服务
docker-compose restart mysql
```

#### 4. Epusdt子模块问题
```bash
# 初始化Git子模块
git submodule update --init --recursive

# 更新子模块
git submodule update --remote
```

### 性能优化

#### 数据库优化
- 定期清理过期订单数据
- 优化数据库索引
- 配置适当的连接池大小

#### 缓存优化
- 合理设置Redis缓存过期时间
- 使用缓存预热策略
- 监控缓存命中率

## API接口概览

### Epusdt支付API

**基础信息**:
- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: API Key (Header: `X-API-Key`)
- **数据格式**: JSON

**核心接口**:
```bash
# 创建支付订单
POST /order
{
  "amount": 10.5,
  "callback_url": "https://your-site.com/callback"
}

# 查询订单状态
GET /order/{order_id}

# 获取钱包余额
GET /wallet/balance
```

### Payment Link API

**基础信息**:
- **Base URL**: `http://localhost:8081/api`
- **认证方式**: JWT Token
- **数据格式**: JSON

**核心接口**:
```bash
# 用户登录
POST /auth/login

# 创建支付链接
POST /payment-links

# 获取支付链接列表
GET /payment-links
```

> 📚 **完整API文档**: 详细的接口参数、响应格式和错误码说明，请参考 [系统说明文档.md](系统说明文档.md) 中的API章节

## 安全考虑

### 核心安全措施
- 🔐 **私钥安全**: 钱包私钥本地存储，不上传云端
- 🛡️ **API认证**: 所有接口均需API Key验证
- 🔒 **HTTPS传输**: 生产环境强制使用HTTPS
- 📝 **签名验证**: 回调数据使用HMAC-SHA256签名
- ⏰ **订单时效**: 支付订单自动过期机制

### 最佳实践
- 定期更换API密钥
- 监控异常交易和API调用
- 备份钱包和数据库
- 使用防火墙限制访问

> 🔒 **详细安全指南**: 完整的安全配置和风险防范措施，请参考 [系统说明文档.md](系统说明文档.md) 中的安全章节

## 技术支持与社区

### 获取帮助
- 📖 **文档中心**: [系统说明文档.md](系统说明文档.md)
- 🐛 **问题反馈**: GitHub Issues
- 💬 **社区讨论**: 开发者交流群
- 📧 **技术支持**: support@example.com

### 贡献指南
欢迎提交Pull Request和Issue，共同完善项目。

## 版本信息

**当前版本**: v1.0.0  
**更新日期**: 2024-01-01  
**兼容性**: 支持TRC20网络，兼容主流浏览器

## 许可证

本项目采用 MIT 许可证，详情请参阅 [LICENSE](LICENSE) 文件。

## 免责声明

本软件仅供学习和研究使用。使用本软件进行任何商业活动或处理真实资金时，请确保遵守当地法律法规。开发者不承担因使用本软件而产生的任何损失或法律责任。

---

**如果本项目对您有帮助，欢迎给我们一个⭐Star！**
