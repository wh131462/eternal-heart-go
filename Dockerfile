# 使用官方 Go 镜像作为基础镜像
FROM golang:1.20-alpine AS builder

# 设置工作目录
WORKDIR /workspace

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制项目代码
COPY . .

# 编译 Go 项目
RUN go build -o eh-go-server .

# 使用轻量级的 Alpine 镜像作为运行时镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /workspace

# 从 builder 镜像中复制编译好的二进制文件
COPY --from=builder /workspace/eh-go-server .

# 暴露端口
EXPOSE 2341

# 运行程序
CMD ["./eh-go-server"]