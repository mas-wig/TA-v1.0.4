package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/initializers"
	"github.com/mas-wig/ta-v1.0.4/routes"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	AbsensiController      controllers.AbsensiController
	AbsensiRouteController routes.UserAbsensiController

	AdminController      controllers.AdminController
	AdminRouteController routes.AdminRouteController

	ProgressController      controllers.ProgressController
	ProgressRouteController routes.ProgressRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal(".env file tidak ditemukan", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	AbsensiController = controllers.NewAbsensiController(initializers.DB)
	AbsensiRouteController = routes.NewAbsensiContoller(AbsensiController)

	AdminController = controllers.NewAdminController(initializers.DB)
	AdminRouteController = routes.NewRouteAdminController(AdminController)

	ProgressController = controllers.NewProgressController(initializers.DB)
	ProgressRouteController = routes.NewRouteProgressController(ProgressController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal(".env file tidak ditemukan", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	server.Static("/src/", "./public/src/")
	server.Static("/vendors/", "./public/vendors/")
	server.Static("/static/", "./public/static/")

	server.LoadHTMLGlob("./public/templates/*.html")

	api := server.Group("/api")
	AuthRouteController.AuthRouters(api)
	UserRouteController.UserRouteProfile(api)
	PostRouteController.PostRoute(api)
	AbsensiRouteController.UserPrensensi(api)

	router := server.Group("/")
	UserRouteController.UserRouteDashboard(router)
	AdminRouteController.AdminRoutes(router)
	ProgressRouteController.ProgressRoutes(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
