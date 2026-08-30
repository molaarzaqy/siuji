package main

import (
	"fmt"
	_ "siuji-backend/docs"
	"siuji-backend/internal/config"
)

// @title Siuji API
// @version 1.0
// @description API documentation for Siuji TOEFL exam application.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT access token.
// @securityDefinitions.apikey TempToken
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and temp JWT token.
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	jwtManager := config.NewJWTManager(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB: db,
		App: app,
		Log: log,
		Validate: validate,
		Config: viperConfig,
		JWTManager: jwtManager,
	})

	port := viperConfig.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	err := app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}