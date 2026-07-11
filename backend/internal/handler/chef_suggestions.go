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

type ChefSuggestionHandler struct {
	db *pgxpool.Pool
}

func NewChefSuggestionHandler(db *pgxpool.Pool) *ChefSuggestionHandler {
	return &ChefSuggestionHandler{db: db}
}

func (h *ChefSuggestionHandler) List(w http.ResponseWriter, r *http.Request) {
	includeExpired := r.URL.Query().Get("all") == "true"
	query := `SELECT s.id, s.name, s.description, s.price_cents, s.shift_date, s.expires_at, s.chef_id, COALESCE(u.name, ''), s.created_at
		FROM chef_suggestions s LEFT JOIN users u ON u.id = s.chef_id`
	if !includeExpired {
		query += ` WHERE s.expires_at > NOW()`
	}
	query += ` ORDER BY s.expires_at`

	rows, err := h.db.Query(r.Context(), query)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var list []domain.ChefSuggestion
	for rows.Next() {
		var s domain.ChefSuggestion
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.PriceCents, &s.ShiftDate, &s.ExpiresAt, &s.ChefID, &s.ChefName, &s.CreatedAt); err != nil {
			respondError(w, "scan error", http.StatusInternalServerError)
			return
		}
		list = append(list, s)
	}
	if list == nil {
		list = []domain.ChefSuggestion{}
	}
	respondJSON(w, list)
}

func (h *ChefSuggestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		PriceCents  int    `json:"price_cents"`
		ExpiresAt   string `json:"expires_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.Name == "" {
		respondError(w, "name is required", http.StatusBadRequest)
		return
	}
	expiresAt := time.Now().Add(8 * time.Hour)
	if input.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, input.ExpiresAt)
		if err == nil {
			expiresAt = t
		}
	}

	var s domain.ChefSuggestion
	err := h.db.QueryRow(r.Context(),
		`INSERT INTO chef_suggestions (name, description, price_cents, shift_date, expires_at, chef_id)
		 VALUES ($1, $2, $3, CURRENT_DATE, $4, $5)
		 RETURNING id, name, description, price_cents, shift_date, expires_at, chef_id, created_at`,
		input.Name, input.Description, input.PriceCents, expiresAt, claims.UserID,
	).Scan(&s.ID, &s.Name, &s.Description, &s.PriceCents, &s.ShiftDate, &s.ExpiresAt, &s.ChefID, &s.CreatedAt)
	if err != nil {
		respondError(w, "could not create suggestion", http.StatusInternalServerError)
		return
	}

	RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "chef_suggestion.created", "chef_suggestion", &s.ID, input)
	respondJSON(w, s)
}

func (h *ChefSuggestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Chef can only delete own, manager can delete any
	if claims != nil && claims.Role == "chef" {
		ct, err := h.db.Exec(r.Context(), `DELETE FROM chef_suggestions WHERE id = $1 AND chef_id = $2`, id, claims.UserID)
		if err != nil || ct.RowsAffected() == 0 {
			respondError(w, "not found or not yours", http.StatusNotFound)
			return
		}
	} else {
		ct, err := h.db.Exec(r.Context(), `DELETE FROM chef_suggestions WHERE id = $1`, id)
		if err != nil || ct.RowsAffected() == 0 {
			respondError(w, "not found", http.StatusNotFound)
			return
		}
	}

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "chef_suggestion.deleted", "chef_suggestion", &id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ChefSuggestionHandler) Renew(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var old domain.ChefSuggestion
	err = h.db.QueryRow(r.Context(),
		`SELECT name, description, price_cents FROM chef_suggestions WHERE id = $1`, id,
	).Scan(&old.Name, &old.Description, &old.PriceCents)
	if err != nil {
		respondError(w, "suggestion not found", http.StatusNotFound)
		return
	}

	expiresAt := time.Now().Add(8 * time.Hour)
	var s domain.ChefSuggestion
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO chef_suggestions (name, description, price_cents, shift_date, expires_at, chef_id)
		 VALUES ($1, $2, $3, CURRENT_DATE, $4, $5)
		 RETURNING id, name, description, price_cents, shift_date, expires_at, chef_id, created_at`,
		old.Name, old.Description, old.PriceCents, expiresAt, claims.UserID,
	).Scan(&s.ID, &s.Name, &s.Description, &s.PriceCents, &s.ShiftDate, &s.ExpiresAt, &s.ChefID, &s.CreatedAt)
	if err != nil {
		respondError(w, "could not renew", http.StatusInternalServerError)
		return
	}

	RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "chef_suggestion.renewed", "chef_suggestion", &s.ID, map[string]interface{}{"from_id": id})
	respondJSON(w, s)
}
