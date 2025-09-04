# Project_template

Это шаблон для решения проектной работы. Структура этого файла повторяет структуру заданий. Заполняйте его по мере работы над решением.

# Задание 1. Анализ и планирование

### 1. Описание функциональности монолитного приложения

**Управление отоплением:**

- Пользователи могут удалённо управлять системой отопления в своём доме
- Пользователи могут включать и выключать отопление
- Система поддерживает синхронное взаимодействие с датчиками отопления
- Система отправляет команды управления от сервера к датчикам отопления
- Для подключения системы требуется выезд специалиста

**Мониторинг температуры:**

- Пользователи могут отслеживать текущую температуру в доме в реальном времени
- Система поддерживает получение данных о температуре через синхронные запросы от сервера к датчику
- Система предоставляет веб-интерфейс для просмотра температурных данных
- Данные о температуре получаются по запросу, а не передаются автоматически

### 2. Анализ архитектуры монолитного приложения

- Язык программирования: Go
- База данных: PostgreSQL
- Архитектурный стиль: Монолитное приложение
- Взаимодействие: Полностью синхронное, без асинхронных вызовов
- Клиенты: 100 веб-клиентов
- Устройства: 100 модулей управления отоплением

### 3. Определение доменов и границы контекстов

1.  Управление Устройствами (Device Management)

- Управление отоплением (включение/выключение)
- Взаимодействие с датчиками температуры
- Взаимодействие с реле управления

2.  Управление пользователями (User Management)

- Аутентификация и авторизация пользователей
- Управление доступом к устройствам
- Привязка пользователей к домам/устройствам

3. Мониторинг и Аналитика (Monitoring & Analytics)

- Сбор данных о температуре
- Мониторинг состояния устройств
- Предоставление данных для отображения
  
4. Интеграция и Коммуникации (Integration & Communications)

- Протоколы взаимодействия с устройствами
- API для внешних систем
- Обработка команд и ответов

### 4. Визуализация контекста системы — диаграмма С4

[Диаграмма контекста системы "Тёплый дом"](./diagrams/images/00_context_diagram_c4_mono.png)

### 5. Проблемы монолитного решения

1. Единая точка отказа
2. Невозможность независимого масштабирования компонентов
3. Синхронные блокировки
4. Сложность изменений и релизов
5. Сложно добавлять новые устройства


# Задание 2. Проектирование микросервисной архитектуры

В этом задании вам нужно предоставить только диаграммы в модели C4. Мы не просим вас отдельно описывать получившиеся микросервисы и то, как вы определили взаимодействия между компонентами To-Be системы. Если вы правильно подготовите диаграммы C4, они и так это покажут.

- **Диаграмма контейнеров [(Containers)](./diagrams/images/02.container_diagram.png)**

- **Диаграмма компонентов [(Components)](./diagrams/images/03_component_diagram_device_service_to_be.png)**

- **Диаграмма кода [(Code)](./diagrams/images/04_code_device_service.png)**

# Задание 3. Разработка ER-диаграммы

**Полная ER-диаграмма [системы](./diagrams/images/er_diagram_complete_system.png):**

**ER-диаграммы по доменам:**

- [Домен пользователей](./diagrams/images/er_diagram_user_domain.png)
- [Домен устройств](./diagrams/images/er_diagram_device_domain.png)
- [Домен телеметрии](./diagrams/images/er_diagram_telemetry_domain.png)
- [Домен автоматизации](./diagrams/images/er_diagram_automation_domain.png)
- [Домен биллинга](./diagrams/images/er_diagram_billing_domain.png)

# Задание 4. Создание и документирование API

### 1. Тип API

**REST API** для синхронного взаимодействия через HTTP/HTTPS:
- CRUD операции для управления сущностями
- Простота интеграции с веб-клиентами
- Стандартизированные HTTP коды ответов
- JSON формат данных

**Event-Driven Architecture** через Apache Kafka:
- Асинхронное взаимодействие между микросервисами
- Обработка событий (device.events, telemetry-events, telemetry-alerts)
- Масштабируемость и отказоустойчивость
- Слабая связанность компонентов

### 2. Документация API

**Основная API документация:**
- **[🌟 API_DOCUMENTATION.html](API_DOCUMENTATION.html)** - **интерактивная консолидированная документация** (открыть в браузере)
- **[🎯 CONSOLIDATED_API.yaml](CONSOLIDATED_API.yaml)** - **OpenAPI спецификация всех микросервисов** (как Yandex API)
- [API Documentation](apps/API_DOCUMENTATION.md) - текстовая документация всех микросервисов

**Swagger/OpenAPI документация по сервисам:**

✅ **Полная Swagger документация реализована:**
- **[API Gateway](http://localhost:8000/swagger/index.html)** - swaggo/gin-swagger
- **[Temperature API](http://localhost:8081/swagger/index.html)** - swaggo/gin-swagger
- **[Device Service](http://localhost:8082/swagger-ui/swagger-ui/index.html)** - SpringDoc OpenAPI
- **[Telemetry Service](http://localhost:8083/docs)** - FastAPI автодокументация  
- **[User Service](http://localhost:8084/swagger/index.html)** - swaggo/gin-swagger
- **[Billing Service](http://localhost:8085/swagger/index.html)** - swaggo/gin-swagger
- **[Automation Service](http://localhost:8087/swagger/index.html)** - swaggo/gin-swagger
- **[Notification Service](http://localhost:8088/api-docs/)** - swagger-ui-express

✅ **Статическая документация:**
- **[Smart Home App](apps/smart_home/swagger.yaml)** - OpenAPI 3.0.3 статический файл (513 строк)

**Event Schemas:**
- [AsyncAPI документация](apps/async-api.yaml) - схемы Kafka событий

# **✅ Задание 5. Работа с docker и docker-compose**

### 1. Микросервисная архитектура в Docker

**Созданы и докеризованы все микросервисы:**

| Сервис | Технология | Порт | Docker контейнер | Swagger документация |
|--------|------------|------|------------------|----------------------|
| **API Gateway** | Go Gin | 8000 | ✅ `api-gateway` | ✅ http://localhost:8000/swagger/index.html |
| **Temperature API** | Go Gin | 8081 | ✅ `temperature-api` | ✅ http://localhost:8081/swagger/index.html |
| **Device Service** | Java Spring Boot | 8082 | ✅ `device-service` | ✅ http://localhost:8082/swagger-ui/swagger-ui/index.html |
| **Telemetry Service** | Python FastAPI | 8083 | ✅ `telemetry-service` | ✅ http://localhost:8083/docs |
| **User Service** | Go Gin | 8084 | ✅ `user-service` | ✅ http://localhost:8084/swagger/index.html |
| **Billing Service** | Go Gin | 8085 | ✅ `billing-service` | ✅ http://localhost:8085/swagger/index.html |
| **Smart Home App** | Go | 8080 | ✅ `smarthome-app` | ✅ [swagger.yaml](apps/smart_home/swagger.yaml) |
| **Automation Service** | Go Gin | 8087 | ✅ `automation-service` | ✅ http://localhost:8087/swagger/index.html |
| **Notification Service** | Node.js Express | 8088 | ✅ `notification-service` | ✅ http://localhost:8088/api-docs/ |

### 2. Полная инфраструктура в Docker Compose

**Базы данных (Database-per-Service pattern):**
- ✅ **PostgreSQL для User Service** - порт 5433
- ✅ **PostgreSQL для Device Service** - порт 5434  
- ✅ **PostgreSQL для Telemetry Service** - порт 5435
- ✅ **PostgreSQL для Automation Service** - порт 5436
- ✅ **PostgreSQL для Billing Service** - порт 5437
- ✅ **PostgreSQL для Smart Home App** - порт 5432
- ✅ **InfluxDB** для временных рядов телеметрии - порт 8086
- ✅ **Redis** для кеширования - порт 6379

**Брокер сообщений:**
- ✅ **Apache Kafka** - порт 9092
- ✅ **Zookeeper** - порт 2181

### 3. Запуск и тестирование

**Запуск всей системы:**
```bash
cd apps
docker-compose up --build

# Проверка статуса всех контейнеров:
docker-compose ps
```

**API тестирование:**
```bash
# API Gateway
curl http://localhost:8000/gateway/status

# Temperature API
curl http://localhost:8081/temperature?location=Kitchen

# Device Service  
curl http://localhost:8082/api/v1/devices

# Telemetry Service
curl http://localhost:8083/telemetry/readings

# User Service
curl http://localhost:8084/api/v1/users

# Billing Service
curl http://localhost:8085/api/v1/plans

# Automation Service  
curl http://localhost:8087/api/v1/stats

# Notification Service
curl http://localhost:8088/health
```

# **✅ Задание 6. Разработка MVP**

## Полностью реализована Smart Home Pro Microservices Architecture

Создана enterprise-grade микросервисная платформа с полной документацией API и готовностью к production.

### **Реализованные микросервисы (9 сервисов):**

#### 🌐 **API Gateway** (Go Gin) - порт 8000
- Центральная точка входа (Single Entry Point)
- Прозрачная маршрутизация к микросервисам
- CORS, Rate Limiting, Request/Response логирование
- **Swagger**: http://localhost:8000/swagger/index.html
- Аутентификация и авторизация (JWT)

#### 🌡️ **Temperature API** (Go Gin) - порт 8081  
- Эмуляция IoT датчиков температуры
- Поддержка множественных локаций
- **Swagger**: http://localhost:8081/swagger/index.html
- Интеграция с Telemetry Service

#### 🔧 **Device Service** (Java Spring Boot) - порт 8082
- Управление IoT устройствами (CRUD)
- Heartbeat мониторинг устройств
- Статистика и типы устройств
- **SpringDoc OpenAPI**: http://localhost:8082/swagger-ui/swagger-ui/index.html
- PostgreSQL с JPA/Hibernate
- Kafka producer для device.events

#### 📊 **Telemetry Service** (Python FastAPI) - порт 8083
- Сбор и агрегация телеметрии в реальном времени
- InfluxDB для временных рядов
- Автоматическая очистка старых данных
- **FastAPI docs**: http://localhost:8083/docs
- Background tasks для сбора данных
- Kafka producer для telemetry-events и telemetry-alerts

#### 👤 **User Service** (Go Gin) - порт 8084
- Управление пользователями и домохозяйствами
- JWT аутентификация и авторизация  
- RBAC (Role-Based Access Control)
- **Swagger**: http://localhost:8084/swagger/index.html
- PostgreSQL с GORM

#### 💳 **Billing Service** (Go Gin) - порт 8085
- Тарифные планы и управление подписками
- Процессинг платежей
- Биллинговая аналитика
- **Swagger**: http://localhost:8085/swagger/index.html
- PostgreSQL для финансовых данных

#### 🔄 **Automation Service** (Go Gin) - порт 8087
- Пользовательские сценарии автоматизации
- Триггеры и действия (Rules Engine)
- История выполнения сценариев
- **Swagger**: http://localhost:8087/swagger/index.html
- PostgreSQL для сценариев

#### 📱 **Notification Service** (Node.js Express) - порт 8088
- Real-time уведомления через WebSocket (Socket.IO)
- Push notifications (планируется)
- Email уведомления (планируется)
- **Swagger UI**: http://localhost:8088/api-docs/
- Kafka consumer для всех событий системы

#### 🏠 **Smart Home App** (Go) - порт 8080
- Монолитное приложение (Legacy) с CRUD для датчиков
- Интеграция с Temperature API и API Gateway
- **OpenAPI 3.0.3**: [swagger.yaml](apps/smart_home/swagger.yaml) (513 строк)
- PostgreSQL база данных с подробной схемой
- Плавная миграция в микросервисы

### **Архитектурные паттерны:**

#### 🔄 **Event-Driven Architecture (Apache Kafka)**
- **Kafka Topics**: `device.events`, `telemetry-events`, `telemetry-alerts`, `processed-events`
- **Asynchronous messaging** между микросервисами
- **Event sourcing** для аудита и воспроизведения состояний
- **CQRS pattern** в Telemetry Service

#### 🌐 **API Gateway Pattern**  
- **Single Entry Point** для всех внешних клиентов
- **Service Discovery** и load balancing
- **Cross-cutting concerns**: аутентификация, логирование, rate limiting
- **Backward compatibility** с монолитным приложением

#### 🗄️ **Database-per-Service Pattern**
- **6 отдельных PostgreSQL** инстансов для каждого домена
- **InfluxDB** для временных рядов телеметрии
- **Redis** для кеширования и сессий
- **Data consistency** через Saga pattern

#### 📡 **Real-time Communication**  
- **WebSocket** уведомления через Socket.IO
- **Server-Sent Events** для стриминга телеметрии
- **Push notifications** (готовность к мобильным приложениям)

### **Полная Production-ready инфраструктура:**

#### 🐳 **Containerization & Orchestration**
```bash
# Запуск всей платформы одной командой
cd apps && docker-compose up --build

# 19 контейнеров работают в связке:
# - 9 микросервисов
# - 6 PostgreSQL баз данных  
# - 1 InfluxDB + 1 Redis
# - 1 Kafka + 1 Zookeeper
```

#### 📊 **Monitoring & Observability**
- **Health checks** на всех endpoints `/health`
- **Metrics endpoints** `/metrics` для Prometheus
- **Structured logging** в JSON формате
- **Distributed tracing** готовность (correlation IDs)

#### 🔒 **Security & Compliance**
- **JWT authentication** в User Service
- **RBAC authorization** (Role-Based Access Control)  
- **API keys** для внешних интеграций
- **Rate limiting** в API Gateway
- **HTTPS/TLS** ready

#### 📚 **Complete API Documentation**
- **9 из 9 сервисов** имеют полную OpenAPI/Swagger документацию
- **8 интерактивных Swagger UI** + 1 статический OpenAPI файл
- **Request/Response examples** с реальными данными  
- **513+ строк** детальной документации API
- **AsyncAPI** схемы для Kafka событий

### **MVP Результаты:**

✅ **Функциональная полнота**: Все основные домены покрыты  
✅ **Scalability**: Горизонтальное масштабирование каждого сервиса  
✅ **Resilience**: Graceful degradation, circuit breakers  
✅ **Developer Experience**: Полная документация + hot reload  
✅ **Production readiness**: Monitoring, logging, security  
✅ **Migration strategy**: Плавный переход от монолита

**Итог**: Enterprise-grade микросервисная платформа Smart Home Pro готова для production deployment с возможностью обслуживать тысячи пользователей и устройств. 