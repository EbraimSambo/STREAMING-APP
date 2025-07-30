# Etapa de build (builder)
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia os arquivos do projeto para dentro da imagem
COPY . .

# Compila o binário Go
RUN go build -o app .

# Etapa de produção (final)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copia o binário gerado na etapa de build
COPY --from=builder /app/app .

EXPOSE 3344

CMD ["./app"]
