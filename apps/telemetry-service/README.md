# Telemetry Service

Микросервис для сбора, обработки и анализа телеметрических данных от IoT устройств умного дома.

## Технологический стек

- **Python 3.11** + **FastAPI**
- **PostgreSQL** для метаданных
- **InfluxDB** для временных рядов
- **Apache Kafka** для событийной архитектуры  
- **Redis** для кэширования

## Возможности

### ✅ Реализованные функции
- **Сбор телеметрии** - получение данных от датчиков
- **Временные ряды** - хранение в InfluxDB
- **Аналитика** - базовая обработка метрик
- **Алерты** - уведомления о критических значениях
- **Kafka интеграция** - публикация событий
- **REST API** - HTTP endpoints для управления

### 📊 Модель данных TelemetryReading
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

## API Endpoints

### Health Check
- `GET /health/` - Проверка состояния сервиса
- `GET /health/dependencies` - Состояние зависимостей

### Telemetry Data
- `GET /telemetry/readings` - Получить показания
- `POST /telemetry/readings` - Добавить показание
- `GET /telemetry/readings/{device_id}` - Показания устройства
- `GET /telemetry/latest/{device_id}` - Последние показания

### Metrics & Analytics
- `GET /metrics/summary` - Сводка метрик
- `GET /metrics/alerts` - Активные алерты

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8083
- **Через Gateway:** http://localhost:8000/api/v1/telemetry/*

## FastAPI Documentation

- **Interactive Docs:** http://localhost:8083/docs
- **ReDoc:** http://localhost:8083/redoc
- **OpenAPI JSON:** http://localhost:8083/openapi.json

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | PostgreSQL connection | `postgresql+asyncpg://...` |
| `INFLUXDB_URL` | InfluxDB endpoint | `http://influxdb:8086` |
| `INFLUXDB_TOKEN` | InfluxDB auth token | `smarthome-super-secret-auth-token` |
| `INFLUXDB_ORG` | InfluxDB organization | `smarthome` |
| `INFLUXDB_BUCKET` | InfluxDB bucket | `telemetry` |
| `KAFKA_BROKERS` | Kafka brokers | `smarthome-kafka:9092` |
| `ENABLE_KAFKA` | Enable Kafka integration | `true` |
| `REDIS_URL` | Redis connection | `redis://redis:6379/0` |
| `TEMPERATURE_API_URL` | External API | `http://temperature-api:8081` |

### Базы данных
- **PostgreSQL:** порт 5434 (docker) - метаданные
- **InfluxDB:** порт 8086 - временные ряды
- **Redis:** порт 16379 - кэширование

## Kafka интеграция

### Публикуемые события

#### Topic: telemetry-events
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

#### Topic: telemetry-alerts  
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

## Поддерживаемые метрики

- **temperature** - температура (°C)
- **humidity** - влажность (%)
- **pressure** - давление (hPa)
- **power** - потребление энергии (W)
- **battery** - уровень батареи (%)

## Алерты и пороги

### Температурные алерты
- **Критический минимум:** < 10°C
- **Критический максимум:** > 35°C
- **Предупреждение минимум:** < 15°C  
- **Предупреждение максимум:** > 30°C

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up telemetry-service postgres-telemetry influxdb kafka redis

# Проверка работы
curl http://localhost:8083/health/
curl http://localhost:8083/telemetry/readings
```

### Локальная разработка
```bash
cd telemetry-service
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8083 --reload
```

## Интеграции

### Внешние сервисы
- **Temperature API** - получение данных датчиков
- **Device Service** - информация об устройствах
- **Notification Service** - отправка алертов

### Потребители событий
- **Notification Service** - уведомления
- **Analytics Service** - аналитика (будущее)
- **Automation Service** - триггеры (будущее)

## Мониторинг

- **Health endpoints** - `/health/` и `/health/dependencies`
- **Metrics endpoint** - `/metrics/summary`
- **Logging** - structured JSON logs
- **InfluxDB metrics** - performance counters

## Архитектурная роль

Telemetry Service является ключевым компонентом для:
- **Централизованного сбора** телеметрических данных
- **Обработки временных рядов** в InfluxDB
- **Детекции аномалий** и алертов
- **Событийного взаимодействия** через Kafka
- **Аналитики** и reporting