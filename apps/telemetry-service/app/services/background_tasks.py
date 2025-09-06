"""Background tasks for telemetry service"""

import asyncio
import logging
import httpx
from datetime import datetime, timedelta

from app.core.config import settings
from app.core.database import async_session_factory
from app.services.telemetry_service import TelemetryService
from app.schemas.telemetry import TelemetryReadingCreate, MetricType

logger = logging.getLogger(__name__)


async def start_background_tasks():
    """Start all background tasks"""
    logger.info("Starting background tasks...")
    
    # Start temperature collection task
    temperature_task = asyncio.create_task(collect_temperature_data())
    
    # Start cleanup task
    cleanup_task = asyncio.create_task(cleanup_old_data())
    
    # Wait for tasks (they run indefinitely)
    await asyncio.gather(temperature_task, cleanup_task, return_exceptions=True)


async def collect_temperature_data():
    """Periodically collect temperature data from temperature API"""
    logger.info("Starting temperature data collection task")
    
    locations = ["Living Room", "Bedroom", "Kitchen"]
    
    while True:
        try:
            async with httpx.AsyncClient() as client:
                for location in locations:
                    try:
                        # Get temperature from temperature API
                        response = await client.get(
                            f"{settings.TEMPERATURE_API_URL}/temperature",
                            params={"location": location},
                            timeout=10.0
                        )
                        
                        if response.status_code == 200:
                            temp_data = response.json()
                            
                            # Create telemetry reading
                            reading = TelemetryReadingCreate(
                                device_id=temp_data.get("sensor_id", "unknown"),
                                device_type="temperature_sensor",
                                location=temp_data.get("location", location),
                                metric_name=MetricType.TEMPERATURE,
                                value=temp_data.get("value", 0.0),
                                unit=temp_data.get("unit", "°C"),
                                timestamp=datetime.utcnow()
                            )
                            
                            # Store in database
                            async with async_session_factory() as session:
                                service = TelemetryService(session)
                                await service.create_reading(reading)
                            
                            logger.debug(f"Collected temperature data for {location}: {temp_data.get('value')}°C")
                        
                        else:
                            logger.warning(f"Failed to get temperature for {location}: HTTP {response.status_code}")
                            
                    except Exception as e:
                        logger.error(f"Error collecting temperature data for {location}: {e}")
                        
                    # Small delay between locations
                    await asyncio.sleep(1)
            
        except Exception as e:
            logger.error(f"Error in temperature collection task: {e}")
        
        # Wait before next collection
        await asyncio.sleep(settings.COLLECTION_INTERVAL)


async def cleanup_old_data():
    """Periodically cleanup old telemetry data"""
    logger.info("Starting data cleanup task")
    
    while True:
        try:
            cutoff_date = datetime.utcnow() - timedelta(days=settings.RETENTION_DAYS)
            
            async with async_session_factory() as session:
                # In a real implementation, you'd delete old data here
                # For MVP, we'll just log the cleanup
                logger.info(f"Would cleanup data older than {cutoff_date}")
                
        except Exception as e:
            logger.error(f"Error in cleanup task: {e}")
        
        # Run cleanup once per day
        await asyncio.sleep(24 * 60 * 60)


async def aggregate_metrics():
    """Aggregate telemetry data into summary metrics"""
    logger.info("Starting metrics aggregation task")
    
    while True:
        try:
            async with async_session_factory() as session:
                # In a real implementation, you'd create hourly/daily aggregates
                # For MVP, we'll just log the aggregation
                logger.info("Would aggregate metrics here")
                
        except Exception as e:
            logger.error(f"Error in aggregation task: {e}")
        
        # Run aggregation every hour
        await asyncio.sleep(60 * 60)