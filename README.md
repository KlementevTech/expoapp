# Expoapp

Сервис предоставляет шаблон для разработки GRPC микросервисов.

## Команды для поднятия сервиса

```shell
# Копирует файл с переменными окружения
cp .env.example .env
```

```shell
# Собирает docker образ
make docker_build
```

```shell
# Запускает сервис в docker контейнере
make up
```

```shell
# Отправляет запрос на proto.ExpoService/GetInfoV1
grpcurl \
  -plaintext \
  -d '{}' \
  localhost:50051 \
  proto.ExpoService/GetInfoV1
```