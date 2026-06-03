BEGIN;

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_roles_name ON roles(name);

INSERT INTO roles (id, name, description) VALUES
    (1, 'regular', 'Звичайний клієнт зі стандартними дозволами'),
    (2, 'wholesale', 'Оптовий клієнт зі спеціальними цінами та доступом до оптових замовлень')
ON CONFLICT (name) DO NOTHING;

COMMIT;