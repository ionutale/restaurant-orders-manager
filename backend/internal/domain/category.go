package domain

import "time"

type Category struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	DisplayOrder int       `json:"display_order"`
	Icon         string    `json:"icon"`
	CreatedAt    time.Time `json:"created_at"`
}
