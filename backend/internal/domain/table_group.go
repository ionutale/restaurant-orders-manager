package domain

type TableGroup struct {
	ID        int64   `json:"id"`
	Name      *string `json:"name"`
	PartySize int     `json:"party_size"`
	Status    string  `json:"status"`
	WaiterID  int64   `json:"waiter_id"`
	OpenedAt  string  `json:"opened_at"`
	ClosedAt  *string `json:"closed_at"`
	TableIDs  []int64 `json:"table_ids,omitempty"`
}
