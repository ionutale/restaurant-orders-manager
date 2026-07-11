package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
	"github.com/ionutale/restaurant-orders-manager/internal/domain"
)

type TableGroupHandler struct {
	db *pgxpool.Pool
}

func NewTableGroupHandler(db *pgxpool.Pool) *TableGroupHandler {
	return &TableGroupHandler{db: db}
}

func (h *TableGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		respondError(w, "at least one table_id is required", http.StatusBadRequest)
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

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, "commit error", http.StatusInternalServerError)
		return
	}

	var g domain.TableGroup
	err = h.db.QueryRow(r.Context(),
		`SELECT id, name, party_size, status, waiter_id, opened_at, closed_at FROM table_groups WHERE id = $1`, groupID,
	).Scan(&g.ID, &g.Name, &g.PartySize, &g.Status, &g.WaiterID, &g.OpenedAt, &g.ClosedAt)
	if err != nil {
		respondError(w, "not found", http.StatusInternalServerError)
		return
	}

	RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table_group.created", "table_group", &groupID, input)
	respondJSON(w, g)
}

func (h *TableGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var g domain.TableGroup
	err = h.db.QueryRow(r.Context(),
		`SELECT id, name, party_size, status, waiter_id, opened_at, closed_at FROM table_groups WHERE id = $1`, id,
	).Scan(&g.ID, &g.Name, &g.PartySize, &g.Status, &g.WaiterID, &g.OpenedAt, &g.ClosedAt)
	if err != nil {
		respondError(w, "not found", http.StatusNotFound)
		return
	}

	tRows, err := h.db.Query(r.Context(), `SELECT table_id FROM table_group_tables WHERE group_id = $1`, id)
	if err == nil {
		defer tRows.Close()
		for tRows.Next() {
			var tid int64
			tRows.Scan(&tid)
			g.TableIDs = append(g.TableIDs, tid)
		}
	}

	respondJSON(w, g)
}

func (h *TableGroupHandler) Close(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	now := time.Now()
	_, err = h.db.Exec(r.Context(),
		`UPDATE table_groups SET status = 'closed', closed_at = $1 WHERE id = $2 AND status != 'closed'`,
		now, id)
	if err != nil {
		respondError(w, "could not close", http.StatusInternalServerError)
		return
	}

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table_group.closed", "table_group", &id, nil)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *TableGroupHandler) UpdateTables(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		AddTableIDs    []int64 `json:"add_table_ids"`
		RemoveTableIDs []int64 `json:"remove_table_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	for _, tid := range input.AddTableIDs {
		_, _ = tx.Exec(r.Context(),
			`INSERT INTO table_group_tables (group_id, table_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, tid)
	}
	for _, tid := range input.RemoveTableIDs {
		_, _ = tx.Exec(r.Context(),
			`DELETE FROM table_group_tables WHERE group_id = $1 AND table_id = $2`, id, tid)
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, "commit error", http.StatusInternalServerError)
		return
	}

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table_group.tables_updated", "table_group", &id, input)
	}
	respondJSON(w, map[string]string{"status": "ok"})
}
