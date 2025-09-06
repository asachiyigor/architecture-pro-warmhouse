# Руководство по работе и тестированию сервисов

## Обзор архитектуры

Система состоит из следующих компонентов:
- **Smart Home API** (Go) - основной монолит на порту 8081
- **Device Service** (Java/Spring Boot) - управление устройствами на порту 8082  
- **Telemetry Service** (Python/FastAPI) - сбор телеметрии на порту 8083
- **Temperature API** (Go) - датчики температуры на порту 8084
- **RabbitMQ** - брокер сообщений на порту 5672
- **PostgreSQL** - базы данных на портах 5432, 5433, 5434
- **Redis** - кэш на порту 6379

## 1. Smart Home API (Go) - Монолит

### Как работает:
- Основной API-шлюз системы
- Интегрируется с новыми микросервисами через HTTP клиенты
- Реализует паттерн Strangler Fig для постепенной миграции
- Предоставляет агрегированные данные из разных сервисов

### API endpoints:
```
GET /health                    - проверка здоровья
GET /api/v1/devices           - получить все устройства (через Device Service)
GET /api/v1/devices/{id}      - получить устройство по ID
POST /api/v1/devices          - создать устройство
GET /api/v1/devices/stats     - статистика устройств
GET /api/v1/telemetry/readings - получить показания телеметрии
GET /api/v1/telemetry/stats   - статистика телеметрии
```

### Тестирование:
```bash
# Проверка здоровья
curl http://localhost:8081/health

# Получить все устройства
curl http://localhost:8081/api/v1/devices

# Создать устройство
curl -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{"name": "Smart Bulb", "type": "SMART_LIGHT", "location": "Bedroom", "userId": 1}'

# Получить статистику устройств
curl http://localhost:8081/api/v1/devices/stats

# Получить показания телеметрии
curl "http://localhost:8081/api/v1/telemetry/readings?limit=10"

# Получить статистику телеметрии
curl http://localhost:8081/api/v1/telemetry/stats
```

## 2. Device Service (Java/Spring Boot)

### Как работает:
- Управляет устройствами умного дома
- Использует PostgreSQL для хранения данных
- Публикует события в RabbitMQ при изменениях
- Предоставляет REST API с Swagger документацией
- Поддерживает статистику и heartbeat устройств

### База данных:
```sql
-- Таблица devices
id BIGINT PRIMARY KEY
name VARCHAR(255) NOT NULL
type VARCHAR(50) NOT NULL  
location VARCHAR(255)
status VARCHAR(50) DEFAULT 'OFFLINE'
mac_address VARCHAR(17)
ip_address VARCHAR(15)
firmware_version VARCHAR(50)
user_id BIGINT
created_at TIMESTAMP
updated_at TIMESTAMP
last_seen TIMESTAMP
```

### События RabbitMQ:
- `device.created` - создание устройства
- `device.updated` - обновление устройства  
- `device.deleted` - удаление устройства
- `device.status.changed` - изменение статуса

### API endpoints:
```
GET /api/v1/devices           - получить все устройства
GET /api/v1/devices/{id}      - получить устройство по ID
POST /api/v1/devices          - создать устройство
PUT /api/v1/devices/{id}      - обновить устройство
DELETE /api/v1/devices/{id}   - удалить устройство
PATCH /api/v1/devices/{id}/status - изменить статус
POST /api/v1/devices/{id}/heartbeat - отправить heartbeat
GET /api/v1/devices/stats     - статистика устройств
GET /swagger-ui.html          - Swagger UI
GET /actuator/health          - health check
```

### Тестирование:
```bash
# Health check
curl http://localhost:8082/actuator/health

# Создать устройство
curl -X POST http://localhost:8082/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smart Thermostat",
    "type": "THERMOSTAT", 
    "location": "Living Room",
    "userId": 1
  }'

# Получить все устройства
curl http://localhost:8082/api/v1/devices

# Получить устройство по ID
curl http://localhost:8082/api/v1/devices/1

# Обновить статус устройства
curl -X PATCH http://localhost:8082/api/v1/devices/1/status \
  -H "Content-Type: application/json" \
  -d '{"status": "ONLINE"}'

# Отправить heartbeat
curl -X POST http://localhost:8082/api/v1/devices/1/heartbeat \
  -H "Content-Type: application/json" \
  -d '{"ipAddress": "192.168.1.100"}'

# Получить статистику
curl http://localhost:8082/api/v1/devices/stats

# Swagger UI (откройте в браузере)
http://localhost:8082/swagger-ui.html
```

## 3. Telemetry Service (Python/FastAPI)

### Как работает:
- Собирает и хранит показания датчиков
- Асинхронная обработка с SQLAlchemy
- Система алертов с пороговыми значениями
- Агрегация данных по локациям
- REST API с автоматической документацией

### База данных:
```sql
-- Таблица telemetry_readings
id BIGINT PRIMARY KEY
device_id VARCHAR(50) NOT NULL
device_type VARCHAR(50)
location VARCHAR(255)
metric_name VARCHAR(100) NOT NULL
value DECIMAL(10,4) NOT NULL
unit VARCHAR(20)
timestamp TIMESTAMP NOT NULL
raw_data TEXT
quality VARCHAR(20) DEFAULT 'good'
created_at TIMESTAMP DEFAULT NOW()

-- Таблица telemetry_alerts  
id BIGINT PRIMARY KEY
device_id VARCHAR(50) NOT NULL
device_type VARCHAR(50)
location VARCHAR(255)
alert_type VARCHAR(50) NOT NULL
severity VARCHAR(20) NOT NULL
current_value DECIMAL(10,4)
threshold_value DECIMAL(10,4)
message TEXT
is_active VARCHAR(10) DEFAULT 'yes'
acknowledged_at TIMESTAMP
resolved_at TIMESTAMP
created_at TIMESTAMP DEFAULT NOW()
updated_at TIMESTAMP DEFAULT NOW()
```

### API endpoints:
```
POST /api/v1/telemetry/readings        - создать показание
GET /api/v1/telemetry/readings         - получить показания (с фильтрами)
GET /api/v1/telemetry/readings/latest/{device_id}/{metric} - последнее показание
GET /api/v1/telemetry/devices/{device_id}/summary - сводка по устройству
GET /api/v1/telemetry/alerts           - получить алерты
GET /api/v1/telemetry/stats            - общая статистика
GET /api/v1/telemetry/stats/locations  - статистика по локациям
GET /docs                              - автодокументация FastAPI
GET /health                            - health check
```

### Тестирование:
```bash
# Health check
curl http://localhost:8083/health

# Создать показание температуры
curl -X POST http://localhost:8083/api/v1/telemetry/readings \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "sensor_001",
    "device_type": "TEMPERATURE_SENSOR",
    "location": "Living Room", 
    "metric_name": "temperature",
    "value": 23.5,
    "unit": "°C",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
    "quality": "good"
  }'

# Создать показание влажности
curl -X POST http://localhost:8083/api/v1/telemetry/readings \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "sensor_001", 
    "device_type": "HUMIDITY_SENSOR",
    "location": "Living Room",
    "metric_name": "humidity",
    "value": 45.2,
    "unit": "%",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
    "quality": "good"
  }'

# Получить все показания
curl "http://localhost:8083/api/v1/telemetry/readings?limit=10"

# Получить показания по устройству
curl "http://localhost:8083/api/v1/telemetry/readings?device_id=sensor_001&limit=5"

# Получить последнее показание
curl http://localhost:8083/api/v1/telemetry/readings/latest/sensor_001/temperature

# Получить сводку по устройству
curl http://localhost:8083/api/v1/telemetry/devices/sensor_001/summary

# Получить алерты  
curl http://localhost:8083/api/v1/telemetry/alerts

# Получить общую статистику
curl http://localhost:8083/api/v1/telemetry/stats

# Получить статистику по локациям
curl http://localhost:8083/api/v1/telemetry/stats/locations

# FastAPI автодокументация (откройте в браузере)
http://localhost:8083/docs
```

## 4. Temperature API (Go)

### Как работает:
- Простой API для показаний температуры
- Использует Gin framework
- Хранит данные в памяти (для демонстрации)

### API endpoints:
```
GET /health                    - health check
GET /api/v1/temperature       - получить текущую температуру
POST /api/v1/temperature      - добавить показание температуры
```

### Тестирование:
```bash
# Health check
curl http://localhost:8084/health

# Получить текущую температуру
curl http://localhost:8084/api/v1/temperature

# Добавить показание температуры
curl -X POST http://localhost:8084/api/v1/temperature \
  -H "Content-Type: application/json" \
  -d '{"value": 24.5, "location": "Kitchen"}'
```

## 5. RabbitMQ

### Как работает:
- Брокер сообщений для асинхронной коммуникации
- Exchange: `device.events` (topic)
- Очереди для каждого типа события устройств

### Мониторинг:
```bash
# Management UI (откройте в браузере)
http://localhost:15672
# Логин: guest, Пароль: guest

# Посмотреть очереди
curl -u guest:guest http://localhost:15672/api/queues

# Посмотреть сообщения в очереди
curl -u guest:guest http://localhost:15672/api/queues/%2F/device.created.queue/get \
  -X POST -d '{"count":5,"ackmode":"ack_requeue_false","encoding":"auto"}'
```

## 6. Базы данных PostgreSQL

### Подключения:
```bash
# Device Service DB
docker exec -it smarthome-postgres-devices psql -U postgres -d smarthome_devices

# Telemetry Service DB  
docker exec -it smarthome-postgres-telemetry psql -U postgres -d smarthome_telemetry

# Smart Home API DB
docker exec -it smarthome-postgres-smarthome psql -U postgres -d smarthome
```

### Полезные SQL запросы:
```sql
-- Device Service
SELECT * FROM devices ORDER BY created_at DESC LIMIT 10;
SELECT status, COUNT(*) FROM devices GROUP BY status;

-- Telemetry Service
SELECT * FROM telemetry_readings ORDER BY timestamp DESC LIMIT 10;
SELECT device_id, AVG(value) as avg_temp FROM telemetry_readings 
WHERE metric_name = 'temperature' GROUP BY device_id;

SELECT * FROM telemetry_alerts WHERE is_active = 'yes';
```

## 7. Redis

### Подключение и проверка:
```bash
# Подключиться к Redis
docker exec -it smarthome-redis redis-cli

# Команды Redis
PING
KEYS *
GET some_key
```

## Комплексное тестирование системы

### 1. Проверка всех сервисов:
```bash
#!/bin/bash
echo "=== Проверка всех сервисов ==="

echo "Smart Home API:"
curl -s http://localhost:8081/health | jq .

echo "Device Service:"
curl -s http://localhost:8082/actuator/health | jq .

echo "Telemetry Service:"
curl -s http://localhost:8083/health | jq .

echo "Temperature API:"  
curl -s http://localhost:8084/health | jq .
```

### 2. Тест полного workflow:
```bash
#!/bin/bash
echo "=== Тест полного workflow ==="

# 1. Создать устройство через монолит
echo "1. Создание устройства..."
DEVICE=$(curl -s -X POST http://localhost:8081/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Sensor", "type": "TEMPERATURE_SENSOR", "location": "Test Room", "userId": 1}')
echo $DEVICE | jq .

DEVICE_ID=$(echo $DEVICE | jq -r .id)

# 2. Добавить показания телеметрии
echo "2. Добавление показаний..."
curl -s -X POST http://localhost:8083/api/v1/telemetry/readings \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "'$DEVICE_ID'",
    "device_type": "TEMPERATURE_SENSOR", 
    "location": "Test Room",
    "metric_name": "temperature",
    "value": 25.0,
    "unit": "°C",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
    "quality": "good"
  }' | jq .

# 3. Получить агрегированные данные через монолит
echo "3. Получение статистики..."
curl -s http://localhost:8081/api/v1/devices/stats | jq .
curl -s http://localhost:8081/api/v1/telemetry/stats | jq .
```

### 3. Проверка логов:
```bash
# Логи всех сервисов
docker-compose logs -f --tail=100

# Логи конкретного сервиса
docker-compose logs device-service --tail=50
docker-compose logs telemetry-service --tail=50
docker-compose logs smart-home-api --tail=50
```

## Документация API

- **Device Service Swagger**: http://localhost:8082/swagger-ui.html
- **Telemetry Service FastAPI Docs**: http://localhost:8083/docs  
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

## Полезные команды Docker

```bash
# Запуск всей системы
docker-compose up -d

# Перезапуск конкретного сервиса
docker-compose restart device-service

# Пересборка и запуск
docker-compose up --build -d

# Остановка системы
docker-compose down

# Просмотр статуса
docker-compose ps

# Очистка всех данных
docker-compose down -v
docker system prune -a
```