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
# Sending GetInfo request
grpcurl \
  -plaintext \
  -d '{}' \
  127.0.0.1:9091 \
  expo.ExpoService/GetInfo
```