package services

import (
	"errors"
	"fmt"

	"github.com/carlon702/Go-Next/backend/internal/database"
	"github.com/carlon702/Go-Next/backend/internal/models"
	"github.com/carlon702/Go-Next/backend/pkg/utils"
	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// GetAll returns all users (excluding soft deleted)
func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	result := database.DB.Find(&users)
	return users, result.Error
}

// GetByID returns a user by ID
func (s *UserService) GetByID(id string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("id = ?", id).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByEmail returns a user by email
func (s *UserService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// GetByRole returns all users with a specific role
func (s *UserService) GetByRole(role models.UserRole) ([]models.User, error) {
	var users []models.User
	result := database.DB.Where("role = ?", role).Find(&users)
	return users, result.Error
}

// Create creates a new user with hashed password
func (s *UserService) Create(req *models.CreateUserRequest) (*models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = models.RoleClient
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	result := database.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// Update updates a user
func (s *UserService) Update(req *models.UpdateUserRequest) (*models.User, error) {
	// Get existing user
	user, err := s.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists
		var existingUser models.User
		if err := database.DB.Where("email = ? AND id != ?", req.Email, req.ID).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	result := database.DB.Save(user)
	return user, result.Error
}

// Delete soft deletes a user
func (s *UserService) Delete(id string) error {
	result := database.DB.Delete(&models.User{}, "id = ?", id)

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return result.Error
}

// Restore restores a soft deleted user
func (s *UserService) Restore(id string) error {
	result := database.DB.Model(&models.User{}).Unscoped().Where("id = ?", id).Update("deleted_at", nil)

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return result.Error
}

// Authenticate checks user credentials and returns user if valid
func (s *UserService) Authenticate(req *models.LoginRequest) (*models.User, error) {
	// Get user by email
	user, err := s.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// Count returns total number of users
func (s *UserService) Count() (int64, error) {
	var count int64
	result := database.DB.Model(&models.User{}).Count(&count)
	return count, result.Error
}

// CountByRole returns number of users by role
func (s *UserService) CountByRole(role models.UserRole) (int64, error) {
	var count int64
	result := database.DB.Model(&models.User{}).Where("role = ?", role).Count(&count)
	return count, result.Error
}
