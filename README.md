# eternal-heart-go
作为go的实践项目,EH架构的核心接口和内容都从此服务提供.

## dev
先构建本地镜像:
```shell
docker build -t eh-go-server:latest .
```
上传到服务端,执行:
```shell
docker run -d -p 2341:2341 --name eh-go-server eh-go-server:latest
```

验证运行:
```shell
curl http://localhost:2341/wx
```

### optional
标记镜像
将本地镜像标记为 Docker Hub 上的镜像：

```Bash
docker tag eh-go-server:latest your_dockerhub_username/eh-go-server:latest
```
your_dockerhub_username：你的 Docker Hub 用户名。

推送镜像:
```shell
docker push your_dockerhub_username/eh-go-server:latest
```
在服务器中拉取镜像并执行:
```shell
# 拉取镜像
docker pull your_dockerhub_username/eh-go-server:latest
# 运行容器
docker run -d -p 2341:2341 --name eh-go-server your_dockerhub_username/eh-go-server:latest
```

