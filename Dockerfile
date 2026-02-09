FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o server ./main.go

FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

RUN mkdir -p /root/models/mobilenetv2

EXPOSE 9090

CMD ["./server"]
