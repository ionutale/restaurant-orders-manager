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

type DishHandler struct {
	db *pgxpool.Pool
}

func NewDishHandler(db *pgxpool.Pool) *DishHandler {
	return &DishHandler{db: db}
}

func (h *DishHandler) List(w http.ResponseWriter, r *http.Request) {
	categoryIDStr := r.URL.Query().Get("category_id")

	query := `SELECT d.id, d.name, d.description, d.price_cents, d.category_id, d.eating_time_min, COALESCE(d.image_url, ''), d.created_at,
		COALESCE(c.name, '') as category_name
		FROM dishes d LEFT JOIN categories c ON c.id = d.category_id`
	args := []interface{}{}
	argIdx := 1

	if categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			query += ` WHERE d.category_id = $` + strconv.Itoa(argIdx)
			args = append(args, categoryID)
			argIdx++
		}
	}
	query += ` ORDER BY c.display_order, d.name`

	rows, err := h.db.Query(r.Context(), query, args...)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dishes []domain.DishWithCategory
	for rows.Next() {
		var d domain.DishWithCategory
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.PriceCents, &d.CategoryID, &d.EatingTimeMin, &d.ImageURL, &d.CreatedAt, &d.CategoryName); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		dishes = append(dishes, d)
	}
	if dishes == nil {
		dishes = []domain.DishWithCategory{}
	}
	respondJSON(w, dishes)
}

func (h *DishHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var d domain.DishWithCategory
	err = h.db.QueryRow(r.Context(),
		`SELECT d.id, d.name, d.description, d.price_cents, d.category_id, d.eating_time_min, COALESCE(d.image_url, ''), d.created_at,
			COALESCE(c.name, '')
		FROM dishes d LEFT JOIN categories c ON c.id = d.category_id WHERE d.id = $1`, id,
	).Scan(&d.ID, &d.Name, &d.Description, &d.PriceCents, &d.CategoryID, &d.EatingTimeMin, &d.ImageURL, &d.CreatedAt, &d.CategoryName)
	if err != nil {
		respondError(w, "dish not found", http.StatusNotFound)
		return
	}
	respondJSON(w, d)
}

func (h *DishHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		PriceCents    int    `json:"price_cents"`
		CategoryID    int64  `json:"category_id"`
		EatingTimeMin int    `json:"eating_time_min"`
		ImageURL      string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		respondError(w, "name is required", http.StatusBadRequest)
		return
	}
	if input.CategoryID == 0 {
		respondError(w, "category_id is required", http.StatusBadRequest)
		return
	}
	if input.EatingTimeMin < 1 {
		input.EatingTimeMin = 10
	}

	var d domain.Dish
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO dishes (name, description, price_cents, category_id, eating_time_min, image_url)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, name, description, price_cents, category_id, eating_time_min, COALESCE(image_url, ''), created_at`,
		input.Name, input.Description, input.PriceCents, input.CategoryID, input.EatingTimeMin, input.ImageURL,
	).Scan(&d.ID, &d.Name, &d.Description, &d.PriceCents, &d.CategoryID, &d.EatingTimeMin, &d.ImageURL, &d.CreatedAt)
	if err != nil {
		respondError(w, "could not create dish", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.created", "dish", &d.ID, input)
	}

	respondJSON(w, d)
}

func (h *DishHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name          *string `json:"name"`
		Description   *string `json:"description"`
		PriceCents    *int    `json:"price_cents"`
		CategoryID    *int64  `json:"category_id"`
		EatingTimeMin *int    `json:"eating_time_min"`
		ImageURL      *string `json:"image_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	var d domain.Dish
	err = h.db.QueryRow(r.Context(),
		`UPDATE dishes SET
			name = COALESCE($1, name),
			description = COALESCE($2, description),
			price_cents = COALESCE($3, price_cents),
			category_id = COALESCE($4, category_id),
			eating_time_min = COALESCE($5, eating_time_min),
			image_url = COALESCE($6, image_url)
		 WHERE id = $7
		 RETURNING id, name, description, price_cents, category_id, eating_time_min, COALESCE(image_url, ''), created_at`,
		input.Name, input.Description, input.PriceCents, input.CategoryID, input.EatingTimeMin, input.ImageURL, id,
	).Scan(&d.ID, &d.Name, &d.Description, &d.PriceCents, &d.CategoryID, &d.EatingTimeMin, &d.ImageURL, &d.CreatedAt)
	if err != nil {
		respondError(w, "dish not found", http.StatusNotFound)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.updated", "dish", &d.ID, input)
	}

	respondJSON(w, d)
}

func (h *DishHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	ct, err := h.db.Exec(r.Context(), `DELETE FROM dishes WHERE id = $1`, id)
	if err != nil {
		respondError(w, "could not delete", http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		respondError(w, "dish not found", http.StatusNotFound)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.deleted", "dish", &id, nil)
	}

	w.WriteHeader(http.StatusNoContent)
}
