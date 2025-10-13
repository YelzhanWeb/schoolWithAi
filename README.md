# Education Platform - Setup Instructions

Образовательная платформа с персонализированным обучением через ИИ и геймификацией.

**Архитектура:** Hexagonal Architecture (Ports & Adapters)
**Backend:** Go + PostgreSQL
**ML Service:** Python (FastAPI)
**Frontend:** Vanilla JavaScript

## 📁 Структура проекта

```
education-platform/
├── backend/                    # Go приложение
│   ├── cmd/
│   │   └── main.go            # Точка входа
│   ├── internal/
│   │   ├── domain/            # 🟢 Бизнес-логика (models, services)
│   │   ├── ports/             # 🟡 Интерфейсы (repositories, services)
│   │   ├── adapters/          # 🔵 Реализации (postgres, http, redis)
│   │   ├── config/
│   │   └── pkg/
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
│
├── ml-service/                # Python ML сервис
├── frontend/                  # Vanilla JS Frontend
├── data/                      # Seed data
├── docker-compose.yml
├── Makefile
└── .env
```

## 🚀 Быстрый старт

### Требования:
- Go 1.21+
- PostgreSQL 15+ (или Docker)
- Python 3.10+ (для ML сервиса)
- Make (опционально, для удобных команд)

### 1. Клонирование и настройка

```bash
# Клонировать репозиторий
git clone <your-repo>
cd education-platform

# Скопировать .env файл
cp .env.example .env

# Отредактировать .env (если нужно)
nano .env
```

### 2. Установка зависимостей

```bash
# Через Makefile (рекомендуется)
make install

# Или вручную
cd backend && go mod download && go mod tidy
cd ../ml-service && pip install -r requirements.txt
```

### 3. Запуск базы данных

**Вариант A: Docker (рекомендуется)**

```bash
# Запустить PostgreSQL и pgAdmin
make docker-up

# Или напрямую
docker-compose up -d
```

**Доступ:**
- PostgreSQL: `localhost:5432`
- pgAdmin: `http://localhost:5050`
  - Email: `admin@education.com`
  - Password: `admin`

**Вариант B: Локальный PostgreSQL**

```bash
# Создать базу данных
createdb education_platform

# Обновить .env с настройками
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
```

### 4. Миграции базы данных

```bash
# Запустить миграции
make migrate-up

# Проверить статус
make status
```

### 5. Запуск сервисов

```bash
# Терминал 1: Backend
make run-backend
# Сервер запустится на http://localhost:8080

# Терминал 2: ML Service (позже)
make run-ml
# ML API на http://localhost:5000

# Терминал 3: Frontend (позже)
make run-frontend
# Frontend на http://localhost:3000
```

## 🎯 Разработка

### Шаг 1: Миграции (✅ Готово)

Миграции созданы и готовы к использованию:
- `000001_init_schema.up.sql` - создание всех таблиц
- `000001_init_schema.down.sql` - откат изменений

**Команды:**
```bash
make migrate-up       # Применить миграции
make migrate-down     # Откатить миграции
make migrate-create name=add_feature  # Создать новую миграцию
```

### Шаг 2: Seed Data (⏳ Следующий)

Создать примеры образовательного контента:
- 10-20 курсов по разным предметам
- 50-100 упражнений на критическое мышление
- Примеры тестов и квизов

### Шаг 3: Backend API (⏳ Следующий)

Реализация REST API endpoints:
```
POST   /api/auth/register
POST   /api/auth/login
GET    /api/courses
GET    /api/courses/:id
POST   /api/courses (teacher only)
GET    /api/recommendations (student) сервис (создадим позже)
│
├── frontend/                  # Frontend (создадим позже)
│
├── docker-compose.yml         # Docker окружение
├── .env                       # Переменные окружения
└── README.md                  # Эта инструкция
```

## 🚀 Шаг 1: Установка зависимостей

### Требования:
- Go 1.21+
- PostgreSQL 15+ (или Docker)
- Git

### Установка Go зависимостей:

```bash
cd backend
go mod download
go mod tidy
```

## 🐳 Шаг 2: Запуск базы данных

### Вариант A: Через Docker (рекомендуется)

```bash
# Запустить PostgreSQL и pgAdmin
docker-compose up -d postgres pgadmin

# Проверить статус
docker-compose ps
```

**Доступ к pgAdmin:**
-