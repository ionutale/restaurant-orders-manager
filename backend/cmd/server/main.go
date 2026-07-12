package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/ionutale/restaurant-orders-manager/internal/auth"
	"github.com/ionutale/restaurant-orders-manager/internal/database"
	"github.com/ionutale/restaurant-orders-manager/internal/handler"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	dbURL := envOrDefault("DATABASE_URL", "postgres://localhost:5432/restaurant_orders?sslmode=disable")
	jwtSecret := envOrDefault("JWT_SECRET", "dev-secret-change-in-production")
	port := envOrDefault("PORT", "8080")

	db, err := database.Connect(context.Background(), dbURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := database.RunMigrations(dbURL); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	jwt := auth.NewJWT(jwtSecret)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(corsMiddleware)

	r.Get("/health", handler.Health(db))

	r.Post("/api/auth/login", handler.Login(db, jwt))
	r.Post("/api/auth/register", handler.Register(db, jwt))

	r.Route("/api", func(r chi.Router) {
		r.Use(auth.Middleware(jwt))

		r.Get("/me", handler.Me(db))

		th := handler.NewTableHandler(db)
		r.Get("/tables", th.List)
		r.Post("/tables", th.Create)
		r.Patch("/tables/{id}", th.Update)
		r.Delete("/tables/{id}", th.Delete)

		ch := handler.NewCategoryHandler(db)
		r.Get("/categories", ch.List)
		r.Post("/categories", ch.Create)
		r.Patch("/categories/{id}", ch.Update)
		r.Delete("/categories/{id}", ch.Delete)
		r.Post("/categories/reorder", ch.Reorder)

		dh := handler.NewDishHandler(db)
		r.Get("/dishes", dh.List)
		r.Get("/dishes/{id}", dh.Get)
		r.Post("/dishes", dh.Create)
		r.Patch("/dishes/{id}", dh.Update)
		r.Delete("/dishes/{id}", dh.Delete)

		ah := handler.NewAllergenHandler(db)
		r.Get("/allergens", ah.ListAllergens)
		r.Post("/allergens", ah.CreateAllergen)
		r.Delete("/allergens/{id}", ah.DeleteAllergen)
		r.Get("/dishes/{dishId}/allergens", ah.GetDishAllergens)
		r.Put("/dishes/{dishId}/allergens", ah.SetDishAllergens)
		r.Get("/dishes/{dishId}/suggestions", ah.GetDishSuggestions)
		r.Post("/dishes/{dishId}/suggestions", ah.CreateDishSuggestion)
		r.Delete("/dish-suggestions/{id}", ah.DeleteDishSuggestion)

		sh := handler.NewChefSuggestionHandler(db)
		r.Get("/chef-suggestions", sh.List)
		r.Post("/chef-suggestions", sh.Create)
		r.Post("/chef-suggestions/{id}/renew", sh.Renew)
		r.Delete("/chef-suggestions/{id}", sh.Delete)

		mh := handler.NewMenuHandler(db)
		r.Get("/menu", mh.GetMenu)

		fh := handler.NewFloorPlanHandler(db)
		r.Get("/floor-plan", fh.GetFloorPlan)

		gh := handler.NewTableGroupHandler(db)
		r.Get("/table-groups", gh.List)
		r.Post("/table-groups", gh.Create)
		r.Get("/table-groups/{id}", gh.Get)
		r.Patch("/table-groups/{id}/tables", gh.UpdateTables)
		r.Post("/table-groups/{id}/close", gh.Close)

		oh := handler.NewOrderHandler(db)
		r.Post("/orders", oh.Create)
		r.Get("/orders/{id}", oh.Get)
		r.Get("/orders", oh.List)
		r.Post("/orders/{id}/courses/{courseId}/items", oh.AddItem)
		r.Post("/orders/{id}/courses", oh.AddCourse)
		r.Post("/orders/{id}/send", oh.Send)
		r.Post("/orders/{id}/advance-course", oh.AdvanceCourse)
		r.Get("/kds/orders", oh.KDSOrders)
		r.Patch("/kds/order-items/{itemId}/ready", oh.MarkItemReady)
		r.Delete("/order-items/{id}", oh.DeleteItem)

		ph := handler.NewPredictionHandler(db)
		r.Get("/predictions", ph.TablePredictions)

		nh := handler.NewNotificationHandler(db)
		r.Get("/notifications", nh.Poll)

		r.Post("/orders/{id}/pay", handler.PayOrder(db))
		r.Post("/orders/{id}/send-invoice", handler.SendInvoice(db))

		auh := handler.NewAuditHandler(db)
		r.Get("/audit-events", auh.List)

		uh := handler.NewUserHandler(db)
		r.Get("/users", uh.List)
		r.Post("/users", uh.Create)
		r.Patch("/users/{id}", uh.Update)
		r.Delete("/users/{id}", uh.Delete)

		uploadHandler := handler.NewUploadHandler(db, envOrDefault("UPLOAD_DIR", "./uploads"))
		r.Post("/upload", uploadHandler.Upload)
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdown
	slog.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
		os.Exit(1)
	}
	slog.Info("server stopped")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
