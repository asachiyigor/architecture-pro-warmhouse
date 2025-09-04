# Billing Service

Микросервис биллинга и подписок - управление тарифными планами, подписками пользователей и платежным циклом Smart Home платформы.

## Технологический стек

- **Go 1.21** + **Gin Framework**
- **PostgreSQL** для хранения биллинговых данных
- **Stripe API** для обработки платежей (планируется)
- **Cron Jobs** для автоматического биллинга (планируется)

## Возможности

### ✅ Реализованные функции
- **Управление тарифными планами** - создание и управление планами
- **Подписки пользователей** - привязка пользователей к тарифам
- **Лимиты устройств** - контроль количества устройств по тарифу
- **Базовый биллинг** - отслеживание статусов подписок
- **REST API** - полный CRUD для планов и подписок

### 📊 Модели данных

#### PricingPlan Model
```json
{
  "id": "uuid",
  "name": "Premium Plan",
  "description": "Unlimited devices and advanced features",
  "price": 39.99,
  "billing_interval": "monthly",
  "device_limit": 50,
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

#### Subscription Model
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "plan_id": "uuid",
  "status": "active",
  "started_at": "2025-01-03T10:00:00Z",
  "ends_at": "2025-02-03T10:00:00Z",
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### Pricing Plans
- `GET /plans` - Получить все тарифные планы
- `GET /plans/{id}` - Получить план по ID
- `POST /plans` - Создать новый план
- `PUT /plans/{id}` - Обновить план
- `DELETE /plans/{id}` - Удалить план

### Subscriptions  
- `GET /subscriptions` - Получить все подписки
- `GET /subscriptions/{id}` - Получить подписку по ID
- `GET /users/{user_id}/subscriptions` - Подписки пользователя
- `POST /subscriptions` - Создать подписку
- `PUT /subscriptions/{id}` - Обновить подписку
- `DELETE /subscriptions/{id}` - Отменить подписку

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8085
- **Через Gateway:** http://localhost:8000/api/v1/billing/*

## Тарифные планы

### 📦 Стандартные тарифы
```json
[
  {
    "name": "Basic",
    "description": "Perfect for small apartments",
    "price": 9.99,
    "device_limit": 5,
    "features": ["Basic automation", "Mobile app", "Email support"]
  },
  {
    "name": "Standard", 
    "description": "Great for medium homes",
    "price": 19.99,
    "device_limit": 15,
    "features": ["Advanced automation", "Voice control", "Priority support"]
  },
  {
    "name": "Premium",
    "description": "For large homes and power users", 
    "price": 39.99,
    "device_limit": 50,
    "features": ["Unlimited automation", "AI insights", "24/7 phone support"]
  }
]
```

## Статусы подписок

### ✅ Основные статусы
- **active** - активная подписка
- **cancelled** - отмененная подписка  
- **expired** - истекшая подписка
- **pending** - ожидание активации
- **suspended** - приостановленная (за неуплату)

### 🔄 Lifecycle подписки
```
pending → active → (cancelled/expired/suspended)
```

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | PostgreSQL connection | `postgres://postgres:postgres@postgres-billing:5432/smarthome_billing?sslmode=disable` |
| `PORT` | Порт сервиса | `8085` |
| `STRIPE_SECRET_KEY` | Stripe API key | (планируется) |
| `WEBHOOK_SECRET` | Stripe webhook secret | (планируется) |

### База данных
- **PostgreSQL** на порту 5436 (docker)
- **Схема:** автоматические миграции через `sql/init.sql`

## Схема базы данных

### Pricing_Plans Table
```sql
CREATE TABLE pricing_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    billing_interval VARCHAR(20) NOT NULL DEFAULT 'monthly',
    device_limit INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Subscriptions Table
```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    plan_id UUID NOT NULL REFERENCES pricing_plans(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    started_at TIMESTAMP,
    ends_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Бизнес-логика

### Device Limit Enforcement
```go
func (s *BillingService) CheckDeviceLimit(userID string, deviceCount int) error {
    subscription := s.GetActiveSubscription(userID)
    if subscription == nil {
        return errors.New("no active subscription")
    }
    
    plan := s.GetPlan(subscription.PlanID)
    if deviceCount > plan.DeviceLimit {
        return fmt.Errorf("device limit exceeded: %d/%d", deviceCount, plan.DeviceLimit)
    }
    
    return nil
}
```

### Subscription Renewal
```go
func (s *BillingService) RenewSubscription(subscriptionID string) error {
    subscription := s.GetSubscription(subscriptionID)
    plan := s.GetPlan(subscription.PlanID)
    
    // Charge customer
    success := s.stripe.ChargeCustomer(subscription.UserID, plan.Price)
    if !success {
        subscription.Status = "suspended"
        return errors.New("payment failed")
    }
    
    // Extend subscription
    subscription.EndsAt = subscription.EndsAt.AddDate(0, 1, 0)
    subscription.Status = "active"
    
    return s.repository.UpdateSubscription(subscription)
}
```

## Интеграции

### 📱 Планируемые интеграции
- **Stripe** - обработка платежей
- **PayPal** - альтернативный способ оплаты  
- **Webhook endpoints** - уведомления о платежах
- **Email service** - уведомления о биллинге

### 🔗 Связи с сервисами
- **User Service** - информация о пользователях
- **Device Service** - проверка лимитов устройств
- **Notification Service** - уведомления о платежах

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up billing-service postgres-billing

# Проверка работы
curl http://localhost:8085/health
curl http://localhost:8085/plans
```

### Локальная разработка
```bash
cd billing-service
go mod tidy
go run cmd/server/main.go
```

## Мониторинг

- **Health endpoint** - `/health`
- **Billing metrics** - доходы, конверсии, churn
- **Payment tracking** - успешные/неуспешные платежи
- **Subscription analytics** - активные подписки по планам

## Безопасность

### ✅ Реализованные меры
- **Input validation** - валидация всех данных
- **Price validation** - проверка корректности цен
- **Plan integrity** - целостность тарифных планов

## Архитектурная роль

Billing Service является критически важным для:
- **Монетизации платформы** - получение доходов
- **Контроля доступа** - лимиты по тарифам  
- **Customer Experience** - прозрачный биллинг
- **Business Intelligence** - аналитика продаж
- **Compliance** - соответствие финансовым требованиям