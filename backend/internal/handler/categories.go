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

type CategoryHandler struct {
	db *pgxpool.Pool
}

func NewCategoryHandler(db *pgxpool.Pool) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, name, display_order, COALESCE(icon, ''), created_at FROM categories ORDER BY display_order`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cats []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.DisplayOrder, &c.Icon, &c.CreatedAt); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		cats = append(cats, c)
	}
	if cats == nil {
		cats = []domain.Category{}
	}
	respondJSON(w, cats)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		respondError(w, "name is required", http.StatusBadRequest)
		return
	}

	var maxOrder int
	h.db.QueryRow(r.Context(), `SELECT COALESCE(MAX(display_order), 0) FROM categories`).Scan(&maxOrder)

	var c domain.Category
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO categories (name, display_order, icon) VALUES ($1, $2, $3)
		 RETURNING id, name, display_order, COALESCE(icon, ''), created_at`,
		input.Name, maxOrder+1, input.Icon,
	).Scan(&c.ID, &c.Name, &c.DisplayOrder, &c.Icon, &c.CreatedAt)
	if err != nil {
		respondError(w, "could not create category", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "category.created", "category", &c.ID, input)
	}

	respondJSON(w, c)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name *string `json:"name"`
		Icon *string `json:"icon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var c domain.Category
	err = h.db.QueryRow(r.Context(),
		`UPDATE categories SET
			name = COALESCE($1, name),
			icon = COALESCE($2, icon)
		 WHERE id = $3
		 RETURNING id, name, display_order, COALESCE(icon, ''), created_at`,
		input.Name, input.Icon, id,
	).Scan(&c.ID, &c.Name, &c.DisplayOrder, &c.Icon, &c.CreatedAt)
	if err != nil {
		respondError(w, "category not found", http.StatusNotFound)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "category.updated", "category", &c.ID, input)
	}

	respondJSON(w, c)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	ct, err := h.db.Exec(r.Context(), `DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		respondError(w, "could not delete", http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		respondError(w, "category not found", http.StatusNotFound)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "category.deleted", "category", &id, nil)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CategoryHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ID    int `json:"id"`
		Delta int `json:"delta"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var currentOrder, total int
	err := h.db.QueryRow(r.Context(),
		`SELECT display_order FROM categories WHERE id = $1`, input.ID).Scan(&currentOrder)
	if err != nil {
		respondError(w, "category not found", http.StatusNotFound)
		return
	}
	h.db.QueryRow(r.Context(), `SELECT COUNT(*) FROM categories`).Scan(&total)

	newOrder := currentOrder + input.Delta
	if newOrder < 1 || newOrder > total {
		respondError(w, "out of bounds", http.StatusBadRequest)
		return
	}

	var otherID int64
	err = h.db.QueryRow(r.Context(),
		`UPDATE categories SET display_order = display_order + $1 WHERE display_order = $2 RETURNING id`,
		-input.Delta, newOrder,
	).Scan(&otherID)
	if err != nil {
		respondError(w, "no category at target position", http.StatusNotFound)
		return
	}

	_, err = h.db.Exec(r.Context(),
		`UPDATE categories SET display_order = $1 WHERE id = $2`, newOrder, input.ID)
	if err != nil {
		respondError(w, "reorder failed", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		id64 := int64(input.ID)
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "category.reordered", "category", &id64, input)
	}

	respondJSON(w, map[string]string{"status": "ok"})
}
