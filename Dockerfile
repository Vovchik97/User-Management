FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOD=linux go build -o api ./cmd/api/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /app/api .
COPY --from=builder /app/.env.example ./.env

# Открываем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./api"]