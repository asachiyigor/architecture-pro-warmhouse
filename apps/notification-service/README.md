# Notification Service

Микросервис уведомлений - обработка событий системы и доставка уведомлений пользователям через различные каналы (WebSocket, Push, Email).

## Технологический стек

- **Node.js 18** + **Express.js**
- **Socket.IO** для WebSocket соединений
- **Kafka** для потребления событий
- **Firebase Cloud Messaging** для push уведомлений (планируется)

## Возможности

### ✅ Реализованные функции
- **Real-time WebSocket** - мгновенные уведомления через Socket.IO
- **Kafka Consumer** - обработка событий от всех сервисов
- **REST API** - управление уведомлениями
- **Multi-event processing** - обработка разных типов событий
- **Configurable routing** - гибкая маршрутизация уведомлений

### 📊 Поддерживаемые события

#### Device Events (device.events)
```json
{
  "eventType": "device.status.changed",
  "deviceName": "Kitchen Sensor",
  "location": "Kitchen",
  "newStatus": "OFFLINE"
}
```

#### Telemetry Events (telemetry-events)
```json
{
  "event_type": "telemetry_received",
  "device_id": "temp_sensor_001",
  "metric_type": "temperature",
  "value": 23.5,
  "location": "Kitchen"
}
```

#### Telemetry Alerts (telemetry-alerts)  
```json
{
  "event_type": "telemetry_alert",
  "device_id": "temp_sensor_001",
  "alert_type": "high_temperature",
  "severity": "critical",
  "message": "Temperature exceeded critical threshold"
}
```

## API Endpoints

### Health Check
- `GET /health` - Проверка состояния сервиса

### Notifications
- `GET /notifications` - Получить уведомления
- `POST /notifications` - Создать уведомление
- `POST /notifications/send` - Отправить уведомление
- `DELETE /notifications/{id}` - Удалить уведомление

### WebSocket Connection
- `ws://localhost:8088/ws` - WebSocket endpoint для real-time уведомлений

### Доступ через API Gateway
- **Прямой доступ:** http://localhost:8088  
- **Через Gateway:** http://localhost:8000/api/v1/notifications/*
- **WebSocket:** ws://localhost:8088/ws

## WebSocket Integration

### Подключение клиента
```javascript
const socket = io('http://localhost:8088');

socket.on('connect', () => {
  console.log('Connected to notification service');
});

socket.on('notification', (data) => {
  console.log('New notification:', data);
  // Display notification to user
});

socket.on('device_alert', (data) => {
  console.log('Device alert:', data);
  // Handle critical device alerts
});
```

### Типы WebSocket событий
- **notification** - общие уведомления
- **device_alert** - критические алерты устройств
- **telemetry_alert** - алерты телеметрии
- **system_notification** - системные уведомления

## Kafka Integration

### Подписка на топики
```javascript
const topics = [
  'device.events',          // События устройств
  'telemetry-events',       // Телеметрические данные
  'telemetry-alerts',       // Критические алерты
  'processed-events'        // Обработанные события
];
```

### Обработка событий
```javascript
async processMessage(topic, message) {
  switch (topic) {
    case 'device.events':
      notification = this.createDeviceNotification(message);
      break;
    case 'telemetry-alerts':
      notification = this.createTelemetryAlertNotification(message);
      break;
    // ... другие топики
  }
  
  if (notification) {
    await this.sendNotification(notification);
  }
}
```

## Типы уведомлений

### 📱 Device Notifications
- **device_event** - события устройств
- **device_alert** - критические состояния
- **device_offline** - устройство недоступно

### 🌡️ Telemetry Notifications
- **telemetry_alert** - пороговые значения
- **sensor_malfunction** - неисправность датчика
- **data_anomaly** - аномальные показания

### 🏠 System Notifications  
- **service_status** - статус сервисов
- **maintenance** - плановые работы
- **security_alert** - безопасность

## Конфигурация

### Переменные окружения
| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `PORT` | Порт сервиса | `8086` |
| `LOG_LEVEL` | Уровень логирования | `info` |
| `KAFKA_BROKERS` | Kafka брокеры | `smarthome-kafka:9092` |
| `CORS_ORIGINS` | CORS origins | `http://localhost:3000,http://localhost:8080` |
| `FCM_SERVER_KEY` | Firebase key | (планируется) |

### Конфигурационный файл (.env.example)
```bash
# Notification Service Configuration
PORT=8086
LOG_LEVEL=info

# Kafka Configuration  
KAFKA_BROKERS=smarthome-kafka:9092

# CORS Configuration
CORS_ORIGINS=http://localhost:3000,http://localhost:8080

# Push Notifications (Future implementation)
FCM_SERVER_KEY=your_fcm_server_key_here
APNS_KEY_ID=your_apns_key_id_here
APNS_TEAM_ID=your_apns_team_id_here
```

## Структура проекта

```
notification-service/
├── src/
│   ├── config/
│   │   └── logger.js          # Логирование
│   ├── services/
│   │   ├── kafkaConsumer.js    # Kafka потребитель
│   │   └── notificationService.js # Бизнес логика
│   ├── routes/
│   │   └── notifications.js    # REST endpoints
│   └── index.js               # Entry point
├── package.json
├── Dockerfile
└── .env.example
```

## Notification Schema

### Базовая структура уведомления
```json
{
  "id": "uuid",
  "userId": "user_123",
  "type": "device_alert",
  "title": "Device Offline",
  "message": "Kitchen sensor has gone offline",
  "data": {
    "deviceId": "temp_sensor_001",
    "location": "Kitchen",
    "severity": "warning"
  },
  "timestamp": "2025-01-03T10:30:00Z",
  "read": false
}
```

## Запуск

### Docker Compose
```bash
# Запуск с зависимостями
docker-compose up notification-service kafka

# Проверка работы
curl http://localhost:8088/health
curl http://localhost:8088/notifications
```

### Локальная разработка
```bash
cd notification-service
npm install
npm start

# Development режим с hot reload
npm run dev
```

## Мониторинг и Логирование

### Структурированное логирование
```json
{
  "level": "info",
  "timestamp": "2025-01-03T12:50:49.793Z",
  "logger": "kafkajs",
  "message": "Connected to Kafka",
  "service": "notification-service"
}
```

### Метрики
- **WebSocket connections** - активные соединения
- **Messages processed** - обработанные сообщения  
- **Kafka lag** - задержка обработки событий
- **Notification delivery** - успешность доставки

## Планируемые функции

### 📱 Push Notifications
- **Firebase Cloud Messaging** - Android/iOS push
- **Web Push API** - браузерные уведомления
- **Apple Push Notification** - iOS нативные уведомления

### 📧 Email Notifications
- **SMTP integration** - email рассылка
- **Template engine** - шаблоны уведомлений
- **Unsubscribe management** - управление подписками

## Архитектурная роль

Notification Service является центральным hub для:
- **Real-time Communication** - мгновенные уведомления пользователям  
- **Multi-channel Delivery** - доставка через разные каналы
- **User Experience** - информирование о состоянии системы
- **System Monitoring** - алерты для администраторов