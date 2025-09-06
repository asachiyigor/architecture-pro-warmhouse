"""Kafka client for message streaming"""

import json
import logging
from typing import Dict, Any, Optional

try:
    from kafka import KafkaProducer, KafkaConsumer
    from kafka.errors import KafkaError
    KAFKA_AVAILABLE = True
except ImportError:
    KAFKA_AVAILABLE = False

from .config import settings

logger = logging.getLogger(__name__)


class KafkaManager:
    """Kafka client manager for event streaming"""
    
    def __init__(self):
        self.producer: Optional[KafkaProducer] = None
        self._initialized = False
        
        # Don't initialize Kafka on startup - do it lazily when first needed
        if not KAFKA_AVAILABLE:
            logger.warning("Kafka client library not available")
        elif not settings.ENABLE_KAFKA:
            logger.info("Kafka is disabled by configuration")
    
    def _initialize_producer(self):
        """Initialize Kafka producer lazily"""
        if self._initialized or not KAFKA_AVAILABLE or not settings.ENABLE_KAFKA:
            return
            
        try:
            # Very aggressive timeouts for instant failure if Kafka not available
            self.producer = KafkaProducer(
                bootstrap_servers=settings.KAFKA_BROKERS.split(','),
                value_serializer=lambda x: json.dumps(x).encode('utf-8'),
                key_serializer=lambda x: x.encode('utf-8') if x else None,
                acks='all',
                retries=0,  # No retries - fail immediately
                retry_backoff_ms=50,  # Minimal backoff
                request_timeout_ms=1000,  # 1 second timeout
                metadata_max_age_ms=1000,  # Refresh metadata every 1 second  
                api_version_auto_timeout_ms=1000,  # API version detection timeout
                # Connection timeouts
                connections_max_idle_ms=1000,
                socket_keepalive=False,
                reconnect_backoff_ms=50,
                reconnect_backoff_max_ms=100
            )
            logger.info("Kafka producer initialized successfully")
        except Exception as e:
            logger.warning(f"Failed to initialize Kafka producer: {e}")
            self.producer = None
        finally:
            self._initialized = True

    def is_available(self) -> bool:
        """Check if Kafka is available"""
        if not self._initialized:
            self._initialize_producer()
        return self.producer is not None
    
    async def publish_telemetry_event(
        self, 
        device_id: str,
        metric_type: str,
        value: float,
        unit: str,
        timestamp: str,
        additional_data: Optional[Dict[str, Any]] = None
    ) -> bool:
        """Publish telemetry event to Kafka"""
        if not self.is_available():
            return False
        
        try:
            event_data = {
                "event_type": "telemetry_received",
                "device_id": device_id,
                "metric_type": metric_type,
                "value": value,
                "unit": unit,
                "timestamp": timestamp
            }
            
            if additional_data:
                event_data.update(additional_data)
            
            # Send to Kafka topic
            future = self.producer.send(
                topic=settings.KAFKA_TOPIC_TELEMETRY,
                key=device_id,
                value=event_data
            )
            
            # Wait for send to complete with very short timeout
            record_metadata = future.get(timeout=1)
            logger.debug(f"Telemetry event sent to Kafka: topic={record_metadata.topic}, partition={record_metadata.partition}")
            return True
            
        except KafkaError as e:
            logger.error(f"Failed to publish telemetry event to Kafka: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error publishing telemetry event: {e}")
            return False
    
    async def publish_alert_event(
        self,
        device_id: str,
        alert_type: str,
        severity: str,
        message: str,
        current_value: float,
        threshold_value: float
    ) -> bool:
        """Publish alert event to Kafka"""
        if not self.is_available():
            return False
        
        try:
            alert_data = {
                "event_type": "telemetry_alert",
                "device_id": device_id,
                "alert_type": alert_type,
                "severity": severity,
                "message": message,
                "current_value": current_value,
                "threshold_value": threshold_value,
                "timestamp": None  # Will be set by receiver
            }
            
            future = self.producer.send(
                topic=settings.KAFKA_TOPIC_ALERTS,
                key=device_id,
                value=alert_data
            )
            
            record_metadata = future.get(timeout=1)
            logger.debug(f"Alert event sent to Kafka: topic={record_metadata.topic}, partition={record_metadata.partition}")
            return True
            
        except KafkaError as e:
            logger.error(f"Failed to publish alert event to Kafka: {e}")
            return False
        except Exception as e:
            logger.error(f"Unexpected error publishing alert event: {e}")
            return False
    
    async def publish_device_status_event(
        self,
        device_id: str,
        old_status: str,
        new_status: str
    ) -> bool:
        """Publish device status change event to Kafka"""
        if not self.is_available():
            return False
        
        try:
            status_data = {
                "event_type": "device_status_changed",
                "device_id": device_id,
                "old_status": old_status,
                "new_status": new_status,
                "timestamp": None  # Will be set by receiver
            }
            
            future = self.producer.send(
                topic=settings.KAFKA_TOPIC_TELEMETRY,
                key=device_id,
                value=status_data
            )
            
            record_metadata = future.get(timeout=1)
            logger.debug(f"Device status event sent to Kafka: topic={record_metadata.topic}")
            return True
            
        except Exception as e:
            logger.error(f"Failed to publish device status event: {e}")
            return False
    
    def close(self):
        """Close Kafka producer"""
        if self.producer:
            self.producer.close()


# Global Kafka manager instance (lazy initialization)
kafka_manager = None

def get_kafka_manager() -> KafkaManager:
    """Get or create Kafka manager instance"""
    global kafka_manager
    if kafka_manager is None:
        kafka_manager = KafkaManager()
    return kafka_manager