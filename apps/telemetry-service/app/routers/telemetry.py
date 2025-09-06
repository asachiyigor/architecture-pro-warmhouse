"""Telemetry API endpoints"""

from typing import List
from fastapi import APIRouter, Depends, HTTPException, BackgroundTasks
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.database import get_db_session
from app.schemas.telemetry import (
    TelemetryReading, TelemetryReadingCreate, TelemetryBatch,
    TelemetryQuery, TelemetryAlert, TelemetryStats, LocationStats
)
from app.services.telemetry_service import TelemetryService

router = APIRouter()


@router.post("/readings", response_model=TelemetryReading)
async def create_reading(
    reading: TelemetryReadingCreate,
    background_tasks: BackgroundTasks,
    db: AsyncSession = Depends(get_db_session)
):
    """Create a single telemetry reading"""
    service = TelemetryService(db)
    
    try:
        result = await service.create_reading(reading)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/readings/batch", response_model=List[TelemetryReading])
async def create_batch_readings(
    batch: TelemetryBatch,
    background_tasks: BackgroundTasks,
    db: AsyncSession = Depends(get_db_session)
):
    """Create multiple telemetry readings in batch"""
    service = TelemetryService(db)
    
    try:
        results = await service.create_batch_readings(batch.readings)
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/readings", response_model=List[TelemetryReading])
async def get_readings(
    device_id: str = None,
    device_type: str = None,
    location: str = None,
    metric_name: str = None,
    start_time: str = None,
    end_time: str = None,
    limit: int = 100,
    offset: int = 0,
    db: AsyncSession = Depends(get_db_session)
):
    """Get telemetry readings with optional filters"""
    from datetime import datetime
    from app.schemas.telemetry import MetricType
    
    # Parse datetime strings
    start_dt = None
    end_dt = None
    
    try:
        if start_time:
            start_dt = datetime.fromisoformat(start_time.replace('Z', '+00:00'))
        if end_time:
            end_dt = datetime.fromisoformat(end_time.replace('Z', '+00:00'))
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid datetime format. Use ISO format.")
    
    # Parse metric name
    metric_enum = None
    if metric_name:
        try:
            metric_enum = MetricType(metric_name)
        except ValueError:
            raise HTTPException(status_code=400, detail=f"Invalid metric name: {metric_name}")
    
    query = TelemetryQuery(
        device_id=device_id,
        device_type=device_type,
        location=location,
        metric_name=metric_enum,
        start_time=start_dt,
        end_time=end_dt,
        limit=min(limit, 1000),  # Cap at 1000
        offset=offset
    )
    
    service = TelemetryService(db)
    
    try:
        results = await service.get_readings(query)
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/readings/latest/{device_id}/{metric_name}", response_model=TelemetryReading)
async def get_latest_reading(
    device_id: str,
    metric_name: str,
    db: AsyncSession = Depends(get_db_session)
):
    """Get the latest reading for a device and metric"""
    service = TelemetryService(db)
    
    result = await service.get_latest_reading(device_id, metric_name)
    if not result:
        raise HTTPException(status_code=404, detail="No readings found for this device and metric")
    
    return result


@router.get("/devices/{device_id}/summary")
async def get_device_summary(
    device_id: str,
    hours: int = 24,
    db: AsyncSession = Depends(get_db_session)
):
    """Get summary of readings for a device"""
    if hours < 1 or hours > 168:  # Max 1 week
        raise HTTPException(status_code=400, detail="Hours must be between 1 and 168")
    
    service = TelemetryService(db)
    
    try:
        result = await service.get_device_readings_summary(device_id, hours)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/alerts", response_model=List[TelemetryAlert])
async def get_alerts(
    active_only: bool = False,
    db: AsyncSession = Depends(get_db_session)
):
    """Get telemetry alerts"""
    service = TelemetryService(db)
    
    try:
        results = await service.get_alerts(active_only)
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.patch("/alerts/{alert_id}/acknowledge", response_model=TelemetryAlert)
async def acknowledge_alert(
    alert_id: int,
    db: AsyncSession = Depends(get_db_session)
):
    """Acknowledge an alert"""
    service = TelemetryService(db)
    
    result = await service.acknowledge_alert(alert_id)
    if not result:
        raise HTTPException(status_code=404, detail="Alert not found")
    
    return result


@router.patch("/alerts/{alert_id}/resolve", response_model=TelemetryAlert)
async def resolve_alert(
    alert_id: int,
    db: AsyncSession = Depends(get_db_session)
):
    """Resolve an alert"""
    service = TelemetryService(db)
    
    result = await service.resolve_alert(alert_id)
    if not result:
        raise HTTPException(status_code=404, detail="Alert not found")
    
    return result


@router.get("/stats", response_model=TelemetryStats)
async def get_telemetry_stats(db: AsyncSession = Depends(get_db_session)):
    """Get overall telemetry statistics"""
    service = TelemetryService(db)
    
    try:
        result = await service.get_telemetry_stats()
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/stats/locations", response_model=List[LocationStats])
async def get_location_stats(db: AsyncSession = Depends(get_db_session)):
    """Get statistics by location"""
    service = TelemetryService(db)
    
    try:
        results = await service.get_location_stats()
        return results
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))