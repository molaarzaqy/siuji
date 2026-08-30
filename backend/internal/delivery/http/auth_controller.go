package http

import (
	"siuji-backend/internal/model"
	"siuji-backend/internal/usecase"
	"siuji-backend/pkg/cookie"
	"siuji-backend/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	UseCase *usecase.AuthUseCase
}

func NewAuthController(useCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{UseCase: useCase}
}

func getUserIDFromContext(c fiber.Ctx) (uint, error) {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "invalid user context")
	}
	return userID, nil
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user with email and password, returns access & refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login credentials"
// @Success      200 {object} response.Response{data=model.AuthResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      401 {object} response.ResponseNoData
// @Router       /auth/login [post]
func (ctrl *AuthController) Login(c fiber.Ctx) error {
	request := new(model.LoginRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := ctrl.UseCase.Login(request)
	if err != nil {
		return err
	}

	if result.AccessToken != "" {
		cookie.SetAuthCookies(c, result.AccessToken, result.RefreshToken,
			ctrl.UseCase.JWTManager.AccessTokenDuration(),
			ctrl.UseCase.JWTManager.RefreshTokenDuration())
	}

	return response.Success(c, "Login successfully. Welcome to the dashboard.", result)
}

// ForgotPassword godoc
// @Summary      Forgot password
// @Description  Request an OTP code to be sent to the registered email for password reset.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body model.ForgotPasswordRequest true "Registered email"
// @Success      200 {object} response.Response{data=model.AuthResponse}
// @Failure      400 {object} response.ResponseNoData
// @Router       /auth/forgot-password [post]
func (ctrl *AuthController) ForgotPassword(c fiber.Ctx) error {
	request := new(model.ForgotPasswordRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := ctrl.UseCase.ForgotPassword(request)
	if err != nil {
		return err
	}

	if result == nil {
		return response.SuccessNoData(c, "If this email is registered, an OTP verification code has been sent.")
	}

	return response.Success(c, "OTP verification code has been sent to your email.", result)
}

// VerifyOTP godoc
// @Summary      Verify OTP
// @Description  Verify the OTP code sent to email. Requires a temp token (from forgot-password) in Authorization header.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     TempToken
// @Param        request body model.VerifyOTPRequest true "OTP code"
// @Success      200 {object} response.Response{data=model.AuthResponse}
// @Failure      400 {object} response.ResponseNoData
// @Failure      401 {object} response.ResponseNoData
// @Router       /auth/verify-otp [post]
func (ctrl *AuthController) VerifyOTP(c fiber.Ctx) error {
	request := new(model.VerifyOTPRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	request.UserID = userID

	result, err := ctrl.UseCase.VerifyOTP(request)
	if err != nil {
		return err
	}

	return response.Success(c, "OTP verified. Please set your new password.", result)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Set a new password after OTP verification. Requires a temp token (from verify-otp) in Authorization header.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     TempToken
// @Param        request body model.ResetPasswordRequest true "New password"
// @Success      200 {object} response.ResponseNoData
// @Failure      400 {object} response.ResponseNoData
// @Failure      401 {object} response.ResponseNoData
// @Router       /auth/reset-password [post]
func (ctrl *AuthController) ResetPassword(c fiber.Ctx) error {
	request := new(model.ResetPasswordRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	request.UserID = userID

	if err := ctrl.UseCase.ResetPassword(request); err != nil {
		return err
	}

	return response.SuccessNoData(c, "Password reset successfully. Please login with your new password.")
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change password for the currently authenticated user.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body model.ChangePasswordRequest true "Old and new password"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.ResponseNoData
// @Failure      401 {object} response.ResponseNoData
// @Router       /auth/change-password [post]
func (ctrl *AuthController) ChangePassword(c fiber.Ctx) error {
	request := new(model.ChangePasswordRequest)
	if err := c.Bind().Body(request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	request.UserID = userID

	if err := ctrl.UseCase.ChangePassword(request); err != nil {
		return err
	}

	return response.Success(c, "Password changed successfully.", nil)
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get the profile of the currently authenticated user.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=model.UserResponse}
// @Failure      401 {object} response.ResponseNoData
// @Router       /auth/me [get]
func (ctrl *AuthController) GetMe(c fiber.Ctx) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}

	result, err := ctrl.UseCase.GetMe(&model.GetMeRequest{UserID: userID})
	if err != nil {
		return err
	}

	return response.Success(c, "User profile retrieved successfully.", result)
}

// Logout godoc
// @Summary      Logout
// @Description  Clear authentication cookies.
// @Tags         Auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.ResponseNoData
// @Router       /auth/logout [post]
func (ctrl *AuthController) Logout(c fiber.Ctx) error {
	cookie.ClearAuthCookies(c)
	return response.SuccessNoData(c, "Logout successfully.")
}