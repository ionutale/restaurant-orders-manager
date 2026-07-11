CREATE TABLE IF NOT EXISTS table_groups (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    party_size INT NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'closed')),
    waiter_id BIGINT NOT NULL REFERENCES users(id),
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    closed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS table_group_tables (
    group_id BIGINT NOT NULL REFERENCES table_groups(id) ON DELETE CASCADE,
    table_id BIGINT NOT NULL REFERENCES tables(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, table_id)
);
