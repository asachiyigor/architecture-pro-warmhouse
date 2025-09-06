-- Automation Service Database Initialization

-- Scenarios table
CREATE TABLE IF NOT EXISTS scenarios (
    id UUID PRIMARY KEY,
    household_id UUID NOT NULL, -- Reference to user-service household
    created_by UUID NOT NULL, -- Reference to user-service user
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Scenario actions table
CREATE TABLE IF NOT EXISTS scenario_actions (
    id UUID PRIMARY KEY,
    scenario_id UUID NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    device_id VARCHAR(100) NOT NULL, -- Reference to device-service device
    action_type VARCHAR(50) NOT NULL, -- turn_on, turn_off, set_temperature, etc.
    parameters TEXT, -- JSON string with action parameters
    order_position INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Scenario triggers table
CREATE TABLE IF NOT EXISTS scenario_triggers (
    id UUID PRIMARY KEY,
    scenario_id UUID NOT NULL REFERENCES scenarios(id) ON DELETE CASCADE,
    trigger_type VARCHAR(50) NOT NULL, -- time, device_state, manual, weather
    device_id VARCHAR(100), -- Reference to device-service device (optional for time triggers)
    condition TEXT NOT NULL, -- JSON string with trigger conditions
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Scenario executions table (audit log)
CREATE TABLE IF NOT EXISTS scenario_executions (
    id UUID PRIMARY KEY,
    scenario_id UUID NOT NULL REFERENCES scenarios(id),
    trigger_type VARCHAR(50) NOT NULL, -- manual, scheduled, device_event
    status VARCHAR(20) NOT NULL DEFAULT 'started', -- started, completed, failed, cancelled
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_scenarios_household_id ON scenarios(household_id);
CREATE INDEX IF NOT EXISTS idx_scenarios_created_by ON scenarios(created_by);
CREATE INDEX IF NOT EXISTS idx_scenarios_active ON scenarios(is_active);

CREATE INDEX IF NOT EXISTS idx_scenario_actions_scenario_id ON scenario_actions(scenario_id);
CREATE INDEX IF NOT EXISTS idx_scenario_actions_device_id ON scenario_actions(device_id);
CREATE INDEX IF NOT EXISTS idx_scenario_actions_order ON scenario_actions(scenario_id, order_position);

CREATE INDEX IF NOT EXISTS idx_scenario_triggers_scenario_id ON scenario_triggers(scenario_id);
CREATE INDEX IF NOT EXISTS idx_scenario_triggers_device_id ON scenario_triggers(device_id);
CREATE INDEX IF NOT EXISTS idx_scenario_triggers_type ON scenario_triggers(trigger_type);
CREATE INDEX IF NOT EXISTS idx_scenario_triggers_active ON scenario_triggers(is_active);

CREATE INDEX IF NOT EXISTS idx_scenario_executions_scenario_id ON scenario_executions(scenario_id);
CREATE INDEX IF NOT EXISTS idx_scenario_executions_status ON scenario_executions(status);
CREATE INDEX IF NOT EXISTS idx_scenario_executions_started_at ON scenario_executions(started_at);

-- Insert sample scenarios (referencing users and households from user-service)
INSERT INTO scenarios (id, household_id, created_by, name, description, is_active) VALUES
('aa0e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'Good Morning', 'Turn on lights and set temperature when waking up', true),
('aa0e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'Leaving Home', 'Turn off all devices when leaving', true),
('aa0e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'Movie Night', 'Dim lights and close blinds', true)
ON CONFLICT (id) DO NOTHING;

-- Insert sample scenario actions
INSERT INTO scenario_actions (id, scenario_id, device_id, action_type, parameters, order_position) VALUES
('bb0e8400-e29b-41d4-a716-446655440001', 'aa0e8400-e29b-41d4-a716-446655440001', '4', 'turn_on', '{"brightness": 80}', 1),
('bb0e8400-e29b-41d4-a716-446655440002', 'aa0e8400-e29b-41d4-a716-446655440001', '1', 'set_temperature', '{"temperature": 22}', 2),
('bb0e8400-e29b-41d4-a716-446655440003', 'aa0e8400-e29b-41d4-a716-446655440002', '4', 'turn_off', '{}', 1),
('bb0e8400-e29b-41d4-a716-446655440004', 'aa0e8400-e29b-41d4-a716-446655440002', '5', 'turn_off', '{}', 2),
('bb0e8400-e29b-41d4-a716-446655440005', 'aa0e8400-e29b-41d4-a716-446655440003', '4', 'set_brightness', '{"brightness": 20}', 1)
ON CONFLICT (id) DO NOTHING;

-- Insert sample scenario triggers
INSERT INTO scenario_triggers (id, scenario_id, trigger_type, device_id, condition, is_active) VALUES
('cc0e8400-e29b-41d4-a716-446655440001', 'aa0e8400-e29b-41d4-a716-446655440001', 'time', NULL, '{"time": "07:00", "days": ["monday", "tuesday", "wednesday", "thursday", "friday"]}', true),
('cc0e8400-e29b-41d4-a716-446655440002', 'aa0e8400-e29b-41d4-a716-446655440002', 'manual', NULL, '{"user_trigger": true}', true),
('cc0e8400-e29b-41d4-a716-446655440003', 'aa0e8400-e29b-41d4-a716-446655440003', 'manual', NULL, '{"user_trigger": true}', true)
ON CONFLICT (id) DO NOTHING;

-- Insert sample executions
INSERT INTO scenario_executions (id, scenario_id, trigger_type, status, started_at, completed_at) VALUES
('dd0e8400-e29b-41d4-a716-446655440001', 'aa0e8400-e29b-41d4-a716-446655440001', 'time', 'completed', NOW() - INTERVAL '1 hour', NOW() - INTERVAL '1 hour' + INTERVAL '30 seconds'),
('dd0e8400-e29b-41d4-a716-446655440002', 'aa0e8400-e29b-41d4-a716-446655440002', 'manual', 'completed', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours' + INTERVAL '15 seconds'),
('dd0e8400-e29b-41d4-a716-446655440003', 'aa0e8400-e29b-41d4-a716-446655440003', 'manual', 'completed', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours' + INTERVAL '20 seconds')
ON CONFLICT (id) DO NOTHING;