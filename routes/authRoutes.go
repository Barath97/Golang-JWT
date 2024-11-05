package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/jwt/controller"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("users/signup", controllers.Signup())
	router.POST("users/login", controllers.Login())
}
