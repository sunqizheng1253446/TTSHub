# 构建阶段
FROM golang:1.23-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev libc-dev pkgconfig sqlite-dev

# 设置工作目录
RUN mkdir /src
WORKDIR /src

# 复制源代码
COPY . /src/

# 下载依赖
RUN go mod download

# 设置CGO环境变量
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV CGO_LDFLAGS_ALLOW=".*"

# 构建应用
RUN go build -v -o ./bin/server ./cmd/ttshub/main.go

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache ca-certificates sqlite-libs

# 创建非root用户
RUN adduser -D -s /bin/sh ttshub

# 从构建阶段复制二进制文件
COPY --from=builder /src/bin/server /app/server

# 复制配置文件
COPY --from=builder /src/configs /app/configs

# 复制静态文件
COPY --from=builder /src/public /app/public

# 设置工作目录
WORKDIR /app

# 更改文件所有者
RUN chown -R ttshub:ttshub /app

# 切换到非root用户
USER ttshub

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./server"]