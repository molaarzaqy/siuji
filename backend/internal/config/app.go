package config

import (
	"siuji-backend/internal/delivery/http"
	"siuji-backend/internal/delivery/http/route"
	"siuji-backend/internal/repository"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/email"
	"siuji-backend/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB         *gorm.DB
	App        *fiber.App
	Log    	   *logrus.Logger
	Validate   *validator.Validate
	Config     *viper.Viper
	JWTManager *jwt.Manager
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.DB)
	otpRepository := repository.NewOTPRepository(config.DB)
	periodRepository := repository.NewPeriodRepository(config.DB)
	sectionRepository := repository.NewSectionRepository(config.DB)
	periodSectionRepository := repository.NewPeriodSectionRepository(config.DB)
	participantPeriodRepository := repository.NewParticipantPeriodRepository(config.DB)

	// setup external services
	emailService := email.NewService()

	// setup usecases
	authUseCase := usecase.NewAuthUseCase(config.Log, config.Validate, userRepository, otpRepository, emailService, config.JWTManager)
	periodUseCase := usecase.NewPeriodUseCase(config.Log, config.Validate, periodRepository, sectionRepository, periodSectionRepository)
	sectionUseCase := usecase.NewSectionUseCase(config.Log, config.Validate, sectionRepository)
	userUseCase := usecase.NewUserUseCase(config.Log, userRepository)
	participantUseCase := usecase.NewParticipantUseCase(config.Log, config.Validate, participantPeriodRepository, userRepository, periodRepository)

	// setup controllers
	authController := http.NewAuthController(authUseCase)
	periodController := http.NewPeriodController(periodUseCase)
	sectionController := http.NewSectionController(sectionUseCase)
	userController := http.NewUserController(userUseCase)
	participantController := http.NewParticipantController(participantUseCase)

	// setup routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		PeriodController:       periodController,
		SectionController:      sectionController,
		UserController:         userController,
		ParticipantController:  participantController,
		JWTManager:     config.JWTManager,
		Log:            config.Log,
	}
	routeConfig.Setup()
}