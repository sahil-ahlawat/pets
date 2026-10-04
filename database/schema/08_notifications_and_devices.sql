-- 08_notifications_and_devices.sql
-- Notification management and delivery audit logs

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pet_id UUID REFERENCES pets(id) ON DELETE SET NULL,
    notification_type notification_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notification_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    channel notification_channel NOT NULL,
    recipient_target TEXT NOT NULL, -- e.g., device token, email, phone
    status VARCHAR(50) NOT NULL, -- 'PENDING', 'SENT', 'FAILED'
    external_message_id VARCHAR(255),
    error_message TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_pet ON notifications(pet_id);
CREATE INDEX idx_notification_deliveries_status ON notification_deliveries(status);
