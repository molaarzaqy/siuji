package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	level := viper.GetString("LOG_LEVEL")
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel // default kalau LOG_LEVEL tidak diset/tidak valid
	}

	log.SetLevel(parsedLevel)
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}