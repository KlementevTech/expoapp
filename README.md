# expoapp

```shell
# Copy environments
cp .env.example .env
```

```shell
# Building docker image
make docker_build
```

```shell
# Starting dev container
make up
```

```shell
# Sending proto.ExpoService/GetInfoV1 request
grpcurl \
  -plaintext \
  -d '{}' \
  localhost:50051 \
  proto.ExpoService/GetInfoV1
```