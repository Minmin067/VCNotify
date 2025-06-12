# Builderステージ
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY bot/go.mod bot/go.sum ./
RUN go mod download
COPY bot/ ./bot/
WORKDIR /app/bot
RUN go build -o /vcnotify

# ランタイムステージ
FROM alpine:latest
COPY --from=builder /vcnotify /vcnotify
WORKDIR /app
CMD ["/vcnotify"]