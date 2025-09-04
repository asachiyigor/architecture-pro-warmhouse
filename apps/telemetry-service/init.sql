-- Telemetry Service Database Initialization
-- This script initializes the telemetry database

-- Create telemetry_readings table
CREATE TABLE IF NOT EXISTS telemetry_readings (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    metric_name VARCHAR(50) NOT NULL,
    value FLOAT NOT NULL,
    unit VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    raw_data TEXT,
    quality VARCHAR(20) DEFAULT 'good',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create telemetry_alerts table
CREATE TABLE IF NOT EXISTS telemetry_alerts (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    current_value FLOAT NOT NULL,
    threshold_value FLOAT NOT NULL,
    message TEXT NOT NULL,
    is_active VARCHAR(10) DEFAULT 'true',
    acknowledged_at TIMESTAMP WITH TIME ZONE,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create device_metrics table for aggregated data
CREATE TABLE IF NOT EXISTS device_metrics (
    id SERIAL PRIMARY KEY,
    device_id VARCHAR(50) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    metric_name VARCHAR(50) NOT NULL,
    avg_value FLOAT,
    min_value FLOAT,
    max_value FLOAT,
    sample_count INTEGER DEFAULT 0,
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    aggregation_type VARCHAR(20) DEFAULT 'hourly',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_device_id ON telemetry_readings(device_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_device_timestamp ON telemetry_readings(device_id, timestamp);
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_location ON telemetry_readings(location);
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_location_timestamp ON telemetry_readings(location, timestamp);
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_metric_timestamp ON telemetry_readings(metric_name, timestamp);
CREATE INDEX IF NOT EXISTS idx_telemetry_readings_timestamp ON telemetry_readings(timestamp);

CREATE INDEX IF NOT EXISTS idx_telemetry_alerts_device_id ON telemetry_alerts(device_id);
CREATE INDEX IF NOT EXISTS idx_telemetry_alerts_is_active ON telemetry_alerts(is_active);
CREATE INDEX IF NOT EXISTS idx_telemetry_alerts_created_at ON telemetry_alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_telemetry_alerts_alert_type ON telemetry_alerts(alert_type);

CREATE INDEX IF NOT EXISTS idx_device_metrics_device_id ON device_metrics(device_id);
CREATE INDEX IF NOT EXISTS idx_device_metrics_device_period ON device_metrics(device_id, period_start);
CREATE INDEX IF NOT EXISTS idx_device_metrics_location_period ON device_metrics(location, period_start);

-- Insert sample telemetry data for testing
INSERT INTO telemetry_readings (device_id, device_type, location, metric_name, value, unit, timestamp) VALUES
('1', 'temperature_sensor', 'Living Room', 'temperature', 22.5, '°C', NOW() - INTERVAL '1 hour'),
('1', 'temperature_sensor', 'Living Room', 'temperature', 23.1, '°C', NOW() - INTERVAL '30 minutes'),
('1', 'temperature_sensor', 'Living Room', 'temperature', 23.8, '°C', NOW()),
('2', 'temperature_sensor', 'Bedroom', 'temperature', 21.2, '°C', NOW() - INTERVAL '1 hour'),
('2', 'temperature_sensor', 'Bedroom', 'temperature', 21.8, '°C', NOW() - INTERVAL '30 minutes'),
('2', 'temperature_sensor', 'Bedroom', 'temperature', 22.3, '°C', NOW()),
('3', 'temperature_sensor', 'Kitchen', 'temperature', 24.1, '°C', NOW() - INTERVAL '1 hour'),
('3', 'temperature_sensor', 'Kitchen', 'temperature', 24.7, '°C', NOW() - INTERVAL '30 minutes'),
('3', 'temperature_sensor', 'Kitchen', 'temperature', 25.2, '°C', NOW());