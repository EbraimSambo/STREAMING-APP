# Etapa de build (builder)
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copia os arquivos do projeto
COPY . .

# Compila o binário Go
RUN go build -o app .

# Etapa de produção (final)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

# Copia o binário compilado
COPY --from=builder /app/app .

EXPOSE 3344

CMD ["./app"]
