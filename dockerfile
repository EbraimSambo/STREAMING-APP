FROM alpine:latest

WORKDIR /app 

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app . 

EXPOSE 3344

CMD ["./app"] 