# Etapa 1: build da aplicação
FROM golang:1.23.7-alpine AS builder

# Instala ferramentas essenciais (ex: swag, se necessário)
RUN apk add --no-cache git

# Define diretório de trabalho dentro do container
WORKDIR /app

# Copia arquivos de dependência e baixa módulos
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o restante do projeto
COPY . .

# Compila a aplicação Go
RUN CGO_ENABLED=0 go build -o main .

# Etapa 2: imagem final enxuta
FROM gcr.io/distroless/static:nonroot

# Define diretório de trabalho final
WORKDIR /

# Copia binário gerado na etapa anterior
COPY --from=builder /app/main /main

# Expõe a porta usada pela aplicação (ajuste se necessário)
EXPOSE 8080

# Comando que será executado ao iniciar o container
CMD ["/main"]
