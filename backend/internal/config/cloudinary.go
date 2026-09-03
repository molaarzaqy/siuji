package config

import (
	"log"
	"siuji-backend/pkg/cloudinary"

	"github.com/spf13/viper"
)

func NewCloudinary(viper *viper.Viper) *cloudinary.Service {
	url := viper.GetString("CLOUDINARY_URL")
	if url == "" {
		log.Fatal("CLOUDINARY_URL is required but not set")
	}
	service, err := cloudinary.NewService(url)
	if err != nil {
		log.Fatalf("failed to initialize Cloudinary: %v", err)
	}
	return service
}