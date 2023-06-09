package models

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	PasswordResetAt    time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
	CreatedAt          time.Time `gorm:"not null"`
	FullName           string    `gorm:"type:varchar(255);not null"`
	ID                 string    `gorm:"type:char(36);primary_key"`
	Gender             string    `gorm:"type:enum('pria','wanita');not null"`
	Address            string    `gorm:"type:text;not null"`
	Photo              string    `gorm:"type:varchar(255);default:'user'"`
	Role               string    `gorm:"type:varchar(255);not null"`
	VerificationCode   string    `gorm:"type:varchar(255)"`
	PasswordResetToken string    `gorm:"not null"`
	Password           string    `gorm:"type:varchar(255);not null"`
	Username           string    `gorm:"type:varchar(255);not null"`
	Email              string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Verified           bool      `gorm:"not null"`
	Acc                bool      `gorm:"not null"`
}

func (note *User) BeforeCreate(tx *gorm.DB) (err error) {
	note.ID = uuid.New().String()
	return nil
}

type SignUpInput struct {
	Email           string                `form:"email" binding:"required"`
	Username        string                `form:"username" binding:"required"`
	Password        string                `form:"password" binding:"required,min=3"`
	PasswordConfirm string                `form:"passwordConfirm" binding:"required"`
	FullName        string                `form:"fullname" binding:"required"`
	Gender          string                `form:"gender" binding:"required"`
	Photo           *multipart.FileHeader `form:"photo" binding:"required"`
	Address         string                `form:"alamat" binding:"required"`
	Acc             bool
}

type SignInInput struct {
	Email    string `form:"email"  binding:"required"`
	Password string `form:"password"  binding:"required"`
	Role     string `form:"role" binding:"required"`
}

type UserResponse struct {
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Role    string `json:"role,omitempty"`
	Photo   string `json:"photo,omitempty"`
	Address string
	ID      uuid.UUID `json:"id,omitempty"`
	Acc     bool
}

type UserResponseProfile struct {
	FullName  string
	PhotoPath string
}

type ForgotPasswordInput struct {
	Email string `form:"email" binding:"required"`
}

type ResetPasswordToken struct {
	Token string `form:"token" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"confirmpassword" binding:"required"`
}
