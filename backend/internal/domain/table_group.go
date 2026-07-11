package domain

import "time"

type TableGroup struct {
	ID        int64      `json:"id"`
	Name      *string    `json:"name"`
	PartySize int        `json:"party_size"`
	Status    string     `json:"status"`
	WaiterID  int64      `json:"waiter_id"`
	OpenedAt  time.Time  `json:"opened_at"`
	ClosedAt  *time.Time `json:"closed_at"`
	TableIDs  []int64    `json:"table_ids,omitempty"`
}
