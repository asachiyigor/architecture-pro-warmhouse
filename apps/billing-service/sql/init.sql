-- Billing Service Database Initialization

-- Pricing plans table
CREATE TABLE IF NOT EXISTS pricing_plans (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    billing_interval VARCHAR(20) NOT NULL DEFAULT 'monthly', -- monthly, yearly
    device_limit INTEGER NOT NULL DEFAULT 10,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Subscriptions table
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL, -- Reference to user-service
    plan_id UUID NOT NULL REFERENCES pricing_plans(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, cancelled, expired, suspended
    started_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ends_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Payments table
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, completed, failed, refunded
    payment_method VARCHAR(50) NOT NULL,
    external_id VARCHAR(255), -- Payment gateway transaction ID
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_pricing_plans_active ON pricing_plans(is_active);
CREATE INDEX IF NOT EXISTS idx_pricing_plans_price ON pricing_plans(price);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_plan_id ON subscriptions(plan_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscriptions_ends_at ON subscriptions(ends_at);

CREATE INDEX IF NOT EXISTS idx_payments_subscription_id ON payments(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_processed_at ON payments(processed_at);

-- Insert sample pricing plans
INSERT INTO pricing_plans (id, name, description, price, billing_interval, device_limit) VALUES
('770e8400-e29b-41d4-a716-446655440001', 'Basic', 'Perfect for small homes', 9.99, 'monthly', 5),
('770e8400-e29b-41d4-a716-446655440002', 'Standard', 'Great for most families', 19.99, 'monthly', 15),
('770e8400-e29b-41d4-a716-446655440003', 'Premium', 'Ultimate smart home experience', 39.99, 'monthly', 50),
('770e8400-e29b-41d4-a716-446655440004', 'Basic Annual', 'Basic plan with annual billing', 99.99, 'yearly', 5),
('770e8400-e29b-41d4-a716-446655440005', 'Standard Annual', 'Standard plan with annual billing', 199.99, 'yearly', 15),
('770e8400-e29b-41d4-a716-446655440006', 'Premium Annual', 'Premium plan with annual billing', 399.99, 'yearly', 50)
ON CONFLICT (id) DO NOTHING;

-- Insert sample subscriptions (referencing users from user-service)
INSERT INTO subscriptions (id, user_id, plan_id, status, started_at) VALUES
('880e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', '770e8400-e29b-41d4-a716-446655440002', 'active', NOW() - INTERVAL '30 days'),
('880e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', '770e8400-e29b-41d4-a716-446655440001', 'active', NOW() - INTERVAL '15 days')
ON CONFLICT (id) DO NOTHING;

-- Insert sample payments
INSERT INTO payments (id, subscription_id, amount, status, payment_method, external_id, processed_at) VALUES
('990e8400-e29b-41d4-a716-446655440001', '880e8400-e29b-41d4-a716-446655440001', 19.99, 'completed', 'credit_card', 'ext_12345', NOW() - INTERVAL '30 days'),
('990e8400-e29b-41d4-a716-446655440002', '880e8400-e29b-41d4-a716-446655440002', 9.99, 'completed', 'paypal', 'ext_12346', NOW() - INTERVAL '15 days'),
('990e8400-e29b-41d4-a716-446655440003', '880e8400-e29b-41d4-a716-446655440001', 19.99, 'completed', 'credit_card', 'ext_12347', NOW() - INTERVAL '1 days')
ON CONFLICT (id) DO NOTHING;