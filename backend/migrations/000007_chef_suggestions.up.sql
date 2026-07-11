CREATE TABLE IF NOT EXISTS chef_suggestions (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price_cents INT NOT NULL DEFAULT 0,
    shift_date DATE NOT NULL DEFAULT CURRENT_DATE,
    expires_at TIMESTAMPTZ NOT NULL,
    chef_id BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
