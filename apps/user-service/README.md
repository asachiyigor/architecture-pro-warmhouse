# User Service

Микросервис управления пользователями и домохозяйствами - аутентификация, авторизация и управление профилями пользователей Smart Home системы.

## Технологический стек

- **Go 1.21** + **Gin Framework**
- **PostgreSQL** для хранения пользователей
- **JWT токены** для аутентификации (планируется)
- **bcrypt** для хэширования паролей

## Возможности

### ✅ Реализованные функции
- **Управление пользователями** - CRUD операции
- **Управление домохозяйствами** - создание и управление домами
- **Профили пользователей** - базовая информация
- **Связи пользователь-дом** - владение домохозяйствами
- **REST API** - HTTP endpoints для всех операций

### 📊 Модели данных

#### User Model
```json
{
  "id": "uuid",
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe", 
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

#### Household Model
```json
{
  "id": "uuid",
  "owner_id": "uuid", 
  "name": "My Smart Home",
  "address": "123 Smart Street, Tech City",
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### User Management
- `GET /users` - Получить всех пользователей
- `GET /users/{id}` - Получить пользователя по ID
- `POST /users` - Создать нового пользователя
- `PUT /users/{id}` - Обновить пользователя
- `DELETE /users/{id}` - Удалить пользователя

### Household Management
- `GET /users/{user_id}/households` - Дома пользователя
- `POST /users/{user_id}/households` - Создать новый дом
- `GET /households/{id}` - Получить дом по ID
- `PUT /households/{id}` - Обновить дом
- `DELETE /households/{id}` - Удалить дом

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8084
- **Через Gateway:** http://localhost:8000/api/v1/users/*

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | PostgreSQL connection | `postgres://postgres:postgres@postgres-users:5432/smarthome_users?sslmode=disable` |
| `PORT` | Порт сервиса | `8084` |
| `JWT_SECRET` | JWT secret key | (планируется) |
| `BCRYPT_COST` | Bcrypt cost factor | `10` |

### База данных
- **PostgreSQL** на порту 5435 (docker)
- **Схема:** автоматические миграции через `sql/init.sql`

## Структура проекта

```
user-service/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # Data models  
│   ├── repository/          # Database layer
│   └── service/             # Business logic
├── sql/
│   └── init.sql            # Database schema
└── Dockerfile
```

## Схема базы данных

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL, 
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Households Table  
```sql
CREATE TABLE households (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up user-service postgres-users

# Проверка работы
curl http://localhost:8084/health
curl http://localhost:8084/users
```

### Локальная разработка
```bash
cd user-service
go mod tidy
go run cmd/server/main.go
```

## Безопасность

### ✅ Реализованные меры
- **Email валидация** - проверка корректности email
- **Уникальность email** - предотвращение дублирования
- **CORS поддержка** - настройка для фронтенда

## Связи с другими сервисами

### Используется сервисами:
- **Device Service** - привязка устройств к пользователям
- **Billing Service** - подписки пользователей
- **Automation Service** - сценарии пользователей
- **Smart Home App** - профили пользователей

### Интеграции:
- **API Gateway** - маршрутизация запросов
- **PostgreSQL** - хранение пользователей
- **Kafka** - события пользователей

## Мониторинг

- **Health endpoint** - `/health`
- **Structured logging** - JSON формат логов
- **Database monitoring** - connection pool metrics
- **HTTP metrics** - request/response статистика

## Архитектурная роль

User Service является фундаментальным сервисом для:
- **Identity Management** - управление идентичностями
- **Authorization** - контроль доступа
- **User Experience** - персонализация интерфейса
- **Multi-tenancy** - разделение данных по пользователям
- **Household Management** - группировка устройств по домам