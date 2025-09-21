# 数字货币支付系统集成项目

## 项目背景

本项目是一个完整的数字货币支付解决方案，专为需要集成USDT支付功能的开发者和企业设计。项目整合了多个核心组件，提供从支付网关到用户界面的完整支付生态系统。

### 核心价值
- **私有化部署**: 完全自主控制，无需依赖第三方支付平台
- **零手续费**: 数字货币直接进入您的钱包，无中间商抽成
- **高安全性**: 开源透明，支持私有化部署，避免资金风险
- **易于集成**: 提供完整的API接口和SDK，快速集成到现有系统
- **多场景支持**: 适用于电商、数字商品销售、服务付费等多种业务场景

## 快速启动指南

### 环境准备

#### 系统要求
- **操作系统**: Linux/macOS/Windows
- **Docker**: 20.10+ 
- **Docker Compose**: 1.29+
- **Git**: 2.0+
- **可用端口**: 3001, 8000, 8080, 8081, 3306, 6379

#### 端口说明
| 端口 | 服务 | 说明 |
|------|------|------|
| 3001 | 前端界面 | React应用主界面 |
| 8000 | Epusdt管理 | USDT支付网关管理后台 |
| 8080 | 后端API | RESTful API服务 |
| 8081 | 数据库管理 | phpMyAdmin数据库管理界面 |
| 3306 | MySQL | 数据库服务 |
| 6379 | Redis | 缓存和会话存储 |

### 一键部署（推荐）

```bash
# 1. 克隆项目
git clone https://github.com/WFHTask/Test_Tasks.git
cd Test_Tasks

# 2. 检查端口占用（可选）
netstat -tulpn | grep -E ':(3001|8000|8080|8081|3306|6379)'

# 3. 启动所有服务
docker-compose up -d

# 4. 等待服务启动完成（约2-3分钟）
docker-compose ps

# 5. 查看服务状态
./status.sh
```

### 分步部署

如果需要更精细的控制，可以分步启动服务：

```bash
# 1. 启动基础服务（数据库和缓存）
docker-compose up -d mysql redis

# 2. 等待数据库初始化完成（约30秒）
docker-compose logs mysql | grep "ready for connections"

# 3. 启动应用服务
docker-compose up -d epusdt backend

# 4. 启动前端服务
docker-compose up -d frontend
```

### 验证部署

部署完成后，通过以下步骤验证系统是否正常运行：

1. **访问前端界面**: http://localhost:3001
2. **访问Epusdt管理后台**: http://localhost:8000
3. **测试API接口**: http://localhost:8080/api/health
4. **访问数据库管理**: http://localhost:8081

### 初始配置

#### Epusdt配置
1. 首次访问 http://localhost:8000 时需要设置管理员账户
2. 配置USDT钱包地址和TRC20网络参数
3. 设置Telegram机器人（可选）

#### 数据库访问
- **地址**: http://localhost:8081
- **用户名**: root
- **密码**: epusdt123456

## 核心组件详解

### 1. Epusdt支付网关

<p align="center">
<img src="epusdt/wiki/img/usdtlogo.png" width="200">
</p>

#### 组件简介
Epusdt（Easy Payment Usdt）是一个由Go语言编写的私有化USDT支付中间件，支持TRC20网络。它为开发者提供了完整的数字货币支付解决方案，无需复杂配置即可实现USDT在线支付和消息回调功能。

#### 核心特性
- **私有化部署**: 完全自主控制，无需担心钱包被篡改和资金安全
- **跨平台支持**: 支持x86和ARM架构的Windows/Linux设备
- **高并发处理**: 多钱包地址轮询，支持高并发订单处理
- **异步队列**: 优雅的异步响应机制，确保高性能
- **零依赖部署**: 仅需一个二进制文件即可运行
- **完整API**: 提供RESTful API，支持任何系统集成
- **Telegram集成**: 支持机器人通知，实时支付消息推送

#### 技术架构
```
Epusdt架构
├── plugins/          # 插件系统（如独角数卡集成）
├── src/             # 核心源代码
│   ├── controller/  # API控制器
│   ├── model/       # 数据模型
│   ├── middleware/  # 中间件
│   ├── route/       # 路由配置
│   └── util/        # 工具函数
├── sql/             # 数据库脚本
└── wiki/            # 文档和教程
```

#### 工作原理
Epusdt通过监听TRC20网络API或节点，实时监控钱包地址的USDT入账事件。系统通过金额差异和时效性来判定交易归属，确保支付的准确性和安全性。

**支付流程**:
1. 客户发起支付请求（如20.05 USDT）
2. 系统分配可用钱包地址和精确金额
3. 客户按指定金额转账到指定地址
4. 系统监听到入账后自动确认支付
5. 触发回调通知商户系统

#### 集成指南
- **宝塔部署教程**: [epusdt/wiki/BT_RUN.md](epusdt/wiki/BT_RUN.md)
- **开发者API文档**: [epusdt/wiki/API.md](epusdt/wiki/API.md)
- **插件开发指南**: [epusdt/plugins/](epusdt/plugins/)

### 2. Payment Link MVP系统

#### 组件简介
Payment Link MVP是一个基于React + Material-UI + Gin的现代化支付链接管理平台，为用户提供直观的支付链接创建和管理功能。

#### 功能特性
- ✅ **访客落地页**: 精美的产品展示页面
- ✅ **用户系统**: 完整的注册、登录、权限管理
- ✅ **支付链接管理**: 创建、编辑、删除支付链接
- ✅ **公开支付页面**: 客户友好的支付界面
- ✅ **响应式设计**: 完美适配桌面和移动设备
- ✅ **订单管理**: 实时订单状态跟踪
- ✅ **数据统计**: 支付数据分析和报表

#### 技术栈
**前端技术**:
- React 18 - 现代化前端框架
- Material-UI - Google Material Design组件库
- TypeScript - 类型安全的JavaScript
- Axios - HTTP客户端
- React Router - 单页应用路由

**后端技术**:
- Gin - 高性能Go Web框架
- GORM - Go语言ORM框架
- JWT - JSON Web Token身份认证
- MySQL - 关系型数据库
- Redis - 缓存和会话存储

#### 项目结构
```
payment-link-mvp/
├── frontend/              # React前端应用
│   ├── src/
│   │   ├── components/    # React组件
│   │   ├── pages/         # 页面组件
│   │   ├── services/      # API服务
│   │   └── utils/         # 工具函数
│   ├── public/            # 静态资源
│   └── package.json       # 依赖配置
├── backend/               # Gin后端服务
│   ├── internal/
│   │   ├── handlers/      # HTTP处理器
│   │   ├── models/        # 数据模型
│   │   ├── services/      # 业务逻辑
│   │   └── middleware/    # 中间件
│   ├── migrations/        # 数据库迁移
│   └── main.go           # 应用入口
└── docker-compose.yml     # Docker配置
```

#### 开发环境搭建
```bash
# 后端开发
cd payment-link-mvp/backend
go mod tidy
go run main.go

# 前端开发
cd payment-link-mvp/frontend
npm install
npm start
```

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

## API接口文档

### 核心API端点

#### Epusdt支付API
```bash
# 创建支付订单
POST /api/v1/order/create-transaction
Content-Type: application/json

{
  "amount": "20.05",
  "order_id": "ORDER_123456",
  "callback_url": "https://your-site.com/callback"
}
```

#### Payment Link API
```bash
# 创建支付链接
POST /api/v1/payment-links
Authorization: Bearer <token>

{
  "title": "商品名称",
  "amount": "100.00",
  "description": "商品描述"
}
```

详细API文档请参考: [API接口文档.md](API接口文档.md)

## 安全考虑

### 部署安全
- 使用HTTPS协议部署生产环境
- 定期更新系统和依赖包
- 配置防火墙规则，仅开放必要端口
- 使用强密码和定期更换

### 数据安全
- 定期备份数据库
- 加密敏感配置信息
- 实施访问控制和权限管理
- 监控异常访问行为

### 钱包安全
- 使用冷钱包存储大额资金
- 定期轮换热钱包地址
- 设置合理的单笔和日限额
- 启用多重签名验证

## 技术支持与社区

### 官方资源
- **项目仓库**: https://github.com/WFHTask/Test_Tasks
- **问题反馈**: GitHub Issues
- **更新日志**: GitHub Releases

### Epusdt社区
- **Telegram频道**: https://t.me/epusdt
- **交流群组**: https://t.me/epusdt_group
- **官方文档**: [epusdt/wiki/](epusdt/wiki/)

### 开发者支持
- 提供完整的API文档和SDK
- 支持多种编程语言集成
- 提供示例代码和最佳实践
- 活跃的开发者社区支持

## 版本信息

### 当前版本
- **Epusdt**: v0.0.2
- **Payment Link MVP**: v1.0.0
- **整体项目**: v1.0.0

### 版本管理策略
项目采用本地版本管理，避免远程仓库依赖风险，确保部署的稳定性和可控性。

## 许可证

本项目遵循以下开源协议：
- **Epusdt**: [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html)
- **Payment Link MVP**: MIT License
- **整体项目**: 遵循各组件相应协议

## 免责声明

⚠️ **重要提醒**

本项目仅供学习和技术交流使用，请遵守当地法律法规：

1. **合规使用**: 请确保在您所在地区合法使用数字货币相关功能
2. **风险自担**: 使用过程中产生的任何法律责任由用户自行承担
3. **学习目的**: 项目中涉及的区块链代币均为学习用途
4. **投资警示**: 不鼓励和支持任何投机性的数字货币交易行为

---

**如果本项目对您有帮助，欢迎给我们一个⭐Star！**
