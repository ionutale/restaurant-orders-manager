package domain

import "time"

type Dish struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PriceCents    int       `json:"price_cents"`
	CategoryID    int64     `json:"category_id"`
	EatingTimeMin int       `json:"eating_time_min"`
	ImageURL      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
}

type DishWithCategory struct {
	Dish
	CategoryName string `json:"category_name"`
}
