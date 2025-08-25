# 支付链接 MVP

基于React + Material-UI + Gin的加密货币支付链接平台

## 项目结构

```
payment-link-mvp/
├── frontend/          # React + Material-UI 前端
├── backend/           # Gin 后端
├── docker-compose.yml # Docker 部署配置
└── README.md
```

## 功能特性

- ✅ 访客落地页展示
- ✅ 用户注册登录
- ✅ 支付链接创建和管理
- ✅ 公开支付页面
- ✅ 响应式设计

## 快速开始

### 开发环境

1. 启动后端服务
```bash
cd backend
go mod tidy
go run main.go
```

2. 启动前端服务
```bash
cd frontend
npm install
npm start
```

### Docker 部署

```bash
docker-compose up -d
```

## 技术栈

- **前端**: React 18 + Material-UI + TypeScript
- **后端**: Gin + GORM + MySQL
- **支付服务**: Epusdt
- **部署**: Docker + Docker Compose
