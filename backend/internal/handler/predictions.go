package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PredictionHandler struct {
	db *pgxpool.Pool
}

func NewPredictionHandler(db *pgxpool.Pool) *PredictionHandler {
	return &PredictionHandler{db: db}
}

func (h *PredictionHandler) TablePredictions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `
		SELECT t.id, t.name, tg.id, tg.name, tg.party_size, o.id,
			(SELECT MAX(oi.ready_at + (d.eating_time_min * interval '1 minute'))
			 FROM order_items oi
			 JOIN courses c ON c.id = oi.course_id
			 JOIN dishes d ON d.id = oi.dish_id
			 WHERE c.order_id = o.id AND oi.ready = true) AS estimated_free
		FROM tables t
		JOIN table_group_tables tgt ON tgt.table_id = t.id
		JOIN table_groups tg ON tg.id = tgt.group_id AND tg.status != 'closed'
		JOIN orders o ON o.table_group_id = tg.id AND o.status = 'sent'
		ORDER BY t.name`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type pred struct {
		TableID       int64   `json:"table_id"`
		TableName     string  `json:"table_name"`
		GroupName     *string `json:"group_name"`
		PartySize     *int    `json:"party_size"`
		EstimatedFree *string `json:"estimated_free"`
		OrderID       int64   `json:"order_id"`
	}
	var predictions []pred
	for rows.Next() {
		var p pred
		var estimatedFree *time.Time
		if err := rows.Scan(&p.TableID, &p.TableName, &p.GroupName, &p.PartySize, &p.OrderID, &estimatedFree); err != nil {
			continue
		}
		if estimatedFree != nil {
			s := estimatedFree.Format(time.RFC3339)
			p.EstimatedFree = &s
		}
		predictions = append(predictions, p)
	}
	if predictions == nil {
		predictions = []pred{}
	}
	respondJSON(w, predictions)
}
