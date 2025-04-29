FROM golang:1.24 AS builder


RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o load_balancer .

EXPOSE 8080

CMD ["./load_balancer"]
