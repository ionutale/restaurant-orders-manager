package domain

type Order struct {
	ID           int64         `json:"id"`
	TableGroupID int64         `json:"table_group_id"`
	WaiterID     int64         `json:"waiter_id"`
	Status       string        `json:"status"`
	CreatedAt    string        `json:"created_at"`
	Courses      []OrderCourse `json:"courses,omitempty"`
}

type OrderCourse struct {
	ID           int64  `json:"id"`
	OrderID      int64  `json:"order_id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"display_order"`
	Status       string `json:"status"`
}
