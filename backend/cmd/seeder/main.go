package main

import (
	"log"
	"siuji-backend/internal/config"
	"siuji-backend/internal/entity"
	"siuji-backend/pkg/password"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	viperConfig := config.NewViper()
	appLog := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, appLog)

	seedAdmin(db, viperConfig)
	seedScoreConversions(db)
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

	log.Println("seeding admin user if not exists for", adminEmail)

	hashed, err := password.Hash(adminPassword)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	admin := entity.User{
		PublicID: uuid.New(),
		Name:     "admin siuji",
		Email:    adminEmail,
		Password: hashed,
		Role:     adminRole,
	}

	// Cek apakah admin sudah ada agar tidak duplikasi
	var existing entity.User
	if err := db.Where("email = ?", adminEmail).First(&existing).Error; err == nil {
		log.Println("admin user already exists, skipping:", adminEmail)
		return
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Println("failed to seed admin:", err)
	} else {
		log.Println("admin user seeded successfully:", adminEmail)
	}
}

func seedScoreConversions(db *gorm.DB) {
	log.Println("seeding TOEFL score conversions...")

	conversions := []entity.ScoreConversion{
		// Section 1: Listening Comprehension (0 - 50)
		{SectionType: "listening", CorrectCount: 0, ScaledScore: 24},
		{SectionType: "listening", CorrectCount: 1, ScaledScore: 25},
		{SectionType: "listening", CorrectCount: 2, ScaledScore: 26},
		{SectionType: "listening", CorrectCount: 3, ScaledScore: 27},
		{SectionType: "listening", CorrectCount: 4, ScaledScore: 28},
		{SectionType: "listening", CorrectCount: 5, ScaledScore: 29},
		{SectionType: "listening", CorrectCount: 6, ScaledScore: 30},
		{SectionType: "listening", CorrectCount: 7, ScaledScore: 31},
		{SectionType: "listening", CorrectCount: 8, ScaledScore: 32},
		{SectionType: "listening", CorrectCount: 9, ScaledScore: 32},
		{SectionType: "listening", CorrectCount: 10, ScaledScore: 33},
		{SectionType: "listening", CorrectCount: 11, ScaledScore: 35},
		{SectionType: "listening", CorrectCount: 12, ScaledScore: 37},
		{SectionType: "listening", CorrectCount: 13, ScaledScore: 37},
		{SectionType: "listening", CorrectCount: 14, ScaledScore: 38},
		{SectionType: "listening", CorrectCount: 15, ScaledScore: 41},
		{SectionType: "listening", CorrectCount: 16, ScaledScore: 41},
		{SectionType: "listening", CorrectCount: 17, ScaledScore: 42},
		{SectionType: "listening", CorrectCount: 18, ScaledScore: 43},
		{SectionType: "listening", CorrectCount: 19, ScaledScore: 44},
		{SectionType: "listening", CorrectCount: 20, ScaledScore: 45},
		{SectionType: "listening", CorrectCount: 21, ScaledScore: 45},
		{SectionType: "listening", CorrectCount: 22, ScaledScore: 46},
		{SectionType: "listening", CorrectCount: 23, ScaledScore: 47},
		{SectionType: "listening", CorrectCount: 24, ScaledScore: 47},
		{SectionType: "listening", CorrectCount: 25, ScaledScore: 48},
		{SectionType: "listening", CorrectCount: 26, ScaledScore: 48},
		{SectionType: "listening", CorrectCount: 27, ScaledScore: 49},
		{SectionType: "listening", CorrectCount: 28, ScaledScore: 49},
		{SectionType: "listening", CorrectCount: 29, ScaledScore: 50},
		{SectionType: "listening", CorrectCount: 30, ScaledScore: 51},
		{SectionType: "listening", CorrectCount: 31, ScaledScore: 52},
		{SectionType: "listening", CorrectCount: 32, ScaledScore: 52},
		{SectionType: "listening", CorrectCount: 33, ScaledScore: 53},
		{SectionType: "listening", CorrectCount: 34, ScaledScore: 53},
		{SectionType: "listening", CorrectCount: 35, ScaledScore: 54},
		{SectionType: "listening", CorrectCount: 36, ScaledScore: 54},
		{SectionType: "listening", CorrectCount: 37, ScaledScore: 55},
		{SectionType: "listening", CorrectCount: 38, ScaledScore: 56},
		{SectionType: "listening", CorrectCount: 39, ScaledScore: 57},
		{SectionType: "listening", CorrectCount: 40, ScaledScore: 57},
		{SectionType: "listening", CorrectCount: 41, ScaledScore: 58},
		{SectionType: "listening", CorrectCount: 42, ScaledScore: 59},
		{SectionType: "listening", CorrectCount: 43, ScaledScore: 60},
		{SectionType: "listening", CorrectCount: 44, ScaledScore: 61},
		{SectionType: "listening", CorrectCount: 45, ScaledScore: 62},
		{SectionType: "listening", CorrectCount: 46, ScaledScore: 63},
		{SectionType: "listening", CorrectCount: 47, ScaledScore: 65},
		{SectionType: "listening", CorrectCount: 48, ScaledScore: 66},
		{SectionType: "listening", CorrectCount: 49, ScaledScore: 67},
		{SectionType: "listening", CorrectCount: 50, ScaledScore: 68},

		// Section 2: Structure and Written Expression (0 - 40)
		{SectionType: "structure", CorrectCount: 0, ScaledScore: 20},
		{SectionType: "structure", CorrectCount: 1, ScaledScore: 20},
		{SectionType: "structure", CorrectCount: 2, ScaledScore: 21},
		{SectionType: "structure", CorrectCount: 3, ScaledScore: 22},
		{SectionType: "structure", CorrectCount: 4, ScaledScore: 23},
		{SectionType: "structure", CorrectCount: 5, ScaledScore: 25},
		{SectionType: "structure", CorrectCount: 6, ScaledScore: 26},
		{SectionType: "structure", CorrectCount: 7, ScaledScore: 27},
		{SectionType: "structure", CorrectCount: 8, ScaledScore: 29},
		{SectionType: "structure", CorrectCount: 9, ScaledScore: 31},
		{SectionType: "structure", CorrectCount: 10, ScaledScore: 33},
		{SectionType: "structure", CorrectCount: 11, ScaledScore: 35},
		{SectionType: "structure", CorrectCount: 12, ScaledScore: 36},
		{SectionType: "structure", CorrectCount: 13, ScaledScore: 37},
		{SectionType: "structure", CorrectCount: 14, ScaledScore: 38},
		{SectionType: "structure", CorrectCount: 15, ScaledScore: 40},
		{SectionType: "structure", CorrectCount: 16, ScaledScore: 40},
		{SectionType: "structure", CorrectCount: 17, ScaledScore: 41},
		{SectionType: "structure", CorrectCount: 18, ScaledScore: 42},
		{SectionType: "structure", CorrectCount: 19, ScaledScore: 43},
		{SectionType: "structure", CorrectCount: 20, ScaledScore: 44},
		{SectionType: "structure", CorrectCount: 21, ScaledScore: 45},
		{SectionType: "structure", CorrectCount: 22, ScaledScore: 46},
		{SectionType: "structure", CorrectCount: 23, ScaledScore: 47},
		{SectionType: "structure", CorrectCount: 24, ScaledScore: 48},
		{SectionType: "structure", CorrectCount: 25, ScaledScore: 49},
		{SectionType: "structure", CorrectCount: 26, ScaledScore: 50},
		{SectionType: "structure", CorrectCount: 27, ScaledScore: 51},
		{SectionType: "structure", CorrectCount: 28, ScaledScore: 52},
		{SectionType: "structure", CorrectCount: 29, ScaledScore: 53},
		{SectionType: "structure", CorrectCount: 30, ScaledScore: 54},
		{SectionType: "structure", CorrectCount: 31, ScaledScore: 55},
		{SectionType: "structure", CorrectCount: 32, ScaledScore: 56},
		{SectionType: "structure", CorrectCount: 33, ScaledScore: 57},
		{SectionType: "structure", CorrectCount: 34, ScaledScore: 58},
		{SectionType: "structure", CorrectCount: 35, ScaledScore: 60},
		{SectionType: "structure", CorrectCount: 36, ScaledScore: 61},
		{SectionType: "structure", CorrectCount: 37, ScaledScore: 63},
		{SectionType: "structure", CorrectCount: 38, ScaledScore: 65},
		{SectionType: "structure", CorrectCount: 39, ScaledScore: 67},
		{SectionType: "structure", CorrectCount: 40, ScaledScore: 68},

		// Section 3: Reading Comprehension (0 - 50)
		{SectionType: "reading", CorrectCount: 0, ScaledScore: 21},
		{SectionType: "reading", CorrectCount: 1, ScaledScore: 22},
		{SectionType: "reading", CorrectCount: 2, ScaledScore: 23},
		{SectionType: "reading", CorrectCount: 3, ScaledScore: 23},
		{SectionType: "reading", CorrectCount: 4, ScaledScore: 24},
		{SectionType: "reading", CorrectCount: 5, ScaledScore: 25},
		{SectionType: "reading", CorrectCount: 6, ScaledScore: 26},
		{SectionType: "reading", CorrectCount: 7, ScaledScore: 27},
		{SectionType: "reading", CorrectCount: 8, ScaledScore: 28},
		{SectionType: "reading", CorrectCount: 9, ScaledScore: 28},
		{SectionType: "reading", CorrectCount: 10, ScaledScore: 29},
		{SectionType: "reading", CorrectCount: 11, ScaledScore: 30},
		{SectionType: "reading", CorrectCount: 12, ScaledScore: 31},
		{SectionType: "reading", CorrectCount: 13, ScaledScore: 32},
		{SectionType: "reading", CorrectCount: 14, ScaledScore: 34},
		{SectionType: "reading", CorrectCount: 15, ScaledScore: 35},
		{SectionType: "reading", CorrectCount: 16, ScaledScore: 36},
		{SectionType: "reading", CorrectCount: 17, ScaledScore: 37},
		{SectionType: "reading", CorrectCount: 18, ScaledScore: 38},
		{SectionType: "reading", CorrectCount: 19, ScaledScore: 39},
		{SectionType: "reading", CorrectCount: 20, ScaledScore: 40},
		{SectionType: "reading", CorrectCount: 21, ScaledScore: 41},
		{SectionType: "reading", CorrectCount: 22, ScaledScore: 42},
		{SectionType: "reading", CorrectCount: 23, ScaledScore: 43},
		{SectionType: "reading", CorrectCount: 24, ScaledScore: 43},
		{SectionType: "reading", CorrectCount: 25, ScaledScore: 44},
		{SectionType: "reading", CorrectCount: 26, ScaledScore: 45},
		{SectionType: "reading", CorrectCount: 27, ScaledScore: 46},
		{SectionType: "reading", CorrectCount: 28, ScaledScore: 46},
		{SectionType: "reading", CorrectCount: 29, ScaledScore: 47},
		{SectionType: "reading", CorrectCount: 30, ScaledScore: 48},
		{SectionType: "reading", CorrectCount: 31, ScaledScore: 48},
		{SectionType: "reading", CorrectCount: 32, ScaledScore: 49},
		{SectionType: "reading", CorrectCount: 33, ScaledScore: 50},
		{SectionType: "reading", CorrectCount: 34, ScaledScore: 51},
		{SectionType: "reading", CorrectCount: 35, ScaledScore: 52},
		{SectionType: "reading", CorrectCount: 36, ScaledScore: 52},
		{SectionType: "reading", CorrectCount: 37, ScaledScore: 53},
		{SectionType: "reading", CorrectCount: 38, ScaledScore: 54},
		{SectionType: "reading", CorrectCount: 39, ScaledScore: 54},
		{SectionType: "reading", CorrectCount: 40, ScaledScore: 55},
		{SectionType: "reading", CorrectCount: 41, ScaledScore: 56},
		{SectionType: "reading", CorrectCount: 42, ScaledScore: 57},
		{SectionType: "reading", CorrectCount: 43, ScaledScore: 58},
		{SectionType: "reading", CorrectCount: 44, ScaledScore: 59},
		{SectionType: "reading", CorrectCount: 45, ScaledScore: 60},
		{SectionType: "reading", CorrectCount: 46, ScaledScore: 61},
		{SectionType: "reading", CorrectCount: 47, ScaledScore: 63},
		{SectionType: "reading", CorrectCount: 48, ScaledScore: 65},
		{SectionType: "reading", CorrectCount: 49, ScaledScore: 66},
		{SectionType: "reading", CorrectCount: 50, ScaledScore: 67},
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "section_type"}, {Name: "correct_count"}},
		DoUpdates: clause.AssignmentColumns([]string{"scaled_score"}),
	}).CreateInBatches(&conversions, 100).Error

	if err != nil {
		log.Printf("failed to seed score conversions: %v\n", err)
	} else {
		log.Printf("successfully seeded %d TOEFL score conversion records\n", len(conversions))
	}
}