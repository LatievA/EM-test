CREATE TABLE subscriptions (
    user_id UUID NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    start_date DATE NOT NULL, -- store as first day of the month
    end_date DATE,
    PRIMARY KEY (user_id, service_name)
);

CREATE INDEX idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX idx_subscriptions_service_name ON subscriptions(service_name);
CREATE INDEX idx_subscriptions_dates ON subscriptions(start_date, end_date);

INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES
-- Подписки для пользователя 1
('60601fee-2bf1-4721-ae6f-7636e79a0cba', 'Yandex Plus', 400, '2025-01-01', '2025-12-31'),
('60601fee-2bf1-4721-ae6f-7636e79a0cba', 'Netflix', 500, '2025-03-01', NULL),
('60601fee-2bf1-4721-ae6f-7636e79a0cba', 'Spotify', 300, '2024-11-01', '2025-05-31'),

-- Подписки для пользователя 2
('a1b2c3d4-5678-90ef-1234-012345678912', 'Yandex Plus', 400, '2025-02-01', NULL),
('a1b2c3d4-5678-90ef-1234-012345678912', 'Apple Music', 250, '2025-01-01', '2025-06-30'),

-- Подписки для пользователя 3
('b2c3d4e5-6789-0123-1234-567890abcdef', 'Disney+', 450, '2025-04-01', NULL),
('b2c3d4e5-6789-0123-1234-567890abcdef', 'Amazon Prime', 350, '2025-01-01', '2025-12-31'),
('b2c3d4e5-6789-0123-1234-567890abcdef', 'Yandex Plus', 400, '2024-09-01', '2025-08-31');