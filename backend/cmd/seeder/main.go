package main

import (
	"log"
	"siuji-backend/internal/config"
	"siuji-backend/internal/entity"
	"siuji-backend/pkg/password"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func main() {
	viperConfig := config.NewViper()
	appLog := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, appLog)

	seedAdmin(db, viperConfig)
}

func seedAdmin(db *gorm.DB, viperConfig *viper.Viper) {
	adminEmail := viperConfig.GetString("ADMIN_EMAIL")
	adminPassword := viperConfig.GetString("ADMIN_PASSWORD")
	adminRole := viperConfig.GetString("ADMIN_ROLE")

	if adminEmail == "" || adminPassword == "" {
		log.Fatal("ADMIN EMAIL and ADMIN PASSWORD must be set in .env")
	}
	if adminRole == "" {
		adminRole = "admin"
	}

	log.Println("no admin found for", adminEmail, "- seeding default admin user")

	hashed, err := password.Hash(adminPassword)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	admin := entity.User{
		PublicID: uuid.New(),
		Name: "admin siuji",
		Email: adminEmail,
		Password: hashed,
		Role: adminRole,
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("failed to seed admin:", err)
	} else {
		log.Println("admin user seeded:", adminEmail)
	}
}