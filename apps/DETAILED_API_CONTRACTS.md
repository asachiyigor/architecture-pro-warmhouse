# Контракты взаимодействия API - Smart Home Pro

## Обзор микросервисной архитектуры

Smart Home Pro использует микросервисную архитектуру с API Gateway и событийной системой на основе Apache Kafka.

### Архитектурные принципы:
- **API Gateway**: Единая точка входа (порт 8000)
- **Database-per-Service**: Каждый сервис имеет собственную БД
- **Event-Driven**: Асинхронная коммуникация через Kafka
- **RESTful APIs**: Стандартизированные HTTP endpoints
- **Technology Diversity**: Go, Java/Spring Boot, Python/FastAPI

---

## 1. Device Service API (Java/Spring Boot)

**Базовый URL**: `http://localhost:8082` (прямой) | `http://localhost:8000/api/v1/devices` (через API Gateway)

### 1.1. Получить список устройств

```http
GET /api/v1/devices
```

#### Параметры запроса:
- `page` (query, optional): Номер страницы (по умолчанию: 0)
- `size` (query, optional): Размер страницы (по умолчанию: 20)
- `userId` (query, optional): Фильтр по ID пользователя

#### Формат ответа:
```json
{
  "content": [
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
  ],
  "pageable": {
    "pageNumber": 0,
    "pageSize": 20,
    "totalElements": 15,
    "totalPages": 1
  }
}
```

#### Коды ответа:
- **200**: Список устройств успешно получен
- **400**: Неверные параметры запроса
- **500**: Внутренняя ошибка сервера

---

### 1.2. Создать устройство

```http
POST /api/v1/devices
```

#### Формат запроса:
```json
{
  "name": "Living Room Smart Light",
  "type": "SMART_LIGHT",
  "location": "Living Room",
  "macAddress": "aa:bb:cc:dd:ee:11",
  "ipAddress": "192.168.1.101",
  "firmwareVersion": "2.1.0",
  "userId": 1
}
```

#### Формат ответа:
```json
{
  "id": 2,
  "name": "Living Room Smart Light",
  "type": "SMART_LIGHT", 
  "location": "Living Room",
  "status": "OFFLINE",
  "macAddress": "aa:bb:cc:dd:ee:11",
  "ipAddress": "192.168.1.101",
  "firmwareVersion": "2.1.0",
  "userId": 1,
  "createdAt": "2025-01-03T11:00:00Z",
  "updatedAt": "2025-01-03T11:00:00Z",
  "lastSeen": null
}
```

#### Коды ответа:
- **201**: Устройство успешно создано
- **400**: Неверные данные в запросе (валидация не прошла)
- **409**: Устройство с таким MAC-адресом уже существует
- **500**: Внутренняя ошибка сервера

---

### 1.3. Получить устройство по ID

```http
GET /api/v1/devices/{id}
```

#### Параметры пути:
- `id` (required): ID устройства

#### Формат ответа:
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

#### Коды ответа:
- **200**: Устройство найдено
- **404**: Устройство не найдено
- **500**: Внутренняя ошибка сервера

---

### 1.4. Обновить статус устройства

```http
PATCH /api/v1/devices/{id}/status
```

#### Формат запроса:
```json
{
  "status": "MAINTENANCE"
}
```

#### Формат ответа:
```json
{
  "id": 1,
  "status": "MAINTENANCE",
  "updatedAt": "2025-01-03T12:00:00Z",
  "lastSeen": "2025-01-03T12:00:00Z"
}
```

#### Коды ответа:
- **200**: Статус устройства успешно обновлен
- **400**: Неверный статус
- **404**: Устройство не найдено
- **500**: Внутренняя ошибка сервера

#### Доступные статусы DeviceStatus:
- `ONLINE`, `OFFLINE`, `ERROR`, `MAINTENANCE`, `CONFIGURING`, `LOW_BATTERY`

---

### 1.5. Статистика устройств

```http
GET /api/v1/devices/stats
```

#### Формат ответа:
```json
{
  "totalDevices": 15,
  "onlineDevices": 12,
  "offlineDevices": 2,
  "errorDevices": 1,
  "devicesByType": {
    "TEMPERATURE_SENSOR": 5,
    "SMART_LIGHT": 4,
    "SMART_PLUG": 3,
    "SMART_LOCK": 2,
    "SECURITY_CAMERA": 1
  },
  "devicesByStatus": {
    "ONLINE": 12,
    "OFFLINE": 2,
    "ERROR": 1,
    "MAINTENANCE": 0,
    "CONFIGURING": 0,
    "LOW_BATTERY": 0
  }
}
```

#### Коды ответа:
- **200**: Статистика успешно получена
- **500**: Внутренняя ошибка сервера

---

## 2. Telemetry Service API (Python/FastAPI)

**Базовый URL**: `http://localhost:8083` (прямой) | `http://localhost:8000/api/v1/telemetry` (через API Gateway)

### 2.1. Создать показание телеметрии

```http
POST /api/v1/telemetry/readings
```

#### Формат запроса:
```json
{
  "device_id": "temp_sensor_001",
  "device_type": "TEMPERATURE_SENSOR",
  "location": "Kitchen",
  "metric_name": "temperature",
  "value": 23.5,
  "unit": "°C",
  "timestamp": "2025-01-03T10:30:00Z"
}
```

#### Формат ответа:
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
  "created_at": "2025-01-03T10:30:05Z"
}
```

#### Коды ответа:
- **201**: Показание телеметрии успешно создано
- **400**: Неверные данные в запросе
- **422**: Ошибка валидации данных
- **500**: Внутренняя ошибка сервера

---

### 2.2. Получить показания телеметрии

```http
GET /api/v1/telemetry/readings
```

#### Параметры запроса:
- `device_id` (query, optional): Фильтр по ID устройства
- `metric_name` (query, optional): Фильтр по названию метрики
- `location` (query, optional): Фильтр по локации
- `start_time` (query, optional): Начальное время (ISO 8601)
- `end_time` (query, optional): Конечное время (ISO 8601)
- `limit` (query, optional): Количество записей (по умолчанию: 100)

#### Формат ответа:
```json
{
  "readings": [
    {
      "id": 1,
      "device_id": "temp_sensor_001",
      "device_type": "TEMPERATURE_SENSOR",
      "location": "Kitchen", 
      "metric_name": "temperature",
      "value": 23.5,
      "unit": "°C",
      "timestamp": "2025-01-03T10:30:00Z",
      "created_at": "2025-01-03T10:30:05Z"
    }
  ],
  "count": 1,
  "total": 150
}
```

#### Коды ответа:
- **200**: Показания успешно получены
- **400**: Неверные параметры запроса
- **500**: Внутренняя ошибка сервера

---

### 2.3. Получить последние показания устройства

```http
GET /api/v1/telemetry/readings/latest/{device_id}
```

#### Параметры пути:
- `device_id` (required): ID устройства

#### Параметры запроса:
- `metric_name` (query, optional): Фильтр по названию метрики

#### Формат ответа:
```json
{
  "device_id": "temp_sensor_001",
  "latest_readings": [
    {
      "metric_name": "temperature",
      "value": 23.5,
      "unit": "°C",
      "timestamp": "2025-01-03T10:30:00Z"
    },
    {
      "metric_name": "humidity", 
      "value": 65.2,
      "unit": "%",
      "timestamp": "2025-01-03T10:29:00Z"
    }
  ]
}
```

#### Коды ответа:
- **200**: Последние показания найдены
- **404**: Устройство не найдено
- **500**: Внутренняя ошибка сервера

---

### 2.4. Статистика телеметрии

```http
GET /api/v1/telemetry/stats
```

#### Параметры запроса:
- `period` (query, optional): Период (hour, day, week, month)
- `device_type` (query, optional): Фильтр по типу устройства

#### Формат ответа:
```json
{
  "total_readings": 1520,
  "active_devices": 15,
  "metrics_by_type": {
    "temperature": 680,
    "humidity": 680,
    "motion": 160
  },
  "period_stats": {
    "period": "day",
    "start_time": "2025-01-03T00:00:00Z",
    "end_time": "2025-01-03T23:59:59Z",
    "readings_count": 288,
    "avg_readings_per_hour": 12
  }
}
```

#### Коды ответа:
- **200**: Статистика успешно получена
- **500**: Внутренняя ошибка сервера

---

## 3. Temperature API Service (Go/Gin)

**Базовый URL**: `http://localhost:8081`

### 3.1. Проверка состояния

```http
GET /health
```

#### Формат ответа:
```json
{
  "status": "ok"
}
```

#### Коды ответа:
- **200**: Сервис работает нормально
- **500**: Сервис недоступен

---

### 3.2. Получить температуру по комнате

```http
GET /temperature?location={room_name}
```

#### Параметры запроса:
- `location` (required): Название комнаты (Living Room, Bedroom, Kitchen)

#### Формат ответа:
```json
{
  "value": 23.45,
  "unit": "°C",
  "timestamp": "2025-01-03T10:30:00Z",
  "location": "Living Room",
  "status": "active",
  "sensor_id": "1",
  "sensor_type": "temperature",
  "description": "Temperature sensor in Living Room"
}
```

#### Коды ответа:
- **200**: Данные о температуре получены
- **400**: Не указана комната
- **500**: Внутренняя ошибка сервера

---

### 3.3. Получить температуру по ID датчика

```http
GET /temperature/{sensor_id}
```

#### Параметры пути:
- `sensor_id` (required): ID датчика (1, 2, 3)

#### Формат ответа:
```json
{
  "value": 24.12,
  "unit": "°C", 
  "timestamp": "2025-01-03T10:35:00Z",
  "location": "Kitchen",
  "status": "active",
  "sensor_id": "3",
  "sensor_type": "temperature",
  "description": "Temperature sensor in Kitchen"
}
```

#### Коды ответа:
- **200**: Данные о температуре получены
- **400**: Неверный ID датчика
- **500**: Внутренняя ошибка сервера

---

## 4. User Service API (Go/Gin)

**Базовый URL**: `http://localhost:8084` (прямой) | `http://localhost:8000/api/v1/users` (через API Gateway)

### 4.1. Создать пользователя

```http
POST /api/v1/users
```

#### Формат запроса:
```json
{
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "password": "secure_password123"
}
```

#### Формат ответа:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

#### Коды ответа:
- **201**: Пользователь успешно создан
- **400**: Неверные данные в запросе
- **409**: Пользователь с таким email уже существует
- **500**: Внутренняя ошибка сервера

---

### 4.2. Получить пользователя

```http
GET /api/v1/users/{id}
```

#### Параметры пути:
- `id` (required): UUID пользователя

#### Формат ответа:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe", 
  "is_active": true,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z"
}
```

#### Коды ответа:
- **200**: Пользователь найден
- **404**: Пользователь не найден
- **500**: Внутренняя ошибка сервера

---

### 4.3. Создать домохозяйство

```http
POST /api/v1/households
```

#### Формат запроса:
```json
{
  "name": "Family Home",
  "address": "123 Main Street, City, Country",
  "owner_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Формат ответа:
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "name": "Family Home",
  "address": "123 Main Street, City, Country", 
  "owner_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2025-01-03T11:00:00Z",
  "updated_at": "2025-01-03T11:00:00Z"
}
```

#### Коды ответа:
- **201**: Домохозяйство успешно создано
- **400**: Неверные данные в запросе
- **404**: Пользователь-владелец не найден
- **500**: Внутренняя ошибка сервера

---

## 5. Automation Service API (Go/Gin)

**Базовый URL**: `http://localhost:8086` (прямой) | `http://localhost:8000/api/v1/automation` (через API Gateway)

### 5.1. Создать сценарий

```http
POST /api/v1/automation/scenarios
```

#### Формат запроса:
```json
{
  "household_id": "660e8400-e29b-41d4-a716-446655440001",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Evening Scene",
  "description": "Turn off all lights and lock doors",
  "actions": [
    {
      "device_id": "light_001",
      "action_type": "turn_off",
      "parameters": "{}",
      "order_position": 1
    },
    {
      "device_id": "lock_001", 
      "action_type": "lock",
      "parameters": "{\"timeout\": 30}",
      "order_position": 2
    }
  ]
}
```

#### Формат ответа:
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "household_id": "660e8400-e29b-41d4-a716-446655440001",
  "created_by": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Evening Scene",
  "description": "Turn off all lights and lock doors",
  "is_active": true,
  "created_at": "2025-01-03T12:00:00Z",
  "updated_at": "2025-01-03T12:00:00Z",
  "actions": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "scenario_id": "770e8400-e29b-41d4-a716-446655440002",
      "device_id": "light_001",
      "action_type": "turn_off",
      "parameters": "{}",
      "order_position": 1
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440004",
      "scenario_id": "770e8400-e29b-41d4-a716-446655440002", 
      "device_id": "lock_001",
      "action_type": "lock", 
      "parameters": "{\"timeout\": 30}",
      "order_position": 2
    }
  ]
}
```

#### Коды ответа:
- **201**: Сценарий успешно создан
- **400**: Неверные данные в запросе
- **404**: Домохозяйство или пользователь не найдены
- **500**: Внутренняя ошибка сервера

---

### 5.2. Выполнить сценарий

```http
POST /api/v1/automation/scenarios/{id}/execute
```

#### Параметры пути:
- `id` (required): UUID сценария

#### Формат ответа:
```json
{
  "scenario_id": "770e8400-e29b-41d4-a716-446655440002",
  "execution_id": "aa0e8400-e29b-41d4-a716-446655440005",
  "status": "executing",
  "started_at": "2025-01-03T20:00:00Z",
  "actions_total": 2,
  "actions_completed": 0,
  "estimated_completion": "2025-01-03T20:00:30Z"
}
```

#### Коды ответа:
- **202**: Сценарий запущен на выполнение
- **404**: Сценарий не найден
- **409**: Сценарий уже выполняется
- **500**: Внутренняя ошибка сервера

---

## 6. Notification Service API (Go/Gin)

**Базовый URL**: `http://localhost:8088` (прямой) | `http://localhost:8000/api/v1/notifications` (через API Gateway)

### 6.1. Отправить уведомление

```http
POST /api/v1/notifications/send
```

#### Формат запроса:
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Device Alert",
  "message": "Kitchen temperature sensor is offline",
  "type": "alert",
  "channels": ["push", "email"],
  "priority": "high"
}
```

#### Формат ответа:
```json
{
  "notification_id": "bb0e8400-e29b-41d4-a716-446655440006",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Device Alert",
  "message": "Kitchen temperature sensor is offline",
  "type": "alert",
  "status": "sent",
  "channels_sent": ["push", "email"],
  "sent_at": "2025-01-03T14:30:00Z"
}
```

#### Коды ответа:
- **201**: Уведомление отправлено
- **400**: Неверные данные в запросе
- **404**: Пользователь не найден
- **500**: Внутренняя ошибка сервера

---

### 6.2. WebSocket подключение

```
WebSocket: ws://localhost:8088/ws?user_id={user_id}
```

#### Параметры подключения:
- `user_id` (required): UUID пользователя

#### Формат сообщений:
```json
{
  "type": "notification",
  "data": {
    "id": "bb0e8400-e29b-41d4-a716-446655440006",
    "title": "Device Alert",
    "message": "Kitchen temperature sensor is offline",
    "type": "alert",
    "timestamp": "2025-01-03T14:30:00Z"
  }
}
```

---

## 7. Общие принципы и коды ответов

### 7.1. Стандартные HTTP коды

#### Успешные ответы:
- **200 OK**: Запрос выполнен успешно
- **201 Created**: Ресурс создан
- **202 Accepted**: Запрос принят к обработке
- **204 No Content**: Запрос выполнен, нет данных для ответа

#### Ошибки клиента:
- **400 Bad Request**: Неверный запрос
- **401 Unauthorized**: Требуется аутентификация
- **403 Forbidden**: Доступ запрещен
- **404 Not Found**: Ресурс не найден
- **409 Conflict**: Конфликт данных
- **422 Unprocessable Entity**: Ошибка валидации

#### Ошибки сервера:
- **500 Internal Server Error**: Внутренняя ошибка сервера
- **502 Bad Gateway**: Ошибка API Gateway
- **503 Service Unavailable**: Сервис недоступен
- **504 Gateway Timeout**: Таймаут API Gateway

### 7.2. Стандартный формат ошибок

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Device name is required",
    "details": {
      "field": "name",
      "value": null,
      "constraint": "NotBlank"
    },
    "timestamp": "2025-01-03T10:30:00Z",
    "path": "/api/v1/devices"
  }
}
```

### 7.3. Общие заголовки

#### Request Headers:
```http
Content-Type: application/json
Accept: application/json
Authorization: Bearer {jwt_token}
X-Request-ID: uuid-для-трассировки
```

#### Response Headers:
```http
Content-Type: application/json
X-Response-Time: 45ms
X-Request-ID: uuid-для-трассировки
```

---

## 8. Примеры интеграции

### 8.1. Создание устройства с телеметрией

```bash
# 1. Создать устройство
curl -X POST http://localhost:8000/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Outdoor Temperature Sensor",
    "type": "TEMPERATURE_SENSOR",
    "location": "Garden",
    "macAddress": "aa:bb:cc:dd:ee:22",
    "userId": 1
  }'

# 2. Отправить показания телеметрии
curl -X POST http://localhost:8000/api/v1/telemetry/readings \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "2",
    "device_type": "TEMPERATURE_SENSOR",
    "location": "Garden",
    "metric_name": "temperature",
    "value": 18.5,
    "unit": "°C"
  }'
```

### 8.2. Создание сценария автоматизации

```bash
# 1. Создать пользователя и домохозяйство
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "Test",
    "last_name": "User"
  }'

# 2. Создать сценарий
curl -X POST http://localhost:8000/api/v1/automation/scenarios \
  -H "Content-Type: application/json" \
  -d '{
    "household_id": "{household_uuid}",
    "created_by": "{user_uuid}",
    "name": "Morning Scene",
    "description": "Turn on lights and start coffee maker",
    "actions": [
      {
        "device_id": "light_001",
        "action_type": "turn_on",
        "parameters": "{\"brightness\": 80}",
        "order_position": 1
      }
    ]
  }'

# 3. Выполнить сценарий
curl -X POST http://localhost:8000/api/v1/automation/scenarios/{scenario_id}/execute
```

---

Эта документация обеспечивает полное понимание контрактов взаимодействия всех микросервисов системы Smart Home Pro, включая форматы запросов и ответов, коды статуса и практические примеры использования.