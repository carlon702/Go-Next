package handlers

import (
	"net/http"

	"github.com/carlon702/Go-Next/backend/internal/models"
	"github.com/carlon702/Go-Next/backend/internal/services"
	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService()

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// HealthCheck handles health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Server is running",
		Data: gin.H{
			"status": "healthy",
		},
	})
}

// GetUsers returns all users
func GetUsers(c *gin.Context) {
	users, err := userService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// GetUser returns a single user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := userService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// GetUsersByRole returns users filtered by role
func GetUsersByRole(c *gin.Context) {
	role := models.UserRole(c.Param("role"))

	if role != models.RoleAdmin && role != models.RoleClient {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid role. Use 'admin' or 'client'",
		})
		return
	}

	users, err := userService.GetByRole(role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	user, err := userService.Create(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	req.ID = id // Set ID from URL parameter

	user, err := userService.Update(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser soft deletes a user
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// RestoreUser restores a soft deleted user
func RestoreUser(c *gin.Context) {
	id := c.Param("id")

	if err := userService.Restore(id); err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User restored successfully",
	})
}

// GetUserStats returns user statistics
func GetUserStats(c *gin.Context) {
	totalCount, _ := userService.Count()
	adminCount, _ := userService.CountByRole(models.RoleAdmin)
	clientCount, _ := userService.CountByRole(models.RoleClient)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "User statistics retrieved successfully",
		Data: gin.H{
			"total":   totalCount,
			"admins":  adminCount,
			"clients": clientCount,
		},
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request body",
		})
		return
	}

	user, err := userService.Authenticate(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Login successful",
		Data:    user,
	})
}
