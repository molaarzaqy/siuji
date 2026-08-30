package config

import (
	"time"

	"siuji-backend/pkg/jwt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewJWTManager(viper *viper.Viper, log *logrus.Logger) *jwt.Manager {
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	accessDuration, err := time.ParseDuration(viper.GetString("JWT_EXPIRES_IN"))
	if err != nil {
		log.Fatalf("invalid JWT_EXPIRES_IN format: %v", err)
	}

	refreshDuration, err := time.ParseDuration(viper.GetString("REFRESH_TOKEN_EXPIRES"))
	if err != nil {
		log.Fatalf("invalid REFRESH_TOKEN_EXPIRES format: %v", err)
	}

	return jwt.NewManager(secret, accessDuration, refreshDuration)
}