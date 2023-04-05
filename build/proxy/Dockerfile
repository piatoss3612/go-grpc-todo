FROM golang:1.20.2-alpine3.17 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./build/proxy/proxy ./cmd/proxy

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/build/proxy /app

CMD ["/app/proxy"]