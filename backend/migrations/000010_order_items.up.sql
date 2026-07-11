CREATE TABLE IF NOT EXISTS order_items (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    dish_id BIGINT REFERENCES dishes(id),
    is_chef_suggestion BOOLEAN NOT NULL DEFAULT false,
    chef_suggestion_id BIGINT REFERENCES chef_suggestions(id),
    quantity INT NOT NULL DEFAULT 1,
    notes TEXT NOT NULL DEFAULT '',
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
