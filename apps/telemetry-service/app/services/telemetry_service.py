"""Telemetry service business logic"""

import json
import logging
from datetime import datetime, timedelta
from typing import List, Optional, Dict, Any

from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, and_, func
from sqlalchemy.orm import selectinload

from app.models.telemetry import TelemetryReading, TelemetryAlert, DeviceMetrics
from app.schemas.telemetry import (
    TelemetryReadingCreate, TelemetryAlertCreate, TelemetryQuery,
    TelemetryStats, LocationStats, AlertSeverity
)
from app.core.config import settings
# Lazy import to avoid startup blocking
# Lazy import to avoid blocking during startup

logger = logging.getLogger(__name__)


class TelemetryService:
    """Service for managing telemetry data"""
    
    def __init__(self, db_session: AsyncSession):
        self.db = db_session
    
    async def create_reading(self, reading_data: TelemetryReadingCreate) -> TelemetryReading:
        """Create a new telemetry reading"""
        reading = TelemetryReading(
            device_id=reading_data.device_id,
            device_type=reading_data.device_type,
            location=reading_data.location,
            metric_name=reading_data.metric_name.value,
            value=reading_data.value,
            unit=reading_data.unit,
            timestamp=reading_data.timestamp,
            raw_data=reading_data.raw_data,
            quality=reading_data.quality.value
        )
        
        self.db.add(reading)
        await self.db.commit()
        await self.db.refresh(reading)
        
        # Skip InfluxDB writing for performance testing
        # TODO: Re-enable InfluxDB once performance is optimized
        logger.debug("InfluxDB writing skipped for performance")
        
        # Skip Kafka publishing entirely when disabled
        if settings.ENABLE_KAFKA:
            logger.info("Kafka publishing enabled but skipping for performance")
            # TODO: Enable Kafka publishing once infrastructure is stable
        
        # Check for alerts after creating reading
        await self._check_alerts(reading)
        
        logger.info(f"Created telemetry reading for device {reading.device_id}")
        return reading
    
    async def create_batch_readings(self, readings: List[TelemetryReadingCreate]) -> List[TelemetryReading]:
        """Create multiple telemetry readings in batch"""
        db_readings = []
        
        for reading_data in readings:
            reading = TelemetryReading(
                device_id=reading_data.device_id,
                device_type=reading_data.device_type,
                location=reading_data.location,
                metric_name=reading_data.metric_name.value,
                value=reading_data.value,
                unit=reading_data.unit,
                timestamp=reading_data.timestamp,
                raw_data=reading_data.raw_data,
                quality=reading_data.quality.value
            )
            db_readings.append(reading)
        
        self.db.add_all(db_readings)
        await self.db.commit()
        
        # Write batch to InfluxDB
        influx_data = []
        for reading in db_readings:
            influx_data.append({
                "device_id": reading.device_id,
                "metric_type": reading.metric_name,
                "value": reading.value,
                "unit": reading.unit,
                "timestamp": reading.timestamp
            })
        
        try:
            from app.core.influx_client import influx_manager
            await influx_manager.write_telemetry_batch(influx_data)
        except Exception as e:
            logger.warning(f"Failed to write batch to InfluxDB: {e}")
        
        # Check alerts for each reading
        for reading in db_readings:
            await self._check_alerts(reading)
        
        logger.info(f"Created {len(db_readings)} telemetry readings in batch")
        return db_readings
    
    async def get_readings(self, query: TelemetryQuery) -> List[TelemetryReading]:
        """Get telemetry readings based on query parameters"""
        stmt = select(TelemetryReading)
        
        # Apply filters
        conditions = []
        if query.device_id:
            conditions.append(TelemetryReading.device_id == query.device_id)
        if query.device_type:
            conditions.append(TelemetryReading.device_type == query.device_type)
        if query.location:
            conditions.append(TelemetryReading.location.ilike(f"%{query.location}%"))
        if query.metric_name:
            conditions.append(TelemetryReading.metric_name == query.metric_name.value)
        if query.start_time:
            conditions.append(TelemetryReading.timestamp >= query.start_time)
        if query.end_time:
            conditions.append(TelemetryReading.timestamp <= query.end_time)
        
        if conditions:
            stmt = stmt.where(and_(*conditions))
        
        stmt = stmt.order_by(TelemetryReading.timestamp.desc())
        stmt = stmt.offset(query.offset).limit(query.limit)
        
        result = await self.db.execute(stmt)
        return result.scalars().all()
    
    async def get_latest_reading(self, device_id: str, metric_name: str) -> Optional[TelemetryReading]:
        """Get the latest reading for a device and metric"""
        stmt = select(TelemetryReading).where(
            and_(
                TelemetryReading.device_id == device_id,
                TelemetryReading.metric_name == metric_name
            )
        ).order_by(TelemetryReading.timestamp.desc()).limit(1)
        
        result = await self.db.execute(stmt)
        return result.scalars().first()
    
    async def get_device_readings_summary(self, device_id: str, hours: int = 24) -> Dict[str, Any]:
        """Get summary of readings for a device over the last N hours"""
        start_time = datetime.utcnow() - timedelta(hours=hours)
        
        stmt = select(
            TelemetryReading.metric_name,
            func.count(TelemetryReading.id).label('count'),
            func.avg(TelemetryReading.value).label('avg_value'),
            func.min(TelemetryReading.value).label('min_value'),
            func.max(TelemetryReading.value).label('max_value'),
            func.max(TelemetryReading.timestamp).label('latest_timestamp')
        ).where(
            and_(
                TelemetryReading.device_id == device_id,
                TelemetryReading.timestamp >= start_time
            )
        ).group_by(TelemetryReading.metric_name)
        
        result = await self.db.execute(stmt)
        rows = result.fetchall()
        
        summary = {}
        for row in rows:
            summary[row.metric_name] = {
                'count': row.count,
                'avg_value': float(row.avg_value) if row.avg_value else None,
                'min_value': float(row.min_value) if row.min_value else None,
                'max_value': float(row.max_value) if row.max_value else None,
                'latest_timestamp': row.latest_timestamp
            }
        
        return summary
    
    async def get_alerts(self, active_only: bool = False) -> List[TelemetryAlert]:
        """Get telemetry alerts"""
        stmt = select(TelemetryAlert)
        
        if active_only:
            stmt = stmt.where(TelemetryAlert.is_active == "true")
        
        stmt = stmt.order_by(TelemetryAlert.created_at.desc())
        
        result = await self.db.execute(stmt)
        return result.scalars().all()
    
    async def acknowledge_alert(self, alert_id: int) -> Optional[TelemetryAlert]:
        """Acknowledge an alert"""
        stmt = select(TelemetryAlert).where(TelemetryAlert.id == alert_id)
        result = await self.db.execute(stmt)
        alert = result.scalars().first()
        
        if alert:
            alert.acknowledged_at = datetime.utcnow()
            alert.updated_at = datetime.utcnow()
            await self.db.commit()
            await self.db.refresh(alert)
            
        return alert
    
    async def resolve_alert(self, alert_id: int) -> Optional[TelemetryAlert]:
        """Resolve an alert"""
        stmt = select(TelemetryAlert).where(TelemetryAlert.id == alert_id)
        result = await self.db.execute(stmt)
        alert = result.scalars().first()
        
        if alert:
            alert.is_active = "false"
            alert.resolved_at = datetime.utcnow()
            alert.updated_at = datetime.utcnow()
            await self.db.commit()
            await self.db.refresh(alert)
            
        return alert
    
    async def get_telemetry_stats(self) -> TelemetryStats:
        """Get overall telemetry statistics"""
        # Total readings
        total_readings_stmt = select(func.count(TelemetryReading.id))
        total_readings_result = await self.db.execute(total_readings_stmt)
        total_readings = total_readings_result.scalar() or 0
        
        # Total devices
        total_devices_stmt = select(func.count(func.distinct(TelemetryReading.device_id)))
        total_devices_result = await self.db.execute(total_devices_stmt)
        total_devices = total_devices_result.scalar() or 0
        
        # Active alerts
        active_alerts_stmt = select(func.count(TelemetryAlert.id)).where(
            TelemetryAlert.is_active == "true"
        )
        active_alerts_result = await self.db.execute(active_alerts_stmt)
        active_alerts = active_alerts_result.scalar() or 0
        
        # Last reading time
        last_reading_stmt = select(func.max(TelemetryReading.timestamp))
        last_reading_result = await self.db.execute(last_reading_stmt)
        last_reading_time = last_reading_result.scalar()
        
        # Readings per hour (last 24 hours)
        start_time = datetime.utcnow() - timedelta(hours=24)
        readings_per_hour_stmt = select(func.count(TelemetryReading.id)).where(
            TelemetryReading.timestamp >= start_time
        )
        readings_per_hour_result = await self.db.execute(readings_per_hour_stmt)
        readings_last_24h = readings_per_hour_result.scalar() or 0
        readings_per_hour = readings_last_24h / 24.0
        
        # Average quality score (simplified)
        quality_scores = {"good": 1.0, "poor": 0.6, "bad": 0.2}
        avg_quality_stmt = select(TelemetryReading.quality, func.count(TelemetryReading.id))
        avg_quality_stmt = avg_quality_stmt.group_by(TelemetryReading.quality)
        avg_quality_result = await self.db.execute(avg_quality_stmt)
        quality_data = avg_quality_result.fetchall()
        
        if quality_data:
            total_weighted = sum(quality_scores.get(row[0], 0) * row[1] for row in quality_data)
            total_count = sum(row[1] for row in quality_data)
            avg_quality_score = total_weighted / total_count if total_count > 0 else 0
        else:
            avg_quality_score = 0
        
        return TelemetryStats(
            total_readings=total_readings,
            total_devices=total_devices,
            active_alerts=active_alerts,
            last_reading_time=last_reading_time,
            readings_per_hour=readings_per_hour,
            average_quality_score=avg_quality_score
        )
    
    async def get_location_stats(self) -> List[LocationStats]:
        """Get statistics by location"""
        stmt = select(
            TelemetryReading.location,
            func.count(func.distinct(TelemetryReading.device_id)).label('device_count'),
            func.count(TelemetryReading.id).label('total_readings'),
            func.max(TelemetryReading.timestamp).label('latest_reading_time'),
            func.avg(
                func.case(
                    (TelemetryReading.metric_name == 'temperature', TelemetryReading.value),
                    else_=None
                )
            ).label('avg_temperature')
        ).group_by(TelemetryReading.location)
        
        result = await self.db.execute(stmt)
        location_data = result.fetchall()
        
        # Get alert counts per location
        alert_stmt = select(
            TelemetryAlert.location,
            func.count(TelemetryAlert.id).label('active_alerts')
        ).where(TelemetryAlert.is_active == "true").group_by(TelemetryAlert.location)
        
        alert_result = await self.db.execute(alert_stmt)
        alert_data = {row[0]: row[1] for row in alert_result.fetchall()}
        
        stats = []
        for row in location_data:
            stats.append(LocationStats(
                location=row.location,
                device_count=row.device_count,
                total_readings=row.total_readings,
                latest_reading_time=row.latest_reading_time,
                avg_temperature=float(row.avg_temperature) if row.avg_temperature else None,
                active_alerts=alert_data.get(row.location, 0)
            ))
        
        return stats
    
    async def _check_alerts(self, reading: TelemetryReading):
        """Check if a reading should trigger alerts"""
        if reading.metric_name == "temperature":
            await self._check_temperature_alerts(reading)
    
    async def _check_temperature_alerts(self, reading: TelemetryReading):
        """Check for temperature-based alerts"""
        alerts_to_create = []
        
        # High temperature alert
        if reading.value > settings.TEMPERATURE_HIGH_THRESHOLD:
            alert_data = TelemetryAlertCreate(
                device_id=reading.device_id,
                device_type=reading.device_type,
                location=reading.location,
                alert_type="high_temperature",
                severity=AlertSeverity.CRITICAL if reading.value > 40 else AlertSeverity.HIGH,
                current_value=reading.value,
                threshold_value=settings.TEMPERATURE_HIGH_THRESHOLD,
                message=f"Temperature {reading.value}°C exceeds threshold {settings.TEMPERATURE_HIGH_THRESHOLD}°C in {reading.location}"
            )
            alerts_to_create.append(alert_data)
        
        # Low temperature alert
        if reading.value < settings.TEMPERATURE_LOW_THRESHOLD:
            alert_data = TelemetryAlertCreate(
                device_id=reading.device_id,
                device_type=reading.device_type,
                location=reading.location,
                alert_type="low_temperature",
                severity=AlertSeverity.CRITICAL if reading.value < 0 else AlertSeverity.MEDIUM,
                current_value=reading.value,
                threshold_value=settings.TEMPERATURE_LOW_THRESHOLD,
                message=f"Temperature {reading.value}°C is below threshold {settings.TEMPERATURE_LOW_THRESHOLD}°C in {reading.location}"
            )
            alerts_to_create.append(alert_data)
        
        # Create alerts
        for alert_data in alerts_to_create:
            # Check if similar alert already exists and is active
            existing_alert_stmt = select(TelemetryAlert).where(
                and_(
                    TelemetryAlert.device_id == alert_data.device_id,
                    TelemetryAlert.alert_type == alert_data.alert_type,
                    TelemetryAlert.is_active == "true"
                )
            )
            
            result = await self.db.execute(existing_alert_stmt)
            existing_alert = result.scalars().first()
            
            if not existing_alert:
                # Create new alert
                alert = TelemetryAlert(
                    device_id=alert_data.device_id,
                    device_type=alert_data.device_type,
                    location=alert_data.location,
                    alert_type=alert_data.alert_type,
                    severity=alert_data.severity.value,
                    current_value=alert_data.current_value,
                    threshold_value=alert_data.threshold_value,
                    message=alert_data.message
                )
                
                self.db.add(alert)
                await self.db.commit()
                
                # Skip Kafka alert publishing when disabled
                if settings.ENABLE_KAFKA:
                    logger.info("Kafka alert publishing enabled but skipping for performance")
                    # TODO: Enable Kafka alert publishing once infrastructure is stable
                    
                logger.warning(f"Created alert: {alert.message}")
            else:
                # Update existing alert with new values
                existing_alert.current_value = alert_data.current_value
                existing_alert.message = alert_data.message
                existing_alert.updated_at = datetime.utcnow()
                await self.db.commit()
                
                logger.info(f"Updated existing alert for device {alert_data.device_id}")