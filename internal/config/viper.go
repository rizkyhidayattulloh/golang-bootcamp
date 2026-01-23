package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	AppPort int `mapstructure:"APP_PORT"`
}

func NewViper() *Config {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.SetDefault("APP_PORT", 8080)

	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return &cfg
}
