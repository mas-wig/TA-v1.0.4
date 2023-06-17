package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type ProgressRouteController struct {
	progressController controllers.ProgressController
}

func NewRouteProgressController(progressController controllers.ProgressController) ProgressRouteController {
	return ProgressRouteController{progressController}
}

func (uc *ProgressRouteController) ProgressRoutes(rg *gin.RouterGroup) {
	router := rg.Group("progress")
	router.POST("/create", middleware.DeserializeUser(), uc.progressController.CreateNewProgress)
	router.GET("/form", middleware.DeserializeUser(), uc.progressController.FormProgress)
	router.GET("/encode", middleware.DeserializeUser(), uc.progressController.GetAllEncodeProgress)
	router.POST("/decode/:progressid", middleware.DeserializeUser(), uc.progressController.DecodeProgressByID)
	router.POST("/update/:updateid", middleware.DeserializeUser(), uc.progressController.UpdateProgressByID)
	router.POST("/delete/:deleteid", middleware.DeserializeUser(), uc.progressController.DeleteProgressByID)
	router.GET("/decode", middleware.DeserializeUser(), uc.progressController.GetAllDecodeProgressData)
}
