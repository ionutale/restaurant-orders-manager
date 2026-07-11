package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
	"github.com/ionutale/restaurant-orders-manager/internal/domain"
)

type OrderHandler struct {
	db *pgxpool.Pool
}

func NewOrderHandler(db *pgxpool.Pool) *OrderHandler {
	return &OrderHandler{db: db}
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var input struct {
		TableGroupID int64    `json:"table_group_id"`
		CourseNames  []string `json:"course_names"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.TableGroupID == 0 {
		respondError(w, "table_group_id is required", http.StatusBadRequest)
		return
	}
	if len(input.CourseNames) == 0 {
		input.CourseNames = []string{"Appetizer", "Main", "Dessert"}
	}

	tx, err := h.db.Begin(r.Context())
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	var orderID int64
	err = tx.QueryRow(r.Context(),
		`INSERT INTO orders (table_group_id, waiter_id) VALUES ($1, $2) RETURNING id`,
		input.TableGroupID, claims.UserID,
	).Scan(&orderID)
	if err != nil {
		respondError(w, "could not create order", http.StatusInternalServerError)
		return
	}

	for i, name := range input.CourseNames {
		status := "pending"
		if i == 0 {
			status = "active"
		}
		_, err = tx.Exec(r.Context(),
			`INSERT INTO courses (order_id, name, display_order, status) VALUES ($1, $2, $3, $4)`,
			orderID, name, i+1, status)
		if err != nil {
			respondError(w, "could not create course", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(r.Context()); err != nil {
		respondError(w, "commit error", http.StatusInternalServerError)
		return
	}

	order := h.loadOrder(r.Context(), orderID)
	RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "order.created", "order", &orderID, input)
	respondJSON(w, order)
}

func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}
	respondJSON(w, h.loadOrder(r.Context(), id))
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, table_group_id, waiter_id, status, created_at FROM orders ORDER BY created_at DESC`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		if err := rows.Scan(&o.ID, &o.TableGroupID, &o.WaiterID, &o.Status, &o.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, o)
	}
	if orders == nil {
		orders = []domain.Order{}
	}
	respondJSON(w, orders)
}

func (h *OrderHandler) loadOrder(ctx context.Context, id int64) domain.Order {
	var o domain.Order
	h.db.QueryRow(ctx,
		`SELECT id, table_group_id, waiter_id, status, created_at FROM orders WHERE id = $1`, id,
	).Scan(&o.ID, &o.TableGroupID, &o.WaiterID, &o.Status, &o.CreatedAt)

	cRows, _ := h.db.Query(ctx,
		`SELECT id, order_id, name, display_order, status FROM courses WHERE order_id = $1 ORDER BY display_order`, id)
	if cRows != nil {
		defer cRows.Close()
		for cRows.Next() {
			var c domain.OrderCourse
			cRows.Scan(&c.ID, &c.OrderID, &c.Name, &c.DisplayOrder, &c.Status)
			o.Courses = append(o.Courses, c)
		}
	}
	if o.Courses == nil {
		o.Courses = []domain.OrderCourse{}
	}
	return o
}
