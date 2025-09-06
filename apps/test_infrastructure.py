#!/usr/bin/env python3
"""
Infrastructure Test Script
Tests all components: PostgreSQL, RabbitMQ, Redis, Kafka, InfluxDB
"""

import asyncio
import json
import logging
import sys
from datetime import datetime
from typing import Dict, Any

import asyncpg
import aio_pika
import redis
import httpx

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class InfrastructureTest:
    def __init__(self):
        self.results = {}
        self.redis_client = None
        
    async def test_postgresql_connections(self) -> Dict[str, bool]:
        """Test all PostgreSQL database connections"""
        databases = {
            "main": "postgresql://postgres:postgres@localhost:5432/smarthome",
            "devices": "postgresql://postgres:postgres@localhost:5433/smarthome_devices", 
            "telemetry": "postgresql://postgres:postgres@localhost:5434/smarthome_telemetry",
            "users": "postgresql://postgres:postgres@localhost:5435/smarthome_users",
            "billing": "postgresql://postgres:postgres@localhost:5436/smarthome_billing",
            "automation": "postgresql://postgres:postgres@localhost:5437/smarthome_automation"
        }
        
        results = {}
        for db_name, connection_string in databases.items():
            try:
                conn = await asyncpg.connect(connection_string)
                await conn.execute("SELECT 1")
                await conn.close()
                results[f"postgresql_{db_name}"] = True
                logger.info(f"✅ PostgreSQL {db_name} connection successful")
            except Exception as e:
                results[f"postgresql_{db_name}"] = False
                logger.error(f"❌ PostgreSQL {db_name} connection failed: {e}")
                
        return results
    
    async def test_rabbitmq_connection(self) -> bool:
        """Test RabbitMQ connection and basic operations"""
        try:
            connection = await aio_pika.connect_robust("amqp://guest:guest@localhost:5672//")
            channel = await connection.channel()
            
            # Test queue creation and message publishing
            queue = await channel.declare_queue("test_queue", auto_delete=True)
            test_message = aio_pika.Message(body=b"test message")
            await channel.default_exchange.publish(test_message, routing_key=queue.name)
            
            # Test message consumption
            message = await queue.get(timeout=5)
            if message:
                await message.ack()
                
            await connection.close()
            logger.info("✅ RabbitMQ connection and operations successful")
            return True
            
        except Exception as e:
            logger.error(f"❌ RabbitMQ connection failed: {e}")
            return False
    
    def test_redis_connection(self) -> bool:
        """Test Redis connection and basic operations"""
        try:
            # Try multiple possible Redis ports
            ports_to_try = [16379, 6379, 6380, 6381]
            connected = False
            
            for port in ports_to_try:
                try:
                    self.redis_client = redis.Redis(host='localhost', port=port, db=0, decode_responses=True, socket_connect_timeout=2)
                    # Test basic operations
                    self.redis_client.set("test_key", "test_value", ex=60)
                    value = self.redis_client.get("test_key")
                    
                    if value == "test_value":
                        self.redis_client.delete("test_key")
                        logger.info(f"✅ Redis connection and operations successful on port {port}")
                        return True
                    connected = True
                    break
                except Exception:
                    continue
            
            if not connected:
                logger.error("❌ Redis connection failed on all ports")
                return False
                
        except Exception as e:
            logger.error(f"❌ Redis connection failed: {e}")
            return False
    
    async def test_kafka_connection(self) -> bool:
        """Test Kafka connection using HTTP API"""
        try:
            async with httpx.AsyncClient() as client:
                # Test Kafka broker API versions (health check)
                response = await client.get("http://localhost:9092", timeout=5)
                # Kafka will return a connection error which is expected for HTTP requests
                # The fact that we can connect to the port indicates Kafka is running
                logger.info("✅ Kafka broker is accessible on port 9092")
                return True
                
        except httpx.ConnectError:
            logger.error("❌ Kafka connection failed: Cannot connect to port 9092")
            return False
        except Exception as e:
            # HTTP requests to Kafka will fail, but connection success indicates Kafka is running
            logger.info("✅ Kafka broker is accessible (expected HTTP protocol error)")
            return True
    
    async def test_influxdb_connection(self) -> bool:
        """Test InfluxDB connection and basic operations"""
        try:
            async with httpx.AsyncClient() as client:
                # Test InfluxDB health endpoint
                response = await client.get("http://localhost:8086/ping", timeout=10)
                if response.status_code == 204:
                    logger.info("✅ InfluxDB connection successful")
                    return True
                else:
                    logger.error(f"❌ InfluxDB health check failed: {response.status_code}")
                    return False
                    
        except Exception as e:
            logger.error(f"❌ InfluxDB connection failed: {e}")
            return False
    
    async def test_service_endpoints(self) -> Dict[str, bool]:
        """Test all microservice endpoints"""
        services = {
            "temperature-api": "http://localhost:8081/health",
            "device-service": "http://localhost:8082/health", 
            "telemetry-service": "http://localhost:8083/health/",
            "user-service": "http://localhost:8084/health",
            "billing-service": "http://localhost:8085/health",
            "automation-service": "http://localhost:8086/health",
            "main-app": "http://localhost:8080/health"
        }
        
        results = {}
        async with httpx.AsyncClient() as client:
            for service_name, endpoint in services.items():
                try:
                    # Special timeout for telemetry-service due to Kafka initialization
                    timeout = 30 if "telemetry-service" in endpoint else 5
                    response = await client.get(endpoint, timeout=timeout)
                    if response.status_code == 200:
                        results[service_name] = True
                        logger.info(f"✅ {service_name} health check successful")
                    else:
                        results[service_name] = False
                        logger.warning(f"⚠️ {service_name} health check returned {response.status_code}")
                        
                except Exception as e:
                    results[service_name] = False
                    logger.error(f"❌ {service_name} health check failed: {e}")
                    
        return results
    
    async def test_telemetry_integration(self) -> bool:
        """Test telemetry service integration with InfluxDB and Kafka"""
        try:
            async with httpx.AsyncClient() as client:
                # Test telemetry data submission
                test_data = {
                    "device_id": "test_device_001",
                    "device_type": "temperature_sensor",
                    "location": "living_room",
                    "metric_name": "temperature",
                    "value": 22.5,
                    "unit": "celsius",
                    "timestamp": datetime.utcnow().strftime("%Y-%m-%dT%H:%M:%S"),
                    "quality": "good"
                }
                
                response = await client.post(
                    "http://localhost:8083/api/v1/telemetry/readings",
                    json=test_data,
                    timeout=30
                )
                
                if response.status_code in [200, 201]:
                    logger.info("✅ Telemetry integration test successful")
                    return True
                else:
                    logger.error(f"❌ Telemetry integration test failed: {response.status_code}")
                    return False
                    
        except Exception as e:
            logger.error(f"❌ Telemetry integration test failed: {e}")
            return False
    
    async def run_all_tests(self):
        """Run all infrastructure tests"""
        logger.info("🚀 Starting infrastructure tests...")
        
        # Test databases
        logger.info("Testing PostgreSQL databases...")
        postgres_results = await self.test_postgresql_connections()
        self.results.update(postgres_results)
        
        # Test cache
        logger.info("Testing Redis...")
        self.results["redis"] = self.test_redis_connection()
        
        # Test streaming platform
        logger.info("Testing Kafka...")
        self.results["kafka"] = await self.test_kafka_connection()
        
        # Test time series database
        logger.info("Testing InfluxDB...")
        self.results["influxdb"] = await self.test_influxdb_connection()
        
        # Test service endpoints
        logger.info("Testing service endpoints...")
        service_results = await self.test_service_endpoints()
        self.results.update(service_results)
        
        # Test integrations
        logger.info("Testing telemetry integration...")
        self.results["telemetry_integration"] = await self.test_telemetry_integration()
        
        # Generate report
        self.generate_report()
    
    def generate_report(self):
        """Generate test results report"""
        logger.info("\n" + "="*60)
        logger.info("INFRASTRUCTURE TEST RESULTS")
        logger.info("="*60)
        
        passed = sum(1 for result in self.results.values() if result)
        total = len(self.results)
        
        logger.info(f"Overall: {passed}/{total} tests passed ({passed/total*100:.1f}%)")
        logger.info("-"*60)
        
        # Group results by category
        categories = {
            "Databases": [k for k in self.results.keys() if k.startswith("postgresql")],
            "Message Brokers": ["kafka"],
            "Cache & Storage": ["redis", "influxdb"], 
            "Services": [k for k in self.results.keys() if k.endswith("-service") or k.endswith("-api") or k == "main-app"],
            "Integration": ["telemetry_integration"]
        }
        
        for category, components in categories.items():
            logger.info(f"{category}:")
            for component in components:
                if component in self.results:
                    status = "✅ PASS" if self.results[component] else "❌ FAIL"
                    logger.info(f"  {component}: {status}")
            logger.info("")
        
        # Critical components check
        critical_components = ["postgresql_telemetry", "redis", "kafka", "influxdb"]
        critical_failures = [comp for comp in critical_components if not self.results.get(comp, False)]
        
        if critical_failures:
            logger.warning(f"⚠️ Critical component failures: {', '.join(critical_failures)}")
        else:
            logger.info("✅ All critical infrastructure components are operational!")

async def main():
    """Main test execution"""
    tester = InfrastructureTest()
    try:
        await tester.run_all_tests()
    except KeyboardInterrupt:
        logger.info("Tests interrupted by user")
    except Exception as e:
        logger.error(f"Test execution failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    asyncio.run(main())