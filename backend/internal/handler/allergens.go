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

type AllergenHandler struct {
	db *pgxpool.Pool
}

func NewAllergenHandler(db *pgxpool.Pool) *AllergenHandler {
	return &AllergenHandler{db: db}
}

func (h *AllergenHandler) ListAllergens(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(), `SELECT id, name, COALESCE(icon, '') FROM allergens ORDER BY name`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []domain.Allergen
	for rows.Next() {
		var a domain.Allergen
		if err := rows.Scan(&a.ID, &a.Name, &a.Icon); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		list = append(list, a)
	}
	if list == nil {
		list = []domain.Allergen{}
	}
	respondJSON(w, list)
}

func (h *AllergenHandler) CreateAllergen(w http.ResponseWriter, r *http.Request) {
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

	var a domain.Allergen
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO allergens (name, icon) VALUES ($1, $2) RETURNING id, name, COALESCE(icon, '')`,
		input.Name, input.Icon,
	).Scan(&a.ID, &a.Name, &a.Icon)
	if err != nil {
		respondError(w, "could not create allergen", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "allergen.created", "allergen", &a.ID, input)
	}
	respondJSON(w, a)
}

func (h *AllergenHandler) DeleteAllergen(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}
	ct, err := h.db.Exec(r.Context(), `DELETE FROM allergens WHERE id = $1`, id)
	if err != nil || ct.RowsAffected() == 0 {
		respondError(w, "not found", http.StatusNotFound)
		return
	}
	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "allergen.deleted", "allergen", &id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AllergenHandler) GetDishAllergens(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.ParseInt(chi.URLParam(r, "dishId"), 10, 64)
	if err != nil {
		respondError(w, "invalid dish id", http.StatusBadRequest)
		return
	}
	rows, err := h.db.Query(r.Context(),
		`SELECT a.id, a.name, COALESCE(a.icon, '')
		 FROM dish_allergens da JOIN allergens a ON a.id = da.allergen_id
		 WHERE da.dish_id = $1 ORDER BY a.name`, dishID)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []domain.Allergen
	for rows.Next() {
		var a domain.Allergen
		if err := rows.Scan(&a.ID, &a.Name, &a.Icon); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		list = append(list, a)
	}
	if list == nil {
		list = []domain.Allergen{}
	}
	respondJSON(w, list)
}

func (h *AllergenHandler) SetDishAllergens(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.ParseInt(chi.URLParam(r, "dishId"), 10, 64)
	if err != nil {
		respondError(w, "invalid dish id", http.StatusBadRequest)
		return
	}

	var input struct {
		AllergenIDs []int64 `json:"allergen_ids"`
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

	_, _ = tx.Exec(r.Context(), `DELETE FROM dish_allergens WHERE dish_id = $1`, dishID)
	for _, aid := range input.AllergenIDs {
		_, _ = tx.Exec(r.Context(), `INSERT INTO dish_allergens (dish_id, allergen_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, dishID, aid)
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, "commit error", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.allergens.updated", "dish", &dishID, input)
	}
	respondJSON(w, map[string]string{"status": "ok"})
}

func (h *AllergenHandler) GetDishSuggestions(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.ParseInt(chi.URLParam(r, "dishId"), 10, 64)
	if err != nil {
		respondError(w, "invalid dish id", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT ds.id, ds.from_dish_id, ds.to_dish_id, ds.suggestion_type, COALESCE(d.name, '')
		 FROM dish_suggestions ds JOIN dishes d ON d.id = ds.to_dish_id
		 WHERE ds.from_dish_id = $1 ORDER BY ds.suggestion_type, d.name`, dishID)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []domain.DishSuggestion
	for rows.Next() {
		var s domain.DishSuggestion
		if err := rows.Scan(&s.ID, &s.FromDishID, &s.ToDishID, &s.SuggestionType, &s.ToDishName); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		list = append(list, s)
	}
	if list == nil {
		list = []domain.DishSuggestion{}
	}
	respondJSON(w, list)
}

func (h *AllergenHandler) CreateDishSuggestion(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.ParseInt(chi.URLParam(r, "dishId"), 10, 64)
	if err != nil {
		respondError(w, "invalid dish id", http.StatusBadRequest)
		return
	}

	var input struct {
		ToDishID       int64  `json:"to_dish_id"`
		SuggestionType string `json:"suggestion_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.SuggestionType != "wine" && input.SuggestionType != "side" {
		respondError(w, "type must be 'wine' or 'side'", http.StatusBadRequest)
		return
	}

	var s domain.DishSuggestion
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO dish_suggestions (from_dish_id, to_dish_id, suggestion_type)
		 VALUES ($1, $2, $3) RETURNING id, from_dish_id, to_dish_id, suggestion_type`,
		dishID, input.ToDishID, input.SuggestionType,
	).Scan(&s.ID, &s.FromDishID, &s.ToDishID, &s.SuggestionType)
	if err != nil {
		respondError(w, "could not create suggestion", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.suggestion.created", "dish", &dishID, input)
	}
	respondJSON(w, s)
}

func (h *AllergenHandler) DeleteDishSuggestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}
	ct, err := h.db.Exec(r.Context(), `DELETE FROM dish_suggestions WHERE id = $1`, id)
	if err != nil || ct.RowsAffected() == 0 {
		respondError(w, "not found", http.StatusNotFound)
		return
	}
	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "dish.suggestion.deleted", "dish_suggestion", &id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}
