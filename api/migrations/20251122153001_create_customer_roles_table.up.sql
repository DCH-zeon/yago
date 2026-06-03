BEGIN;

CREATE TABLE IF NOT EXISTS customer_roles (
    customer_id UUID NOT NULL,
    role_id INTEGER NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (customer_id, role_id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
);

CREATE INDEX idx_customer_roles_customer_id ON customer_roles(customer_id);
CREATE INDEX idx_customer_roles_role_id ON customer_roles(role_id);

COMMIT;