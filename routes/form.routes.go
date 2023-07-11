package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func FormUsers(rg *gin.RouterGroup) {
	rg.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/login")
	})

	rg.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{"endpoint": "/api/auth/login"})
	})
	rg.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "register.html", gin.H{})
	})
	rg.GET("/lupasandi", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "forgot-password.html", gin.H{"Action": "/api/auth/forgotpassword"})
	})

}
