docker buildx build --platform linux/amd64,linux/arm64 -t eh-go-server:latest  .
docker save -o eh-go-server.tar eh-go-server:latest
