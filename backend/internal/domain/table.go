package domain

type Table struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Capacity  int     `json:"capacity"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Label     *string `json:"label"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
