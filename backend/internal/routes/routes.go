package routes

import (
	"github.com/carlon702/Go-Next/backend/internal/handlers"
	"github.com/carlon702/Go-Next/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes with Gin
func SetupRouter() *gin.Engine {
	r := gin.Default() // Includes Logger and Recovery middleware

	// Apply CORS middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// API group
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.CreateUser)
		}

		// User routes
		users := api.Group("/users")
		{
			users.GET("", handlers.GetUsers)
			users.GET("/stats", handlers.GetUserStats)
			users.GET("/role/:role", handlers.GetUsersByRole)
			users.GET("/:id", handlers.GetUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.PATCH("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
			users.POST("/:id/restore", handlers.RestoreUser)
		}
	}

	return r
}
