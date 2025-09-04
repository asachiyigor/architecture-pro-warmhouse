# Device Service

Микросервис управления устройствами IoT - управление регистрацией, конфигурацией и мониторингом умных устройств.

## Технологический стек

- **Java 17** + **Spring Boot 3**
- **PostgreSQL** для хранения данных устройств
- **Apache Kafka** для событийной архитектуры
- **Swagger/OpenAPI** для документации API

## Возможности

### ✅ Реализованные функции
- **Регистрация устройств** - добавление новых IoT устройств
- **CRUD операции** - создание, чтение, обновление, удаление
- **Мониторинг статуса** - отслеживание online/offline состояния
- **Сетевые настройки** - управление IP, MAC адресами
- **Версии прошивки** - отслеживание firmware versions
- **Kafka события** - публикация device.events

### 📊 Модель данных Device
```json
{
  "id": 1,
  "name": "Kitchen Temperature Sensor",
  "type": "TEMPERATURE_SENSOR",
  "location": "Kitchen",
  "status": "online",
  "mac_address": "aa:bb:cc:dd:ee:ff",
  "ip_address": "192.168.1.100",
  "firmware_version": "1.2.3",
  "user_id": 1,
  "created_at": "2025-01-03T10:00:00Z",
  "updated_at": "2025-01-03T10:00:00Z",
  "last_seen": "2025-01-03T10:30:00Z"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### Device Management
- `GET /devices` - Получить все устройства
- `GET /devices/{id}` - Получить устройство по ID
- `POST /devices` - Создать новое устройство
- `PUT /devices/{id}` - Обновить устройство
- `DELETE /devices/{id}` - Удалить устройство
- `PATCH /devices/{id}/status` - Обновить статус устройства

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8082
- **Через Gateway:** http://localhost:8000/api/v1/devices/*

## Swagger Documentation

- **Interactive UI:** http://localhost:8082/swagger-ui.html
- **OpenAPI JSON:** http://localhost:8082/api-docs

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_HOST` | Хост PostgreSQL | `localhost` |
| `DATABASE_PORT` | Порт PostgreSQL | `5432` |
| `DATABASE_NAME` | Имя базы данных | `smarthome_devices` |
| `DATABASE_USER` | Пользователь БД | `postgres` |
| `DATABASE_PASSWORD` | Пароль БД | `postgres` |
| `KAFKA_BROKERS` | Kafka брокеры | `localhost:9092` |

### База данных
- **PostgreSQL** на порту 5433 (docker) или 5432 (native)
- **Схема:** автоматическое создание через Hibernate DDL

## Kafka интеграция

### Публикуемые события
- **Topic:** `device.events`
- **Event Types:**
  - `device_created` - новое устройство
  - `device_updated` - обновление устройства
  - `device_deleted` - удаление устройства  
  - `device_status_changed` - изменение статуса

### Event Schema
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

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up device-service postgres-devices kafka

# Проверка работы
curl http://localhost:8082/health
curl http://localhost:8082/devices
```

### Локальная разработка
```bash
cd device-service
./mvnw clean install
./mvnw spring-boot:run
```

## Типы устройств

- `TEMPERATURE_SENSOR` - датчики температуры
- `SMART_LIGHT` - умные лампы
- `SMART_PLUG` - умные розетки  
- `SMART_LOCK` - умные замки

## Мониторинг

- **Health checks:** `/actuator/health`
- **Metrics:** `/actuator/metrics`
- **Prometheus:** `/actuator/prometheus`
- **Database:** Hibernate статистика

## Архитектурная роль

Device Service является центральным компонентом для:
- **Регистрации IoT устройств**
- **Управления жизненным циклом**
- **Интеграции с Telemetry Service**
- **Событийного взаимодействия**