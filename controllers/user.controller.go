package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mas-wig/ta-v1.0.4/models"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	userResponse := &models.UserResponse{
		ID:      userID,
		Name:    currentUser.FullName,
		Email:   currentUser.Email,
		Photo:   currentUser.Photo,
		Role:    currentUser.Role,
		Address: currentUser.Address,
	}
	ctx.HTML(http.StatusOK, "profile.html", gin.H{"Profile": userResponse})
}

func (uc *UserController) UserDashboard(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userResponse := &models.UserResponseProfile{
		FullName:  currentUser.FullName,
		PhotoPath: currentUser.Photo,
	}
	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"Fullname":  userResponse.FullName,
		"PhotoPath": userResponse.PhotoPath,
	})
}
