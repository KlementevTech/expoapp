ARG ALPINE_VERSION=3.23

FROM alpine:$ALPINE_VERSION

RUN apk add --no-cache ca-certificates

RUN adduser -D -u 1000 appuser

COPY --chmod=755 ./bin/expoapp /usr/local/bin/expoapp

USER appuser

CMD ["expoapp"]