package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRouters(rg *gin.RouterGroup) {
	router := rg.Group("auth")

	router.POST("/register", rc.authController.SignUpUser)
	router.POST("/login", rc.authController.SignInUser)
	router.GET("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)
	router.GET("/verifyemail/:verificationCode", rc.authController.VerifyEmail)
	router.POST("/forgotpassword", rc.authController.ForgotPassword)
	router.POST("/resetpassword/:resetToken", rc.authController.ResetPassword)

	router.GET("/signin", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{"endpoint": "/api/auth/login"})
	})

	router.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", gin.H{})
	})

	router.GET("/lupasandi", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "forgot-password.html", gin.H{"Action": "/api/auth/forgotpassword"})
	})

}
