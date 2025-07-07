# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем модуль и зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходники
COPY . .

# Копируем конфиг из папки config
COPY config.env .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# Stage 2: Run
FROM alpine:3.18

# Установка сертификатов (если нужны HTTPS-запросы)
RUN apk add --no-cache ca-certificates

WORKDIR /
COPY --from=builder /app/app .
COPY --from=builder /app/config.env .

CMD ["./app"]