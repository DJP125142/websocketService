package routers

import (
	"github.com/gin-gonic/gin"
	"websocketService/controller"
)

func WsRouter(Router *gin.Engine) {
	Router.GET("/", controller.CreateConn)
	Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}
