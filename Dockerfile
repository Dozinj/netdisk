FROM golang:latest  AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ENV GOPROXY=https://goproxy.cn,direct

# 移动到工作目录 build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将代码编译为可执行文件到app
RUN go build -o app .

# 创建一个小镜像
# FROM debian:stretch-slim

FROM scratch

# 将静态文件或者配置文件拷贝到项目中
COPY ./conf /conf

# 从builder镜像中将拉取 build/app 到当前目录
COPY --from=builder /build/app /


ENTRYPOINT ["/app", "conf/config.ini"]



