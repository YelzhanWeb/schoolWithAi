# 🏗️ Education Platform - Hexagonal Architecture

## 📂 Полная структура проекта

```
education-platform/
│
├── backend/                           # Go Backend
│   ├── cmd/
│   │   └── main.go                   # Entry point, dependency injection
│   │
│   ├── internal/
│   │   │
│   │   ├── domain/                   # 🟢 DOMAIN LAYER (Core Business Logic)
│   │   │   ├── models/              # Domain entities
│   │   │   │   ├── user.go          # User entity
│   │   │   │   ├── course.go        # Course entity
│   │   │   │   ├── resource.go      # Resource entity
│   │   │   │   ├── progress.go      # Progress tracking
│   │   │   │   ├── recommendation.go # Recommendation entity
│   │   │   │   ├── achievement.go   # Achievement entity
│   │   │   │   └── rating.go        # Rating entity
│   │   │   │
│   │   │   ├── services/            # Use cases (business rules)
│   │   │   │   ├── auth_service.go        # Authentication logic
│   │   │   │   ├── student_service.go     # Student operations
│   │   │   │   ├── teacher_service.go     # Teacher operations
│   │   │   │   ├── course_service.go      # Course management
│   │   │   │   ├── progress_service.go    # Progress tracking
│   │   │   │   ├── recommendation_service.go # Recommendation logic
│   │   │   │   └── gamification_service.go   # Achievements, points
│   │   │   │
│   │   │   └── errors/              # Domain-specific errors
│   │   │       └── errors.go
│   │   │
│   │   ├── ports/                   # 🟡 PORTS (Interfaces)
│   │   │   ├── repositories/       # Database interfaces
│   │   │   │   ├── user_repository.go
│   │   │   │   ├── course_repository.go
│   │   │   │   ├── resource_repository.go
│   │   │   │   ├── progress_repository.go
│   │   │   │   ├── recommendation_repository.go
│   │   │   │   ├── rating_repository.go
│   │   │   │   └── achievement_repository.go
│   │   │   │
│   │   │   └── services/           # External service interfaces
│   │   │       ├── ml_service.go   # ML recommendation interface
│   │   │       ├── cache_service.go # Cache interface
│   │   │       └── email_service.go # Email notifications
│   │   │
│   │   ├── adapters/               # 🔵 ADAPTERS (Implementations)
│   │   │   │
│   │   │   ├── postgres/          # PostgreSQL adapter
│   │   │   │   ├── postgres.go    # DB connection setup
│   │   │   │   ├── user_repo.go
│   │   │   │   ├── course_repo.go
│   │   │   │   ├── resource_repo.go
│   │   │   │   ├── progress_repo.go
│   │   │   │   ├── recommendation_repo.go
│   │   │   │   ├── rating_repo.go
│   │   │   │   └── achievement_repo.go
│   │   │   │
│   │   │   ├── redis/             # Redis cache adapter
│   │   │   │   ├── redis.go
│   │   │   │   └── cache_service.go
│   │   │   │
│   │   │   ├── http/              # HTTP adapter (REST API)
│   │   │   │   ├── server.go      # HTTP server setup
│   │   │   │   ├── router.go      # Routes configuration
│   │   │   │   │
│   │   │   │   ├── middleware/    # HTTP middlewares
│   │   │   │   │   ├── auth.go
│   │   │   │   │   ├── cors.go
│   │   │   │   │   ├── logger.go
│   │   │   │   │   └── rate_limit.go
│   │   │   │   │
│   │   │   │   └── handlers/      # HTTP handlers
│   │   │   │       ├── auth_handler.go       # /api/auth/*
│   │   │   │       ├── student_handler.go    # /api/students/*
│   │   │   │       ├── teacher_handler.go    # /api/teachers/*
│   │   │   │       ├── course_handler.go     # /api/courses/*
│   │   │   │       ├── resource_handler.go   # /api/resources/*
│   │   │   │       ├── progress_handler.go   # /api/progress/*
│   │   │   │       └── recommendation_handler.go # /api/recommendations/*
│   │   │   │
│   │   │   └── ml_client/         # ML service HTTP client
│   │   │       └── ml_client.go
│   │   │
│   │   ├── config/                # Configuration
│   │   │   └── config.go
│   │   │
│   │   └── pkg/                   # Shared utilities
│   │       ├── jwt/               # JWT token generation/validation
│   │       │   └── jwt.go
│   │       ├── logger/            # Logging utilities
│   │       │   └── logger.go
│   │       ├── validator/         # Input validation
│   │       │   └── validator.go
│   │       └── response/          # HTTP response helpers
│   │           └── response.go
│   │
│   ├── migrations/                # Database migrations
│   │   ├── 000001_init_schema.up.sql
│   │   └── 000001_init_schema.down.sql
│   │
│   ├── go.mod
│   └── go.sum
│
├── ml-service/                    # Python ML Service
│   ├── app/
│   │   ├── __init__.py
│   │   ├── main.py               # FastAPI application
│   │   │
│   │   ├── models/               # ML models
│   │   │   ├── collaborative.py  # Collaborative filtering
│   │   │   ├── content_based.py  # Content-based filtering
│   │   │   └── hybrid.py         # Hybrid recommender
│   │   │
│   │   ├── services/             # Business logic
│   │   │   ├── recommendation_service.py
│   │   │   └── skill_assessment_service.py
│   │   │
│   │   └── utils/
│   │       ├── data_loader.py    # Load data from PostgreSQL
│   │       └── preprocessing.py
│   │
│   ├── requirements.txt
│   ├── Dockerfile
│   └── .env
│
├── frontend/                      # Vanilla JS Frontend
│   ├── public/
│   │   ├── index.html
│   │   ├── login.html
│   │   ├── dashboard.html
│   │   └── assets/
│   │       ├── images/
│   │       └── icons/
│   │
│   ├── src/
│   │   ├── js/
│   │   │   ├── app.js            # Main application
│   │   │   ├── api.js            # API client
│   │   │   ├── auth.js           # Authentication logic
│   │   │   │
│   │   │   ├── components/       # UI components
│   │   │   │   ├── course-card.js
│   │   │   │   ├── progress-bar.js
│   │   │   │   └── achievement-badge.js
│   │   │   │
│   │   │   └── pages/            # Page-specific logic
│   │   │       ├── dashboard.js
│   │   │       ├── courses.js
│   │   │       └── profile.js
│   │   │
│   │   └── css/
│   │       ├── main.css
│   │       ├── components.css
│   │       └── responsive.css
│   │
│   └── package.json
│
├── data/                          # Seed data
│   ├── courses.json              # Sample courses
│   ├── exercises.json            # Sample exercises
│   └── achievements.json         # Achievement definitions
│
├── docker-compose.yml            # Development environment
├── Dockerfile                    # Production build
├── Makefile                      # Build commands
├── .env.example                  # Environment variables template
├── .gitignore
└── README.md
```

## 🎯 Принципы Hexagonal Architecture

### 1. **Domain Layer** (Центр)
- **Не зависит** от внешних библиотек
- Содержит только бизнес-логику
- Определяет интерфейсы (ports)

```go
// ✅ ПРАВИЛЬНО - domain не знает о БД
type UserRepository interface {
    Create(ctx context.Context, user *User) error
}

// ❌ НЕПРАВИЛЬНО - domain зависит от конкретной реализации
import "github.com/lib/pq"
func SaveUser(db *sql.DB, user *User) error
```

### 2. **Ports** (Интерфейсы)
- Определяют контракты между слоями
- Позволяют менять реализацию без изменения бизнес-логики

### 3. **Adapters** (Реализация)
- Конкретные имплементации портов
- Могут меняться независимо от domain

## 🔄 Поток данных (Data Flow)

```
HTTP Request
     ↓
[HTTP Handler] ← adapter
     ↓
[Service] ← domain (use case)
     ↓
[Repository Interface] ← port
     ↓
[PostgreSQL Repo] ← adapter
     ↓
Database
```

**Пример:**
```
POST /api/auth/register
     ↓
auth_handler.go → вызывает
     ↓
auth_service.Register() → вызывает
     ↓
UserRepository.Create() → вызывает
     ↓
postgres/user_repo.go → пишет в
     ↓
PostgreSQL
```

## 📦 Зависимости между слоями

```
┌─────────────────────────────────┐
│   Adapters (HTTP, DB, Redis)   │
│            ↓ depends on          │
├─────────────────────────────────┤
│   Ports (Interfaces)            │
│            ↑ implements          │
├─────────────────────────────────┤
│   Domain (Business Logic)       │
│   (NO external dependencies!)   │
└─────────────────────────────────┘
```

**Правило:** Зависимости идут **внутрь** (к domain), никогда наружу!

## 🧪 Преимущества для тестирования

```go
// Легко создать mock для тестирования
type MockUserRepository struct {
    users map[int64]*models.User
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
    m.users[user.ID] = user
    return nil
}

// Теперь можно тестировать сервис БЕЗ реальной БД
func TestAuthService_Register(t *testing.T) {
    mockRepo := &MockUserRepository{users: make(map[int64]*models.User)}
    authService := services.NewAuthService(mockRepo, "secret")
    
    // Test logic...
}
```

## 🚀 Следующие шаги

1. ✅ Миграции БД созданы
2. ✅ Базовая структура определена
3. 🔄 **Текущий этап:** Создание seed data (примеры курсов)
4. ⏳ Реализация HTTP server
5. ⏳ Создание ML сервиса
6. ⏳ Frontend разработка

## 📝 Команды для разработки

```bash
# Запуск БД
docker-compose up -d postgres

# Запуск миграций
go run cmd/main.go -migrate up

# Откат миграций
go run cmd/main.go -migrate down

# Запуск backend
go run cmd/main.go

# Тесты
go test ./...
```

## 🎓 Дополнительные ресурсы

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Ports and Adapters Pattern](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/)