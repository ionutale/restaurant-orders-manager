package domain

import "time"

type ChefSuggestion struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PriceCents  int       `json:"price_cents"`
	ShiftDate   time.Time `json:"shift_date"`
	ExpiresAt   time.Time `json:"expires_at"`
	ChefID      int64     `json:"chef_id"`
	ChefName    string    `json:"chef_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
