FROM alpine:latest

RUN mkdir /app

COPY server /app

CMD ["/app/server"]