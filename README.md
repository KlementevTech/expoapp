# Expo gRPC Service

Сервис предоставляет шаблон для разработки gRPC сервисов.

Для работы с сервисом вам необходимо установить [Task](https://taskfile.dev/).

```shell
# Список доступных команд в Taskfile
task
```

```shell
# Копирует файл с переменными окружения
task cp-env
```

```shell
# Собирает docker образ
task build-docker
```

```shell
# Запускает сервис в docker контейнере
task compose-up
```

```shell
# Отправляет запрос на proto.ExpoService/GetInfoV1
    ./bin/grpcurl \
    -plaintext \
    -d '{}' \
    localhost:50051 \
    expo.v1.ExpoService/GetVersion
```

```shell
# Запускает нагрузочный тест
task test-perf
```