package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type AdminRouteController struct {
	userController controllers.AdminController
}

func NewRouteAdminController(adminController controllers.AdminController) AdminRouteController {
	return AdminRouteController{adminController}
}

func (uc *AdminRouteController) AdminRoutes(rg *gin.RouterGroup) {
	router := rg.Group("admin")
	router.GET("/dashboard", middleware.DeserializeUser(), uc.userController.GetAllNewUser)
	router.GET("/access", middleware.DeserializeUser(), uc.userController.GetAllUserACC)
	router.GET("/list", middleware.DeserializeUser(), uc.userController.GetAllUserList)
	router.POST("/keyaccess/:updateId", middleware.DeserializeUser(), uc.userController.GetUserDecodeKey)
}
