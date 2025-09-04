"""Pydantic schemas for telemetry service"""

from datetime import datetime
from typing import Optional, List
from pydantic import BaseModel, Field, validator
from enum import Enum


class MetricType(str, Enum):
    """Available metric types"""
    TEMPERATURE = "temperature"
    HUMIDITY = "humidity"
    PRESSURE = "pressure"
    MOTION = "motion"
    DOOR_STATE = "door_state"
    LIGHT_LEVEL = "light_level"
    POWER_CONSUMPTION = "power_consumption"
    BATTERY_LEVEL = "battery_level"


class QualityLevel(str, Enum):
    """Data quality levels"""
    GOOD = "good"
    POOR = "poor"
    BAD = "bad"


class AlertSeverity(str, Enum):
    """Alert severity levels"""
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"


class TelemetryReadingCreate(BaseModel):
    """Schema for creating telemetry readings"""
    device_id: str = Field(..., min_length=1, max_length=50)
    device_type: str = Field(..., min_length=1, max_length=50)
    location: str = Field(..., min_length=1, max_length=100)
    metric_name: MetricType
    value: float
    unit: str = Field(..., min_length=1, max_length=20)
    timestamp: datetime
    raw_data: Optional[str] = None
    quality: QualityLevel = QualityLevel.GOOD
    
    @validator('value')
    def validate_value(cls, v):
        if not isinstance(v, (int, float)):
            raise ValueError('Value must be a number')
        return float(v)


class TelemetryReading(BaseModel):
    """Schema for telemetry reading response"""
    id: int
    device_id: str
    device_type: str
    location: str
    metric_name: str
    value: float
    unit: str
    timestamp: datetime
    raw_data: Optional[str] = None
    quality: str
    created_at: datetime
    
    class Config:
        from_attributes = True


class TelemetryBatch(BaseModel):
    """Schema for batch telemetry readings"""
    readings: List[TelemetryReadingCreate] = Field(..., min_items=1, max_items=100)


class TelemetryQuery(BaseModel):
    """Schema for querying telemetry data"""
    device_id: Optional[str] = None
    device_type: Optional[str] = None
    location: Optional[str] = None
    metric_name: Optional[MetricType] = None
    start_time: Optional[datetime] = None
    end_time: Optional[datetime] = None
    limit: int = Field(default=100, ge=1, le=1000)
    offset: int = Field(default=0, ge=0)


class TelemetryAlertCreate(BaseModel):
    """Schema for creating telemetry alerts"""
    device_id: str = Field(..., min_length=1, max_length=50)
    device_type: str = Field(..., min_length=1, max_length=50)
    location: str = Field(..., min_length=1, max_length=100)
    alert_type: str = Field(..., min_length=1, max_length=50)
    severity: AlertSeverity
    current_value: float
    threshold_value: float
    message: str = Field(..., min_length=1)


class TelemetryAlert(BaseModel):
    """Schema for telemetry alert response"""
    id: int
    device_id: str
    device_type: str
    location: str
    alert_type: str
    severity: str
    current_value: float
    threshold_value: float
    message: str
    is_active: str
    acknowledged_at: Optional[datetime] = None
    resolved_at: Optional[datetime] = None
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True


class DeviceMetrics(BaseModel):
    """Schema for device metrics response"""
    id: int
    device_id: str
    device_type: str
    location: str
    metric_name: str
    avg_value: Optional[float] = None
    min_value: Optional[float] = None
    max_value: Optional[float] = None
    sample_count: int
    period_start: datetime
    period_end: datetime
    aggregation_type: str
    created_at: datetime
    
    class Config:
        from_attributes = True


class TelemetryStats(BaseModel):
    """Schema for telemetry statistics"""
    total_readings: int
    total_devices: int
    active_alerts: int
    last_reading_time: Optional[datetime] = None
    readings_per_hour: float
    average_quality_score: float


class LocationStats(BaseModel):
    """Schema for location-based statistics"""
    location: str
    device_count: int
    total_readings: int
    latest_reading_time: Optional[datetime] = None
    avg_temperature: Optional[float] = None
    active_alerts: int