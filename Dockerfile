# STAGE 1: Build
FROM golang:1.24-alpine AS builder

# Install dependencies yang dibutuhkan untuk build
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# 1. Cache dependencies (Sangat penting agar build cepat)
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy source code
COPY . .

# 3. Compile API
# Pastikan folder cmd/api memang ada. Jika namanya cmd/main.go, sesuaikan pathnya.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o seido-api ./cmd/api

# 4. Compile Worker
# Saya sesuaikan dengan path yang umum kamu pakai. 
# Jika nama foldernya 'worker', ganti menjadi ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux go build -v -o seido-worker ./cmd/ocr-worker

# STAGE 2: Runtime (Gunakan Alpine yang super ringan)
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary hasil build
COPY --from=builder /app/seido-api .
COPY --from=builder /app/seido-worker .

# Pastikan folder configs ikut terbawa karena aplikaCOPY --from=builder /app/.env ./.envsi butuh .env atau yaml
COPY --from=builder /app/configs ./configs
# Jika kamu punya file .env di root, copy juga
#COPY --from=builder /app/.env ./.env 

# Fiber biasanya pakai 8080 atau 3000, sesuaikan dengan settingan API kamu
EXPOSE 8080

CMD ["./seido-api"]