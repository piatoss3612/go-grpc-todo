FROM alpine:latest

RUN mkdir /app

COPY proxy /app

CMD ["/app/proxy"]