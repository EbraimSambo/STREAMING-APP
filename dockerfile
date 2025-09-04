FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copia os arquivos do projeto
COPY . .

# Compila o binário Go a partir do diretório cmd
RUN go build -o app ./cmd

# Etapa de produção (final)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

RUN apk add --no-cache ffmpeg

# Copia o binário gerado
COPY --from=builder /app/app .

EXPOSE 3344

CMD ["./app"]
