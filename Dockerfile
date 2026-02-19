ARG GOLANG_VERSION=1.26
FROM golang:${GOLANG_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG CGO_ENABLED=0

RUN CGO_ENABLED=${CGO_ENABLED} GOOS=linux \
    go build \
    -ldflags="-X main.version=${VERSION} -s -w" \
    -o /tmp/expo \
    ./cmd/server

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /tmp/expo /bin/expo

USER 65532:65532

ENTRYPOINT ["/bin/expo"]
CMD ["--config", "/etc/expo/config.toml"]