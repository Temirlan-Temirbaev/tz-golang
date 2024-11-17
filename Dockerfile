# Устанавливаем базовый образ с Go
FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy
COPY . .

RUN go build -o main .

CMD ["./main"]