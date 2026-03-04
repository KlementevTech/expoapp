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
task docker-build
```

```shell
# Запускает сервис в docker контейнере
task compose-up
```

```shell
# Устанавливает grpcurl
task install-grpcurl
```

```shell
# Отправляет запрос на proto.ExpoService/GetVersion
grpcurl \
  -plaintext \
  -d '{}' \
  localhost:50051 \
  expo.v1.ExpoService/GetVersion
```

```shell
grpcurl -plaintext \
    -d '{"id": "019cbea4-e202-731f-9206-966c34048cae", "name": "Turbo"}' \
    localhost:50051 \
    expo.v1.ExpoService/CreatePart
```

```shell
# Отправляет запрос на api/version
curl -i \
  -X GET \
  "http://localhost:8080/api/version"
```

[Pprof server](http://localhost:6060/debug/pprof/)

```shell
# Запускает нагрузочный тест
task test-perf
```

```shell
# Останавливает контейнер
task compose-stop
```