package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
)

type userResp struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, name, email, role FROM users ORDER BY role, name`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []userResp
	for rows.Next() {
		var u userResp
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			continue
		}
		users = append(users, u)
	}
	if users == nil {
		users = []userResp{}
	}
	respondJSON(w, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || input.Email == "" || input.Password == "" {
		respondError(w, "name, email, password required", http.StatusBadRequest)
		return
	}
	if input.Role != "waiter" && input.Role != "chef" && input.Role != "manager" {
		respondError(w, "role must be waiter, chef, or manager", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, "internal error", http.StatusInternalServerError)
		return
	}

	var u userResp
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, name, email, role`,
		input.Name, input.Email, string(hash), input.Role,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err != nil {
		if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
			respondError(w, "email already registered", http.StatusConflict)
			return
		}
		respondError(w, "could not create user", http.StatusInternalServerError)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "user.created", "user", &u.ID, input)
	}
	respondJSON(w, u)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
		Role  *string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Only managers can change roles
	if input.Role != nil && claims.Role != "manager" {
		respondError(w, "only managers can change roles", http.StatusForbidden)
		return
	}

	// Non-managers can only edit their own profile
	if claims.Role != "manager" && claims.UserID != id {
		respondError(w, "cannot edit other users", http.StatusForbidden)
		return
	}
	// Managers cannot change their own role (prevent self-demotion)
	if claims.Role == "manager" && input.Role != nil && claims.UserID == id {
		respondError(w, "cannot change your own role", http.StatusBadRequest)
		return
	}

	var u userResp
	err = h.db.QueryRow(r.Context(),
		`UPDATE users SET
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			role = COALESCE($3, role)
		 WHERE id = $4
		 RETURNING id, name, email, role`,
		input.Name, input.Email, input.Role, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Role)
	if err != nil {
		respondError(w, "user not found", http.StatusNotFound)
		return
	}

	if claims := auth.ClaimsFromCtx(r.Context()); claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "user.updated", "user", &u.ID, input)
	}
	respondJSON(w, u)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	claims := auth.ClaimsFromCtx(r.Context())
	if claims != nil && claims.UserID == id {
		respondError(w, "cannot delete yourself", http.StatusBadRequest)
		return
	}

	ct, err := h.db.Exec(r.Context(), `DELETE FROM users WHERE id = $1`, id)
	if err != nil || ct.RowsAffected() == 0 {
		respondError(w, "user not found", http.StatusNotFound)
		return
	}

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "user.deleted", "user", &id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}
