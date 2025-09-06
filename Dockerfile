# Билд образ
FROM 172.16.99.6:8085/golang:1.24.0-alpine AS builder
LABEL stage=gobuilder

RUN apk update --no-cache


ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/main cmd/http/main.go

# Финальный образ
FROM 172.16.99.6:8085/alpine:3.19.1

RUN apk update --no-cache && apk add --no-cache ca-certificates

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/main /app/main

CMD ["./main"]
