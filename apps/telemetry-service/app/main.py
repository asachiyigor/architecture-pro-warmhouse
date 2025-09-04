"""
Telemetry Service - MVP
Collects and processes telemetry data from smart home devices
"""

import asyncio
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

from app.core.config import settings
from app.core.database import init_db
from app.routers import telemetry, health, metrics
from app.services.background_tasks import start_background_tasks

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager"""
    logger.info("Starting Telemetry Service...")
    
    # Initialize database
    await init_db()
    logger.info("Database initialized")
    
    # RabbitMQ messaging removed - using Kafka only
    logger.info("Messaging skipped - using Kafka")
    
    # Start background tasks for data collection
    logger.info("Starting background tasks for data collection")
    asyncio.create_task(start_background_tasks())
    
    logger.info("Telemetry Service started successfully")
    
    yield
    
    # Cleanup
    logger.info("Shutting down Telemetry Service...")
    # RabbitMQ messaging removed
    logger.info("Telemetry Service stopped")


app = FastAPI(
    title="Telemetry Service",
    description="Smart Home Telemetry Collection and Processing Service",
    version="1.0.0",
    lifespan=lifespan
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(health.router, prefix="/health", tags=["Health"])
app.include_router(telemetry.router, prefix="/api/v1/telemetry", tags=["Telemetry"])
app.include_router(metrics.router, prefix="/api/v1/metrics", tags=["Metrics"])


@app.get("/")
async def root():
    """Root endpoint"""
    return {
        "service": "Telemetry Service",
        "version": "1.0.0",
        "status": "running"
    }


if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=settings.PORT,
        reload=settings.DEBUG,
        log_level="info"
    )