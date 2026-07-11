CREATE TABLE IF NOT EXISTS dishes (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price_cents INT NOT NULL DEFAULT 0,
    category_id BIGINT NOT NULL REFERENCES categories(id),
    eating_time_min INT NOT NULL DEFAULT 10,
    image_url TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO dishes (name, description, price_cents, category_id, eating_time_min) VALUES
    ('Bruschetta', 'Toasted bread with tomato, basil, olive oil', 850, 1, 8),
    ('Calamari', 'Fried squid with lemon aioli', 1200, 1, 10),
    ('Caesar Salad', 'Romaine, parmesan, croutons, caesar dressing', 1050, 1, 8),
    ('Ribeye Steak', '300g prime ribeye with roasted vegetables', 2800, 2, 15),
    ('Grilled Salmon', 'Atlantic salmon with asparagus and hollandaise', 2200, 2, 14),
    ('Pasta Carbonara', 'Spaghetti, pancetta, egg yolk, pecorino', 1600, 2, 12),
    ('Tiramisu', 'Classic Italian coffee dessert', 950, 3, 5),
    ('Panna Cotta', 'Vanilla cream with berry compote', 850, 3, 5),
    ('Espresso', 'Single origin espresso', 350, 4, 5),
    ('Sparkling Water', 'San Pellegrino 750ml', 450, 4, 3),
    ('Chianti Classico', 'DOCG 2021, Italy', 3800, 5, 10),
    ('Sauvignon Blanc', 'Marlborough, New Zealand 2023', 3200, 5, 10);
