package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carlon702/Go-Next/backend/internal/config"
	"github.com/carlon702/Go-Next/backend/internal/database"
	"github.com/carlon702/Go-Next/backend/internal/models"
	"github.com/carlon702/Go-Next/backend/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()
	log.Printf("🚀 Starting server in %s mode", cfg.Environment)

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Test database connection
	if err := database.Ping(); err != nil {
		log.Fatal("❌ Failed to ping database:", err)
	}

	// Run migrations
	log.Println("🔄 Running database migrations...")
	if err := database.Migrate(&models.User{}); err != nil {
		log.Fatal("❌ Failed to migrate database:", err)
	}
	log.Println("✅ Database migrations completed")

	// Setup router
	router := routes.SetupRouter()

	// Configure server
	port := cfg.Port
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("✅ Server running on http://localhost:%s", port)
		log.Println("📚 Available endpoints:")
		log.Println("   GET    /health")
		log.Println("   POST   /api/auth/login")
		log.Println("   POST   /api/auth/register")
		log.Println("   GET    /api/users")
		log.Println("   GET    /api/users/stats")
		log.Println("   GET    /api/users/role/:role")
		log.Println("   GET    /api/users/:id")
		log.Println("   PUT    /api/users/:id")
		log.Println("   DELETE /api/users/:id")
		log.Println("   POST   /api/users/:id/restore")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown:", err)
	}

	// Close database connection
	database.Close()
	log.Println("✅ Server stopped gracefully")
}
