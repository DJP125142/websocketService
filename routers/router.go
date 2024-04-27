package routers

import (
	"github.com/gin-gonic/gin"
	"websocketService/controller"
)

func WsRouter(Router *gin.Engine) {
	Router.GET("/", controller.CreateConn)
}
