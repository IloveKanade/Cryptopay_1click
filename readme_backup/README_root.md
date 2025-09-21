# 支付系统集成项目

本项目集成了多种支付解决方案，包括 epusdt 数字货币收款网关等。

## 快速启动

### 一键启动（推荐）
```bash
# 克隆项目
git clone https://github.com/WFHTask/Test_Tasks.git
cd Test_Tasks

# 启动所有服务
docker-compose up -d
```

### 手动启动
```bash
# 1. 启动基础服务（MySQL、Redis）
docker-compose up -d mysql redis

# 2. 等待数据库就绪后启动应用服务
docker-compose up -d epusdt backend frontend
```

## 重要说明

### 服务访问地址
- **前端界面**: http://localhost:3001
- **后端API**: http://localhost:8080
- **epusdt管理**: http://localhost:8000
- **数据库管理**: http://localhost:8081 (phpMyAdmin)

### 默认账户
- **数据库**: root / epusdt123456
- **epusdt**: 首次访问时需要设置管理员账户

### 版本管理
项目使用本地版本的 epusdt v0.0.2，避免远程仓库依赖风险。

## 核心功能

### Epusdt支付系统
- 私有化USDT支付网关
- 支持TRC20网络
- 多钱包地址轮询
- Telegram机器人通知
- 完整的支付回调机制

### Payment Link MVP
- 支付链接生成和管理
- 用户认证和授权
- 订单管理系统
- RESTful API接口

## 管理命令

```bash
# 启动所有服务
./start-all.sh

# 停止所有服务
./stop-all.sh

# 重启所有服务
./restart-all.sh

# 查看服务日志
./logs-all.sh

# 查看服务状态
./status.sh
```

## 故障排除

### 1. 端口冲突
```bash
# 确保 3001、8000、8080、8081、3306、6379 端口未被占用
netstat -tulpn | grep :3001
netstat -tulpn | grep :8000
netstat -tulpn | grep :8080
netstat -tulpn | grep :8081
netstat -tulpn | grep :3306
netstat -tulpn | grep :6379
```

### 2. epusdt目录为空
```bash
# 初始化子模块
git submodule update --init --recursive
```

### 2. Docker服务启动失败
```bash
# 检查Docker是否运行
docker info

# 检查端口占用
netstat -tulpn | grep :8000
netstat -tulpn | grep :3001
netstat -tulpn | grep :8080
```

### 3. 数据库连接失败
```bash
# 检查MySQL服务状态
docker-compose ps mysql

# 查看MySQL日志
docker-compose logs mysql
```

## 开发说明

### 项目结构
```
Test_Tasks/
├── epusdt/                 # USDT支付网关（Git子模块）
├── payment-link-mvp/       # 支付链接系统
│   ├── backend/           # Go后端服务
│   └── frontend/          # React前端应用
├── docker-compose.yml     # Docker编排配置
├── start-all.sh          # 一键启动脚本
└── README.md             # 项目文档
```

### 环境要求
- Docker & Docker Compose
- Git
- Linux/macOS/Windows

## 技术栈

### 后端技术
- **Go** - 高性能后端服务
- **Gin** - Web框架
- **GORM** - ORM框架
- **JWT** - 身份认证
- **MySQL** - 关系型数据库
- **Redis** - 缓存和会话存储

### 前端技术
- **React** - 用户界面框架
- **Material-UI** - UI组件库
- **Axios** - HTTP客户端
- **React Router** - 路由管理

### 基础设施
- **Docker** - 容器化部署
- **Docker Compose** - 服务编排
- **Nginx** - 反向代理（生产环境）

## 许可证

本项目遵循相关开源协议，详见各组件的LICENSE文件。
