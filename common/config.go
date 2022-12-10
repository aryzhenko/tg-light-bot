package common

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	BotToken string     `mapstructure:"LIGHT_BOT_TOKEN"`
	Db       DbConfig   `mapstructure:",squash"`
	Ping     PingConfig `mapstructure:",squash"`
}

type DbConfig struct {
	Host     string `mapstructure:"LIGHT_BOT_DB_HOST"`
	User     string `mapstructure:"LIGHT_BOT_DB_USER"`
	Password string `mapstructure:"LIGHT_BOT_DB_PASSWORD"`
	DbName   string `mapstructure:"LIGHT_BOT_DB_NAME"`
}

type PingConfig struct {
	Port      int32         `mapstructure:"LIGHT_BOT_PING_PORT"`
	Threshold time.Duration `mapstructure:"LIGHT_BOT_PING_THRESHOLD"`
}

func LoadConfig(file string) (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(file)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
