package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jwt/routes"
)

func main() {
	//Retrieve port from environment variable
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000" //default to 8000 if not specified
	}

	router := gin.New()      //initalize the gin router
	router.Use(gin.Logger()) //logging in middleware

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Access granted for api-1"})
	})

	router.GET("api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
