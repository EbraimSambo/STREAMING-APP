# Etapa final
FROM alpine:latest

WORKDIR /app 

RUN apk --no-cache add ca-certificates

# Copia o binário gerado na etapa anterior
COPY --from=builder /app/app . 

# Expõe a porta usada pelo Fiber
EXPOSE 3344

# Comando para iniciar a aplicação
CMD ["./app"] # <--- Isso executará /app/app porque o WORKDIR é /a