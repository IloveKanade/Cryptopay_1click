# 使用Go官方镜像作为构建环境
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制源代码
COPY src/ .

# 设置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o epusdt main.go

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S epusdt && \
    adduser -u 1001 -S epusdt -G epusdt

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/epusdt .

# 复制静态文件
COPY --from=builder /app/static ./static

# 创建必要的目录
RUN mkdir -p /app/logs /app/runtime && \
    chown -R epusdt:epusdt /app

# 切换到非root用户
USER epusdt

# 暴露端口
EXPOSE 8000

# 启动应用
CMD ["./epusdt"]