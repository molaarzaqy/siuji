package route

import (
	"siuji-backend/internal/delivery/http"
	"siuji-backend/internal/delivery/http/middleware"
	"siuji-backend/pkg/jwt"

	"github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	App            *fiber.App
	AuthController *http.AuthController
	JWTManager     *jwt.Manager
	Log            *logrus.Logger
}

func (r *RouteConfig) Setup() {
	r.setupSwaggerRoutes()
	r.setupAuthRoutes()
}

func (r *RouteConfig) setupSwaggerRoutes() {
	r.App.Get("/swagger/*", swaggo.New(swaggo.Config{}))
}

func (r *RouteConfig) setupAuthRoutes() {
	auth := r.App.Group("/api/v1/auth")

	auth.Post("/login", r.AuthController.Login)
	auth.Post("/forgot-password", r.AuthController.ForgotPassword)
	auth.Post("/verify-otp", middleware.TempAuth(r.JWTManager, r.Log, jwt.PurposeVerifyEmail), r.AuthController.VerifyOTP)
	auth.Post("/reset-password", middleware.TempAuth(r.JWTManager, r.Log, jwt.PurposeResetPassword), r.AuthController.ResetPassword)
	auth.Post("/change-password", middleware.JWTAuth(r.JWTManager, r.Log), r.AuthController.ChangePassword)
	auth.Get("/me", middleware.JWTAuth(r.JWTManager, r.Log), r.AuthController.GetMe)
	auth.Post("/logout", middleware.JWTAuth(r.JWTManager, r.Log), r.AuthController.Logout)
}