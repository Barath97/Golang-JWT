package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/jwt/controller"
	"github.com/jwt/middleware"
)

func UserRoutes(router *gin.Engine) {
	router.Use(middleware.Authenticate())
	router.GET("/users", controllers.GetUsers())
	router.GET("/users/:user_id", controllers.GetUser())
}
