-- Device Service Database Initialization
-- This script initializes the device management database

-- Create devices table (will also be created by JPA, this is for reference)
CREATE TABLE IF NOT EXISTS devices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'offline',
    mac_address VARCHAR(17) UNIQUE,
    ip_address VARCHAR(15),
    firmware_version VARCHAR(50),
    user_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    last_seen TIMESTAMP WITH TIME ZONE
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_devices_type ON devices(type);
CREATE INDEX IF NOT EXISTS idx_devices_location ON devices(location);
CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_mac_address ON devices(mac_address);

-- Insert sample devices for testing
INSERT INTO devices (name, type, location, status, mac_address, user_id) VALUES
('Living Room Temperature Sensor', 'TEMPERATURE_SENSOR', 'Living Room', 'online', '00:1B:44:11:3A:B7', 1),
('Bedroom Temperature Sensor', 'TEMPERATURE_SENSOR', 'Bedroom', 'online', '00:1B:44:11:3A:B8', 1),
('Kitchen Temperature Sensor', 'TEMPERATURE_SENSOR', 'Kitchen', 'online', '00:1B:44:11:3A:B9', 1),
('Living Room Smart Light', 'SMART_LIGHT', 'Living Room', 'offline', '00:1B:44:11:3A:C1', 1),
('Kitchen Smart Plug', 'SMART_PLUG', 'Kitchen', 'online', '00:1B:44:11:3A:C2', 1),
('Front Door Smart Lock', 'SMART_LOCK', 'Entrance', 'online', '00:1B:44:11:3A:C3', 1)
ON CONFLICT (mac_address) DO NOTHING;