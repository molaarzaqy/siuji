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
	PeriodController      *http.PeriodController
	SectionController     *http.SectionController
	ParticipantController *http.ParticipantController
	UserController        *http.UserController
	JWTManager     *jwt.Manager
	Log            *logrus.Logger
}

func (r *RouteConfig) Setup() {
	r.setupSwaggerRoutes()
	r.setupAuthRoutes()
	r.setupPeriodRoutes()
	r.setupSectionRoutes()
	r.setupUserRoutes()
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

func (r *RouteConfig) setupPeriodRoutes() {
	periods := r.App.Group("/api/v1/periods",
		middleware.JWTAuth(r.JWTManager, r.Log),
		middleware.RequireRole("admin"),
	)

	periods.Post("/", r.PeriodController.Create)
	periods.Get("/", r.PeriodController.GetAll)
	periods.Get("/:period_public_id", r.PeriodController.GetDetail)
	periods.Put("/:period_public_id", r.PeriodController.Update)
	periods.Delete("/:period_public_id", r.PeriodController.Delete)

	periods.Post("/:period_public_id/sections", r.PeriodController.AddSection)
	periods.Delete("/:period_public_id/sections/:section_public_id", r.PeriodController.RemoveSection)
	periods.Put("/:period_public_id/sections/reorder", r.PeriodController.ReorderSections)

	periods.Post("/:period_public_id/participants", r.ParticipantController.Add)
	periods.Get("/:period_public_id/participants", r.ParticipantController.GetAll)
	periods.Get("/:period_public_id/participants/:user_public_id", r.ParticipantController.GetDetail)
	periods.Put("/:period_public_id/participants/:user_public_id", r.ParticipantController.Update)
	periods.Delete("/:period_public_id/participants/:user_public_id", r.ParticipantController.Remove)
}

func (r *RouteConfig) setupSectionRoutes() {
	sections := r.App.Group("/api/v1/sections",
		middleware.JWTAuth(r.JWTManager, r.Log),
		middleware.RequireRole("admin"),
	)

	sections.Post("/", r.SectionController.Create)
	sections.Get("/", r.SectionController.GetAll)
	sections.Get("/:section_public_id", r.SectionController.GetDetail)
	sections.Put("/:section_public_id", r.SectionController.Update)
	sections.Delete("/:section_public_id", r.SectionController.Delete)
}

func (r *RouteConfig) setupUserRoutes() {
	users := r.App.Group("/api/v1/users",
		middleware.JWTAuth(r.JWTManager, r.Log),
		middleware.RequireRole("admin"),
	)

	users.Get("/", r.UserController.GetAll)
	users.Get("/:user_public_id", r.UserController.GetDetail)
	users.Delete("/:user_public_id", r.UserController.Delete)
}