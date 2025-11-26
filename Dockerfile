# Build stage
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Разрешаем Go автоматически скачать нужный toolchain
ENV GOTOOLCHAIN=auto

# Скачиваем зависимости
RUN go mod download

# Копируем остальной код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o task-service ./cmd/server

# Runtime stage
FROM alpine:3.19

# Устанавливаем ca-certificates для HTTPS и tzdata для таймзон
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/task-service .

# Создаём непривилегированного пользователя
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8082

CMD ["./task-service"]
