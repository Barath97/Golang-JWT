package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jwt/controller"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
