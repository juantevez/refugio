# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Instalamos git por si alguna dependencia de go mod lo requiere
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilación optimizada
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /shelter-api ./cmd/api/main.go

# Stage 2: Final image (Distroless o Alpine para mayor seguridad)
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /shelter-api .

EXPOSE 8080

CMD ["./shelter-api"]
