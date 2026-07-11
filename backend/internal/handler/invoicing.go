package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
)

func PayOrder(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.ClaimsFromCtx(r.Context())
		orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			respondError(w, "invalid id", http.StatusBadRequest)
			return
		}

		var input struct {
			PaymentMethod string `json:"payment_method"`
		}
		json.NewDecoder(r.Body).Decode(&input)
		if input.PaymentMethod == "" {
			input.PaymentMethod = "cash"
		}

		_, err = db.Exec(r.Context(),
			`UPDATE orders SET status = 'paid' WHERE id = $1 AND status = 'completed'`, orderID)
		if err != nil {
			respondError(w, "could not pay", http.StatusBadRequest)
			return
		}

		if claims != nil {
			RecordAudit(r.Context(), db, claims.UserID, claims.Name, "order.paid", "order", &orderID, input)
		}
		respondJSON(w, map[string]string{"status": "paid"})
	}
}

func SendInvoice(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.ClaimsFromCtx(r.Context())
		orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			respondError(w, "invalid id", http.StatusBadRequest)
			return
		}

		var input struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Email == "" {
			respondError(w, "email is required", http.StatusBadRequest)
			return
		}

		// Build invoice data
		var o struct {
			ID int64
			TableGroupID int64
			CreatedAt time.Time
		}
		db.QueryRow(r.Context(), `SELECT id, table_group_id, created_at FROM orders WHERE id = $1`, orderID).
			Scan(&o.ID, &o.TableGroupID, &o.CreatedAt)

		type line struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
			Price    int    `json:"price_cents"`
		}
		var items []line
		rows, _ := db.Query(r.Context(), `
			SELECT COALESCE(d.name, cs.name, 'Item'), oi.quantity, COALESCE(d.price_cents, cs.price_cents, 0)
			FROM order_items oi
			JOIN courses c ON c.id = oi.course_id
			LEFT JOIN dishes d ON d.id = oi.dish_id
			LEFT JOIN chef_suggestions cs ON cs.id = oi.chef_suggestion_id
			WHERE c.order_id = $1`, orderID)
		if rows != nil {
			defer rows.Close()
			for rows.Next() {
				var l line
				rows.Scan(&l.Name, &l.Quantity, &l.Price)
				items = append(items, l)
			}
		}

		total := 0
		for _, l := range items {
			total += l.Price * l.Quantity
		}

		invoiceHTML := fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Invoice #%d</title>
<style>body{font-family:sans-serif;max-width:600px;margin:40px auto;padding:20px}
h1{color:#333}table{width:100%%;border-collapse:collapse}th,td{padding:8px;text-align:left;border-bottom:1px solid #ddd}
.total{font-weight:bold;font-size:1.2em;text-align:right;margin-top:16px}
.footer{margin-top:32px;color:#666;font-size:0.9em}</style></head>
<body>
<h1>Invoice #%d</h1>
<p>Date: %s</p>
<p>Order ID: #%d</p>
<table><thead><tr><th>Item</th><th>Qty</th><th>Price</th></tr></thead><tbody>
%s
</tbody></table>
<div class="total">Total: €%s</div>
<div class="footer">Thank you for dining with us!</div>
</body></html>`,
			orderID, orderID, o.CreatedAt.Format("2006-01-02 15:04"), orderID,
			func() string {
				var b strings.Builder
				for _, l := range items {
					fmt.Fprintf(&b, "<tr><td>%s</td><td>%d</td><td>€%.2f</td></tr>", l.Name, l.Quantity, float64(l.Price)/100)
				}
				return b.String()
			}(),
			fmt.Sprintf("%.2f", float64(total)/100))

		// Log the invoice send
		db.Exec(r.Context(),
			`INSERT INTO audit_events (actor_id, actor_name, action, entity, entity_id, payload)
			 VALUES ($1, $2, 'invoice.sent', 'order', $3, $4)`,
			claims.UserID, claims.Name, &orderID,
			json.RawMessage(fmt.Sprintf(`{"email":"%s","total_cents":%d}`, input.Email, total)))

		respondJSON(w, map[string]interface{}{
			"status":       "sent",
			"email":        input.Email,
			"total":        fmt.Sprintf("%.2f", float64(total)/100),
			"invoice_html": invoiceHTML,
		})
	}
}
