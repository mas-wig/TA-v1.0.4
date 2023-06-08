package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRouteProfile(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)
}

func (uc *UserRouteController) UserRouteDashboard(rg *gin.RouterGroup) {
	router := rg.Group("users")
	router.GET("/dashboard", middleware.DeserializeUser(), uc.userController.UserDashboard)
}
