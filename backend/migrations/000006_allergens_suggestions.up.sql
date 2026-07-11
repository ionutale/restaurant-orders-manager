CREATE TABLE IF NOT EXISTS allergens (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    icon TEXT NOT NULL DEFAULT ''
);

INSERT INTO allergens (name, icon) VALUES
    ('Gluten', '🌾'),
    ('Lactose', '🥛'),
    ('Eggs', '🥚'),
    ('Fish', '🐟'),
    ('Shellfish', '🦐'),
    ('Tree Nuts', '🥜'),
    ('Peanuts', '🥜'),
    ('Soy', '🫘'),
    ('Sesame', '🌱'),
    ('Sulfites', '🧪'),
    ('Celery', '🥬'),
    ('Mustard', '🌭'),
    ('Molluscs', '🐌'),
    ('Lupin', '🌼');

CREATE TABLE IF NOT EXISTS dish_allergens (
    dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    allergen_id BIGINT NOT NULL REFERENCES allergens(id) ON DELETE CASCADE,
    PRIMARY KEY (dish_id, allergen_id)
);

CREATE TABLE IF NOT EXISTS dish_suggestions (
    id BIGSERIAL PRIMARY KEY,
    from_dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    to_dish_id BIGINT NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    suggestion_type TEXT NOT NULL CHECK (suggestion_type IN ('wine', 'side')),
    UNIQUE (from_dish_id, to_dish_id, suggestion_type)
);
