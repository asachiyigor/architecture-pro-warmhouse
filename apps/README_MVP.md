# Smart Home Microservices MVP

## Архитектура

Создан MVP микросервисной архитектуры умного дома с плавным переходом от монолита:

### Компоненты системы:

1. **Smart Home App** (Go) - Существующий монолит на порту 8080
2. **Temperature API** (Go) - Существующий сервис на порту 8081  
3. **Device Management Service** (Go) - Сервис управления устройствами на порту 8082
4. **Telemetry Service** (Go) - Сервис телеметрии на порту 8083
5. **User Service** (Go) - Сервис управления пользователями на порту 8084
6. **Billing Service** (Go) - Сервис биллинга на порту 8085
7. **Automation Service** (Go) - Сервис автоматизации на порту 8086
8. **Notification Service** (Go) - Сервис уведомлений на порту 8088
9. **API Gateway** (Go) - Шлюз API на порту 8000
10. **Apache Kafka** - Брокер сообщений на порту 9092 (+ Zookeeper)
11. **InfluxDB** - База данных временных рядов на порту 8086
12. **PostgreSQL** - Отдельные базы данных для каждого сервиса
13. **Redis** - Кэш на порту 16379

### Диаграмма взаимодействия:

```
                    ┌────────────────────┐
                    │   API Gateway      │
                    │   (Go) :8000       │
                    └──────────┬─────────┘
                               │
    ┌───────────────────────────────────────────────────────────────────┐
    │                                                                  │
┌───┬─────────────┌─────────────┌────────────┌────────────┌────────────┐
│User│Device      │Telemetry   │Billing    │Automation │Notification│
│:8084│Service     │Service     │Service    │Service    │Service     │
│    │:8082       │:8083       │:8085      │:8086      │:8087       │
└────└─────────────└─────────────└────────────└────────────└─────────────┘
                           ▲
                           │
            ┌──────────────┼──────────────┐
            │              │              │
┌────────────────────┐    ┌────────────────────┐
│  Smart Home App     │◄──►│  Temperature API    │
│  (Go) :8080         │    │  (Go) :8081         │
│ [Telemetry proxy]   │    │                    │
└────────────────────┘    └────────────────────┘

┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Apache Kafka  │  │   PostgreSQL    │  │    InfluxDB     │
│   + Zookeeper   │  │   Multi-DB      │  │   Time Series   │
│    :9092        │  │    + Redis      │  │     :8086       │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

## Запуск системы

### Предварительные требования

- Docker и Docker Compose
- Минимум 4GB RAM
- Свободные порты: 8080-8083, 5432, 5433, 5434, 5672, 15672, 6379

### Команды запуска

```bash
# Перейти в директорию apps
cd apps

# Инициализация системы
./init.sh

# Запустить всю систему
docker-compose up --build

# Запустить в фоне
docker-compose up --build -d

# Посмотреть логи
docker-compose logs -f

# Остановить систему
docker-compose down

# Очистить данные и перезапустить
docker-compose down -v
docker-compose up --build

# Запустить только один сервис
docker-compose up device-service

# Масштабирование сервиса
docker-compose up --scale device-service=3
```

## API Endpoints

### API Gateway - http://localhost:8000
- `GET /health` - Общая проверка состояния всех сервисов
- `GET /api/v1/devices/*` - Прокси к device-service
- `GET /api/v1/users/*` - Прокси к user-service
- `GET /api/v1/telemetry/*` - Прокси к telemetry-service
- `GET /api/v1/billing/*` - Прокси к billing-service
- `GET /api/v1/automation/*` - Прокси к automation-service
- `GET /api/v1/notifications/*` - Прокси к notification-service

### Smart Home App (монолит) - http://localhost:8080
- `GET /health` - Проверка состояния
- `GET /api/v1/sensors` - Список датчиков
- `POST /api/v1/sensors` - Создание датчика
- `GET /api/v1/telemetry/readings` - Получение показаний телеметрии (через API Gateway)
- `POST /api/v1/telemetry/readings` - Создание показания телеметрии (через API Gateway)
- `GET /api/v1/telemetry/readings/latest/:deviceId/:metricName` - Последние показания устройства
- `GET /api/v1/telemetry/stats` - Статистика телеметрии (через API Gateway)

### Temperature API - http://localhost:8081  
- `GET /health` - Проверка состояния
- `GET /temperature?location={room}` - Температура по комнате
- `GET /temperature/{sensor_id}` - Температура по ID датчика
- `GET /swagger/index.html` - Swagger UI

### User Service - http://localhost:8084
- `GET /health` - Проверка состояния
- `GET /api/v1/users` - Список пользователей
- `POST /api/v1/users` - Создание пользователя
- `GET /api/v1/users/{id}` - Пользователь по ID
- `GET /api/v1/households` - Список домов
- `POST /api/v1/households` - Создание дома

### Device Management Service - http://localhost:8082
- `GET /health` - Проверка состояния
- `GET /api/v1/devices` - Список устройств
- `POST /api/v1/devices` - Создание устройства
- `GET /api/v1/devices/{id}` - Устройство по ID
- `PUT /api/v1/devices/{id}` - Обновление устройства
- `DELETE /api/v1/devices/{id}` - Удаление устройства
- `PATCH /api/v1/devices/{id}/status` - Обновление статуса
- `GET /api/v1/devices/stats` - Статистика устройств

### Telemetry Service - http://localhost:8083
- `GET /health` - Проверка состояния
- `POST /api/v1/telemetry` - Создание показания телеметрии
- `GET /api/v1/telemetry` - Получение показаний (с фильтрами)
- `GET /api/v1/telemetry/stats` - Статистика телеметрии

### Billing Service - http://localhost:8085
- `GET /health` - Проверка состояния
- `GET /api/v1/subscriptions` - Список подписок
- `POST /api/v1/subscriptions` - Создание подписки
- `GET /api/v1/pricing-plans` - Список тарифных планов

### Automation Service - http://localhost:8086
- `GET /health` - Проверка состояния
- `GET /api/v1/scenarios` - Список сценариев
- `POST /api/v1/scenarios` - Создание сценария
- `POST /api/v1/scenarios/{id}/execute` - Выполнение сценария
- `GET /api/v1/scenarios/stats` - Статистика автоматизации

### Notification Service - http://localhost:8088
- `GET /health` - Проверка состояния
- `POST /api/v1/notifications/send` - Отправка уведомления
- `GET /api/v1/notifications` - Получение списка уведомлений

### Apache Kafka - http://localhost:9092
- **Zookeeper:** порт 2181
- **Kafka Broker:** порт 9092
- **Topics:** device-events, telemetry-alerts, user-events

### InfluxDB - http://localhost:8086
- **Username:** admin
- **Password:** password123
- **Organization:** smarthome
- **Bucket:** telemetry
- **Token:** smarthome-super-secret-auth-token


## Интеграция сервисов

### Синхронное взаимодействие (REST API):
- API Gateway ↔ Все микросервисы
- Smart Home App ↔ Temperature API (основная интеграция)
- Smart Home App → API Gateway → Telemetry Service (интегрированная телеметрия)
- Smart Home App ↔ Device Management Service (напрямую)
- Smart Home App ↔ User Service (напрямую)
- Telemetry Service → Temperature API (автоматический сбор данных)
- Device Service ↔ User Service
- Billing Service ↔ User Service
- Automation Service ↔ Device Service
- Notification Service ↔ User Service

### Асинхронное взаимодействие (Apache Kafka):
- Device Service → События создания/обновления устройств
- Telemetry Service → Критические алерты и метрики
- Notification Service ← Потребляет события для отправки уведомлений

### Хранение данных:
- **PostgreSQL**: Реляционные данные каждого сервиса
- **InfluxDB**: Временные ряды телеметрии
- **Redis**: Кэширование (порт 16379)
- Прямые HTTP вызовы между сервисами

## Базы данных

### PostgreSQL основная (порт 5432):
- База: `smarthome`
- Таблицы: sensors (существующая)

### PostgreSQL для устройств (порт 5433):
- База: `smarthome_devices`  
- Таблицы: devices

### PostgreSQL для телеметрии (порт 5434):
- База: `smarthome_telemetry`
- Таблицы: telemetry_readings, telemetry_alerts, device_metrics

## Мониторинг и отладка

### Health Checks:
```bash
# Общая проверка через API Gateway
curl http://localhost:8000/health

# Отдельные сервисы
curl http://localhost:8080/health   # Smart Home App
curl http://localhost:8081/health   # Temperature API
curl http://localhost:8082/health   # Device Service
curl http://localhost:8083/health   # Telemetry Service
curl http://localhost:8084/health   # User Service
curl http://localhost:8085/health   # Billing Service
curl http://localhost:8086/health   # Automation Service (порт конфликт с InfluxDB)
curl http://localhost:8088/health   # Notification Service
```

### Просмотр логов:
```bash
# Все сервисы
docker-compose logs

# Конкретный сервис
docker-compose logs device-service
docker-compose logs telemetry-service
docker-compose logs temperature-api
docker-compose logs app
```

### Доступ к базам данных:
```bash
# Основная БД
docker exec -it smarthome-postgres psql -U postgres -d smarthome

# БД устройств  
docker exec -it smarthome-postgres-devices psql -U postgres -d smarthome_devices

# БД телеметрии
docker exec -it smarthome-postgres-telemetry psql -U postgres -d smarthome_telemetry
```

## Тестирование MVP

### 1. Создание устройства:
```bash
curl -X POST http://localhost:8082/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Temperature Sensor",
    "type": "TEMPERATURE_SENSOR", 
    "location": "Test Room",
    "userId": 1
  }'
```

### 2. Создание показания телеметрии:
```bash
curl -X POST http://localhost:8083/api/v1/telemetry/readings \
  -H "Content-Type: application/json" \
  -d '{
    "device_id": "1",
    "device_type": "temperature_sensor",
    "location": "Test Room",
    "metric_name": "temperature",
    "value": 25.5,
    "unit": "°C",
    "timestamp": "2024-01-15T10:30:00Z"
  }'
```

### 3. Получение статистики:
```bash
# Статистика устройств
curl http://localhost:8082/api/v1/devices/stats

# Статистика телеметрии  
curl http://localhost:8083/api/v1/telemetry/stats
```

## Градуальная миграция

MVP демонстрирует паттерн **Strangler Fig** для постепенного перехода от монолита:

1. **Фаза 1 (Текущая)**: Монолит + новые микросервисы параллельно
2. **Фаза 2**: Перенос функций управления датчиками в Device Service
3. **Фаза 3**: Перенос телеметрии в Telemetry Service
4. **Фаза 4**: Декомпозиция оставшихся функций монолита

## Планы развития

- [x] API Gateway 
- [ ] Service Discovery (Consul)
- [ ] Distributed Tracing (Jaeger)
- [ ] Metrics (Prometheus + Grafana)
- [ ] Authentication/Authorization (JWT)
- [ ] Circuit Breakers (Hystrix/Resilience4j)
- [ ] Event Sourcing и CQRS
- [ ] Kubernetes deployment
- [ ] Load Balancing
- [ ] Rate Limiting
- [ ] API Versioning
- [ ] Service Mesh (Istio)
- [ ] Centralized Logging (ELK Stack)
- [ ] Config Management (Vault)

## Документация

- **API документация**: `API_DOCUMENTATION.md`
- **Тестирование**: `SERVICE_TESTING_GUIDE.md`
- **Микросервисы**: `MICROSERVICES_README.md`
- **Postman коллекция**: `smarthome-api.postman_collection.json`
- **AsyncAPI**: `async-api.yaml`
- **Архитектурные диаграммы**: `../diagrams/images/`