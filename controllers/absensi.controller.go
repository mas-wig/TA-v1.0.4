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

	newAbsensi := models.EncodePresensi{
		NamaSiswa:      base64.StdEncoding.EncodeToString([]byte(currentUser.FullName)),
		IDSiswa:        userID,
		Hari:           base64.StdEncoding.EncodeToString([]byte(payload.Hari)),
		Lokasi:         base64.StdEncoding.EncodeToString([]byte(payload.Lokasi)),
		TanggalWaktu:   base64.StdEncoding.EncodeToString([]byte(payload.TanggalWaktu)),
		Kehadiran:      base64.StdEncoding.EncodeToString([]byte(payload.Kehadiran)),
		InformasiMedis: base64.StdEncoding.EncodeToString([]byte(payload.InformasiMedis)),
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

func (abc *AbsensiController) GetAllEncodeAbsensi(ctx *gin.Context) {
	var allPresensi []models.EncodePresensi
	results := abc.DB.Find(&allPresensi)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "encoded-absensi.html", gin.H{
		"Absensi": allPresensi,
	})
}

func (abc *AbsensiController) GetAllDecodeAbsensi(ctx *gin.Context) {
	var allPresensi []models.DecodePresensi
	results := abc.DB.Find(&allPresensi)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "decoded-absensi.html", gin.H{
		"Absensi": allPresensi,
	})
}

func (abc *AbsensiController) DecodeByID(ctx *gin.Context) {
	absenID := ctx.Param("absenId")

	var payload *models.DecodeKey
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	var absensiSiswa models.EncodePresensi
	result := abc.DB.First(&absensiSiswa, "id = ?", absenID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	if payload.Key != currentUser.ID {
		ctx.Redirect(http.StatusFound, "/api/absen/encode")
		return
	}

	fullname, _ := base64.StdEncoding.DecodeString(absensiSiswa.NamaSiswa)
	hari, _ := base64.StdEncoding.DecodeString(absensiSiswa.Hari)
	lokasi, _ := base64.StdEncoding.DecodeString(absensiSiswa.Lokasi)
	date, _ := base64.StdEncoding.DecodeString(absensiSiswa.TanggalWaktu)
	kehadiran, _ := base64.StdEncoding.DecodeString(absensiSiswa.Kehadiran)
	catatanMedis, _ := base64.StdEncoding.DecodeString(absensiSiswa.InformasiMedis)

	decodeData := models.DecodePresensi{
		NamaSiswa:      string(fullname),
		IDSiswa:        userID,
		Hari:           string(hari),
		Lokasi:         string(lokasi),
		TanggalWaktu:   string(date),
		Kehadiran:      string(kehadiran),
		InformasiMedis: string(catatanMedis),
	}

	result = abc.DB.Create(&decodeData)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	abc.DB.Delete(&models.EncodePresensi{}, "id = ?", absenID)

	ctx.Redirect(http.StatusFound, "/api/absen/decode")

}

func (abc *AbsensiController) DeleteByID(ctx *gin.Context) {
	deleteID := ctx.Param("deleteId")
	var absensiSiswa models.DecodePresensi
	result := abc.DB.First(&absensiSiswa, "id = ?", deleteID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	result = abc.DB.Delete(&absensiSiswa, "id = ?", deleteID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	ctx.Redirect(http.StatusFound, "/api/absen/decode")
}

func (abc *AbsensiController) UpdatePresensiByID(ctx *gin.Context) {
	updateID := ctx.Param("updateId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.CreatePresensi
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var newUpdatePresensi models.DecodePresensi
	result := abc.DB.First(&newUpdatePresensi, "id = ?", updateID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	updateDecodeData := models.DecodePresensi{
		NamaSiswa:      newUpdatePresensi.NamaSiswa,
		IDSiswa:        userID,
		Hari:           payload.Hari,
		Lokasi:         payload.Lokasi,
		TanggalWaktu:   payload.TanggalWaktu,
		Kehadiran:      payload.Kehadiran,
		InformasiMedis: payload.InformasiMedis,
	}
	abc.DB.Model(&newUpdatePresensi).Updates(updateDecodeData)
	ctx.Redirect(http.StatusFound, "/api/absen/decode")
}
