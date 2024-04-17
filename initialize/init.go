package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"websocketService/config"
	"websocketService/global"
	"websocketService/middlewares"
	"websocketService/routers"
	"websocketService/utils"
)

func InitConfig() {
	// 实例化viper
	v := viper.New()
	v.SetConfigName("env")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	// 声明一个ServerConfig类型的实例
	serverConfig := config.ServerConfig{}
	//给serverConfig初始值
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	// 传递给全局变量
	global.Settings = serverConfig
}

func InitRouter() *gin.Engine {
	Router := gin.Default()
	// 注册中间件
	Router.Use(middlewares.CORS(), middlewares.GinLogger(), middlewares.GinRecovery(true))

	routers.WsRouter(Router)
	return Router
}

// InitLogger 初始化Logger
func InitLogger() {
	// 实例化zap配置
	cfg := zap.NewDevelopmentConfig()
	// 配置日志的输出地址
	cfg.OutputPaths = []string{
		fmt.Sprintf("%slog_%s.log", global.Settings.LogsAddress, utils.GetNowFormatTodayTime()),
		"stdout", // "stdout" 表示同时将日志输出到标准输出流（控制台）。这样就可以将日志同时输出到文件和控制台
	}
	// 创建logger实例
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	global.Lg = logger         // 注册到全局变量中
}
