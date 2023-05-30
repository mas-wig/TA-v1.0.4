package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               string    `gorm:"type:varchar(36);default:uuid();primary_key"`
	Name             string    `gorm:"type:varchar(255);not null"`
	Email            string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password         string    `gorm:"type:varchar(255);not null"`
	Role             string    `gorm:"type:varchar(255);not null"`
	Provider         string    `gorm:"type:varchar(255);not null"`
	Photo            string    `gorm:"type:varchar(255);default:'default.png'"`
	VerificationCode string    `gorm:"type:varchar(255)"`
	Verified         bool      `gorm:"not null"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

type SignUpInput struct {
	Name            string                `form:"name" binding:"required"`
	Email           string                `form:"email" binding:"required"`
	Password        string                `form:"password" binding:"required,min=8"`
	PasswordConfirm string                `form:"passwordConfirm" binding:"required"`
	Photo           *multipart.FileHeader `form:"photo"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
