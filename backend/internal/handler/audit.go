package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditEvent struct {
	ID        int64           `json:"id"`
	ActorID   int64           `json:"actor_id"`
	ActorName string          `json:"actor_name"`
	Action    string          `json:"action"`
	Entity    string          `json:"entity"`
	EntityID  *int64          `json:"entity_id"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

func RecordAudit(ctx context.Context, db *pgxpool.Pool, actorID int64, actorName, action, entity string, entityID *int64, payload interface{}) {
	var raw json.RawMessage
	if payload != nil {
		b, err := json.Marshal(payload)
		if err == nil {
			raw = b
		}
	}
	db.Exec(ctx,
		`INSERT INTO audit_events (actor_id, actor_name, action, entity, entity_id, payload) VALUES ($1, $2, $3, $4, $5, $6)`,
		actorID, actorName, action, entity, entityID, raw,
	)
}

type AuditHandler struct {
	db *pgxpool.Pool
}

func NewAuditHandler(db *pgxpool.Pool) *AuditHandler {
	return &AuditHandler{db: db}
}

func (h *AuditHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, actor_id, actor_name, action, entity, entity_id, created_at
		 FROM audit_events ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []AuditEvent
	for rows.Next() {
		var e AuditEvent
		if err := rows.Scan(&e.ID, &e.ActorID, &e.ActorName, &e.Action, &e.Entity, &e.EntityID, &e.CreatedAt); err != nil {
			continue
		}
		events = append(events, e)
	}
	if events == nil {
		events = []AuditEvent{}
	}
	respondJSON(w, events)
}
