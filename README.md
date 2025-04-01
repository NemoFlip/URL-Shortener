# URL Shortener

Этот проект — сервис сокращения URL, написанный на Go. Он позволяет пользователям сокращать длинные ссылки и получать удобные короткие URL.

## Возможности
- Генерация коротких ссылок
- Перенаправление по короткому URL
- API для работы с сервисом
- Документация Swagger
- Поддержка Docker и Docker Compose
- Логирование и конфигурация через YAML

## Установка и запуск
### Через Docker Compose
```sh
docker-compose up --build
```

### Локальный запуск
1. Установите зависимости:
```sh
go mod tidy
```
2. Запустите приложение:
```sh
go run cmd/app/main.go
```

## Конфигурация
Конфигурация хранится в файлах `config/local.yaml` и `config/prod.yaml`. Можно указать настройки базы данных, порта сервера и другие параметры.

## Тестовый .env файл
```env
CONFIG_PATH="./config/local.yaml"
DB_USER=admin
DB_PASSWORD=admin
DB_URL=postgres://admin:admin@localhost:5432/postgres?sslmode=disable
DB_PORT=5432
DB_NAME=postgres
```

## API
Документация API доступна по адресу `/swagger/index.html` при запущенном сервере.

## Деплой
Проект поддерживает деплой через systemd. Пример сервиса:
```
[Unit]
Description=URL Shortener Service
After=network.target

[Service]
ExecStart=/path/to/binary
Restart=always
User=your-user

[Install]
WantedBy=multi-user.target
```

## Лицензия
Этот проект распространяется под лицензией MIT.

