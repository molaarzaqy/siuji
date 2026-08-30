package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// NewViper loads configuration from .env file and environment variables.
// Viper acts as a centralized config accessor over godotenv-loaded env vars.
func NewViper() *viper.Viper {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system environment variables")
	}

	config := viper.New()
	config.AutomaticEnv() // let viper read from OS environment variables

	return config
}