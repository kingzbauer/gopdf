FROM alpine:latest

LABEL maintainer="jackmwangi@gmail.com"

WORKDIR /app

COPY browserless_linux-amd64 browserless

RUN apk update; \
    apk add chromium;

ENV CHROME_BIN_PATH=/usr/bin/chromium-browser

EXPOSE 8089

ENTRYPOINT ["/app/browserless"]
