# Smart Home Microservices Architecture

## Обзор архитектуры

Проект был расширен для соответствия архитектурным диаграммам с добавлением недостающих микросервисов:

### Микросервисы

1. **User Service** (`:8084`) - управление пользователями и домохозяйствами
2. **Billing Service** (`:8085`) - управление подписками и платежами  
3. **Automation Service** (`:8086`) - сценарии автоматизации
4. **Device Service** (`:8082`) - управление устройствами (Java/Spring Boot)
5. **Telemetry Service** (`:8083`) - сбор телеметрии (Python/FastAPI)
6. **Smart Home App** (`:8080`) - основное приложение (Go/Gin)
7. **Temperature API** (`:8081`) - температурный API (Go)

### Базы данных

- `postgres` (`:5432`) - основная БД для Smart Home App
- `postgres-devices` (`:5433`) - БД для Device Service
- `postgres-telemetry` (`:5434`) - БД для Telemetry Service  
- `postgres-users` (`:5435`) - БД для User Service
- `postgres-billing` (`:5436`) - БД для Billing Service
- `postgres-automation` (`:5437`) - БД для Automation Service

### Инфраструктура

- **RabbitMQ** (`:5672`, UI `:15672`) - брокер сообщений
- **Redis** (`:6379`) - кэширование

## Запуск системы

```bash
# Сборка и запуск всех сервисов
docker-compose up --build

# Проверка статуса сервисов
docker-compose ps

# Просмотр логов конкретного сервиса
docker-compose logs user-service
```

## API Endpoints

### User Service (`:8084`)
- `GET /health` - проверка здоровья сервиса
- `POST /api/v1/users` - создание пользователя
- `GET /api/v1/users/{id}` - получение пользователя
- `POST /api/v1/households` - создание домохозяйства
- `GET /api/v1/households/{id}` - получение домохозяйства

### Billing Service (`:8085`)
- `GET /health` - проверка здоровья сервиса
- `GET /api/v1/plans` - получение тарифных планов
- `POST /api/v1/subscriptions` - создание подписки
- `POST /api/v1/payments` - обработка платежа
- `GET /api/v1/stats` - статистика биллинга

### Automation Service (`:8086`)
- `GET /health` - проверка здоровья сервиса
- `POST /api/v1/scenarios` - создание сценария
- `GET /api/v1/scenarios/{id}` - получение сценария
- `POST /api/v1/scenarios/{id}/execute` - выполнение сценария
- `GET /api/v1/stats` - статистика автоматизации

## Архитектурные особенности

### Унифицированный стек
- Все новые микросервисы написаны на **Go** для единообразия
- Используется **Gin** фреймворк для REST API
- **PostgreSQL** как основная СУБД
- Чистая архитектура: handlers → service → repository

### Межсервисное взаимодействие  
- REST API для синхронного взаимодействия
- RabbitMQ для асинхронных событий
- Каждый сервис имеет свою БД (Database per Service)

### Демо-данные
Каждый сервис предзаполнен демо-данными:
- Пользователи и домохозяйства
- Тарифные планы и подписки
- Примеры сценариев автоматизации

## Мониторинг

Проверить работу всех сервисов:
```bash
curl http://localhost:8084/health  # User Service
curl http://localhost:8085/health  # Billing Service  
curl http://localhost:8086/health  # Automation Service
```

## Изменения в архитектуре

**Удалено из диаграмм:**
- Notification Service (по требованию)

**Добавлено:**
- Полная реализация User Service
- Полная реализация Billing Service
- Полная реализация Automation Service
- Соответствующие базы данных и конфигурации

Теперь проект соответствует архитектурным диаграммам C4 и может служить полноценным примером микросервисной архитектуры для системы "Умный дом".