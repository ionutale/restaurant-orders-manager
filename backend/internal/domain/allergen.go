package domain

type Allergen struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type DishAllergen struct {
	DishID     int64 `json:"dish_id"`
	AllergenID int64 `json:"allergen_id"`
}

type DishSuggestion struct {
	ID             int64  `json:"id"`
	FromDishID     int64  `json:"from_dish_id"`
	ToDishID       int64  `json:"to_dish_id"`
	SuggestionType string `json:"suggestion_type"`
	ToDishName     string `json:"to_dish_name,omitempty"`
}
