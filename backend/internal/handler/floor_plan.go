package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type floorPlanTable struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Capacity  int     `json:"capacity"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Label     *string `json:"label"`
	Status    string  `json:"status"`
	GroupID   *int64  `json:"group_id,omitempty"`
	GroupName *string `json:"group_name,omitempty"`
	PartySize *int    `json:"party_size,omitempty"`
}

type FloorPlanHandler struct {
	db *pgxpool.Pool
}

func NewFloorPlanHandler(db *pgxpool.Pool) *FloorPlanHandler {
	return &FloorPlanHandler{db: db}
}

func (h *FloorPlanHandler) GetFloorPlan(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT t.id, t.name, t.capacity, t.x, t.y, t.label,
			CASE WHEN tg.id IS NOT NULL THEN 'occupied' ELSE 'free' END AS status,
			tg.id AS group_id, tg.name AS group_name, tg.party_size
		FROM tables t
		LEFT JOIN table_group_tables tgt ON tgt.table_id = t.id
		LEFT JOIN table_groups tg ON tg.id = tgt.group_id AND tg.status != 'closed'
		ORDER BY t.name`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []floorPlanTable
	for rows.Next() {
		var t floorPlanTable
		if err := rows.Scan(&t.ID, &t.Name, &t.Capacity, &t.X, &t.Y, &t.Label, &t.Status, &t.GroupName, &t.PartySize); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		tables = append(tables, t)
	}
	if tables == nil {
		tables = []floorPlanTable{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}
