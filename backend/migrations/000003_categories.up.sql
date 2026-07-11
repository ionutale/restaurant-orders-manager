CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    display_order INT NOT NULL DEFAULT 0,
    icon TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO categories (name, display_order) VALUES
    ('Appetizers', 1),
    ('Mains', 2),
    ('Desserts', 3),
    ('Drinks', 4),
    ('Wines', 5);
