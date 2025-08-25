# 支付系统一键启动指南

## 🚀 项目概述

本项目包含两个主要系统：
1. **Epusdt支付系统** - USDT支付处理系统
2. **Payment Link MVP** - 支付链接管理系统（包含中间人钱包功能）

## 📁 项目结构

```
payment/
├── docker-compose.yml          # 主系统配置（Epusdt + 数据库）
├── epusdt/                     # Epusdt支付系统
├── payment-link-mvp/           # Payment Link MVP系统
│   ├── docker-compose.yml      # MVP系统配置
│   ├── backend/                # 后端服务
│   └── frontend/               # 前端服务
├── start.bat                   # Windows一键启动
├── start-all.sh                # Linux/Mac一键启动
├── stop-all.bat                # Windows一键停止
├── stop-all.sh                 # Linux/Mac一键停止
├── restart-all.bat             # Windows一键重启
├── restart-all.sh              # Linux/Mac一键重启
├── logs-all.bat                # Windows查看日志
├── logs-all.sh                 # Linux/Mac查看日志
├── status.bat                  # Windows状态检查
└── status.sh                   # Linux/Mac状态检查
```

## 🎯 一键启动

### Windows 用户
```bash
# 双击运行或在命令行中执行
start.bat
```

### Linux/Mac 用户
```bash
# 设置执行权限（首次使用）
chmod +x *.sh

# 启动所有服务
./start-all.sh
```

## 🛑 一键停止

### Windows 用户
```bash
stop-all.bat
```

### Linux/Mac 用户
```bash
./stop-all.sh
```

## 🔄 一键重启

### Windows 用户
```bash
restart-all.bat
```

### Linux/Mac 用户
```bash
./restart-all.sh
```

## 📊 状态检查

### Windows 用户
```bash
status.bat
```

### Linux/Mac 用户
```bash
./status.sh
```

## 📋 服务地址

启动成功后，可以通过以下地址访问：

| 服务 | 地址 | 说明 |
|------|------|------|
| Epusdt管理后台 | http://localhost:8000 | USDT支付系统管理界面 |
| Payment Link前端 | http://localhost:3000 | 支付链接管理前端 |
| Payment Link后端API | http://localhost:8080 | 后端API接口 |
| MySQL数据库 | localhost:3306 | 数据库服务 |
| Redis缓存 | localhost:6379 | 缓存服务 |

## 🔐 默认账户

### Epusdt管理后台
- **用户名**: admin
- **密码**: admin123

### Payment Link系统

#### 管理员账户
- **邮箱**: admin@example.com
- **密码**: admin123
- **角色**: admin
- **登录后**: 直接进入管理后台

#### 普通用户
- 需要先注册账户
- 登录后进入dashboard界面

## 👨‍💼 管理员功能

管理员用户登录后可以直接访问管理后台，包含以下功能：

### 用户管理
- 查看所有注册用户
- 管理用户权限
- 用户状态管理

### 订单管理
- 查看所有支付订单
- 订单状态跟踪
- 支付记录查询

### 中间人钱包管理
- 添加/删除中间人钱包地址
- 钱包状态管理
- 钱包余额查看

### 转账记录
- 查看所有转账记录
- 转账状态跟踪
- 手续费统计

### 费率配置
- 设置中间人手续费率
- 最小/最大手续费限制
- 费率历史记录

## 🔄 用户权限说明

### 管理员用户
- 登录后直接进入管理后台
- 可以访问所有管理功能
- 无法访问普通用户的dashboard页面
- 管理后台不显示"返回用户面板"按钮

### 普通用户
- 登录后进入dashboard页面
- 无法访问管理后台
- 只能管理自己的支付链接

## 🛠️ 手动命令

如果不想使用脚本，也可以直接使用Docker Compose命令：

```bash
# 启动Epusdt系统
docker-compose up -d mysql redis epusdt

# 启动Payment Link MVP系统
cd payment-link-mvp
docker-compose up -d

# 停止所有服务
docker-compose down
cd ..
docker-compose down

# 查看日志
docker-compose logs -f
cd payment-link-mvp
docker-compose logs -f
```

## 🔧 故障排除

### 1. Docker未安装
- **Windows/Mac**: 安装Docker Desktop
- **Linux**: 安装Docker Engine

### 2. Docker未运行
- **Windows/Mac**: 启动Docker Desktop
- **Linux**: `sudo systemctl start docker`

### 3. 端口被占用
如果端口被占用，可以修改docker-compose.yml文件中的端口映射：
- 3000 → 3001 (前端)
- 8000 → 8001 (Epusdt)
- 8080 → 8081 (后端API)

### 4. 权限问题（Linux/Mac）
```bash
chmod +x *.sh
```

### 5. 服务启动失败
```bash
# 查看详细日志
./logs-all.sh

# 检查服务状态
./status.sh

# 重新构建镜像
docker-compose build --no-cache
```

## 📝 开发模式

### 前端开发
```bash
cd payment-link-mvp/frontend
npm install
npm start
```

### 后端开发
```bash
cd payment-link-mvp/backend
go mod tidy
go run main.go
```

## 🔄 数据持久化

- **MySQL数据**: 保存在 `./epusdt/data/mysql/`
- **Redis数据**: 保存在 `./epusdt/data/redis/`
- **应用日志**: 保存在 `./logs/`

## ⚠️ 注意事项

1. **首次启动**: 需要下载Docker镜像，可能需要较长时间
2. **系统要求**: 确保有足够的磁盘空间和内存
3. **网络连接**: 需要网络连接下载Docker镜像
4. **防火墙**: 确保端口3000、8000、8080、3306、6379未被防火墙阻止
5. **数据备份**: 定期备份数据库数据
6. **管理员密码**: 首次登录后建议修改默认管理员密码
7. **用户权限**: 管理员和普通用户权限完全分离

## 🆘 常见问题

### Q: 启动后无法访问服务
A: 检查服务状态：`./status.sh`，确保所有容器都正常运行

### Q: 数据库连接失败
A: 等待数据库完全启动（通常需要30-60秒）

### Q: 前端页面空白
A: 检查后端API是否正常运行，查看浏览器控制台错误信息

### Q: 支付功能异常
A: 检查Epusdt系统状态，确保钱包地址配置正确

### Q: 管理员登录失败
A: 确保使用正确的管理员账户：admin@example.com / admin123

### Q: 返回用户面板按钮不工作
A: 管理员用户不需要访问普通用户面板，该按钮已被移除

## 📞 技术支持

如果遇到问题，请：
1. 查看服务日志：`./logs-all.sh`
2. 检查服务状态：`./status.sh`
3. 查看本文档的故障排除部分
4. 检查Docker和系统资源使用情况
