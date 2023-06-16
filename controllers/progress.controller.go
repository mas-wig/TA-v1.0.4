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
	ctx.HTML(http.StatusOK, "progress-form.html", gin.H{
		"Action": "/progress/create",
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

	ctx.Redirect(http.StatusOK, "/progress/encode")
}

func (pc *ProgressController) GetAllEncodeProgress(ctx *gin.Context) {
	var allProgress []models.EncodeProgressLatihan
	results := pc.DB.Find(&allProgress)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.HTML(http.StatusOK, "encoded-progress.html", gin.H{
		"Progress": allProgress,
	})
}
