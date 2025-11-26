# Task Service

Микросервис для управления задачами. Часть системы управления проектами.

📋 Содержание

- [Описание](#описание)
- [Технологии](#технологии)
- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [API Reference](#api-reference)
- [Аутентификация](#аутентификация)
- [Модели данных](#модели-данных)
- [Примеры запросов](#примеры-запросов)

---

## Описание

Task Service предоставляет REST API для:
- Создания, получения и удаления задач
- Назначения исполнителей
- Управления статусами задач
- Контроля доступа на основе ролей (создатель/исполнитель)

### Бизнес-логика

| Действие | Кто может выполнить |
|----------|---------------------|
| Создать задачу | Любой авторизованный пользователь |
| Просмотреть задачу | Любой авторизованный пользователь |
| Назначить исполнителя | Только создатель задачи |
| Сменить исполнителя | Только создатель задачи |
| Изменить статус | Создатель или исполнитель |
| Удалить задачу | Только создатель задачи |

---

## Технологии

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.23+ |
| Веб-фреймворк | Gin |
| ORM | GORM |
| База данных | PostgreSQL 15 |
| Аутентификация | JWT (HS256) |
| Валидация | go-playground/validator |
| Контейнеризация | Docker, Docker Compose |

---

## Быстрый старт

### Требования

- Docker и Docker Compose
- Go 1.23+ (для локальной разработки)

### Запуск через Docker

```bash
# Клонировать репозиторий
git clone <repository-url>
cd hc-task-service

# Запустить сервисы
docker-compose up -d

# Проверить статус
docker-compose ps

# Проверить работоспособность
curl http://localhost:8082/health
```

### Запуск локально

```bash
# Запустить только БД
docker-compose up -d task-db

# Запустить сервис
go run ./cmd/server
```

### Остановка

```bash
docker-compose down

# С удалением данных
docker-compose down -v
```

---

## Конфигурация

Сервис конфигурируется через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | Порт сервиса | `8082` |
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `password` |
| `DB_NAME` | Имя базы данных | `taskdb` |
| `JWT_SECRET` | Секрет для подписи JWT | `supersecretkey` |

### Пример .env файла

PORT=8082

DB_HOST=localhost 

DB_PORT=5432

DB_USER=postgres

DB_PASSWORD=mypassword

DB_NAME=taskdb

JWT_SECRET=your-super-secret-key-change-in-production


---

## API Reference

**Base URL:** `http://localhost:8082`

### Endpoints

| Метод | Endpoint | Описание | Авторизация |
|-------|----------|----------|-------------|
| GET | `/health` | Проверка работоспособности | ❌ Нет |
| POST | `/api/v1/tasks` | Создать задачу | ✅ Да |
| GET | `/api/v1/tasks` | Получить все задачи пользователя | ✅ Да |
| GET | `/api/v1/tasks/:id` | Получить задачу по ID | ✅ Да |
| PUT | `/api/v1/tasks/:id/assign` | Назначить исполнителя | ✅ Да |
| PUT | `/api/v1/tasks/:id/assignee` | Сменить исполнителя | ✅ Да |
| PUT | `/api/v1/tasks/:id/status` | Изменить статус | ✅ Да |
| DELETE | `/api/v1/tasks/:id` | Удалить задачу | ✅ Да |

---

### GET /health

Проверка работоспособности сервиса.

**Ответ:**
```json
{
    "status": "ok"
}
```

---

### POST /api/v1/tasks

Создание новой задачи.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "title": "string (обязательно, 3-100 символов)",
    "description": "string (опционально, до 1000 символов)",
    "team_id": "uuid (обязательно)",
    "assignee_id": "uuid (опционально)"
}
```

**Response (201 Created):**
```json
{
    "id": "uuid",
    "title": "string",
    "description": "string",
    "status": "created",
    "assignee_id": "uuid | null",
    "creator_id": "uuid",
    "team_id": "uuid",
    "created_at": "2025-11-26T12:00:00Z",
    "updated_at": "2025-11-26T12:00:00Z"
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 400 | Невалидные данные |
| 401 | Не авторизован |

---

### GET /api/v1/tasks

Получение всех задач текущего пользователя.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "tasks": [
        {
            "id": "uuid",
            "title": "string",
            "description": "string",
            "status": "created | in_process | done",
            "assignee_id": "uuid | null",
            "creator_id": "uuid",
            "team_id": "uuid",
            "created_at": "2025-11-26T12:00:00Z",
            "updated_at": "2025-11-26T12:00:00Z"
        }
    ],
    "count": 1
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 401 | Не авторизован |

---

### GET /api/v1/tasks/:id

Получение задачи по ID.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "id": "uuid",
    "title": "string",
    "description": "string",
    "status": "created | in_process | done",
    "assignee_id": "uuid | null",
    "creator_id": "uuid",
    "team_id": "uuid",
    "created_at": "2025-11-26T12:00:00Z",
    "updated_at": "2025-11-26T12:00:00Z"
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 400 | Невалидный ID |
| 401 | Не авторизован |
| 404 | Задача не найдена |

---

### PUT /api/v1/tasks/:id/assign

Назначение исполнителя задачи. Только создатель.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "assignee_id": "uuid (обязательно)"
}
```

**Response (200 OK):**
```json
{
    "message": "assigned"
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 400 | Невалидные данные |
| 401 | Не авторизован |
| 403 | Нет прав (не создатель) |
| 404 | Задача не найдена |

---

### PUT /api/v1/tasks/:id/assignee

Смена исполнителя задачи. Только создатель.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "assignee_id": "uuid (обязательно)"
}
```

**Response (200 OK):**
```json
{
    "message": "assignee changed"
}
```

---

### PUT /api/v1/tasks/:id/status

Изменение статуса задачи. Создатель или исполнитель.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "status": "created | in_process | done"
}
```

**Response (200 OK):**
```json
{
    "message": "status updated",
    "status": "in_process"
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 400 | Невалидный статус |
| 401 | Не авторизован |
| 403 | Нет прав (не создатель и не исполнитель) |
| 404 | Задача не найдена |

---

### DELETE /api/v1/tasks/:id

Удаление задачи. Только создатель.

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
    "message": "task deleted"
}
```

**Ошибки:**
| Код | Описание |
|-----|----------|
| 400 | Невалидный ID |
| 401 | Не авторизован |
| 403 | Нет прав (не создатель) |
| 404 | Задача не найдена |

---

## Аутентификация

Сервис использует JWT (JSON Web Token) для аутентификации.

### Формат токена

```
Header: Authorization: Bearer <jwt_token>
```

### Структура JWT payload

```json
{
    "user_id": "22222222-2222-2222-2222-222222222222",
    "exp": 1764260241,
    "iat": 1764173841
}
```

### Генерация тестового токена

```bash
go run scripts/generate_token.go
```

### Пример токена для тестирования

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQyNjAyNDEsImlhdCI6MTc2NDE3Mzg0MSwidXNlcl9pZCI6IjIyMjIyMjIyLTIyMjItMjIyMi0yMjIyLTIyMjIyMjIyMjIyMiJ9.2nK0BXVa9Wpq7TUdZSspDCdKCtGhWt2TPrKBobldBTE
```

**User ID:** `22222222-2222-2222-2222-222222222222`

---

## Модели данных

### Task

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | UUID | Уникальный идентификатор |
| `title` | string | Название задачи (3-100 символов) |
| `description` | string | Описание задачи (до 1000 символов) |
| `status` | enum | Статус: `created`, `in_process`, `done` |
| `assignee_id` | UUID | ID исполнителя (может быть null) |
| `creator_id` | UUID | ID создателя задачи |
| `team_id` | UUID | ID команды |
| `created_at` | datetime | Дата создания |
| `updated_at` | datetime | Дата обновления |

### TaskStatus

| Значение | Описание |
|----------|----------|
| `created` | Задача создана |
| `in_process` | Задача в работе |
| `done` | Задача выполнена |

---

## Примеры запросов

### cURL

```bash
# Health check
curl http://localhost:8082/health

# Получить все задачи пользователя
curl http://localhost:8082/api/v1/tasks \
  -H "Authorization: Bearer <token>"

# Создать задачу
curl -X POST http://localhost:8082/api/v1/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Моя задача",
    "description": "Описание",
    "team_id": "11111111-1111-1111-1111-111111111111"
  }'

# Получить задачу
curl http://localhost:8082/api/v1/tasks/<task_id> \
  -H "Authorization: Bearer <token>"

# Назначить исполнителя
curl -X PUT http://localhost:8082/api/v1/tasks/<task_id>/assign \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"assignee_id": "33333333-3333-3333-3333-333333333333"}'

# Изменить статус
curl -X PUT http://localhost:8082/api/v1/tasks/<task_id>/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_process"}'

# Удалить задачу
curl -X DELETE http://localhost:8082/api/v1/tasks/<task_id> \
  -H "Authorization: Bearer <token>"
```

---


---

## Тестовые данные

При первом запуске создаётся тестовая команда:

| ID | Название |
|----|----------|
| `11111111-1111-1111-1111-111111111111` | Моя команда |

---

## License

MIT
