package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardHandler struct {
	db *pgxpool.Pool
}

func NewDashboardHandler(db *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{db: db}
}

type topItem struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Total int    `json:"total_cents"`
}

type statsResp struct {
	PeopleToday   int       `json:"people_today"`
	RevenueToday  int       `json:"revenue_today_cents"`
	TopDishes     []topItem `json:"top_dishes"`
	TopWines      []topItem `json:"top_wines"`
	OrdersToday   int       `json:"orders_today"`
	TablesOccupied int      `json:"tables_occupied"`
}

func (h *DashboardHandler) Stats(w http.ResponseWriter, r *http.Request) {
	var s statsResp

	// People served today
	h.db.QueryRow(r.Context(), `
		SELECT COALESCE(SUM(tg.party_size), 0)
		FROM orders o
		JOIN table_groups tg ON tg.id = o.table_group_id
		WHERE o.status IN ('completed', 'paid') AND o.created_at >= CURRENT_DATE
	`).Scan(&s.PeopleToday)

	// Revenue today
	h.db.QueryRow(r.Context(), `
		SELECT COALESCE(SUM(oi.quantity * COALESCE(d.price_cents, cs.price_cents, 0)), 0)
		FROM order_items oi
		JOIN courses c ON c.id = oi.course_id
		JOIN orders o ON o.id = c.order_id
		LEFT JOIN dishes d ON d.id = oi.dish_id
		LEFT JOIN chef_suggestions cs ON cs.id = oi.chef_suggestion_id
		WHERE o.status IN ('completed', 'paid') AND o.created_at >= CURRENT_DATE
	`).Scan(&s.RevenueToday)

	// Orders today
	h.db.QueryRow(r.Context(), `
		SELECT COUNT(*) FROM orders WHERE status IN ('completed', 'paid') AND created_at >= CURRENT_DATE
	`).Scan(&s.OrdersToday)

	// Tables currently occupied
	h.db.QueryRow(r.Context(), `
		SELECT COUNT(DISTINCT tgt.table_id) FROM table_groups tg
		JOIN table_group_tables tgt ON tgt.group_id = tg.id
		WHERE tg.status != 'closed'
	`).Scan(&s.TablesOccupied)

	// Top dishes (excluding wines)
	rows, err := h.db.Query(r.Context(), `
		SELECT COALESCE(d.name, cs.name, 'Unknown'), COUNT(*), SUM(oi.quantity * COALESCE(d.price_cents, cs.price_cents, 0))
		FROM order_items oi
		JOIN courses c ON c.id = oi.course_id
		JOIN orders o ON o.id = c.order_id
		LEFT JOIN dishes d ON d.id = oi.dish_id AND (d.category_id NOT IN (SELECT id FROM categories WHERE LOWER(name) = 'wines'))
		LEFT JOIN chef_suggestions cs ON cs.id = oi.chef_suggestion_id
		WHERE o.status IN ('completed', 'paid') AND o.created_at >= CURRENT_DATE
		GROUP BY d.name, cs.name
		ORDER BY COUNT(*) DESC
		LIMIT 5
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t topItem
			rows.Scan(&t.Name, &t.Count, &t.Total)
			s.TopDishes = append(s.TopDishes, t)
		}
	}
	if s.TopDishes == nil {
		s.TopDishes = []topItem{}
	}

	// Top wines
	wRows, err := h.db.Query(r.Context(), `
		SELECT d.name, COUNT(*), SUM(oi.quantity * d.price_cents)
		FROM order_items oi
		JOIN courses c ON c.id = oi.course_id
		JOIN orders o ON o.id = c.order_id
		JOIN dishes d ON d.id = oi.dish_id
		JOIN categories cat ON cat.id = d.category_id AND LOWER(cat.name) = 'wines'
		WHERE o.status IN ('completed', 'paid') AND o.created_at >= CURRENT_DATE
		GROUP BY d.name
		ORDER BY COUNT(*) DESC
		LIMIT 5
	`)
	if err == nil {
		defer wRows.Close()
		for wRows.Next() {
			var t topItem
			wRows.Scan(&t.Name, &t.Count, &t.Total)
			s.TopWines = append(s.TopWines, t)
		}
	}
	if s.TopWines == nil {
		s.TopWines = []topItem{}
	}

	respondJSON(w, s)
}
