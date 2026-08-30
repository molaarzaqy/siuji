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

	// setup external services
	emailService := email.NewService()

	// setup usecases
	authUseCase := usecase.NewAuthUseCase(
		config.Log,
		config.Validate,
		userRepository,
		otpRepository,
		emailService,
		config.JWTManager,
	)

	// setup controllers
	authController := http.NewAuthController(authUseCase)

	// setup routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		AuthController: authController,
		JWTManager:     config.JWTManager,
		Log:            config.Log,
	}
	routeConfig.Setup()
}