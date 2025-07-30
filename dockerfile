FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ffmpeg

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app .

RUN adduser -D -g '' appuser

USER appuser

EXPOSE 3344

CMD ["./app"]
