# Build stage
FROM golang:1.25-trixie AS builder

# Single target build for the main bot

WORKDIR /app

# Copia os arquivos de módulo.
COPY go.mod go.sum ./

# Baixa as dependências. Esta camada será armazenada em cache
# e só será reconstruída se go.mod ou go.sum mudarem.
RUN go mod download

# (Opcional, mas boa prática) Verifica a integridade das dependências.
RUN go mod verify

# Copia o restante do código-fonte.
# Esta camada será invalidada com mais frequência.
COPY . .

# Build the binaries for the main bot and auxiliary tools
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/bot
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o healthcheck ./cmd/healthcheck

# Final stage
# Usamos a tag :nonroot para uma imagem que já roda sem privilégios.
FROM gcr.io/distroless/static-debian13:nonroot

# Definimos um diretório de trabalho limpo e sem privilégios.
WORKDIR /app

# Copy the pre-built binaries from the build stage
# Os binários são copiados para o /app definido acima.
COPY --from=builder /app/app .
COPY --from=builder /app/healthcheck .

# Command to run the executable
# O comando é executado a partir do WORKDIR /app.
CMD ["./app"]