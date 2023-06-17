package controllers

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/mas-wig/ta-v1.0.4/models"
	"github.com/mas-wig/ta-v1.0.4/utils"
)

type ProgressController struct {
	DB *gorm.DB
}

func NewProgressController(DB *gorm.DB) ProgressController {
	return ProgressController{DB}
}

func (pc *ProgressController) FormProgress(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	userResponse := &models.UserResponseProfile{
		FullName:  currentUser.FullName,
		PhotoPath: currentUser.Photo,
	}

	ctx.HTML(http.StatusOK, "progress-form.html", gin.H{
		"Action":    "/progress/create",
		"Fullname":  userResponse.FullName,
		"PhotoPath": userResponse.PhotoPath,
	})
}

func (pc *ProgressController) CreateNewProgress(ctx *gin.Context) {
	var payload *models.CreateProgress
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	var mediaURL string
	if payload.Media != nil {
		url, err := utils.SaveUploadedFile(payload.Media)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error 2", "message": err.Error()})
			return
		}
		mediaURL = url
	}

	now := time.Now()

	newProgress := models.EncodeProgressLatihan{
		SiswaID:             userID,
		TingkatKemajuan:     base64.StdEncoding.EncodeToString([]byte(payload.TingkatKemajuan)),
		CatatanPerkembangan: base64.StdEncoding.EncodeToString([]byte(payload.CatatanPerkembangan)),
		PenilaianTeknik:     base64.StdEncoding.EncodeToString([]byte(payload.PenilaianTeknik)),
		PenilaianDisiplin:   base64.StdEncoding.EncodeToString([]byte(payload.PenilaianDisiplin)),
		CatatanEvaluasi:     base64.StdEncoding.EncodeToString([]byte(payload.CatatanEvaluasi)),
		Media:               mediaURL,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	result := pc.DB.Create(&newProgress)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.Redirect(http.StatusFound, "/progress/encode")
}

func (pc *ProgressController) GetAllEncodeProgress(ctx *gin.Context) {
	var allProgress []models.EncodeProgressLatihan
	currentUser := ctx.MustGet("currentUser").(models.User)
	results := pc.DB.Find(&allProgress)
	userResponse := &models.UserResponseProfile{
		FullName:  currentUser.FullName,
		PhotoPath: currentUser.Photo,
	}

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "encoded-progress.html", gin.H{
		"Fullname":  userResponse.FullName,
		"Progress":  allProgress,
		"PhotoPath": userResponse.PhotoPath,
	})
}

func (pc *ProgressController) DecodeProgressByID(ctx *gin.Context) {
	progressID := ctx.Param("progressid")

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

	var progressSiswa models.EncodeProgressLatihan
	result := pc.DB.First(&progressSiswa, "id = ?", progressID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	if payload.Key != currentUser.ID {
		ctx.Redirect(http.StatusFound, "/progress/encode")
		return
	}

	tingkatKemajuan, _ := base64.StdEncoding.DecodeString(progressSiswa.TingkatKemajuan)
	catatanPerkembangan, _ := base64.StdEncoding.DecodeString(progressSiswa.CatatanPerkembangan)
	penilaianTeknik, _ := base64.StdEncoding.DecodeString(progressSiswa.PenilaianTeknik)
	penilaianDisiplin, _ := base64.StdEncoding.DecodeString(progressSiswa.PenilaianDisiplin)
	catatanEvaluasi, _ := base64.StdEncoding.DecodeString(progressSiswa.CatatanEvaluasi)

	decodeData := models.DecodeProgressLatihan{
		SiswaID:             userID,
		TingkatKemajuan:     string(tingkatKemajuan),
		CatatanPerkembangan: string(catatanPerkembangan),
		PenilaianTeknik:     string(penilaianTeknik),
		PenilaianDisiplin:   string(penilaianDisiplin),
		CatatanEvaluasi:     string(catatanEvaluasi),
		Media:               progressSiswa.Media,
	}

	result = pc.DB.Create(&decodeData)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	pc.DB.Delete(&models.EncodePresensi{}, "id = ?", progressID)
	ctx.Redirect(http.StatusFound, "/progress/decode")

}
func (pc *ProgressController) GetAllDecodeProgressData(ctx *gin.Context) {
	var allProgress []models.DecodeProgressLatihan
	results := pc.DB.Find(&allProgress)

	currentUser := ctx.MustGet("currentUser").(models.User)
	userResponse := &models.UserResponseProfile{
		FullName:  currentUser.FullName,
		PhotoPath: currentUser.Photo,
	}

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "decoded-progress.html", gin.H{
		"Progress":  allProgress,
		"Fullname":  userResponse.FullName,
		"PhotoPath": userResponse.PhotoPath,
	})
}

func (pc *ProgressController) DeleteProgressByID(ctx *gin.Context) {
	deleteID := ctx.Param("deleteid")
	var progressSiswa models.DecodeProgressLatihan
	result := pc.DB.First(&progressSiswa, "id = ?", deleteID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	result = pc.DB.Delete(&progressSiswa, "id = ?", deleteID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	ctx.Redirect(http.StatusFound, "/progress/decode")
}

func (pc *ProgressController) UpdateProgressByID(ctx *gin.Context) {
	updateID := ctx.Param("updateid")
	var payload *models.CreateProgress
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID, err := uuid.Parse(currentUser.ID)
	if err != nil {
		panic("UUID kosong ")
	}

	var newUpdateProgress models.DecodeProgressLatihan
	result := pc.DB.First(&newUpdateProgress, "id = ?", updateID)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	var mediaURL string
	if payload.Media != nil {
		url, err := utils.SaveUploadedFile(payload.Media)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error 2", "message": err.Error()})
			return
		}
		mediaURL = url
	}

	now := time.Now()
	updateDecodeData := models.DecodeProgressLatihan{
		SiswaID:             userID,
		TingkatKemajuan:     payload.TingkatKemajuan,
		CatatanPerkembangan: payload.CatatanPerkembangan,
		PenilaianTeknik:     payload.PenilaianTeknik,
		PenilaianDisiplin:   payload.PenilaianDisiplin,
		CatatanEvaluasi:     payload.CatatanEvaluasi,
		Media:               mediaURL,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	pc.DB.Model(&newUpdateProgress).Updates(updateDecodeData)
	ctx.Redirect(http.StatusFound, "/progress/decode")
}
