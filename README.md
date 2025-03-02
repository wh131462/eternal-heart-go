# eternal-heart-go
作为go的实践项目,EH架构的核心接口和内容都从此服务提供.

## dev
先构建本地镜像:
```shell
docker build -t eh-go-server:latest .
```
上传到服务端,执行:
```shell
docker run -d -p 9999:9999 --name eh-go-server eh-go-server:latest
```

验证运行:
```shell
curl http://localhost:9999/wx
```

打包,构建本地镜像并打标签:
```shell
docker buildx build --platform linux/amd64 -t eh-go-server:latest . --load
```

保存到tar文件:
```shell
docker save -o eh-go-server.tar eh-go-server:latest
```

加载镜像(需要先上传到服务器):
```shell
docker load -i /workspace/go/eh-go-server.tar
```
运行即可:
```shell
docker run -d --name eh-go-server -p 9999:9999 eh-go-server:latest
```

### 多平台构建

```shell
# 启用 buildx 并创建多平台构建器
docker buildx create --use --name multi-platform-builder
# 安装跨平台模拟器（非必需，但推荐）
docker run --privileged --rm tonistiigi/binfmt --install all
# 构建并推送多架构镜像（AMD64 + ARM64）按需构建
docker buildx build --platform linux/amd64 -t eh-go-server:latest . --load
docker buildx build --platform linux/arm64 -t eh-go-server:latest . --load
docker manifest create eh-go-server:latest eh-go-server:amd64 eh-go-server:arm64
docker manifest push eh-go-server:latest
```


### optional
标记镜像
将本地镜像标记为 Docker Hub 上的镜像：

```Bash
docker tag eh-go-server:latest eternalheart/eh-go-server:latest
```

推送镜像:
```shell
docker push eternalheart/eh-go-server:latest
```
在服务器中拉取镜像并执行:
```shell
# 拉取镜像
docker pull eternalheart/eh-go-server:latest
# 运行容器
docker run -d -p 9999:9999 --name eh-go-server eternalheart/eh-go-server:latest
```

