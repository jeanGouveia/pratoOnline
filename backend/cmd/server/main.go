package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/seu-usuario/my-app/backend/internal/handler"
	"github.com/seu-usuario/my-app/backend/internal/infra/database"
	"github.com/seu-usuario/my-app/backend/internal/infra/repository"
	"github.com/seu-usuario/my-app/backend/internal/middleware"
	"github.com/seu-usuario/my-app/backend/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using defaults")
	}

	// Database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate (for development - use Goose migrations in production)
	if os.Getenv("AUTO_MIGRATE") == "true" {
		// Auto-migrate schema
		// db.AutoMigrate(&domain.User{})
	}

	// Initialize repository
	userRepo := repository.NewGormUserRepository(db)

	// Initialize service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	authService := service.NewAuthService(userRepo, jwtSecret)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup router
	r := chi.NewRouter()

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public routes
		authHandler.RegisterRoutes(r)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Auth)
			userHandler.RegisterRoutes(r)
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
