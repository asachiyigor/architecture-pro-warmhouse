# Automation Service

Микросервис автоматизации умного дома - создание и выполнение пользовательских сценариев для управления IoT устройствами.

## Технологический стек

- **Go 1.21** + **Gin Framework**
- **PostgreSQL** для хранения сценариев
- **Kafka интеграция** для событийной архитектуры (планируется)
- **REST API** для управления сценариями

## Возможности

### ✅ Реализованные функции
- **Создание сценариев** - пользовательские automation скрипты
- **Управление действиями** - последовательность команд устройствам
- **Ручное выполнение** - запуск сценариев по требованию
- **CRUD операции** - полное управление сценариями
- **Связь с домохозяйствами** - сценарии привязаны к домам

### 📊 Модели данных

#### Scenario Model
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

#### ScenarioAction Model
```json
{
  "id": "uuid",
  "scenario_id": "uuid",
  "device_id": "light_living_room",
  "action_type": "turn_off",
  "parameters": "{\"brightness\": 0}",
  "order_position": 1,
  "created_at": "2025-01-03T10:00:00Z"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### Scenario Management
- `GET /scenarios` - Получить все сценарии
- `GET /scenarios/{id}` - Получить сценарий по ID
- `POST /scenarios` - Создать новый сценарий
- `PUT /scenarios/{id}` - Обновить сценарий
- `DELETE /scenarios/{id}` - Удалить сценарий
- `POST /scenarios/{id}/execute` - Выполнить сценарий

### Action Management
- `GET /scenarios/{scenario_id}/actions` - Действия сценария
- `POST /scenarios/{scenario_id}/actions` - Добавить действие
- `PUT /actions/{id}` - Обновить действие
- `DELETE /actions/{id}` - Удалить действие

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8087
- **Через Gateway:** http://localhost:8000/api/v1/automation/*

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | PostgreSQL connection | `postgres://postgres:postgres@postgres-automation:5432/smarthome_automation?sslmode=disable` |
| `PORT` | Порт сервиса | `:8087` |
| `KAFKA_BROKERS` | Kafka брокеры | (планируется) |

### База данных
- **PostgreSQL** на порту 5437 (docker)
- **Схема:** автоматические миграции через `sql/init.sql`

## Схема базы данных

### Scenarios Table
```sql
CREATE TABLE scenarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    household_id UUID NOT NULL,
    created_by UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### Scenario_Actions Table
```sql
CREATE TABLE scenario_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scenario_id UUID NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    device_id VARCHAR(100) NOT NULL,
    action_type VARCHAR(50) NOT NULL,
    parameters TEXT,
    order_position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

## Типы действий

### ✅ Базовые действия (MVP)
- **turn_on** - включить устройство
- **turn_off** - выключить устройство  
- **set_brightness** - установить яркость (для ламп)
- **set_temperature** - установить температуру (для термостатов)
- **lock** - заблокировать (для замков)
- **unlock** - разблокировать (для замков)

## Примеры сценариев

### "Вечерний режим"
```json
{
  "name": "Evening Scene",
  "description": "Prepare house for evening",
  "actions": [
    {
      "device_id": "living_room_light",
      "action_type": "set_brightness", 
      "parameters": "{\"brightness\": 30}",
      "order_position": 1
    },
    {
      "device_id": "front_door_lock",
      "action_type": "lock",
      "parameters": "{}",
      "order_position": 2
    }
  ]
}
```

### "Уход из дома"
```json
{
  "name": "Leaving Home",
  "description": "Turn off everything when leaving",
  "actions": [
    {
      "device_id": "all_lights",
      "action_type": "turn_off",
      "parameters": "{}",
      "order_position": 1
    },
    {
      "device_id": "all_plugs", 
      "action_type": "turn_off",
      "parameters": "{}",
      "order_position": 2
    }
  ]
}
```

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up automation-service postgres-automation

# Проверка работы
curl http://localhost:8087/health
curl http://localhost:8087/scenarios
```

### Локальная разработка
```bash
cd automation-service
go mod tidy
go run cmd/server/main.go
```

## Kafka интеграция (Планируется)

### Планируемые события
- **Topic:** `automation.events`
- **Event Types:**
  - `scenario_executed` - сценарий выполнен
  - `scenario_failed` - ошибка выполнения
  - `action_completed` - действие завершено

### Event Schema (Планируется)
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

## Связи с другими сервисами

### Зависимости:
- **User Service** - создатели сценариев
- **Household Service** - привязка к домам
- **Device Service** - целевые устройства (планируется)

### Интеграции:
- **API Gateway** - маршрутизация запросов
- **Notification Service** - уведомления о выполнении (планируется)

## Мониторинг

- **Health endpoint** - `/health`
- **Execution logging** - логи выполнения сценариев
- **Performance metrics** - время выполнения
- **Error tracking** - отслеживание ошибок

## Архитектурная роль

Automation Service является ключевым для:
- **Smart Home UX** - автоматизация рутинных задач
- **Device Orchestration** - координация устройств
- **User Personalization** - персональные сценарии
- **Event-driven Actions** - реакция на события системы