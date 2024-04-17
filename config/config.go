package config

type ServerConfig struct {
	Name        string `mapstructure:"appName"`
	Port        int    `mapstructure:"port"`
	LogsAddress string `mapstructure:"logsAddress"`
}
