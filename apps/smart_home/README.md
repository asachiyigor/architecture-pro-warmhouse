# Smart Home App

Монолитное приложение умного дома - MVP версия системы для управления IoT устройствами. Работает как основной фронтенд и бэкенд для Smart Home системы с интеграцией в микросервисную архитектуру через API Gateway.

## Технологический стек

- **Go 1.22** + **Gin Framework**
- **PostgreSQL** для хранения данных датчиков
- **REST API** для управления устройствами
- **Swagger/OpenAPI** для документации API
- **Docker** для контейнеризации

## Возможности

### ✅ Реализованные функции
- **Управление датчиками** - CRUD операции для температурных датчиков
- **Интеграция с Temperature API** - получение актуальных данных
- **Proxy для Telemetry Service** - маршрутизация через API Gateway
- **Swagger документация** - интерактивная API документация
- **Graceful shutdown** - корректное завершение работы
- **Health check endpoint** - мониторинг состояния

### 📊 Модель Sensor
```json
{
  "id": 1,
  "name": "Kitchen Temperature Sensor",
  "type": "temperature",
  "location": "Kitchen", 
  "value": 23.45,
  "unit": "°C",
  "status": "active",
  "last_updated": "2024-01-15T10:30:00Z",
  "created_at": "2024-01-15T09:00:00Z"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### Sensor Management
- `GET /api/v1/sensors` - Получить все датчики
- `GET /api/v1/sensors/{id}` - Получить датчик по ID
- `POST /api/v1/sensors` - Создать новый датчик
- `PUT /api/v1/sensors/{id}` - Обновить датчик
- `DELETE /api/v1/sensors/{id}` - Удалить датчик
- `PATCH /api/v1/sensors/{id}/value` - Обновить значение датчика

### Temperature API Integration
- `GET /api/v1/sensors/temperature/{location}` - Температура по местоположению

### Telemetry Proxy (через API Gateway)
- `GET /api/v1/telemetry/*` - Проксирование к Telemetry Service
- `POST /api/v1/telemetry/*` - Проксирование к Telemetry Service

### Доступ
- **Прямой доступ:** http://localhost:8080
- **OpenAPI документация:** [swagger.yaml](swagger.yaml) (статический файл)

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | PostgreSQL connection | `postgres://postgres:postgres@localhost:5432/smarthome` |
| `PORT` | Порт приложения | `:8080` |
| `TEMPERATURE_API_URL` | Temperature API endpoint | `http://temperature-api:8081` |
| `API_GATEWAY_URL` | API Gateway endpoint | `http://api-gateway:8000` |

### База данных
- **PostgreSQL** на порту 5432 (docker)
- **Схема:** автоматические миграции через `init.sql`
- **Индексы:** по типу, местоположению и статусу датчиков

## Структура проекта

```
smart_home/
├── main.go                    # Entry point и сервер
├── handlers/
│   ├── sensors.go            # Обработчики датчиков
│   └── telemetry.go          # Проксирование телеметрии
├── services/
│   ├── temperature_service.go # Интеграция с Temperature API
│   ├── device_client.go       # Клиент Device Service (планируется)
│   ├── telemetry_client.go    # Прямой клиент Telemetry Service
│   └── telemetry_gateway_client.go # Клиент через API Gateway
├── models/
│   └── sensor.go             # Модели данных
├── db/
│   └── db.go                 # Подключение к базе данных
├── init.sql                  # Схема базы данных
├── swagger.yaml              # OpenAPI спецификация
├── Dockerfile               # Docker конфигурация
├── go.mod                   # Go modules
└── go.sum                   # Зависимости
```

## Схема базы данных

### Sensors Table
```sql
CREATE TABLE sensors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    value FLOAT DEFAULT 0,
    unit VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'inactive',
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_sensors_type ON sensors(type);
CREATE INDEX idx_sensors_location ON sensors(location);
CREATE INDEX idx_sensors_status ON sensors(status);
```

## Поддерживаемые типы датчиков

### ✅ Текущие типы
- **temperature** - температурные датчики с интеграцией Temperature API

### 🔄 Планируемые типы  
- **humidity** - датчики влажности
- **pressure** - барометры
- **motion** - датчики движения
- **light** - датчики освещенности

## Статусы датчиков

- **active** - активный и работающий
- **inactive** - неактивный
- **error** - ошибка датчика
- **maintenance** - техническое обслуживание

## Интеграции

### Внешние API
- **Temperature API** (port 8081) - получение данных о температуре
- **API Gateway** (port 8000) - маршрутизация к микросервисам

### Связи с микросервисами (через API Gateway)
- **Telemetry Service** - отправка данных телеметрии
- **Device Service** - управление устройствами (планируется)
- **User Service** - аутентификация (планируется)

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up smarthome-app postgres temperature-api api-gateway

# Проверка работы
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/sensors
```

### Локальная разработка
```bash
cd smart_home
go mod tidy
go run main.go
```

## Swagger/OpenAPI

### Документация API
- **Полная спецификация:** `swagger.yaml` (513 строк)
- **Статический OpenAPI файл:** доступен локально и в репозитории
- **Версия:** OpenAPI 3.0.3
- **Поддерживаемые форматы:** YAML, JSON

### Основные эндпоинты в Swagger
- Sensor management (CRUD)
- Temperature API integration
- Health checks
- Error handling schemas

## Мониторинг

### Встроенный мониторинг
- **Health endpoint** - `/health`
- **Graceful shutdown** - обработка SIGINT/SIGTERM
- **Structured logging** - логирование операций
- **Database health** - проверка подключения к PostgreSQL

## Архитектурная роль

Smart Home App служит как:
- **MVP Frontend** - пользовательский интерфейс управления устройствами
- **Legacy API** - обратная совместимость со старыми клиентами
- **Gateway Adapter** - адаптер между монолитной и микросервисной архитектурой
- **Development Environment** - быстрый запуск для разработки
- **Proof of Concept** - демонстрация базовых функций умного дома

## Миграция к микросервисам

### ✅ Уже интегрировано
- Проксирование телеметрии через API Gateway
- Использование Temperature API
- Подготовка к интеграции с другими сервисами