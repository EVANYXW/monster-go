# FROM 基于 golang:1.16-alpine
FROM golang:1.19-alpine AS builder

WORKDIR /opt/repo/monster-go

# COPY 源路径 目标路径
COPY . .


RUN go env -w GO111MODULE=on \
    && go env -w GOPATH=/opt/repo \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o monsterGo .


# FROM 基于 alpine:latest
FROM alpine:latest

# RUN 设置代理镜像
RUN echo -e http://mirrors.ustc.edu.cn/alpine/v3.13/main/ > /etc/apk/repositories

# RUN 设置 Asia/Shanghai 时区
RUN apk --no-cache add tzdata  && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

# COPY 源路径 目标路径 从镜像中 COPY
COPY --from=builder /opt/repo/monster-go /opt

# EXPOSE 设置端口映射
EXPOSE 8023
EXPOSE 8024
EXPOSE 6060

# WORKDIR 设置工作目录
WORKDIR /opt

# CMD 设置启动命令
CMD ["./monsterGo","run","--server_name","world", "--env", "dev"]