package routes

import (
	"github.com/gin-gonic/gin"
)

func StatusRoutes(apiRouter *gin.RouterGroup) {
	route := apiRouter.Group("/status")
	{
		route.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "success",
			})
		})
	}
}
