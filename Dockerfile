# ========================================
# Build Stage
# ========================================
FROM golang:1.24-alpine AS builder

# Instalar dependências de build
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /build

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download de dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o zpwoot \
    ./cmd/zpwoot/main.go

# ========================================
# Runtime Stage
# ========================================
FROM alpine:latest

# Instalar certificados SSL e timezone data
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1000 zpwoot && \
    adduser -D -u 1000 -G zpwoot zpwoot

# Criar diretório de dados
RUN mkdir -p /app/data && \
    chown -R zpwoot:zpwoot /app

WORKDIR /app

# Copiar binário do builder
COPY --from=builder /build/zpwoot .

# Mudar para usuário não-root
USER zpwoot

# Expor porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando de execução
CMD ["./zpwoot"]

