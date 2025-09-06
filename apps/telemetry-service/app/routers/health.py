"""Health check endpoints"""

from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import text

from app.core.database import get_db_session
from app.core.config import settings

router = APIRouter()


@router.get("/")
async def health_check():
    """Basic health check"""
    return {
        "status": "healthy",
        "service": "telemetry-service",
        "version": settings.VERSION
    }


@router.get("/detailed")
async def detailed_health_check(db: AsyncSession = Depends(get_db_session)):
    """Detailed health check with dependencies"""
    health_status = {
        "status": "healthy",
        "service": "telemetry-service",
        "version": settings.VERSION,
        "database": "unknown",
        "messaging": "unknown"
    }
    
    # Check database connectivity
    try:
        result = await db.execute(text("SELECT 1"))
        health_status["database"] = "healthy"
    except Exception as e:
        health_status["database"] = f"unhealthy: {str(e)}"
        health_status["status"] = "unhealthy"
    
    # Check messaging (simplified - in real implementation you'd ping RabbitMQ)
    try:
        # This is a simplified check - in production you'd check RabbitMQ connectivity
        health_status["messaging"] = "healthy"
    except Exception as e:
        health_status["messaging"] = f"unhealthy: {str(e)}"
        health_status["status"] = "degraded"
    
    return health_status