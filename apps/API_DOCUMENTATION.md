# Комплексная API Документация - Smart Home Pro

## 1. Микросервисная архитектура с API Gateway

### Современная архитектура: REST + Event-Driven + API Gateway

**Архитектурные компоненты:**

#### 🚪 **API Gateway** - Единая точка входа
- **Назначение**: Маршрутизация запросов, авторизация, мониторинг
- **Технология**: Go (Gin framework) 
- **Порт**: 8000
- **URL**: http://localhost:8000

**Преимущества API Gateway:**
1. **Единая точка входа**: Все запросы проходят через один endpoint
2. **Прозрачная маршрутизация**: Клиенты не знают о внутренней архитектуре
3. **Централизованная авторизация**: Security на уровне gateway
4. **Мониторинг**: Логирование всех запросов
5. **Версионирование**: Управление версиями API

#### 🔄 **REST API** для синхронного взаимодействия
- **Назначение**: CRUD операции, получение данных, управление ресурсами
- **Микросервисы**: User, Device, Telemetry, Automation, Billing, Notification
- **Протокол**: HTTP/JSON через API Gateway

#### ⚡ **Event-Driven Architecture** для асинхронного взаимодействия  
- **Назначение**: Уведомления, алерты, автоматизация, межсервисная коммуникация
- **Технология**: Apache Kafka
- **Паттерны**: Publish/Subscribe, Event Sourcing

**Преимущества Kafka:**
1. **Высокая производительность**: Миллионы сообщений в секунду
2. **Distributed**: Горизонтальное масштабирование
3. **Persistent**: Долговременное хранение событий
4. **Fault-tolerant**: Репликация и восстановление
5. **Ordering guarantee**: Порядок сообщений в партициях

### Когда используется каждый тип:

| Сценарий | API Тип | Маршрут | Обоснование |
|----------|---------|---------|-------------|
| Получить список устройств | REST | API Gateway → Device Service | Синхронная операция |
| Создать пользователя | REST | API Gateway → User Service | CRUD с валидацией |
| Критический перегрев | Kafka Event | Telemetry → Notification | Асинхронное уведомление |
| Сценарий автоматизации | REST + Kafka | API Gateway → Automation + Events | Гибридный подход |
| Метрики телеметрии | Kafka Stream | Telemetry → InfluxDB | Потоковая обработка |

## 2. REST API через API Gateway

### 🚪 API Gateway - Единая точка входа
**Базовый URL:** `http://localhost:8000` (development)

Все запросы проходят через API Gateway и маршрутизируются к соответствующим микросервисам.

#### Общие Endpoints

##### 🏥 Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "api-gateway",
  "timestamp": 1735896000
}
```

### 👥 User Service API (Go/Gin)
**Внутренний URL:** `http://user-service:8084`
**Через Gateway:** `http://localhost:8000/api/v1/users/*`
**Технологии:** Gin HTTP, GORM

##### Пользователи
```http
GET /api/v1/users
POST /api/v1/users
PUT /api/v1/users/{id}
DELETE /api/v1/users/{id}
```

**User Object:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

##### Домохозяйства
```http
GET /api/v1/users/{user_id}/households
POST /api/v1/users/{user_id}/households
PUT /api/v1/households/{id}
DELETE /api/v1/households/{id}
```

### 📱 Device Service API (Java/Spring Boot)
**Внутренний URL:** `http://device-service:8082`
**Через Gateway:** `http://localhost:8000/api/v1/devices/*`
**Технологии:** Spring Boot, Spring Data JPA, Spring REST Controllers

##### Управление устройствами
```http
GET /api/v1/devices
GET /api/v1/devices/{id}
POST /api/v1/devices
PUT /api/v1/devices/{id}
DELETE /api/v1/devices/{id}
```

**Device Object:**
```json
{
  "id": 1,
  "name": "Kitchen Temperature Sensor",
  "type": "TEMPERATURE_SENSOR",
  "location": "Kitchen",
  "status": "ONLINE",
  "macAddress": "aa:bb:cc:dd:ee:ff",
  "ipAddress": "192.168.1.100",
  "firmwareVersion": "1.2.3",
  "userId": 1,
  "createdAt": "2025-01-03T10:00:00Z",
  "updatedAt": "2025-01-03T10:00:00Z",
  "lastSeen": "2025-01-03T10:30:00Z"
}
```

### 📊 Telemetry Service API (Python/FastAPI)
**Внутренний URL:** `http://telemetry-service:8083`
**Через Gateway:** `http://localhost:8000/api/v1/telemetry/*`
**Технологии:** FastAPI, SQLAlchemy, Pydantic, InfluxDB, asyncio

##### Телеметрия
```http
GET /api/v1/telemetry/readings
POST /api/v1/telemetry/readings
GET /api/v1/telemetry/readings/{device_id}
GET /api/v1/telemetry/latest/{device_id}
```

**Telemetry Reading Object:**
```json
{
  "id": 1,
  "device_id": "temp_sensor_001",
  "device_type": "TEMPERATURE_SENSOR",
  "location": "Kitchen",
  "metric_name": "temperature",
  "value": 23.5,
  "unit": "°C",
  "timestamp": "2025-01-03T10:30:00Z",
  "created_at": "2025-01-03T10:30:00Z"
}
```

### 🏠 Automation Service API (Go/Gin)
**Внутренний URL:** `http://automation-service:8087`
**Через Gateway:** `http://localhost:8000/api/v1/automation/*`
**Технологии:** Gin HTTP, GORM, Go goroutines

##### Сценарии автоматизации
```http
GET /api/v1/automation/scenarios
POST /api/v1/automation/scenarios
PUT /api/v1/automation/scenarios/{id}
DELETE /api/v1/automation/scenarios/{id}
POST /api/v1/automation/scenarios/{id}/execute
```

**Scenario Object:**
```json
{
  "id": "uuid",
  "household_id": "uuid",
  "created_by": "uuid",
  "name": "Evening Scene",
  "description": "Turn off all lights and lock doors",
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z",
  "actions": [
    {
      "id": "uuid",
      "device_id": "light_001",
      "action_type": "turn_off",
      "parameters": "{}",
      "order_position": 1
    }
  ]
}
```

### 💳 Billing Service API (Go/Gin)
**Внутренний URL:** `http://billing-service:8085`
**Через Gateway:** `http://localhost:8000/api/v1/billing/*`
**Технологии:** Gin HTTP, GORM, платёжные системы интеграция

##### Тарифные планы и подписки
```http
GET /api/v1/billing/plans
GET /api/v1/billing/subscriptions
POST /api/v1/billing/subscriptions
PUT /api/v1/billing/subscriptions/{id}
```

### 🔔 Notification Service API (Node.js Express)
**Внутренний URL:** `http://notification-service:8088`
**Через Gateway:** `http://localhost:8000/api/v1/notifications/*`
**Технологии:** Node.js Express, Socket.IO WebSocket, swagger-ui-express

##### Уведомления
```http
GET /api/v1/notifications
POST /api/v1/notifications
WebSocket: ws://localhost:8088/ws
```

### 🌡️ Temperature API Service (Go/Gin)

**Базовый URL:** `http://localhost:8081` (development) / `http://temperature-api:8081` (docker)
**Технологии:** Gin HTTP, Swagger/OpenAPI
**Интеграция:** Автоматический сбор данных Telemetry Service

#### Endpoints

##### 🏥 Health Check
```http
GET /health
```

##### 🌡️ Получение температуры
```http
GET /temperature?location={room_name}
GET /temperature/{sensor_id}
```

**Примеры ответов:**
```json
{
  "value": 23.45,
  "unit": "°C", 
  "timestamp": "2024-01-15T10:30:00Z",
  "location": "Living Room",
  "status": "active",
  "sensor_id": "1",
  "sensor_type": "temperature",
  "description": "Temperature sensor in Living Room"
}
```

### 📋 Mapping Location ↔ Sensor ID

| Location     | Sensor ID | Description                    |
|-------------|-----------|--------------------------------|
| Living Room | 1         | Датчик в гостиной             |
| Bedroom     | 2         | Датчик в спальне               |
| Kitchen     | 3         | Датчик на кухне                |

## 3. Event-Driven Architecture (Kafka)

### 📨 Apache Kafka Event Streaming

**Message Broker:** Apache Kafka
**Servers:** 
- Development: `localhost:9092`
- Docker network: `smarthome-kafka:9092`

#### Kafka Topics и Events

##### 📱 Device Events
```
Topic: device-events
Producer: Device Service (Java/Spring Boot)
Consumers: Notification Service, Automation Service
```

**Event Schema:**
```json
{
  "event_type": "device_status_changed",
  "device_id": "temp_sensor_001",
  "device_name": "Kitchen Sensor",
  "location": "Kitchen",
  "old_status": "online",
  "new_status": "offline",
  "reason": "network_timeout",
  "timestamp": "2025-01-03T10:30:00Z",
  "user_id": "user_123"
}
```

##### 📊 Telemetry Events
```
Topic: telemetry-events
Producer: Telemetry Service
Consumers: Notification Service, Analytics Service
```

**Event Schema:**
```json
{
  "event_type": "telemetry_received",
  "device_id": "temp_sensor_001",
  "metric_type": "temperature",
  "value": 23.5,
  "unit": "°C",
  "timestamp": "2025-01-03T10:30:00Z",
  "location": "Kitchen",
  "device_type": "TEMPERATURE_SENSOR"
}
```

##### 🚨 Telemetry Alerts
```
Topic: telemetry-alerts
Producer: Telemetry Service
Consumers: Notification Service, Automation Service
```

**Alert Schema:**
```json
{
  "event_type": "telemetry_alert",
  "device_id": "temp_sensor_001",
  "alert_type": "high_temperature",
  "severity": "critical",
  "message": "Temperature exceeded critical threshold",
  "current_value": 36.8,
  "threshold_value": 35.0,
  "timestamp": "2025-01-03T14:45:00Z"
}
```

##### 🔧 Device Status Changes
```
Topic: device.events (same as above)
Event Type: device_status_changed
Purpose: Мониторинг состояния устройств
```

##### 🏠 Automation Events
```
Topic: automation-events
Producer: Automation Service (Go/Gin:8087)
Consumers: Device Service, Notification Service
```

**Automation Event:**
```json
{
  "event_type": "scenario_executed",
  "scenario_id": "evening_scene_001",
  "scenario_name": "Evening Scene",
  "household_id": "house_123",
  "executed_by": "user_456",
  "actions_count": 5,
  "execution_status": "completed",
  "timestamp": "2025-01-03T20:00:00Z"
}
```

**Доступные типы устройств (DeviceType ENUM):**
- TEMPERATURE_SENSOR, HUMIDITY_SENSOR, MOTION_SENSOR
- DOOR_SENSOR, SMOKE_DETECTOR
- SMART_LIGHT, SMART_PLUG, SMART_LOCK
- THERMOSTAT, SECURITY_CAMERA, GARAGE_DOOR, SMART_SPEAKER

**Статусы устройств (DeviceStatus ENUM):**
- ONLINE, OFFLINE, ERROR, MAINTENANCE, CONFIGURING, LOW_BATTERY

## 4. Интеграционные Паттерны

### 🔄 Синхронное взаимодействие через API Gateway
1. **Client** → **API Gateway** → **Microservices**
   - Единая точка входа для всех запросов
   - Маршрутизация к соответствующим сервисам
   - CORS и security на уровне gateway

2. **Межсервисное взаимодействие**
   - **Telemetry Service** → **Temperature API** (данные датчиков)
   - **Smart Home App** → **All Services** (legacy интеграция)

### ⚡ Асинхронное взаимодействие через Kafka
1. **Device Service** → **Kafka** → **Notification Service**
   - События статуса устройств
   - Уведомления о подключении/отключении
   
2. **Telemetry Service** → **Kafka** → **Multiple Consumers**
   - Метрики и показания датчиков
   - Критические алерты
   - Аналитика данных

3. **Automation Service** → **Kafka** → **Device Service**
   - Выполнение сценариев автоматизации
   - Команды управления устройствами

### 🏗️ Современная микросервисная архитектура

```
┌─────────────┐    HTTP/REST    ┌──────────────────┐
│   Clients   │◄───────────────►│   API Gateway    │
│ (Web, Mobile│                 │   (Go/Gin:8000)  │
│  3rd Party) │                 └──────────┬───────┘
└─────────────┘                            │
                                          │ Route
                   ┌──────────────────────┼──────────────────────┐
                   │                      │                      │
                   ▼                      ▼                      ▼
       ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐
       │   User Service    │  │  Device Service   │  │ Telemetry Service │
       │   (Go/Gin:8084)   │  │(Java/Spring:8082) │  │(Python/FastAPI:   │
       │                   │  │                   │  │      8083)        │
       └─────────┬─────────┘  └─────────┬─────────┘  └─────────┬─────────┘
                 │                      │                      │
                 │                      │                      │
                 └──────────────────────┼──────────────────────┘
                                        │
                           Publish      ▼      Consume
                   ┌─────────────────────────────────────────────┐
                   │          Apache Kafka Event Bus            │
                   │  Topics:                                    │
                   │  • device-events                           │
                   │  • telemetry-events                        │
                   │  • telemetry-alerts                        │
                   │  • automation-events                       │
                   └─────────────────┬───────────────────────────┘
                                     │ Subscribe
                                     ▼
       ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐
       │ Automation Service│  │Notification Service│  │   Smart Home App  │
       │   (Go/Gin:8087)   │  │ (Node.js/Exp:8088) │  │ (Legacy Go:8080)  │
       └───────────────────┘  └───────────────────┘  └───────────────────┘

                          ┌───────────────────┐
                          │ Billing Service   │
                          │   (Go/Gin:8085)   │
                          └───────────────────┘

                          ┌───────────────────┐
                          │ Temperature API   │
                          │   (Go/Gin:8081)   │
                          └───────────────────┘
```

## 5. Контракты взаимодействия

### 🔒 HTTP Status Codes
- **200** OK - Успешный запрос
- **201** Created - Ресурс создан  
- **400** Bad Request - Неверный запрос
- **404** Not Found - Ресурс не найден
- **500** Internal Server Error - Ошибка сервера

### 📝 Error Response Format
```json
{
  "error": "Описание ошибки"
}
```

### 🔐 Security & Validation
- **Input Validation**: Все параметры валидируются на стороне сервера
- **Type Safety**: Строгая типизация через JSON Schema/Go structs
- **Rate Limiting**: Планируется на уровне API Gateway
- **Authentication**: JWT токены (планируется)

## 6. Инфраструктура и Инструменты

### 🏗️ Развертывание и Окружение
```bash
# Запуск всей микросервисной системы
docker-compose up --build

# Запуск отдельных сервисов
docker-compose up api-gateway user-service device-service telemetry-service
```

### 🌐 Доступные endpoints (Development):
- **API Gateway:** http://localhost:8000 (Go/Gin)
- **Smart Home App (Legacy):** http://localhost:8080 (Go/Gin)
- **Temperature API:** http://localhost:8081 (Go/Gin + Swagger)
- **Device Service:** http://localhost:8082 (Java/Spring Boot)
- **Telemetry Service:** http://localhost:8083 (Python/FastAPI)
- **User Service:** http://localhost:8084 (Go/Gin)
- **Billing Service:** http://localhost:8085 (Go/Gin)
- **Automation Service:** http://localhost:8087 (Go/Gin)
- **Notification Service:** http://localhost:8088 (Node.js Express + Socket.IO)

### 🗄️ Databases:
- **PostgreSQL Instances (Database-per-Service):**
  - Main DB: `localhost:5432` (smarthome) - Smart Home App
  - User DB: `localhost:5433` (smarthome_users) - Go/GORM
  - Device DB: `localhost:5434` (smarthome_devices) - Java/JPA
  - Telemetry DB: `localhost:5435` (smarthome_telemetry) - Python/SQLAlchemy
  - Automation DB: `localhost:5436` (smarthome_automation) - Go/GORM
  - Billing DB: `localhost:5437` (smarthome_billing) - Go/GORM
- **InfluxDB:** `localhost:8086` (admin/password123) - Time Series для телеметрии
- **Redis:** `localhost:16379` - Кеширование

### 🚌 Message Broker:
- **Apache Kafka:** `localhost:9092`
- **Zookeeper:** `localhost:2181`

### 📖 API Documentation
- **OpenAPI/Swagger:** Доступно на каждом сервисе с `/swagger` endpoint
- **Kafka Topics:** Документированы в данном файле
- **Postman Collection:** [`smarthome-api.postman_collection.json`](./smarthome-api.postman_collection.json)

### 🧪 Тестирование
- **Health Checks:** `/health` на каждом сервисе
- **Unit Tests:** Специфичные для каждой технологии (Go, Java, Python, Node.js)
- **Integration Tests:** Через API Gateway
- **Load Testing:** Apache Kafka performance тесты

### 📊 Мониторинг и Observability
- **Health checks** ✅ Реализованы на всех сервисах
- **Docker health checks** ✅ Настроены для всех контейнеров
- **Kafka monitoring** ✅ Через Kafka built-in метрики

### 🔧 Development Tools
```bash
# Проверка состояния всех сервисов
docker-compose ps

# Логи конкретного сервиса
docker logs [service-name]

# Kafka topics
docker exec -it smarthome-kafka kafka-topics --bootstrap-server localhost:9092 --list

# PostgreSQL подключение
docker exec -it smarthome-postgres-[service] psql -U postgres -d [database]
```