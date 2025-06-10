# Flask API Exam Project

REST API приложение на Flask с PostgreSQL для экзаменационной работы.

**Автор:** Метельский Павел  
**Дата:** 10 июня 2025  
**Сервер:** Ubuntu 24.04.1 LTS (37.9.53.237)

## Функциональность

- **POST /submit** - сохранение данных (name, score) в БД
- **GET /results** - получение всех записей из БД
- **GET /ping** - health check endpoint

## Endpoint для проверки

```
http://37.9.53.237:5000/results
```

## Запуск проекта локально

### Предварительные требования

- Docker
- Docker Compose
- Git

### Шаги запуска

1. Клонирование репозитория:
```bash
git clone <repository-url>
cd flask-api-exam
```

2. Создание .env файла:
```bash
cp .env.example .env
```

3. Запуск через docker-compose:
```bash
docker-compose up -d
```

4. Проверка работы:
```bash
# Health check
curl http://localhost:5000/ping

# Отправка данных
curl -X POST http://localhost:5000/submit \
  -H "Content-Type: application/json" \
  -d '{"name": "Kirill", "score": 88}'

# Получение результатов
curl http://localhost:5000/results
```

## Настройка Jenkins

### 1. Установка необходимых плагинов

- SSH Agent Plugin
- Docker Pipeline Plugin
- Pipeline Plugin

### 2. Настройка credentials

1. Добавить SSH ключ для доступа к удаленному серверу:
   - Manage Jenkins → Manage Credentials
   - Add Credentials → SSH Username with private key
   - ID: `ssh-key-id`

### 3. Создание Pipeline Job

1. New Item → Pipeline
2. В разделе Pipeline выбрать "Pipeline script from SCM"
3. Указать репозиторий и путь к Jenkinsfile

### 4. Настройка Jenkins для Docker

На Jenkins сервере выполнить:
```bash
# Добавить пользователя jenkins в группу docker
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
```

## CI/CD Pipeline

### Стадии pipeline:

1. **Checkout** - клонирование репозитория
2. **Build** - сборка Docker образа
3. **Test/Lint** - проверка кода с помощью flake8
4. **Push** - сохранение образа
5. **Deploy** - развертывание на удаленном сервере

### Автоматическое развертывание

При успешном прохождении всех стадий приложение автоматически развертывается на сервере 37.9.53.237.

## API Endpoints

### POST /submit

Сохранение данных в БД.

**Request:**
```bash
curl -X POST http://localhost:5000/submit \
  -H "Content-Type: application/json" \
  -d '{"name": "Kirill", "score": 88}'
```

**Response:**
```json
{
  "message": "Data saved successfully",
  "id": 1
}
```

### GET /ping

Health check endpoint.

**Request:**
```bash
curl http://localhost:5000/ping
```

**Response:**
```json
{"status": "ok"}
```

### GET /results

Получение всех записей из БД.

**Request:**
```bash
curl http://localhost:5000/results
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Kirill",
    "score": 88,
    "timestamp": "2025-06-10T10:25:43"
  },
  {
    "id": 2,
    "name": "Roman",
    "score": 92,
    "timestamp": "2025-06-10T10:30:10"
  }
]
```

## Структура базы данных

### Таблица: results

| Поле | Тип | Описание |
|------|-----|----------|
| id | INTEGER | Первичный ключ, автоинкремент |
| name | VARCHAR(100) | Имя пользователя |
| score | INTEGER | Количество баллов |
| timestamp | DATETIME | Время создания записи |

## Troubleshooting

### Проблема с правами Docker в Jenkins

```bash
sudo usermod -aG docker jenkins
sudo systemctl restart jenkins
```

### Проблема с подключением к БД

Проверить переменные окружения в .env файле и убедиться, что PostgreSQL контейнер запущен:
```bash
docker-compose ps
docker-compose logs db
```

### Очистка Docker

```bash
docker-compose down -v
docker system prune -a
```