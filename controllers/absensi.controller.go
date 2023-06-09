package controllers

import (
	"encoding/base64"
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
		Hari:           base64.StdEncoding.EncodeToString([]byte(payload.Hari)),
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

	ctx.Redirect(http.StatusFound, "/api/absen/encode")

}

func (abc *AbsensiController) GetAllAbsensi(ctx *gin.Context) {
	var allPresensi []models.Presensi
	results := abc.DB.Find(&allPresensi)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "encoded-absensi.html", gin.H{
		"Absensi": allPresensi,
	})
}

func (abc *AbsensiController) DecodeByID(ctx *gin.Context) {
	absenID := ctx.Param("absenId")

	var absensiSiswa models.Presensi
	result := abc.DB.First(&absensiSiswa, "id = ?", absenID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	hari, _ := base64.StdEncoding.DecodeString(absensiSiswa.Hari)

	decodeData := models.Presensi{
		Hari: string(hari),
	}
	abc.DB.Model(&absensiSiswa).Updates(decodeData)

	ctx.Redirect(http.StatusFound, "/api/absen/encode")
}
