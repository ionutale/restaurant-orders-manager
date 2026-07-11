package domain

type ChefSuggestion struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriceCents  int    `json:"price_cents"`
	ShiftDate   string `json:"shift_date"`
	ExpiresAt   string `json:"expires_at"`
	ChefID      int64  `json:"chef_id"`
	ChefName    string `json:"chef_name,omitempty"`
	CreatedAt   string `json:"created_at"`
}
