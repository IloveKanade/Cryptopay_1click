# 启动脚本使用说明

## 脚本文件

本项目提供了以下启动和管理脚本：

### Windows 系统
- `start.bat` - 启动所有服务
- `stop.bat` - 停止所有服务
- `logs.bat` - 查看服务日志

### Linux/Mac 系统
- `start.sh` - 启动所有服务
- `stop.sh` - 停止所有服务
- `logs.sh` - 查看服务日志

## 使用方法

### 启动服务

#### Windows
```bash
# 双击运行或在命令行中执行
start.bat
```

#### Linux/Mac
```bash
# 首先设置执行权限（仅首次需要）
chmod +x start.sh

# 运行启动脚本
./start.sh
```

### 停止服务

#### Windows
```bash
stop.bat
```

#### Linux/Mac
```bash
chmod +x stop.sh
./stop.sh
```

### 查看日志

#### Windows
```bash
logs.bat
```

#### Linux/Mac
```bash
chmod +x logs.sh
./logs.sh
```

## 服务地址

启动成功后，可以通过以下地址访问：

- **前端界面**: http://localhost:3000
- **后端API**: http://localhost:8080

## 手动命令

如果不想使用脚本，也可以直接使用Docker Compose命令：

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down

# 查看日志
docker-compose logs -f

# 重新构建并启动
docker-compose up -d --build
```

## 故障排除

### 1. Docker未安装
确保已安装Docker Desktop（Windows/Mac）或Docker Engine（Linux）

### 2. Docker未运行
- Windows/Mac: 启动Docker Desktop
- Linux: 启动Docker服务 `sudo systemctl start docker`

### 3. 端口被占用
如果3000或8080端口被占用，可以修改docker-compose.yml文件中的端口映射

### 4. 权限问题（Linux/Mac）
确保脚本有执行权限：
```bash
chmod +x *.sh
```

## 开发模式

如果需要开发模式（代码热重载），可以修改docker-compose.yml文件，将构建模式改为开发模式。

## 注意事项

1. 首次启动可能需要较长时间，因为需要下载Docker镜像
2. 确保有足够的磁盘空间和内存
3. 如果修改了代码，需要重新构建镜像
4. 数据库数据会持久化保存在本地
