# Отчет по домашнему заданию к лекции 7. Docker

**Выполнил:**  Метельский Павел  
**Дата выполнения:** 2 июня 2025  
**Сервер:** Ubuntu 24.04.1 LTS (37.9.53.237)

- **Репозиторий с кодом:** https://github.com/PavelMetelsky/HW-lesta/blob/master/hw-Docker/report.md
- **Рабочая директория:** ~/hw-docker

## 1. Задание 1: Multi-stage сборка простого Go приложения
![alt text](telegram-cloud-photo-size-2-5328065789099635107-y.jpg)

### Описание задания
При помощи multi-stage собрать Docker образ простого web приложения на Go с тегом `lesta-start:7.1`, работающего на порту 8080.

### Структура файлов
```
~/hw-docker/task1/
├── main.go
└── Dockerfile
```

### Исходный код приложения

<details>
<summary>main.go</summary>

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello World!")
    })
    http.ListenAndServe(":8080", nil)
}
```
</details>

### Multi-stage Dockerfile

<details>
<summary>Dockerfile</summary>

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY main.go .

RUN go mod init task1
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```
</details>

### Сборка и запуск

```bash
# Сборка образа
docker build -t lesta-start:7.1 .

# Запуск контейнера
docker run -d -p 8080:8080 --name task1-app lesta-start:7.1

# Проверка работы
curl http://localhost:8080
# Вывод: Hello World!
```

### Результаты
- ✅ Образу присвоен тег `lesta-start:7.1`
- ✅ Приложение работает на порту 8080

## 2. Задание 2: Docker Compose для Go API с PostgreSQL и Nginx

![alt text](telegram-cloud-photo-size-2-5328065789099635165-y.jpg)

### Описание задания
Подготовить сборку и запуск Go API приложения при помощи Docker Compose с PostgreSQL и Nginx.

### Структура проекта
```
~/hw-docker/task2/
├── Dockerfile
├── docker-compose.yml
├── nginx/
│   └── nginx.conf
├── api.tar.gz          # Исходный архив
├── cmd/                # Распакованные файлы
│   └── demo/
│       └── main.go
├── internal/
├── docs/
├── go.mod
└── go.sum
```

### Docker конфигурация

<details>
<summary>Dockerfile</summary>

```dockerfile
FROM golang:1.21.5 as builder

WORKDIR /go-app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd /go-app/cmd
COPY internal /go-app/internal
COPY docs /go-app/docs

RUN cd /go-app/cmd/demo/ && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/web ./main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates

ENV API_PORT="${API_PORT}" \
    DB_URL="${DB_URL}"

COPY --from=builder /bin/web /bin/web

RUN chmod +x /bin/web

CMD ["/bin/sh", "-c", "/bin/web"]
```
</details>

### Docker Compose конфигурация

<details>
<summary>docker-compose.yml</summary>

```yaml
version: '3.8'

services:
  postgres_db:
    image: postgres:16.1-alpine3.19
    container_name: postgres_db
    environment:
      POSTGRES_DB: DB
      POSTGRES_USER: USER_DB
      POSTGRES_PASSWORD: PWD_DB
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U USER_DB -d DB"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - app-network

  lesta_start_demo:
    build: .
    container_name: lesta_start_demo
    environment:
      DB_URL: postgres://USER_DB:PWD_DB@postgres_db:5432/DB
      API_PORT: 8887
    ports:
      - "8887:8887"
    depends_on:
      postgres_db:
        condition: service_healthy
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    container_name: nginx
    ports:
      - "8088:8088"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
    depends_on:
      - lesta_start_demo
    networks:
      - app-network

volumes:
  postgres_data:

networks:
  app-network:
    driver: bridge
```
</details>

### Nginx конфигурация

<details>
<summary>nginx/nginx.conf</summary>

```nginx
events {
    worker_connections 1024;
}

http {
    upstream backend {
        server lesta_start_demo:8887;
    }
    
    server {
        listen 8088;
        
        location / {
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
    }
}
```
</details>

### Запуск и проверка

```bash
# Запуск всех сервисов
docker-compose up -d

# Проверка статуса контейнеров
docker-compose ps

# Вывод:
#       Name                   Command                State                 Ports           
# ------------------------------------------------------------------------------------------
# lesta_start_demo   /bin/sh -c /bin/web           Up      0.0.0.0:8887->8887/tcp
# nginx              /docker-entrypoint.sh nginx   Up      0.0.0.0:8088->8088/tcp
# postgres_db        docker-entrypoint.sh postgres Up      0.0.0.0:5432->5432/tcp
```

### Тестирование API

```bash
# Проверка health endpoint
curl http://localhost:8088/ping
# Ответ: {"result":"I am alive"}

# POST запрос - бросок кости
curl -X POST http://localhost:8088/api/v1/roll_dice \
  -H "Content-Type: application/json" \
  -d '{}'
# Ответ: {"result":6,"error":null}

# GET запрос - статистика
curl http://localhost:8088/api/v1/roll_statistic
# Ответ: {"result":{"all_count_roll":1,"all_awg_roll":6},"error":null}

# Swagger UI доступен по адресу:
# http://localhost:8088/swagger/index.html
```

### Архитектура решения

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
┌──────▼──────┐
│    Nginx    │ :8088
└──────┬──────┘
       │ Proxy
┌──────▼──────┐
│   Go API    │ :8887
└──────┬──────┘
       │ DB Connection
┌──────▼──────┐
│ PostgreSQL  │ :5432
└─────────────┘
```

### Особенности реализации

1. **Multi-stage build** для минимизации размера образа
2. **Health check** для PostgreSQL обеспечивает правильный порядок запуска
3. **Environment variables** для конфигурации:
   - `DB_URL`: connection string для БД
   - `API_PORT`: порт приложения
4. **Nginx** как reverse proxy для API
5. **Docker networks** для изоляции сервисов

### Проверка артефакта

```bash
# Проверка бинарного файла в контейнере
docker-compose exec lesta_start_demo ls -la /bin/web
# -rwxr-xr-x    1 root     root      37391561 Jun  2 15:09 /bin/web
```

## Результаты
- ✅ Dockerfile с multi-stage сборкой
- ✅ Docker Compose с тремя сервисами (backend, PostgreSQL, Nginx)
- ✅ PostgreSQL с заданными ENV переменными (DB, USER_DB, PWD_DB)
- ✅ Health check для контроля статуса БД
- ✅ Передача порта и connection string через ENV
- ✅ Nginx конфигурация с проксированием на backend
- ✅ Правильный порядок запуска сервисов через depends_on
- ✅ Swagger UI доступен по адресу http://localhost:8088/swagger/index.html
- ✅ API методы работают корректно (POST /api/v1/roll_dice, GET /api/v1/roll_statistic)

## Команды для остановки

```bash
# Остановка всех сервисов
docker-compose down

# Остановка с удалением volumes
docker-compose down -v

# Остановка контейнера из задания 1
docker stop task1-app
docker rm task1-app
```