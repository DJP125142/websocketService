package global

import (
	"go.uber.org/zap"
	"websocketService/config"
)

var (
	Settings config.ServerConfig
	Lg       *zap.Logger
)
