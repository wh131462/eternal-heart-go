# 启用 BuildKit 参数（必须放在文件开头）
# syntax=docker/dockerfile:1.4

# 多阶段构建：编译阶段（自动适配宿主机架构）
FROM golang:latest AS builder

# 接收构建时传入的平台参数
ARG TARGETOS
ARG TARGETARCH

# 配置Go编译环境
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH}

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 静态编译优化（移除调试信息）
RUN go build -trimpath -ldflags="-s -w" -o eh-go-server .

# 多阶段构建：运行时镜像（自动匹配目标平台）
FROM alpine:latest

# 时区配置（国内服务器建议）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /workspace
COPY --from=builder /workspace/eh-go-server .

# 暴露端口（与您的监听端口一致）
EXPOSE 9999

# 容器启动命令
CMD ["./eh-go-server"]