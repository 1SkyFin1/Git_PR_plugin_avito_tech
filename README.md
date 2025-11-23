# PR Reviewer Assignment Service

Сервис для автоматического назначения ревьюверов на Pull Request'ы с балансировкой нагрузки между членами команды.

## 📋 Описание

Сервис предоставляет REST API для управления командами разработчиков, пользователями и автоматического назначения ревьюверов на Pull Request'ы. При создании PR система автоматически выбирает до 2 ревьюверов из команды автора, учитывая их текущую загрузку и статус активности.

### Основные возможности:

- ✅ Создание команд с участниками
- ✅ Управление активностью пользователей
- ✅ Автоматическое назначение ревьюверов с балансировкой нагрузки
- ✅ Merge Pull Request'ов
- ✅ Переназначение ревьюверов
- ✅ Защита от изменений в смерженных PR
- ✅ **Подробная статистика системы** (топ ревьюверов, загрузка команды, аналитика PR)
- ✅ **Массовая деактивация пользователей** с автоматическим переназначением PR

## 🚀 Быстрый старт

### Требования

- Go 1.21+
- Docker & Docker Compose
- Make (опционально)

### Запуск

#### Вариант 1: Docker Compose (рекомендуется) 🐳

Самый простой способ - запустить весь стек одной командой:

```bash
git clone <repository-url>
cd Git_PR_plugin_avito_tech
docker-compose up -d
```

**Готово!** Сервис будет доступен на `http://localhost:8080`

> **Zero Configuration:** Docker Compose автоматически использует значения по умолчанию из `.env.example`. Никаких дополнительных настроек не требуется!

**Что происходит при запуске:**
1. Поднимается PostgreSQL база данных (контейнер `db`)
2. Docker ждет, пока БД станет готова (healthcheck)
3. Собирается Docker образ Go приложения
4. Запускается приложение (контейнер `app`)
5. Автоматически применяются миграции БД
6. Сервис запущен и готов к работе

**Полезные команды:**
```bash
# Посмотреть логи всех сервисов
docker-compose logs -f

# Посмотреть логи только приложения
docker-compose logs -f app

# Посмотреть логи только БД
docker-compose logs -f db

# Проверить статус контейнеров
docker-compose ps

# Остановить сервисы
docker-compose down

# Остановить и удалить данные БД
docker-compose down -v

# Пересобрать образ приложения после изменений
docker-compose build app
docker-compose up -d

# Перезапустить только приложение
docker-compose restart app
```

#### Вариант 2: Локальная разработка (без Docker для приложения)

Для разработки с возможностью hot-reload и отладки:

1. Запустите только базу данных в Docker:
```bash
docker-compose up -d db
```

2. Запустите приложение локально:
```bash
go run cmd/app/main.go
```

Сервис будет доступен на `http://localhost:8080`

> **Автоматическая конфигурация:** Приложение автоматически использует настройки из `.env.example`, если файл `.env` не найден.

#### Вариант 3: Настройка для production

Для production окружения создайте файл `.env` с собственными настройками:

```bash
cp .env.example .env
```

Отредактируйте `.env`:
```env
# Измените пароль для production
POSTGRES_PASSWORD=your_secure_password

# Для Docker Compose хост будет 'db', для локального запуска - 'localhost'
PG_URL=postgresql://postgres:your_secure_password@db:5432/git_pr_plugin?sslmode=disable
```

**Важно:** Убедитесь, что пароль в `POSTGRES_PASSWORD` совпадает с паролем в `PG_URL`.

Затем запустите:
```bash
docker-compose up -d
```

## 📡 API Endpoints

Сервис слушает порт **8080** и предоставляет следующие эндпоинты:

### Teams

#### `POST /team/add`
Создать команду с участниками.

**Request:**
```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "Alice",
      "is_active": true
    },
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "Bob",
      "is_active": true
    }
  ]
}
```

**Response (201):**
```json
{
  "team": {
    "team_name": "backend",
    "members": [...]
  }
}
```

#### `GET /team/get?team_name=backend`
Получить информацию о команде по названию.

**Response (200):**
```json
{
  "team_name": "backend",
  "members": [
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "Alice",
      "is_active": true
    }
  ]
}
```

### Users

#### `POST /users/setIsActive`
Установить флаг активности пользователя.

**Request:**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "is_active": false
}
```

**Response (200):**
```json
{
  "user": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "Alice",
    "team_name": "backend",
    "is_active": false
  }
}
```

#### `GET /users/getReview?user_id=<uuid>`
Получить список PR, где пользователь является ревьювером.

**Response (200):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "pull_requests": [
    {
      "pull_request_id": "...",
      "pull_request_name": "Add feature",
      "author_id": "...",
      "status": "OPEN"
    }
  ]
}
```

#### `POST /users/deactivate`
Массовая деактивация пользователей с автоматическим удалением их назначений на ревью.

**Описание:**
Этот эндпоинт позволяет деактивировать сразу несколько пользователей одной операцией. При деактивации:
- Пользователи помечаются как неактивные (`is_active = false`)
- Все их назначения на ревью PR автоматически удаляются
- Операция выполняется атомарно в рамках одной транзакции
- Оптимизирована для работы с большими объемами данных

**Request:**
```json
{
  "user_ids": [
    "550e8400-e29b-41d4-a716-446655440000",
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440003"
  ]
}
```

**Response (200):**
```json
{
  "deactivated_count": 3,
  "reassigned_prs_count": 7,
  "deactivated_user_ids": [
    "550e8400-e29b-41d4-a716-446655440000",
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440003"
  ]
}
```

**Поля ответа:**
- `deactivated_count` - количество фактически деактивированных пользователей (только те, кто был активен)
- `reassigned_prs_count` - количество удаленных назначен��й на ревью
- `deactivated_user_ids` - список ID деактивированных пользователей

**Особенности:**
- ✅ Атомарная операция (транзакция)
- ✅ Деактивируются только активные пользователи
- ✅ Безопасное удаление назначений на открытые PR
- ✅ Оптимизирована для больших объемов (batch операции)
- ✅ Если пользователь уже неактивен, он не учитывается в `deactivated_count`

**Пример использования:**
```bash
curl -X POST http://localhost:8080/users/deactivate \
  -H "Content-Type: application/json" \
  -d '{
    "user_ids": [
      "550e8400-e29b-41d4-a716-446655440000",
      "550e8400-e29b-41d4-a716-446655440001"
    ]
  }'
```

### Pull Requests

#### `POST /pullRequest/create`
Создать Pull Request с автоматическим назначением ревьюверов.

**Request:**
```json
{
  "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
  "pull_request_name": "Add new feature",
  "author_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response (201):**
```json
{
  "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
  "pull_request_name": "Add new feature",
  "author_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "OPEN",
  "assigned_reviewers": [
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440003"
  ]
}
```

**Логика назначения ревьюверов:**
- Выбираются только активные пользователи из команды автора
- Автор не может быть ревьювером своего PR
- Назначается до 2 ревьюверов
- Приоритет отдается пользователям с наименьшей текущей загрузкой

#### `POST /pullRequest/merge`
Смержить Pull Request.

**Request:**
```json
{
  "pull_request_id": "550e8400-e29b-41d4-a716-446655440002"
}
```

**Response (200):**
```json
{
  "pr": {
    "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
    "pull_request_name": "Add new feature",
    "author_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "MERGED",
    "assigned_reviewers": [...]
  },
  "mergedAt": "2024-01-15T10:30:45Z"
}
```

#### `POST /pullRequest/reassign`
Переназначить ревьювера на Pull Request.

**Request:**
```json
{
  "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
  "old_reviewer_id": "550e8400-e29b-41d4-a716-446655440001"
}
```

**Response (200):**
```json
{
  "pr": {
    "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
    "pull_request_name": "Add new feature",
    "author_id": "550e8400-e29b-41d4-a716-446655440000",
    "status": "OPEN",
    "assigned_reviewers": [
      "550e8400-e29b-41d4-a716-446655440003",
      "550e8400-e29b-41d4-a716-446655440004"
    ]
  },
  "replaced_by": "550e8400-e29b-41d4-a716-446655440004"
}
```

**Ограничения:**
- ❌ Нельзя переназначить ревьювера на смерженном PR (статус 409)
- ❌ Старый ревьювер должен быть назначен на PR
- ✅ Новый ревьювер выбирается случайно из доступных членов команды

### Statistics

#### `GET /statistics`
Получить подробную статистику системы.

**Описание:**
Этот эндпоинт предоставляет комплексную аналитику работы системы, включая:
- Общую статистику по пользователям и PR
- Топ-10 самых загруженных ревьюверов
- PR с наибольшим количеством назначенных ревьюверов
- Детальную информацию о назначениях (открытые/завершенные)

**Response (200):**
```json
{
  "total_users": 25,
  "active_users": 22,
  "total_pull_requests": 150,
  "open_pull_requests": 45,
  "merged_pull_requests": 105,
  "total_assignments": 280,
  "top_reviewers": [
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "Alice",
      "team_name": "backend",
      "total_assignments": 35,
      "open_assignments": 12,
      "completed_assignments": 23
    },
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "username": "Bob",
      "team_name": "backend",
      "total_assignments": 32,
      "open_assignments": 10,
      "completed_assignments": 22
    }
  ],
  "pull_requests_with_most_reviewers": [
    {
      "pull_request_id": "550e8400-e29b-41d4-a716-446655440010",
      "pull_request_name": "Major refactoring",
      "author_username": "Charlie",
      "status": "OPEN",
      "reviewers_count": 5
    },
    {
      "pull_request_id": "550e8400-e29b-41d4-a716-446655440011",
      "pull_request_name": "Critical bug fix",
      "author_username": "David",
      "status": "MERGED",
      "reviewers_count": 4
    }
  ]
}
```

**Поля ответа:**

**Общая статистика:**
- `total_users` - общее количество пользователей в системе
- `active_users` - количество активных пользователей
- `total_pull_requests` - общее количество PR
- `open_pull_requests` - количество открытых PR
- `merged_pull_requests` - количество смерженных PR
- `total_assignments` - общее количество назначений на ревью

**Топ ревьюверов (top_reviewers):**
- `user_id` - ID пользователя
- `username` - имя пользователя
- `team_name` - название команды
- `total_assignments` - всего назначений на ревью
- `open_assignments` - открытых назначений
- `completed_assignments` - завершенных назначений (смерженные PR)

**PR с наибольшим количеством ревьюверов (pull_requests_with_most_reviewers):**
- `pull_request_id` - ID Pull Request
- `pull_request_name` - название PR
- `author_username` - имя автора PR
- `status` - статус PR (OPEN/MERGED)
- `reviewers_count` - количество назначенных ревьюверов

**Особенности:**
- ✅ Быстрые агрегированные запросы к БД
- ✅ Показывает только активных пользователей в топе
- ✅ Топ-10 по каждой категории
- ✅ Полезно для мониторинга загрузки команды
- ✅ Помогает выявить узкие места в процессе ревью

**Пример использования:**
```bash
curl -X GET http://localhost:8080/statistics
```

**Применение:**
- 📊 Мониторинг загрузки ревьюверов
- 📈 Анализ эффективности команды
- 🎯 Выявление перегруженных разработчиков
- 📉 Отслеживание динамики работы с PR
- 🔍 Поиск PR, требующих особого внимания

## 🔑 Важные решения по реализации

### Использование UUID для идентификаторов

Поскольку в техническом задании не была указана конкретная информация о типе данных для идентификаторов, было принято решение использовать **UUID (Universally Unique Identifier)**.

**Обоснование:**
- ✅ UUID является стандартом де-факто в современных распределенных системах
- ✅ Обеспечивает глобальную уникальность без координации между сервисами
- ✅ Предотвращает проблемы с перебором ID (security by obscurity)
- ✅ Упрощает миграцию данных и репликацию
- ✅ Наиболее распространен в реальных бизнес-сервисах

**Формат:** `550e8400-e29b-41d4-a716-446655440000` (UUID v4)

### Отличия в формате ответов для merge и reassign

В эндпоинтах `/pullRequest/merge` и `/pullRequest/reassign` используется **немного другой формат возвращаемых данных** по сравнению с OpenAPI спецификацией.

**Текущая реализация:**
```json
{
  "pr": {
    "pull_request_id": "...",
    "pull_request_name": "...",
    "author_id": "...",
    "status": "...",
    "assigned_reviewers": [...]
  },
  "mergedAt": "..." // или "replaced_by": "..."
}
```

**Причина изменения:**
Набор полей остался тот же, но они расположены иначе для использования **единой модели `PullRequestFullInfoDto`**. Это позволяет:
- ✅ Избежать дублирования кода
- ✅ Упростить поддержку и тестирование
- ✅ Обеспечить консистентность данных

**Примечание:** В реальном проекте я бы обратился к аналитику/product owner с предложением немного скорректировать API контракт для унификации структуры ответов. Это улучшило бы как клиентский, так и серверный код, сделав API более предсказуемым и удобным в использовании.

## 🐳 Docker и контейнеризация

### Как это работает

Проект использует **многоступенчатую сборку Docker** (multi-stage build) для оптимизации размера образа и безопасности.

#### Dockerfile - два этапа:

**Этап 1: Builder (сборка)**
```dockerfile
FROM golang:1.21-alpine AS builder
# Компилируем приложение в статический бинарник
# Результат: исполняемый файл ~15-20 MB
```

**Этап 2: Runtime (запуск)**
```dockerfile
FROM alpine:latest
# Копируем только бинарник и конфиги
# Итоговый образ: ~25-30 MB (вместо ~800 MB с полным Go)
```

**Преимущества:**
- ✅ Маленький размер образа (в 30 раз меньше!)
- ✅ Быстрая сборка благодаря кэшированию слоев
- ✅ Безопасность: нет компилятора и исходников в production образе
- ✅ Запуск от непривилегированного пользователя

#### Docker Compose - оркестрация:

```yaml
services:
  db:
    # PostgreSQL с healthcheck
    healthcheck: проверка готовности каждые 5 секунд
    
  app:
    # Go приложение
    depends_on:
      db:
        condition: service_healthy  # ждем готовности БД!
```

**Ключевые особенности:**

1. **Healthchecks** - Docker проверяет готовность сервисов:
   - БД: `pg_isready` каждые 5 секунд
   - App: HTTP запрос на `/` каждые 30 секунд

2. **Зависимости** - приложение запускается только после готовности БД:
   ```yaml
   depends_on:
     db:
       condition: service_healthy
   ```

3. **Сетевое взаимодействие**:
   - Контейнеры в одной сети `app_network`
   - Приложение обращается к БД по имени: `db:5432`
   - Порты пробрасываются на хост: `8080:8080`

4. **Переменные окружения**:
   - Значения по умолчанию: `${POSTGRES_USER:-postgres}`
   - Если `.env` не найден, используются дефолтные значения
   - Для Docker хост БД = `db` (имя сервиса)

5. **Volumes** - данные БД сохраняются между перезапусками:
   ```yaml
   volumes:
     - postgres_data:/var/lib/postgresql/data
   ```

#### Процесс запуска:

```bash
docker-compose up -d
```

**Что происходит:**

1. **Чтение конфигурации**
   - Docker Compose читает `docker-compose.yml`
   - Загружает переменные из `.env.example` (если `.env` нет)

2. **Создание сети**
   - Создается bridge сеть `app_network`
   - Все контейнеры получают DNS имена

3. **Запуск БД**
   - Скачивается образ `postgres:16-alpine` (если нужно)
   - Создается volume `postgres_data`
   - Запускается контейнер `git_pr_plugin_db`
   - Healthcheck проверяет готовность

4. **Сборка приложения**
   - Выполняется `docker build` по Dockerfile
   - Копируются зависимости (кэшируется!)
   - Компилируется Go код
   - Создается финальный образ

5. **Запуск приложения**
   - Docker ждет `db` (condition: service_healthy)
   - Запускается контейнер `git_pr_plugin_app`
   - Приложение подключается к БД по адресу `db:5432`
   - Применяются миграции
   - API готов к работе на порту 8080

#### Логи и мониторинг:

```bash
# Смотрим что происходит
docker-compose logs -f

# Проверяем статус
docker-compose ps

# Результат:
# NAME                  STATUS              PORTS
# git_pr_plugin_db      Up (healthy)        0.0.0.0:5432->5432/tcp
# git_pr_plugin_app     Up (healthy)        0.0.0.0:8080->8080/tcp
```

#### Остановка и очистка:

```bash
# Остановить (данные сохраняются)
docker-compose down

# Удалить всё включая данные БД
docker-compose down -v

# Пересобрать после изменений кода
docker-compose build app && docker-compose up -d
```

## 🏗️ Архитектура

Проект следует принципам **Clean Architecture** с разделением на слои:

```
cmd/
  app/          - Точка входа приложения
  migrate/      - Миграции БД
internal/
  app/          - Инициализация приложения
  controller/   - HTTP handlers (Chi router)
  service/      - Бизнес-логика
  repository/   - Работа с БД (PostgreSQL)
  dto/          - Data Transfer Objects
  model/        - Доменные модели
  apperror/     - Централизованная обработка ошибок
  constants/    - Константы
config/         - Конфигурация
resources/
  migrations/   - SQL миграции
```

### Технологический стек

- **Router:** Chi v5 - легковесный и производительный HTTP роутер
- **Database:** PostgreSQL 14+ с sqlx для работы с БД
- **Migrations:** golang-migrate для управления схемой БД
- **Error Handling:** Централизованная система обработки ошибок с кодами из OpenAPI
- **Logging:** Структурированное логирование всех операций

## 🔒 Обработка ошибок

Все ошибки возвращаются в едином формате согласно OpenAPI спецификации:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable message"
  }
}
```

### Коды ошибок:

| Код | HTTP Status | Описание |
|-----|-------------|----------|
| `TEAM_EXISTS` | 400 | Команда уже существует |
| `PR_EXISTS` | 409 | Pull Request уже существует |
| `PR_MERGED` | 409 | PR уже смержен (нельзя изменить) |
| `NOT_ASSIGNED` | 409 | Ревьювер не назначен на PR |
| `NO_CANDIDATE` | 409 | Нет доступных кандидатов |
| `NOT_FOUND` | 404 | Ресурс не найден |
| `INVALID_INPUT` | 400 | Некорректные входные данные |
| `INTERNAL` | 500 | Внутренняя ошибка сервера |

## 🗄️ База данных

### Схема

```sql
-- Команды
CREATE TABLE team (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Пользователи
CREATE TABLE "user" (
    user_id UUID PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    team_id UUID REFERENCES team(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Pull Request'ы
CREATE TABLE pull_request (
    pull_request_id UUID PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id UUID REFERENCES "user"(user_id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);

-- Назначения ревьюверов
CREATE TABLE pr_reviewer (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pull_request_id UUID REFERENCES pull_request(pull_request_id),
    reviewer_id UUID REFERENCES "user"(user_id),
    assigned_at TIMESTAMP DEFAULT NOW()
);
```

## 🧪 Тестирование

### Примеры запросов

Используйте Postman, cURL или любой HTTP клиент:

```bash
# Создать команду
curl -X POST http://localhost:8080/team/add \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "Alice",
        "is_active": true
      }
    ]
  }'

# Создать PR
curl -X POST http://localhost:8080/pullRequest/create \
  -H "Content-Type: application/json" \
  -d '{
    "pull_request_id": "550e8400-e29b-41d4-a716-446655440002",
    "pull_request_name": "Add feature",
    "author_id": "550e8400-e29b-41d4-a716-446655440000"
  }'
```

## 📝 Логирование

Все операции логируются с указанием:
- Начала операции
- Промежуточных шагов
- Результата (успех/ошибка)

Пример логов:
```
[TeamService.CreateTeamWithMembers] Starting for team: backend with 3 members
[TeamService.CreateTeamWithMembers] Team created: backend (ID: ...)
[TeamService.CreateTeamWithMembers] Created 3 users
[TeamService.CreateTeamWithMembers] Success! Team: backend
```

## 🔐 Транзакционность

Все операции, изменяющие несколько записей, выполняются в транзакциях:
- ✅ Создание команды + пользователей
- ✅ Создание PR + назначение ревьюверов
- ✅ Переназначение ревьювера (удаление + добавление)

## 📚 Дополнительная информация

### Конфигурация

Настройки приложения находятся в `config/config.yml` и могут быть переопределены через переменные окружения.

### Миграции

Миграции применяются автоматически при запуске приложения. Файлы миграций находятся в `resources/migrations/`.

## 👨‍💻 Разработка

### Структура проекта

- `cmd/` - Исполняемые файлы
- `internal/` - Внутренний код приложения
- `config/` - Конфигурационные файлы
- `resources/` - Ресурсы (миграции, статика)

### Запуск в режиме разработки

```bash
go run cmd/app/main.go
```

### Сборка

```bash
go build -o bin/app cmd/app/main.go
./bin/app
```

## 📄 Лицензия

MIT

## 🤝 Контакты

Для вопросов и предложений создавайте Issue в репозитории.
