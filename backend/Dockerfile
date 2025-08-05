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

FROM alpine:3.20

RUN apk add --no-cache \
    postgresql15-client \
    tzdata \
    bash

WORKDIR /app

COPY --from=builder /app/wait-for-postgres.sh /app/
RUN chmod +x /app/wait-for-postgres.sh

COPY --from=builder /app/migrator /app/app ./
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/order.json .

ENV CONFIG_PATH=/app/configs/config.yaml