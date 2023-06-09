package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mas-wig/ta-v1.0.4/controllers"
	"github.com/mas-wig/ta-v1.0.4/middleware"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewRoutePostController(postController controllers.PostController) PostRouteController {
	return PostRouteController{postController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {
	router := rg.Group("posts")
	router.Use(middleware.DeserializeUser())
	router.POST("/create", pc.postController.CreatePost)
	router.GET("/", pc.postController.GetAllPosts)
	router.PUT("/update/:postId", pc.postController.UpdatePost)
	router.GET("/find/:postId", pc.postController.FindPostByID)
	router.POST("/delete/:postId", pc.postController.DeletePost)
}
