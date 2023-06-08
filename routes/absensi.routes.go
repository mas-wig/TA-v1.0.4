package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type UserAbsensiController struct {
	absensiController controllers.AbsensiController
}

func NewAbsensiContoller(absensiController controllers.AbsensiController) UserAbsensiController {
	return UserAbsensiController{absensiController}
}

func (abc *UserAbsensiController) UserPrensensi(rg *gin.RouterGroup) {
	router := rg.Group("absen")
	router.Use(middleware.DeserializeUser())
	router.POST("/create", abc.absensiController.CreateAbsensi)
	router.GET("/", abc.absensiController.GetInputPresensi)
}
