package main

import (
	"fmt"
	"websocketService/global"
	"websocketService/initialize"
	"websocketService/service"
)

func main() {
	// 1.初始化配置
	initialize.InitConfig()
	// 2.初始化路由
	Router := initialize.InitRouter()
	// 3.初始化日志信息
	initialize.InitLogger()

	// 启动一个协程来创建一个空的消息通道
	go service.NewChatRoomThread().Start()

	Router.Run(fmt.Sprintf(":%d", global.Settings.Port))
}
