package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mas-wig/ta-v1.0.4/models"
)

type AbsensiController struct {
	DB *gorm.DB
}

func NewAbsensiController(DB *gorm.DB) AbsensiController {
	return AbsensiController{DB}
}

func (abc *AbsensiController) GetInputPresensi(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userResponse := &models.UserResponseProfile{
		FullName:  currentUser.FullName,
		PhotoPath: currentUser.Photo,
	}
	ctx.HTML(http.StatusOK, "absensi-form.html", gin.H{
		"Action":    "/api/absen/create",
		"Fullname":  userResponse.FullName,
		"PhotoPath": userResponse.PhotoPath,
	})

}

func (abc *AbsensiController) CreateAbsensi(ctx *gin.Context) {
	var payload *models.CreatePresensi
	currentUser := ctx.MustGet("currentUser").(models.User)
	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newAbsensi := models.Presensi{
		NamaSiswa:      currentUser.FullName,
		IDSiswa:        userID,
		Hari:           payload.Hari,
		Lokasi:         payload.Lokasi,
		TanggalWaktu:   payload.TanggalWaktu,
		Kehadiran:      payload.Kehadiran,
		InformasiMedis: payload.InformasiMedis,
	}

	result := abc.DB.Create(&newAbsensi)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newAbsensi})

}
