CREATE TABLE IF NOT EXISTS tables (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INT NOT NULL DEFAULT 4,
    x FLOAT NOT NULL DEFAULT 0,
    y FLOAT NOT NULL DEFAULT 0,
    label TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO tables (name, capacity, x, y) VALUES
    ('T1', 4, 100, 100),
    ('T2', 4, 300, 100),
    ('T3', 6, 100, 300),
    ('T4', 2, 300, 300),
    ('T5', 8, 200, 500);
