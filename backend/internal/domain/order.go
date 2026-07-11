package domain

type Order struct {
	ID           int64             `json:"id"`
	TableGroupID int64             `json:"table_group_id"`
	WaiterID     int64             `json:"waiter_id"`
	Status       string            `json:"status"`
	CreatedAt    string            `json:"created_at"`
	Courses      []OrderCourse     `json:"courses,omitempty"`
}

type OrderCourse struct {
	ID           int64       `json:"id"`
	OrderID      int64       `json:"order_id"`
	Name         string      `json:"name"`
	DisplayOrder int         `json:"display_order"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID                int64   `json:"id"`
	CourseID          int64   `json:"course_id"`
	DishID            *int64  `json:"dish_id"`
	IsChefSuggestion  bool    `json:"is_chef_suggestion"`
	ChefSuggestionID  *int64  `json:"chef_suggestion_id"`
	Quantity          int     `json:"quantity"`
	Notes             string  `json:"notes"`
	Ready             bool    `json:"ready"`
	ReadyAt           *string `json:"ready_at,omitempty"`
	AddedAt           string  `json:"added_at"`
	DishName          string  `json:"dish_name,omitempty"`
}
