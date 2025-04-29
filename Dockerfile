FROM golang:1.24 AS builder

# Установим git через apt, не apk
RUN apt-get update && apt-get install -y git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o load_balancer main.go

# Финальный образ
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/load_balancer .

COPY config.yaml .

CMD ["./load_balancer", "--config=config.yaml"]
