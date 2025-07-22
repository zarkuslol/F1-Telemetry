# --- Estágio 1: Build ---
# Usamos a imagem oficial do Go para compilar nossa aplicação.
# A tag 'alpine' resulta em uma imagem de build menor.
FROM golang:1.24-alpine AS builder

# Define o diretório de trabalho.
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências.
# Isso aproveita o cache do Docker: as dependências só são baixadas novamente se o go.mod/sum mudar.
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código-fonte da aplicação.
COPY . .

# Compila a aplicação.
# CGO_ENABLED=0 cria um executável estático, sem depender de bibliotecas C do sistema.
# -o /app/f1-collector define o nome do arquivo de saída.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/f1-collector /app/src/

# --- Estágio 2: Final ---
# Usamos uma imagem base mínima. 'alpine' é pequena e segura.
FROM alpine:latest

# Define o diretório de trabalho.
WORKDIR /app

# Copia APENAS o executável compilado do estágio de build.
COPY --from=builder /app/f1-collector .

# Expõe a porta UDP para o jogo.
EXPOSE 20777/udp

# Define o comando para executar nossa aplicação.
CMD ["./f1-collector"]
