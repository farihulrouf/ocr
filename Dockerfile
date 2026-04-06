# STAGE 1: Build (Pake Alpine biar ringan)
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files dulu (agar caching layer efisien)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Compile API dan Worker
# CGO_ENABLED=0 agar binary bersifat static dan bisa jalan di Alpine mana pun
RUN CGO_ENABLED=0 GOOS=linux go build -o seido-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o seido-worker ./cmd/ocr-worker

# STAGE 2: Run (Hanya berisi binary)
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /app/seido-api .
COPY --from=builder /app/seido-worker .

# Copy folder configs/docs jika aplikasi kamu membutuhkannya saat runtime
COPY --from=builder /app/configs ./configs

# Default port untuk API
EXPOSE 3000

# Default command (bisa dioverride saat running)
CMD ["./seido-api"]
