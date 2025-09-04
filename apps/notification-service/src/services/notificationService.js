const logger = require('../config/logger');

class NotificationService {
  constructor() {
    this.io = null;
    this.stats = {
      totalSent: 0,
      webSocketSent: 0,
      pushSent: 0,
      errors: 0,
      startTime: new Date()
    };
  }

  initialize(socketIO) {
    this.io = socketIO;
    logger.info('Notification service initialized with Socket.IO');
  }

  async sendNotification(notification) {
    try {
      const { userId, type, title, message, data } = notification;
      
      logger.info('Sending notification:', { userId, type, title });
      
      // Create notification object
      const notificationPayload = {
        id: this.generateNotificationId(),
        type,
        title,
        message,
        data: data || {},
        timestamp: new Date().toISOString(),
        read: false
      };

      // Send via WebSocket to web applications
      await this.sendWebSocketNotification(userId, notificationPayload);
      
      // TODO: Send push notifications to mobile apps (FCM/APNS)
      // await this.sendPushNotification(userId, notificationPayload);
      
      this.stats.totalSent++;
      logger.info('Notification sent successfully', { notificationId: notificationPayload.id });
      
    } catch (error) {
      this.stats.errors++;
      logger.error('Failed to send notification:', error);
      throw error;
    }
  }

  async sendWebSocketNotification(userId, notification) {
    if (!this.io) {
      logger.warn('Socket.IO not initialized, skipping WebSocket notification');
      return;
    }

    try {
      if (userId === 'broadcast') {
        // Send to all connected clients
        this.io.emit('notification', notification);
        logger.info('Broadcast notification sent via WebSocket');
      } else {
        // Send to specific user room
        this.io.to(`user:${userId}`).emit('notification', notification);
        logger.info(`User notification sent via WebSocket to user: ${userId}`);
      }
      
      this.stats.webSocketSent++;
      
    } catch (error) {
      logger.error('WebSocket notification failed:', error);
      throw error;
    }
  }

  async sendPushNotification(userId, notification) {
    // TODO: Implement FCM for Android and APNS for iOS
    // This would require user device tokens stored in a database
    
    logger.info(`Push notification would be sent to user ${userId}:`, {
      title: notification.title,
      message: notification.message
    });
    
    this.stats.pushSent++;
    
    // Mock implementation for now
    return Promise.resolve();
  }

  generateNotificationId() {
    return `notif_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  getStats() {
    return {
      ...this.stats,
      uptime: Date.now() - this.stats.startTime.getTime(),
      connectedClients: this.io ? this.io.engine.clientsCount : 0
    };
  }

  // Test notification endpoint
  async sendTestNotification(userId = 'broadcast') {
    const testNotification = {
      userId,
      type: 'test',
      title: 'Test Notification',
      message: 'This is a test notification from the Notification Service',
      data: {
        test: true,
        timestamp: new Date().toISOString()
      }
    };

    await this.sendNotification(testNotification);
    return testNotification;
  }
}

module.exports = new NotificationService();