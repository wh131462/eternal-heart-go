# 显式指定平台（以 AMD64 为例）
# 若需支持多平台，建议使用 docker buildx 构建
ARG TARGETPLATFORM=linux/amd64
FROM --platform=$TARGETPLATFORM golang:latest AS builder

# 设置 Go 模块代理和私有库配置
ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 设置工作目录
WORKDIR /workspace

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖（缓存优化）
RUN go mod download

# 复制项目代码
COPY . .

# 编译 Go 项目（禁用 CGO）
RUN go build -ldflags="-s -w" -o eh-go-server .

# 使用与构建平台一致的运行时镜像
FROM --platform=$TARGETPLATFORM alpine:latest

# 设置时区（可选）
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# 设置工作目录
WORKDIR /workspace

# 从 builder 镜像中复制编译好的二进制文件
COPY --from=builder /workspace/eh-go-server .

# 暴露端口
EXPOSE 9999

# 运行程序
CMD ["./eh-go-server"]