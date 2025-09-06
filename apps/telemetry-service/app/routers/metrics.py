"""Prometheus metrics endpoints"""

from fastapi import APIRouter
from fastapi.responses import PlainTextResponse
import time

router = APIRouter()

# Simplified metrics for MVP - in production you'd use prometheus_client
start_time = time.time()
request_count = 0
error_count = 0


@router.get("/prometheus", response_class=PlainTextResponse)
async def get_prometheus_metrics():
    """Prometheus metrics endpoint"""
    global request_count, error_count
    
    uptime = time.time() - start_time
    
    metrics = f"""# HELP telemetry_service_uptime_seconds Service uptime in seconds
# TYPE telemetry_service_uptime_seconds counter
telemetry_service_uptime_seconds {uptime}

# HELP telemetry_service_requests_total Total number of requests
# TYPE telemetry_service_requests_total counter
telemetry_service_requests_total {request_count}

# HELP telemetry_service_errors_total Total number of errors
# TYPE telemetry_service_errors_total counter
telemetry_service_errors_total {error_count}
"""
    
    return metrics


@router.get("/")
async def get_basic_metrics():
    """Basic metrics endpoint"""
    global request_count
    request_count += 1
    
    uptime = time.time() - start_time
    
    return {
        "uptime_seconds": uptime,
        "total_requests": request_count,
        "total_errors": error_count,
        "service": "telemetry-service"
    }