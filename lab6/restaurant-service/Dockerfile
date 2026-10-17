# === Этап 1: сборка ===
FROM golang:1.22-alpine AS builder

WORKDIR /app

# go.mod и go.sum копируются отдельно: слой с зависимостями
# переиспользуется из кэша, если изменился только код
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# === Этап 2: рантайм ===
FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/server .

# HTTP и gRPC
EXPOSE 8080 9090

ENTRYPOINT ["./server"]
