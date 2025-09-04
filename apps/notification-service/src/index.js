const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const { createServer } = require('http');
const { Server } = require('socket.io');
const { v4: uuidv4 } = require('uuid');
const swaggerUi = require('swagger-ui-express');
require('dotenv').config();

const logger = require('./config/logger');
const swaggerSpecs = require('./config/swagger');
const kafkaConsumer = require('./services/kafkaConsumer');
const notificationService = require('./services/notificationService');

const app = express();
const httpServer = createServer(app);

// Socket.IO setup for WebSocket connections
const io = new Server(httpServer, {
  cors: {
    origin: process.env.CORS_ORIGINS?.split(',') || ["http://localhost:3000"],
    methods: ["GET", "POST"]
  }
});

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());

// Swagger UI setup
app.use('/api-docs', swaggerUi.serve, swaggerUi.setup(swaggerSpecs, {
  customCss: '.swagger-ui .topbar { display: none }',
  customSiteTitle: 'Smart Home Notification Service API',
}));

/**
 * @swagger
 * /:
 *   get:
 *     summary: Service information
 *     description: Получить информацию о сервисе уведомлений
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Service information
 *         content:
 *           application/json:
 *             schema:
 *               type: object
 *               properties:
 *                 service:
 *                   type: string
 *                   example: "Smart Home Notification Service"
 *                 version:
 *                   type: string
 *                   example: "1.0.0"
 *                 description:
 *                   type: string
 *                   example: "WebSocket and Kafka-based notification service"
 *                 endpoints:
 *                   type: object
 *                   properties:
 *                     docs:
 *                       type: string
 *                       example: "/api-docs"
 *                     health:
 *                       type: string
 *                       example: "/health"
 *                     notifications:
 *                       type: string
 *                       example: "/api/v1/notifications"
 *                     stats:
 *                       type: string
 *                       example: "/api/v1/notifications/stats"
 */
app.get('/', (req, res) => {
  res.json({
    service: 'Smart Home Notification Service',
    version: '1.0.0',
    description: 'WebSocket and Kafka-based notification service for smart home system',
    endpoints: {
      docs: '/api-docs',
      health: '/health',
      notifications: '/api/v1/notifications',
      stats: '/api/v1/notifications/stats'
    }
  });
});

/**
 * @swagger
 * /health:
 *   get:
 *     summary: Health check endpoint
 *     description: Проверить состояние сервиса уведомлений
 *     tags: [Health]
 *     responses:
 *       200:
 *         description: Service is healthy
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/HealthStatus'
 */
app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    service: 'notification-service',
    timestamp: new Date().toISOString()
  });
});

// WebSocket connection handling
io.on('connection', (socket) => {
  const sessionId = uuidv4();
  logger.info(`Client connected: ${socket.id}, session: ${sessionId}`);
  
  // Join client to a room for targeted notifications
  socket.on('join-user', (userId) => {
    socket.join(`user:${userId}`);
    logger.info(`Client ${socket.id} joined user room: ${userId}`);
  });
  
  socket.on('disconnect', () => {
    logger.info(`Client disconnected: ${socket.id}`);
  });
});

// Initialize notification service with Socket.IO
notificationService.initialize(io);

/**
 * @swagger
 * /api/v1/notifications:
 *   post:
 *     summary: Send a notification
 *     description: Отправить уведомление пользователю через WebSocket и сохранить в системе
 *     tags: [Notifications]
 *     requestBody:
 *       required: true
 *       content:
 *         application/json:
 *           schema:
 *             $ref: '#/components/schemas/NotificationRequest'
 *           examples:
 *             device_alert:
 *               summary: Device alert notification
 *               value:
 *                 userId: "user123"
 *                 type: "device_alert"
 *                 title: "Устройство отключено"
 *                 message: "Датчик температуры в гостиной перестал отвечать"
 *                 data:
 *                   deviceId: "temp_sensor_001"
 *                   location: "Living Room"
 *             automation_event:
 *               summary: Automation event notification
 *               value:
 *                 userId: "user123"
 *                 type: "automation_event"
 *                 title: "Сценарий выполнен"
 *                 message: "Сценарий 'Уходим из дома' успешно выполнен"
 *                 data:
 *                   scenarioId: "scenario_001"
 *                   executionTime: "2025-09-04T10:30:00Z"
 *     responses:
 *       200:
 *         description: Notification sent successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Success'
 *       400:
 *         description: Missing required fields
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Error'
 *       500:
 *         description: Failed to send notification
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/Error'
 */
app.post('/api/v1/notifications', async (req, res) => {
  try {
    const { userId, type, title, message, data } = req.body;
    
    if (!userId || !type || !title || !message) {
      return res.status(400).json({ 
        error: 'Missing required fields: userId, type, title, message' 
      });
    }

    await notificationService.sendNotification({
      userId,
      type,
      title,
      message,
      data: data || {}
    });

    res.status(200).json({ 
      success: true,
      message: 'Notification sent successfully' 
    });
    
  } catch (error) {
    logger.error('Error sending notification:', error);
    res.status(500).json({ 
      error: 'Failed to send notification' 
    });
  }
});

/**
 * @swagger
 * /api/v1/notifications/stats:
 *   get:
 *     summary: Get notification statistics
 *     description: Получить статистику отправленных уведомлений и активных соединений
 *     tags: [Statistics]
 *     responses:
 *       200:
 *         description: Statistics retrieved successfully
 *         content:
 *           application/json:
 *             schema:
 *               $ref: '#/components/schemas/NotificationStats'
 *             example:
 *               totalSent: 156
 *               totalByType:
 *                 device_alert: 45
 *                 automation_event: 38
 *                 info: 73
 *               activeConnections: 12
 *               uptime: "2h 34m 12s"
 */
app.get('/api/v1/notifications/stats', (req, res) => {
  const stats = notificationService.getStats();
  res.json(stats);
});

// Start Kafka consumer for processing events
kafkaConsumer.start().catch((error) => {
  logger.error('Failed to start Kafka consumer:', error);
  process.exit(1);
});

// Start HTTP server
const PORT = process.env.PORT || 8086;
httpServer.listen(PORT, () => {
  logger.info(`Notification Service started on port ${PORT}`);
  logger.info('WebSocket server ready for connections');
});

// Graceful shutdown
process.on('SIGTERM', async () => {
  logger.info('SIGTERM received, shutting down gracefully');
  await kafkaConsumer.stop();
  httpServer.close(() => {
    logger.info('Server stopped');
    process.exit(0);
  });
});