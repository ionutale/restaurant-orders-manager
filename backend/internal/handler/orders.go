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
	claims := auth.ClaimsFromCtx(r.Context())
	query := `SELECT id, table_group_id, waiter_id, status, created_at FROM orders`
	var args []interface{}
	if claims != nil && claims.Role == "waiter" {
		query += ` WHERE waiter_id = $1`
		args = append(args, claims.UserID)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := h.db.Query(r.Context(), query, args...)
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

func (h *OrderHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	courseID, err := strconv.ParseInt(chi.URLParam(r, "courseId"), 10, 64)
	if err != nil {
		respondError(w, "invalid course id", http.StatusBadRequest)
		return
	}

	var input struct {
		DishID           *int64 `json:"dish_id"`
		IsChefSuggestion bool   `json:"is_chef_suggestion"`
		ChefSuggestionID *int64 `json:"chef_suggestion_id"`
		Quantity         int    `json:"quantity"`
		Notes            string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if input.Quantity < 1 {
		input.Quantity = 1
	}

	var item domain.OrderItem
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO order_items (course_id, dish_id, is_chef_suggestion, chef_suggestion_id, quantity, notes)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, course_id, dish_id, is_chef_suggestion, chef_suggestion_id, quantity, notes, added_at`,
		courseID, input.DishID, input.IsChefSuggestion, input.ChefSuggestionID, input.Quantity, input.Notes,
	).Scan(&item.ID, &item.CourseID, &item.DishID, &item.IsChefSuggestion, &item.ChefSuggestionID, &item.Quantity, &item.Notes, &item.AddedAt)
	if err != nil {
		respondError(w, "could not add item", http.StatusInternalServerError)
		return
	}

	// Get dish name
	if item.DishID != nil {
		h.db.QueryRow(r.Context(), `SELECT name FROM dishes WHERE id = $1`, *item.DishID).Scan(&item.DishName)
	} else if item.ChefSuggestionID != nil {
		h.db.QueryRow(r.Context(), `SELECT name FROM chef_suggestions WHERE id = $1`, *item.ChefSuggestionID).Scan(&item.DishName)
	}

	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "order_item.added", "order_item", &item.ID, input)
	}
	respondJSON(w, item)
}

func (h *OrderHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}
	_, err = h.db.Exec(r.Context(), `DELETE FROM order_items WHERE id = $1`, id)
	if err != nil {
		respondError(w, "not found", http.StatusNotFound)
		return
	}
	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "order_item.deleted", "order_item", &id, nil)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) Send(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		respondError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var currentStatus string
	h.db.QueryRow(r.Context(), `SELECT status FROM orders WHERE id = $1`, orderID).Scan(&currentStatus)
	if currentStatus != "pending" {
		respondError(w, "order already sent", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec(r.Context(), `UPDATE orders SET status = 'sent' WHERE id = $1`, orderID)
	if err != nil {
		respondError(w, "could not send", http.StatusInternalServerError)
		return
	}

	// First course should already be active from creation
	if claims != nil {
		RecordAudit(r.Context(), h.db, claims.UserID, claims.Name, "order.sent", "order", &orderID, nil)
	}
	respondJSON(w, h.loadOrder(r.Context(), orderID))
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

			iRows, _ := h.db.Query(ctx,
				`SELECT oi.id, oi.course_id, oi.dish_id, oi.is_chef_suggestion, oi.chef_suggestion_id, oi.quantity, oi.notes, oi.added_at,
					COALESCE(d.name, cs.name, '') as dish_name
				 FROM order_items oi
				 LEFT JOIN dishes d ON d.id = oi.dish_id
				 LEFT JOIN chef_suggestions cs ON cs.id = oi.chef_suggestion_id
				 WHERE oi.course_id = $1 ORDER BY oi.added_at`, c.ID)
			if iRows != nil {
				for iRows.Next() {
					var item domain.OrderItem
					iRows.Scan(&item.ID, &item.CourseID, &item.DishID, &item.IsChefSuggestion, &item.ChefSuggestionID, &item.Quantity, &item.Notes, &item.AddedAt, &item.DishName)
					c.Items = append(c.Items, item)
				}
				iRows.Close()
			}
			if c.Items == nil {
				c.Items = []domain.OrderItem{}
			}
			o.Courses = append(o.Courses, c)
		}
	}
	if o.Courses == nil {
		o.Courses = []domain.OrderCourse{}
	}
	return o
}
