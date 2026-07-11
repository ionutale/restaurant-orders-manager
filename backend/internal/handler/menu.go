package handler

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ionutale/restaurant-orders-manager/internal/domain"
)

type MenuHandler struct {
	db *pgxpool.Pool
}

func NewMenuHandler(db *pgxpool.Pool) *MenuHandler {
	return &MenuHandler{db: db}
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	type categoryGroup struct {
		Category domain.Category `json:"category"`
		Dishes   []domain.Dish   `json:"dishes"`
	}

type menuResponse struct {
	Categories   []categoryGroup          `json:"categories"`
	Suggestions  []domain.ChefSuggestion  `json:"suggestions"`
	Allergens    []domain.Allergen        `json:"allergens"`
	DishAllergens map[int64][]domain.Allergen `json:"dish_allergens"`
}

	// Categories with dishes
	catRows, err := h.db.Query(r.Context(),
		`SELECT id, name, display_order, COALESCE(icon, ''), created_at FROM categories ORDER BY display_order`)
	if err != nil {
		respondError(w, "database error", http.StatusInternalServerError)
		return
	}
	defer catRows.Close()

	var cats []domain.Category
	for catRows.Next() {
		var c domain.Category
		catRows.Scan(&c.ID, &c.Name, &c.DisplayOrder, &c.Icon, &c.CreatedAt)
		cats = append(cats, c)
	}

	var categories []categoryGroup
	for _, c := range cats {
		dRows, err := h.db.Query(r.Context(),
			`SELECT id, name, description, price_cents, category_id, eating_time_min, COALESCE(image_url, ''), created_at
			 FROM dishes WHERE category_id = $1 ORDER BY name`, c.ID)
		if err != nil {
			continue
		}
		var dishes []domain.Dish
		for dRows.Next() {
			var d domain.Dish
			dRows.Scan(&d.ID, &d.Name, &d.Description, &d.PriceCents, &d.CategoryID, &d.EatingTimeMin, &d.ImageURL, &d.CreatedAt)
			dishes = append(dishes, d)
		}
		dRows.Close()
		if dishes == nil {
			dishes = []domain.Dish{}
		}
		categories = append(categories, categoryGroup{Category: c, Dishes: dishes})
	}

	// Active chef suggestions
	sRows, err := h.db.Query(r.Context(),
		`SELECT s.id, s.name, s.description, s.price_cents, s.shift_date, s.expires_at, s.chef_id, COALESCE(u.name, ''), s.created_at
		 FROM chef_suggestions s LEFT JOIN users u ON u.id = s.chef_id
		 WHERE s.expires_at > NOW() ORDER BY s.created_at DESC`)
	var suggestions []domain.ChefSuggestion
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var s domain.ChefSuggestion
			sRows.Scan(&s.ID, &s.Name, &s.Description, &s.PriceCents, &s.ShiftDate, &s.ExpiresAt, &s.ChefID, &s.ChefName, &s.CreatedAt)
			suggestions = append(suggestions, s)
		}
	}
	if suggestions == nil {
		suggestions = []domain.ChefSuggestion{}
	}

	// All allergens
	aRows, err := h.db.Query(r.Context(), `SELECT id, name, COALESCE(icon, '') FROM allergens ORDER BY name`)
	var allergens []domain.Allergen
	if err == nil {
		defer aRows.Close()
		for aRows.Next() {
			var a domain.Allergen
			aRows.Scan(&a.ID, &a.Name, &a.Icon)
			allergens = append(allergens, a)
		}
	}
	if allergens == nil {
		allergens = []domain.Allergen{}
	}

	// Per-dish allergen mapping
	daRows, err := h.db.Query(r.Context(),
		`SELECT da.dish_id, a.id, a.name, COALESCE(a.icon, '')
		 FROM dish_allergens da JOIN allergens a ON a.id = da.allergen_id ORDER BY da.dish_id, a.name`)
	dishAllergens := map[int64][]domain.Allergen{}
	if err == nil {
		defer daRows.Close()
		for daRows.Next() {
			var dishID int64
			var a domain.Allergen
			daRows.Scan(&dishID, &a.ID, &a.Name, &a.Icon)
			dishAllergens[dishID] = append(dishAllergens[dishID], a)
		}
	}

	respondJSON(w, menuResponse{
		Categories:    categories,
		Suggestions:   suggestions,
		Allergens:     allergens,
		DishAllergens: dishAllergens,
	})
}
