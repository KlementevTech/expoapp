# expoapp

```shell
# Copy environments
cp .env.example .env
```

```shell
# Starting Expoapp application
make run
```

```shell
# Sending proto.ExpoService/GetInfoV1 request
grpcurl \
  -plaintext \
  -d '{}' \
  127.0.0.1:50051 \
  proto.ExpoService/GetInfoV1
```