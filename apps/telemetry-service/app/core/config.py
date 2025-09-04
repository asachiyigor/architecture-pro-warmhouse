"""Configuration settings for Telemetry Service"""

from typing import List
from pydantic import Field
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """Application settings"""
    
    # Application
    APP_NAME: str = "Telemetry Service"
    VERSION: str = "1.0.0"
    DEBUG: bool = False
    PORT: int = 8083
    
    # Database
    DATABASE_URL: str = Field(
        default="postgresql+asyncpg://postgres:postgres@localhost:5432/smarthome_telemetry",
        env="DATABASE_URL"
    )
    
   # Redis
    REDIS_URL: str = Field(
        default="redis://redis:6379/0",
        env="REDIS_URL"
    )
    
    # InfluxDB (Time Series Database)
    INFLUXDB_URL: str = Field(
        default="http://smarthome-influxdb:8086",
        env="INFLUXDB_URL"
    )
    INFLUXDB_TOKEN: str = Field(
        default="smarthome-super-secret-auth-token",
        env="INFLUXDB_TOKEN"
    )
    INFLUXDB_ORG: str = Field(
        default="smarthome",
        env="INFLUXDB_ORG"
    )
    INFLUXDB_BUCKET: str = Field(
        default="telemetry",
        env="INFLUXDB_BUCKET"
    )
    
    # Kafka (Message Streaming)
    KAFKA_BROKERS: str = Field(
        default="smarthome-kafka:9092",
        env="KAFKA_BROKERS"
    )
    KAFKA_TOPIC_TELEMETRY: str = "telemetry-events"
    KAFKA_TOPIC_ALERTS: str = "telemetry-alerts"
    ENABLE_KAFKA: bool = Field(
        default=True,
        env="ENABLE_KAFKA"
    )
    
    # External Services
    TEMPERATURE_API_URL: str = Field(
        default="http://temperature-api:8081",
        env="TEMPERATURE_API_URL"
    )
    
    DEVICE_SERVICE_URL: str = Field(
        default="http://device-service:8082",
        env="DEVICE_SERVICE_URL"
    )
    
    # CORS
    CORS_ORIGINS: List[str] = ["*"]
    
    # Telemetry Collection
    COLLECTION_INTERVAL: int = 30  # seconds
    BATCH_SIZE: int = 100
    RETENTION_DAYS: int = 365
    
    # Alert Thresholds
    TEMPERATURE_HIGH_THRESHOLD: float = 35.0
    TEMPERATURE_LOW_THRESHOLD: float = 5.0
    
    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()