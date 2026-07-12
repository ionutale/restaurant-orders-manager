package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
	"github.com/ionutale/restaurant-orders-manager/internal/domain"
)

type ServeOrderHandler struct {
	db *pgxpool.Pool
}

func NewServeOrderHandler(db *pgxpool.Pool) *ServeOrderHandler {
	return &ServeOrderHandler{db: db}
}

func (h *ServeOrderHandler) Start(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		TableIDs  []int64 `json:"table_ids"`
		PartySize int     `json:"party_size"`
		Name      string  `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(input.TableIDs) == 0 {
		respondError(w, "table_ids required", http.StatusBadRequest)
		return
	}
	if input.PartySize < 1 {
		input.PartySize = 1
	}

	var name *string
	if input.Name != "" {
		name = &input.Name
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Create table group
	var groupID int64
	err = tx.QueryRow(r.Context(),
		`INSERT INTO table_groups (name, party_size, waiter_id) VALUES ($1, $2, $3) RETURNING id`,
		name, input.PartySize, claims.UserID,
	).Scan(&groupID)
	if err != nil {
		respondError(w, "could not create group", http.StatusInternalServerError)
		return
	}

	for _, tid := range input.TableIDs {
		_, err = tx.Exec(r.Context(),
			`INSERT INTO table_group_tables (group_id, table_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, groupID, tid)
		if err != nil {
			respondError(w, "could not link table", http.StatusInternalServerError)
			return
		}
	}

	// Create order with a single default course
	var orderID int64
	err = tx.QueryRow(r.Context(),
		`INSERT INTO orders (table_group_id, waiter_id) VALUES ($1, $2) RETURNING id`,
		groupID, claims.UserID,
	).Scan(&orderID)
	if err != nil {
		respondError(w, "could not create order", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(r.Context(),
		`INSERT INTO courses (order_id, name, display_order, status) VALUES ($1, 'Course 1', 1, 'active')`,
		orderID)
	if err != nil {
		respondError(w, "could not create course", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, "commit error", http.StatusInternalServerError)
		return
	}

	RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "start_order", "order", &orderID, input)

	// Return the order
	var o domain.Order
	h.db.QueryRow(r.Context(),
		`SELECT id, table_group_id, waiter_id, status, created_at FROM orders WHERE id = $1`, orderID,
	).Scan(&o.ID, &o.TableGroupID, &o.WaiterID, &o.Status, &o.CreatedAt)
	respondJSON(w, o)
}
