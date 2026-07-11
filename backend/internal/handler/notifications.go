package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
)

type NotificationHandler struct {
	db *pgxpool.Pool
}

func NewNotificationHandler(db *pgxpool.Pool) *NotificationHandler {
	return &NotificationHandler{db: db}
}

func (h *NotificationHandler) Poll(w http.ResponseWriter, r *http.Request) {
	claims := auth.ClaimsFromCtx(r.Context())
	if claims == nil {
		respondError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := h.db.Query(r.Context(), `
		SELECT oi.id, oi.dish_id, COALESCE(d.name, cs.name, ''), oi.quantity, o.id AS order_id, t.name AS table_name
		FROM order_items oi
		JOIN courses c ON c.id = oi.course_id
		JOIN orders o ON o.id = c.order_id
		JOIN table_groups tg ON tg.id = o.table_group_id
		JOIN table_group_tables tgt ON tgt.group_id = tg.id
		JOIN tables t ON t.id = tgt.table_id
		LEFT JOIN dishes d ON d.id = oi.dish_id
		LEFT JOIN chef_suggestions cs ON cs.id = oi.chef_suggestion_id
		WHERE o.waiter_id = $1 AND oi.ready = true
		  AND NOT EXISTS (SELECT 1 FROM order_items oi2 WHERE oi2.course_id = oi.course_id AND oi2.ready = false)
		  AND c.status = 'active'
		GROUP BY oi.id, d.name, cs.name, o.id, t.name
		ORDER BY oi.ready_at DESC
		LIMIT 20`, claims.UserID)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type notif struct {
		ItemID    int64  `json:"item_id"`
		DishName  string `json:"dish_name"`
		Quantity  int    `json:"quantity"`
		OrderID   int64  `json:"order_id"`
		TableName string `json:"table_name"`
	}
	var notifs []notif
	for rows.Next() {
		var n notif
		if err := rows.Scan(&n.ItemID, &n.DishName, &n.Quantity, &n.OrderID, &n.TableName); err != nil {
			continue
		}
		notifs = append(notifs, n)
	}
	if notifs == nil {
		notifs = []notif{}
	}
	respondJSON(w, notifs)
}
