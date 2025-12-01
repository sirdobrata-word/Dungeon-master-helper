# Многоэтапная сборка для минимального образа
# Этап 1: Сборка приложения
FROM golang:1.21-alpine AS builder

# Устанавливаем зависимости для PostgreSQL (libpq требует CGO)
RUN apk add --no-cache gcc musl-dev postgresql-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости (кэшируется, если go.mod/go.sum не изменились)
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник (CGO_ENABLED=1 для работы с PostgreSQL)
RUN CGO_ENABLED=1 GOOS=linux go build -a -o server ./cmd/server

# Этап 2: Финальный минимальный образ
FROM alpine:3.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из этапа сборки
COPY --from=builder /app/server /app/server

# Копируем веб-интерфейс
COPY --from=builder /app/web /app/web

# Порт по умолчанию (можно переопределить через переменную окружения)
ENV PORT=9190

# Открываем порт
EXPOSE 9190

# Запускаем сервер
CMD ["/app/server"]

