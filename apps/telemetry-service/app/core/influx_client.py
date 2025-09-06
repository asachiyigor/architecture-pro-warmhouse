"""InfluxDB client for time series data"""

import logging
from datetime import datetime
from typing import Dict, List, Optional

try:
    from influxdb_client import InfluxDBClient, Point, WritePrecision
    from influxdb_client.client.write_api import SYNCHRONOUS
    INFLUXDB_AVAILABLE = True
except ImportError:
    INFLUXDB_AVAILABLE = False

from .config import settings

logger = logging.getLogger(__name__)


class InfluxDBManager:
    """InfluxDB client manager for telemetry data"""
    
    def __init__(self):
        self.client: Optional[InfluxDBClient] = None
        self.write_api = None
        self.query_api = None
        
        if INFLUXDB_AVAILABLE:
            try:
                self.client = InfluxDBClient(
                    url=settings.INFLUXDB_URL,
                    token=settings.INFLUXDB_TOKEN,
                    org=settings.INFLUXDB_ORG
                )
                self.write_api = self.client.write_api(write_options=SYNCHRONOUS)
                self.query_api = self.client.query_api()
                logger.info("InfluxDB client initialized successfully")
            except Exception as e:
                logger.warning(f"Failed to initialize InfluxDB client: {e}")
                self.client = None
        else:
            logger.warning("InfluxDB client library not available, falling back to PostgreSQL")
    
    def is_available(self) -> bool:
        """Check if InfluxDB is available"""
        return self.client is not None
    
    async def write_telemetry_point(
        self, 
        device_id: str, 
        metric_type: str, 
        value: float, 
        unit: str,
        timestamp: Optional[datetime] = None,
        additional_tags: Optional[Dict[str, str]] = None
    ) -> bool:
        """Write a telemetry data point to InfluxDB"""
        if not self.is_available():
            return False
        
        try:
            point = Point("telemetry") \
                .tag("device_id", device_id) \
                .tag("metric_type", metric_type) \
                .tag("unit", unit) \
                .field("value", float(value))
            
            if additional_tags:
                for key, value in additional_tags.items():
                    point = point.tag(key, value)
            
            if timestamp:
                point = point.time(timestamp, WritePrecision.S)
            
            self.write_api.write(bucket=settings.INFLUXDB_BUCKET, record=point)
            return True
        except Exception as e:
            logger.error(f"Failed to write telemetry point to InfluxDB: {e}")
            return False
    
    async def write_telemetry_batch(self, telemetry_data: List[Dict]) -> bool:
        """Write multiple telemetry points to InfluxDB"""
        if not self.is_available():
            return False
        
        try:
            points = []
            for data in telemetry_data:
                point = Point("telemetry") \
                    .tag("device_id", data["device_id"]) \
                    .tag("metric_type", data["metric_type"]) \
                    .tag("unit", data["unit"]) \
                    .field("value", float(data["value"]))
                
                if "timestamp" in data and data["timestamp"]:
                    point = point.time(data["timestamp"], WritePrecision.S)
                
                points.append(point)
            
            self.write_api.write(bucket=settings.INFLUXDB_BUCKET, record=points)
            return True
        except Exception as e:
            logger.error(f"Failed to write telemetry batch to InfluxDB: {e}")
            return False
    
    async def query_telemetry_data(
        self, 
        device_id: str, 
        metric_type: str,
        start_time: str = "-1h",
        stop_time: str = "now()"
    ) -> List[Dict]:
        """Query telemetry data from InfluxDB"""
        if not self.is_available():
            return []
        
        try:
            query = f'''
                from(bucket: "{settings.INFLUXDB_BUCKET}")
                |> range(start: {start_time}, stop: {stop_time})
                |> filter(fn: (r) => r._measurement == "telemetry")
                |> filter(fn: (r) => r.device_id == "{device_id}")
                |> filter(fn: (r) => r.metric_type == "{metric_type}")
                |> sort(columns: ["_time"], desc: true)
            '''
            
            tables = self.query_api.query(query)
            results = []
            
            for table in tables:
                for record in table.records:
                    results.append({
                        "device_id": record.values.get("device_id"),
                        "metric_type": record.values.get("metric_type"),
                        "value": record.values.get("_value"),
                        "unit": record.values.get("unit"),
                        "timestamp": record.values.get("_time")
                    })
            
            return results
        except Exception as e:
            logger.error(f"Failed to query telemetry data from InfluxDB: {e}")
            return []
    
    def close(self):
        """Close InfluxDB connection"""
        if self.client:
            self.client.close()


# Global InfluxDB manager instance
influx_manager = InfluxDBManager()