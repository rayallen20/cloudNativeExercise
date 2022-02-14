# 编译用镜像
# 指定GO编译器版本
ARG GO_VERSION=1.17.6
FROM golang:${GO_VERSION} as builder
# 安装git
RUN apt-get update && apt-get install -y git
WORKDIR /root/
# 拉取代码 此处直接clone也是不正规的 应该用pull
RUN git clone https://github.com/rayallen20/httpserver.git
WORKDIR /root/httpserver/
# 拉取依赖
RUN go mod vendor
# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o httpserver

# 运行用镜像
FROM alpine:latest
RUN apk --no-cache add ca-certificates
# 指定环境变量version
ENV VERSION=1.0
WORKDIR /root/
# 从编译用镜像处复制二进制文件
COPY --from=builder /root/httpserver/httpserver .
CMD ["./httpserver"]