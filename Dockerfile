# --- ETAPA 1: Construcción (Builder) ---
FROM golang:1.25-alpine AS builder




# Instalamos certificados y herramientas básicas
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copiamos los archivos de dependencias primero para aprovechar el cache de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY . .

# Compilamos el binario (CGO_ENABLED=0 para que sea estático y ligero)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- ETAPA 2: Producción (Final) ---
FROM alpine:latest

# Instalamos ca-certificates para que las peticiones HTTPS funcionen
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiamos solo el binario desde la etapa anterior
COPY --from=builder /app/main .

# Exponemos el puerto de tu servidor GraphQL
EXPOSE 8080

# Comando para arrancar la app
CMD ["./main"]