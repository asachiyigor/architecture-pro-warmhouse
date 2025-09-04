"""Telemetry data models"""

from datetime import datetime
from typing import Optional
from sqlalchemy import Column, Integer, String, Float, DateTime, Text, Index
from sqlalchemy.dialects.postgresql import TIMESTAMP
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class TelemetryReading(Base):
    """Telemetry reading model"""
    __tablename__ = "telemetry_readings"
    
    id = Column(Integer, primary_key=True, index=True)
    device_id = Column(String(50), nullable=False, index=True)
    device_type = Column(String(50), nullable=False)
    location = Column(String(100), nullable=False, index=True)
    metric_name = Column(String(50), nullable=False)  # temperature, humidity, etc.
    value = Column(Float, nullable=False)
    unit = Column(String(20), nullable=False)
    timestamp = Column(TIMESTAMP(timezone=True), nullable=False, index=True)
    
    # Additional metadata
    raw_data = Column(Text, nullable=True)  # JSON string for additional data
    quality = Column(String(20), default="good")  # good, poor, bad
    
    created_at = Column(TIMESTAMP(timezone=True), default=datetime.utcnow)
    
    # Indexes for performance
    __table_args__ = (
        Index('ix_device_timestamp', 'device_id', 'timestamp'),
        Index('ix_location_timestamp', 'location', 'timestamp'),
        Index('ix_metric_timestamp', 'metric_name', 'timestamp'),
    )
    
    def __repr__(self):
        return f"<TelemetryReading(device_id={self.device_id}, metric={self.metric_name}, value={self.value})>"


class TelemetryAlert(Base):
    """Telemetry alert model"""
    __tablename__ = "telemetry_alerts"
    
    id = Column(Integer, primary_key=True, index=True)
    device_id = Column(String(50), nullable=False, index=True)
    device_type = Column(String(50), nullable=False)
    location = Column(String(100), nullable=False)
    alert_type = Column(String(50), nullable=False)  # high_temperature, low_temperature, etc.
    severity = Column(String(20), nullable=False)  # low, medium, high, critical
    
    current_value = Column(Float, nullable=False)
    threshold_value = Column(Float, nullable=False)
    message = Column(Text, nullable=False)
    
    is_active = Column(String(10), default="true")  # true, false
    acknowledged_at = Column(TIMESTAMP(timezone=True), nullable=True)
    resolved_at = Column(TIMESTAMP(timezone=True), nullable=True)
    
    created_at = Column(TIMESTAMP(timezone=True), default=datetime.utcnow, index=True)
    updated_at = Column(TIMESTAMP(timezone=True), default=datetime.utcnow, onupdate=datetime.utcnow)
    
    def __repr__(self):
        return f"<TelemetryAlert(device_id={self.device_id}, type={self.alert_type}, severity={self.severity})>"


class DeviceMetrics(Base):
    """Aggregated device metrics"""
    __tablename__ = "device_metrics"
    
    id = Column(Integer, primary_key=True, index=True)
    device_id = Column(String(50), nullable=False, index=True)
    device_type = Column(String(50), nullable=False)
    location = Column(String(100), nullable=False)
    
    # Aggregated values (e.g., hourly averages)
    metric_name = Column(String(50), nullable=False)
    avg_value = Column(Float)
    min_value = Column(Float)
    max_value = Column(Float)
    sample_count = Column(Integer, default=0)
    
    # Time period
    period_start = Column(TIMESTAMP(timezone=True), nullable=False, index=True)
    period_end = Column(TIMESTAMP(timezone=True), nullable=False)
    aggregation_type = Column(String(20), default="hourly")  # hourly, daily, weekly
    
    created_at = Column(TIMESTAMP(timezone=True), default=datetime.utcnow)
    
    __table_args__ = (
        Index('ix_device_period', 'device_id', 'period_start'),
        Index('ix_location_period', 'location', 'period_start'),
    )
    
    def __repr__(self):
        return f"<DeviceMetrics(device_id={self.device_id}, metric={self.metric_name}, avg={self.avg_value})>"