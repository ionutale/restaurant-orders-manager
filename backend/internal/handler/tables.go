package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
	"github.com/ionutale/restaurant-orders-manager/internal/domain"
)

type TableHandler struct {
	db *pgxpool.Pool
}

func NewTableHandler(db *pgxpool.Pool) *TableHandler {
	return &TableHandler{db: db}
}

func (h *TableHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, name, capacity, x, y, label, created_at, updated_at FROM tables ORDER BY name`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []domain.Table
	for rows.Next() {
		var t domain.Table
		if err := rows.Scan(&t.ID, &t.Name, &t.Capacity, &t.X, &t.Y, &t.Label, &t.CreatedAt, &t.UpdatedAt); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		tables = append(tables, t)
	}
	if tables == nil {
		tables = []domain.Table{}
	}
	respondJSON(w, tables)
}

func (h *TableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string   `json:"name"`
		Capacity int      `json:"capacity"`
		X        *float64 `json:"x"`
		Y        *float64 `json:"y"`
		Label    *string  `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		respondError(w, "name is required", http.StatusBadRequest)
		return
	}
	if input.Capacity < 1 {
		input.Capacity = 4
	}
	x := 0.0
	y := 0.0
	if input.X != nil && input.Y != nil {
		x = *input.X
		y = *input.Y
	} else {
		// Auto-place in a grid based on existing table count
		var count int
		h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM tables`).Scan(&count)
		cols := 5
		row := count / cols
		col := count % cols
		x = 30.0 + float64(col)*160.0
		y = 30.0 + float64(row)*130.0
	}

	var t domain.Table
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO tables (name, capacity, x, y, label) VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, name, capacity, x, y, label, created_at, updated_at`,
		input.Name, input.Capacity, x, y, input.Label,
	).Scan(&t.ID, &t.Name, &t.Capacity, &t.X, &t.Y, &t.Label, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		respondError(w, "could not create table", http.StatusInternalServerError)
		return
	}

	claims := auth.ClaimsFromCtx(r.Context())
	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table.created", "table", &t.ID, input)
	}

	respondJSON(w, t)
}

func (h *TableHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name     *string  `json:"name"`
		Capacity *int     `json:"capacity"`
		X        *float64 `json:"x"`
		Y        *float64 `json:"y"`
		Label    *string  `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var t domain.Table
	err = h.db.QueryRow(r.Context(),
		`UPDATE tables SET
			name = COALESCE($1, name),
			capacity = COALESCE($2, capacity),
			x = COALESCE($3, x),
			y = COALESCE($4, y),
			label = COALESCE($5, label),
			updated_at = NOW()
		 WHERE id = $6
		 RETURNING id, name, capacity, x, y, label, created_at, updated_at`,
		input.Name, input.Capacity, input.X, input.Y, input.Label, id,
	).Scan(&t.ID, &t.Name, &t.Capacity, &t.X, &t.Y, &t.Label, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		respondError(w, "table not found", http.StatusNotFound)
		return
	}

	claims := auth.ClaimsFromCtx(r.Context())
	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table.updated", "table", &t.ID, input)
	}

	respondJSON(w, t)
}

func (h *TableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	ct, err := h.db.Exec(r.Context(), `DELETE FROM tables WHERE id = $1`, id)
	if err != nil {
		respondError(w, "could not delete", http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		respondError(w, "table not found", http.StatusNotFound)
		return
	}

	claims := auth.ClaimsFromCtx(r.Context())
	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "table.deleted", "table", &id, nil)
	}

	w.WriteHeader(http.StatusNoContent)
}
