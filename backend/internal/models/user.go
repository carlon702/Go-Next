package models

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleClient UserRole = "client"
)

type User struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      UserRole       `gorm:"type:text;not null;default:'client';check:role IN ('admin','client')" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook - runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleClient
	}
	return nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsClient checks if user has client role
func (u *User) IsClient() bool {
	return u.Role == RoleClient
}

// ========== REQUEST/RESPONSE STRUCTS ==========

// CreateUserRequest is the request payload for creating a user
type CreateUserRequest struct {
	Name     string   `json:"name" binding:"required,min=2,max=100"`
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	Role     UserRole `json:"role" binding:"omitempty,oneof=admin client"`
}

// UpdateUserRequest is the request payload for updating a user
type UpdateUserRequest struct {
	ID       string   `json:"id"`
	Name     string   `json:"name" binding:"omitempty,min=2,max=100"`
	Email    string   `json:"email" binding:"omitempty,email"`
	Password string   `json:"password" binding:"omitempty,min=6"`
	Role     UserRole `json:"role" binding:"omitempty,oneof=admin client"`
}

// LoginRequest is the request payload for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
