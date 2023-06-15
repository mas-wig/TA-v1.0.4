package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/mas-wig/ta-v1.0.4/models"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(DB *gorm.DB) AdminController {
	return AdminController{DB}
}

func (ac *AdminController) GetAllNewUser(ctx *gin.Context) {
	var allUser []models.User

	var (
		verifiedCount   int64
		allUserCount    int64
		unverifiedCount int64
	)

	ac.DB.Model(&models.User{}).Where("role=? AND verified=?", "user", true).Count(&verifiedCount)
	ac.DB.Model(&models.User{}).Where("role=? AND verified=?", "user", false).Count(&unverifiedCount)
	ac.DB.Model(&models.User{}).Where("role=?", "user").Count(&allUserCount)

	results := ac.DB.Where("role = ?", "user").Find(&allUser)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "admin-dashboard.html", gin.H{
		"JumlahUser":      allUserCount,
		"JumlahVerified":  verifiedCount,
		"BelumVerifikasi": unverifiedCount,
	})
}
