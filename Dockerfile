# ── Stage 1: Build ──
FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GOTOOLCHAIN=auto

# Dependencias del sistema necesarias para compilar
RUN apk add --no-cache git ca-certificates tzdata

# Copiar módulos primero para aprovechar cache de Docker
COPY go.mod go.sum ./
RUN GOTOOLCHAIN=go1.24.13 go mod download

# Copiar el resto del código
COPY . .

# Compilar el binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/api/main.go

# ── Stage 2: Runtime ──
FROM alpine:3.20

WORKDIR /app

# Certificados para HTTPS (necesario para llamadas a AWS S3)
RUN apk add --no-cache ca-certificates tzdata

# Copiar solo el binario compilado
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
