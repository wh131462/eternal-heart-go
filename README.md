# eternal-heart-go
EH架构的的核心后台代码,支持微信公众号系列接口,全平台的用户管理,权限服务管理,应用服务管理.

## ToDo
微信公众号部分
- [x] 整体菜单逻辑
- [x] 黄历查看
- [x] 复读机模式
- [ ] 生日管理
- [ ] AI api-key管理
用户管理
- [ ] 用户表建立维护
- [ ] 登录注册校验模块
权限管理
- [ ] [openfga](https://github.com/openfga/openfga)对接
博客管理
- [ ] blog manager
- [ ] website manager
## dev
构建本地镜像(暂时不用):
```shell
docker build -t eh-go-server:latest .
```

打包,使用build x构建多平台(当前为单个平台)本地镜像并打标签:
```shell
docker buildx build --platform linux/amd64 -t eh-go-server:latest . --load
```
保存到tar文件:
```shell
docker save -o eh-go-server.tar eh-go-server:latest
```
删除服务器中的镜像:
```shell
docker rm -f  eh-go-server
```
加载镜像(需要先上传到服务器):
```shell
docker load -i /workspace/go/eh-go-server.tar
```
运行即可:
```shell
docker run -d --name eh-go-server -p 9999:9999 eh-go-server:latest
```

### 多平台构建细节

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

