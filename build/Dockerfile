ARG ALPINE_VERSION=3.23

FROM alpine:$ALPINE_VERSION

RUN apk add --no-cache ca-certificates

RUN adduser -D -u 1000 expouser

COPY --chmod=755 ./bin/expo /usr/local/bin/expo

USER expouser

CMD ["expo"]