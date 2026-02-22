FROM m.daocloud.io/docker.io/library/golang:1.25 AS builder

WORKDIR /app

# 安装必要的构建依赖
RUN apt-get update && apt-get install -y \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./main.go

FROM m.daocloud.io/docker.io/library/ubuntu:22.04

RUN apt-get update && apt-get install -y \
    ca-certificates \
    libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# 复制二进制文件和配置
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
COPY --from=builder /app/imagenet_classes.txt /root/imagenet_classes.txt

# 复制 ONNX 库和模型
COPY --from=builder /app/libs/libonnxruntime.so /usr/lib/libonnxruntime.so
RUN ln -sf /usr/lib/libonnxruntime.so /usr/lib/onnxruntime.so
COPY --from=builder /app/models /root/models

EXPOSE 9090

ENV ORT_DYLIB_PATH=/usr/lib/libonnxruntime.so
ENV IS_DOCKER=true

CMD ["./server"]
