# Etapa 1: Compilar o Go
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copia os arquivos de definição do módulo
COPY go.mod ./

# Baixa as dependências (cria o go.sum automaticamente)
RUN go mod tidy

# Copia o código fonte
COPY *.go ./

# Compila o executável
RUN go build -o weather-worker main.go

# Etapa 2: Imagem final leve para rodar
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/weather-worker .

# Comando para rodar
CMD ["./weather-worker"]