const swaggerJSDoc = require('swagger-jsdoc');

const options = {
  definition: {
    openapi: '3.0.0',
    info: {
      title: 'Smart Home Notification Service API',
      version: '1.0.0',
      description: 'Микросервис уведомлений для умного дома с поддержкой WebSocket и Kafka',
      contact: {
        name: 'Smart Home Team',
        email: 'support@smarthome.com'
      }
    },
    servers: [
      {
        url: 'http://localhost:8088',
        description: 'Development server'
      }
    ],
    components: {
      schemas: {
        Notification: {
          type: 'object',
          properties: {
            id: {
              type: 'string',
              format: 'uuid',
              description: 'Unique notification ID'
            },
            userId: {
              type: 'string',
              description: 'User ID who will receive the notification'
            },
            type: {
              type: 'string',
              enum: ['info', 'warning', 'error', 'success', 'device_alert', 'automation_event'],
              description: 'Type of notification'
            },
            title: {
              type: 'string',
              description: 'Notification title'
            },
            message: {
              type: 'string',
              description: 'Notification message content'
            },
            data: {
              type: 'object',
              description: 'Additional data payload',
              additionalProperties: true
            },
            timestamp: {
              type: 'string',
              format: 'date-time',
              description: 'When the notification was created'
            },
            read: {
              type: 'boolean',
              description: 'Whether the notification has been read'
            }
          },
          required: ['userId', 'type', 'title', 'message']
        },
        NotificationRequest: {
          type: 'object',
          properties: {
            userId: {
              type: 'string',
              description: 'User ID who will receive the notification'
            },
            type: {
              type: 'string',
              enum: ['info', 'warning', 'error', 'success', 'device_alert', 'automation_event'],
              description: 'Type of notification'
            },
            title: {
              type: 'string',
              description: 'Notification title'
            },
            message: {
              type: 'string',
              description: 'Notification message content'
            },
            data: {
              type: 'object',
              description: 'Additional data payload',
              additionalProperties: true
            }
          },
          required: ['userId', 'type', 'title', 'message']
        },
        NotificationStats: {
          type: 'object',
          properties: {
            totalSent: {
              type: 'integer',
              description: 'Total notifications sent'
            },
            totalByType: {
              type: 'object',
              additionalProperties: {
                type: 'integer'
              },
              description: 'Count of notifications by type'
            },
            activeConnections: {
              type: 'integer',
              description: 'Number of active WebSocket connections'
            },
            uptime: {
              type: 'string',
              description: 'Service uptime'
            }
          }
        },
        HealthStatus: {
          type: 'object',
          properties: {
            status: {
              type: 'string',
              enum: ['ok', 'error'],
              description: 'Service status'
            },
            service: {
              type: 'string',
              description: 'Service name'
            },
            timestamp: {
              type: 'string',
              format: 'date-time',
              description: 'Health check timestamp'
            }
          }
        },
        Error: {
          type: 'object',
          properties: {
            error: {
              type: 'string',
              description: 'Error message'
            }
          }
        },
        Success: {
          type: 'object',
          properties: {
            success: {
              type: 'boolean'
            },
            message: {
              type: 'string'
            }
          }
        }
      }
    },
    tags: [
      {
        name: 'Health',
        description: 'Health check endpoints'
      },
      {
        name: 'Notifications',
        description: 'Notification management endpoints'
      },
      {
        name: 'Statistics',
        description: 'Service statistics endpoints'
      }
    ]
  },
  apis: ['./src/index.js', './src/routes/*.js']
};

const specs = swaggerJSDoc(options);
module.exports = specs;