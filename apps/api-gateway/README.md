# API Gateway

API Gateway для платформы "Умный дом" - единая точка входа для всех клиентских приложений.

## Архитектура

API Gateway реализует паттерн **Backend for Frontend (BFF)** и обеспечивает:

- **Единую точку входа** для всех клиентов
- **Маршрутизацию** запросов к соответствующим микросервисам
- **Прозрачное проксирование** во время миграции
- **CORS, логирование, мониторинг**

## Поэтапное внедрение

### Фаза 1: Прозрачный прокси (текущая)
Все запросы проксируются в Smart Home App для обеспечения обратной совместимости:

```
Client → API Gateway:8000 → Smart Home App:8080
```

### Фаза 2: Постепенная миграция микросервисов
Поэтапно мигрируем маршруты к соответствующим микросервисам:

```
Client → API Gateway:8000 → {
    /api/v1/sensors/*     → Smart Home App:8080
    /api/v1/telemetry/*   → Telemetry Service:8083  
    /api/v1/devices/*     → Device Service:8082
    /api/v1/users/*       → User Service:8084
    /api/v1/billing/*     → Billing Service:8085
    /api/v1/automation/*  → Automation Service:8087
}
```

## Конфигурация

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | Порт API Gateway | `:8000` |
| `SMART_HOME_URL` | URL Smart Home App | `http://app:8080` |
| `DEVICE_SERVICE_URL` | URL Device Service | `http://device-service:8082` |
| `TELEMETRY_SERVICE_URL` | URL Telemetry Service | `http://telemetry-service:8083` |
| `USER_SERVICE_URL` | URL User Service | `http://user-service:8084` |
| `BILLING_SERVICE_URL` | URL Billing Service | `http://billing-service:8085` |
| `AUTOMATION_SERVICE_URL` | URL Automation Service | `http://automation-service:8087` |

## Endpoints

### Health Check
- `GET /health` - Проверка состояния API Gateway

### API Routes
- `GET|POST|PUT|DELETE /api/v1/*` - Все API маршруты проксируются к соответствующим сервисам

## Запуск

### Docker Compose
```bash
# Запуск всей системы с API Gateway
docker-compose up --build api-gateway

# Проверка работы
curl http://localhost:8000/health
curl http://localhost:8000/api/v1/sensors
```

### Локальная разработка
```bash
cd api-gateway
go mod tidy
go run main.go
```

## Миграция сервисов

Для миграции сервиса из Smart Home App в отдельный микросервис:

1. Обновить маршрутизацию в `internal/gateway/gateway.go`:
```go
// Заменить:
api.Any("/telemetry/*path", gw.proxyToSmartHome)

// На:
api.Any("/telemetry/*path", gw.proxyToTelemetryService)
```

2. Перезапустить API Gateway:
```bash
docker-compose restart api-gateway
```

## Функции

- ✅ **Прозрачное проксирование** запросов
- ✅ **CORS** поддержка
- ✅ **Логирование** запросов
- ✅ **Health check**
- ⏳ **Аутентификация** (планируется)
- ⏳ **Rate limiting** (планируется)
- ⏳ **Circuit breaker** (планируется)