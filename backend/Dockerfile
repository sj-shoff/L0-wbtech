FROM golang:1.24.3-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o app ./cmd/app/main.go
RUN go build -ldflags="-w -s" -o migrator ./cmd/migrator/main.go

FROM alpine:latest

RUN apk add --no-cache \
    postgresql15-client \
    tzdata \
    bash

WORKDIR /app

COPY --from=builder /app/wait-for-postgres.sh /app/
COPY --from=builder /app/app /app/
COPY --from=builder /app/migrator /app/
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/order.json .

RUN chmod +x /app/wait-for-postgres.sh \
    && chmod +x /app/app \
    && chmod +x /app/migrator

ENV CONFIG_PATH=/app/configs/config.yaml