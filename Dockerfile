ARG ALPINE_VERSION=3.23
FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 expouser

WORKDIR /usr/local/bin

COPY --chmod=755 bin/expo expo

USER expouser

CMD ["expo"]