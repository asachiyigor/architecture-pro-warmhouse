const { Kafka } = require('kafkajs');
const logger = require('../config/logger');
const notificationService = require('./notificationService');

class KafkaConsumer {
  constructor() {
    this.kafka = new Kafka({
      clientId: 'notification-service',
      brokers: (process.env.KAFKA_BROKERS || 'localhost:9092').split(','),
      retry: {
        retries: 3,
        initialRetryTime: 300,
        maxRetryTime: 30000
      }
    });
    
    this.consumer = this.kafka.consumer({ 
      groupId: 'notification-service-group',
      sessionTimeout: 30000,
      heartbeatInterval: 3000,
    });
    
    this.isRunning = false;
  }

  async start() {
    try {
      logger.info('Starting Kafka consumer...');
      
      await this.consumer.connect();
      logger.info('Connected to Kafka');

      // Subscribe to relevant topics for notifications
      await this.consumer.subscribe({
        topics: [
          'device.events',
          'telemetry-events', 
          'telemetry-alerts',
          'processed-events' // From Event Processor (if implemented)
        ],
        fromBeginning: false
      });

      await this.consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
          try {
            const value = JSON.parse(message.value.toString());
            logger.info(`Received message from topic ${topic}:`, value);
            
            await this.processMessage(topic, value);
            
          } catch (error) {
            logger.error('Error processing Kafka message:', error);
          }
        }
      });

      this.isRunning = true;
      logger.info('Kafka consumer is running');
      
    } catch (error) {
      logger.error('Failed to start Kafka consumer:', error);
      throw error;
    }
  }

  async processMessage(topic, message) {
    try {
      let notification = null;

      switch (topic) {
        case 'device.events':
          notification = this.createDeviceNotification(message);
          break;
          
        case 'telemetry-alerts':
          notification = this.createTelemetryAlertNotification(message);
          break;
          
        case 'telemetry-events':
          // Only create notifications for significant telemetry events
          if (message.event_type === 'device_status_changed') {
            notification = this.createStatusChangeNotification(message);
          }
          break;
          
        case 'processed-events':
          notification = this.createProcessedEventNotification(message);
          break;
          
        default:
          logger.debug(`Unhandled topic: ${topic}`);
      }

      if (notification) {
        await notificationService.sendNotification(notification);
        logger.info('Notification sent for event:', { topic, eventType: message.event_type });
      }
      
    } catch (error) {
      logger.error('Error processing message:', error);
    }
  }

  createDeviceNotification(message) {
    const { eventType, deviceName, location } = message;
    
    switch (eventType) {
      case 'device.created':
        return {
          userId: message.userId || 'broadcast',
          type: 'device_event',
          title: 'New Device Added',
          message: `Device "${deviceName}" has been added to ${location}`,
          data: { deviceId: message.deviceId, location }
        };
        
      case 'device.deleted':
        return {
          userId: message.userId || 'broadcast',
          type: 'device_event',
          title: 'Device Removed',
          message: `Device "${deviceName}" has been removed from ${location}`,
          data: { deviceId: message.deviceId, location }
        };
        
      case 'device.status.changed':
        if (message.newStatus === 'OFFLINE') {
          return {
            userId: message.userId || 'broadcast',
            type: 'device_alert',
            title: 'Device Offline',
            message: `Device "${deviceName}" in ${location} has gone offline`,
            data: { deviceId: message.deviceId, location, status: message.newStatus }
          };
        }
        break;
    }
    
    return null;
  }

  createTelemetryAlertNotification(message) {
    return {
      userId: 'broadcast', // Send to all users for now
      type: 'telemetry_alert',
      title: `${message.alert_type} Alert`,
      message: message.message,
      data: {
        deviceId: message.device_id,
        severity: message.severity,
        currentValue: message.current_value,
        thresholdValue: message.threshold_value
      }
    };
  }

  createStatusChangeNotification(message) {
    if (message.new_status === 'offline') {
      return {
        userId: 'broadcast',
        type: 'device_status',
        title: 'Device Status Change',
        message: `Device ${message.device_id} is now ${message.new_status}`,
        data: {
          deviceId: message.device_id,
          oldStatus: message.old_status,
          newStatus: message.new_status
        }
      };
    }
    return null;
  }

  createProcessedEventNotification(message) {
    // Handle processed events from Event Processor
    return {
      userId: 'broadcast',
      type: 'processed_event',
      title: 'System Notification',
      message: message.description || 'System event processed',
      data: message
    };
  }

  async stop() {
    if (this.isRunning) {
      logger.info('Stopping Kafka consumer...');
      await this.consumer.disconnect();
      this.isRunning = false;
      logger.info('Kafka consumer stopped');
    }
  }
}

module.exports = new KafkaConsumer();