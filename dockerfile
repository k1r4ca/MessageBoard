# 基于官方的Golang镜像构建
FROM golang:1.21

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件拷贝到工作目录
COPY . .

# 下载依赖并构建应用
RUN go mod download
RUN go build -o main .

# 暴露端口
EXPOSE 10070

# 启动应用
CMD ["./main"]