package usecase

import (
	"crypto/subtle"
	"fmt"
	"siuji-backend/internal/entity"
	"siuji-backend/internal/model"
	"siuji-backend/internal/model/converter"
	"siuji-backend/internal/repository"
	"siuji-backend/pkg/email"
	"siuji-backend/pkg/jwt"
	"siuji-backend/pkg/otp"
	"siuji-backend/pkg/password"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

const (
	OTPValidDuration = 5
	MaxOTPAttempts = 5
)

type AuthUseCase struct {
	Log 		   *logrus.Logger
	Validate	   *validator.Validate
	UserRepository repository.UserRepository
	OTPRepository  repository.OTPRepository
	EmailService   email.Service
	JWTManager	   *jwt.Manager
}

func NewAuthUseCase(
	log *logrus.Logger,
	validate *validator.Validate,
	userRepository repository.UserRepository,
	otpRepository repository.OTPRepository,
	emailService email.Service,
	jwtManager *jwt.Manager,
) *AuthUseCase {
	return &AuthUseCase{
		Log: log,
		Validate: validate,
		UserRepository: userRepository,
		OTPRepository: otpRepository,
		EmailService: emailService,
		JWTManager: jwtManager,
	}
}


func (c *AuthUseCase) Login(request *model.LoginRequest) (*model.AuthResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid login request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByEmail(request.Email)
	if err != nil || !password.Check(request.Password, user.Password) {
		c.Log.Warnf("login failed for email %s", request.Email)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid credentials")
	}

	accessToken, err := c.JWTManager.GenerateAccessToken(user.ID, user.PublicID, user.Email, user.Role)
	if err != nil {
		c.Log.Errorf("failed to generate access token: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to generate access token")
	}

	refreshToken, err := c.JWTManager.GenerateRefreshToken(user.ID, jwt.GenerateTokenFamily())
	if err != nil {
		c.Log.Errorf("failed to generate refresh token: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to generate refresh token")
	}

	return &model.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(c.JWTManager.AccessTokenDuration().Seconds()),
		User:         converter.UserToResponse(user),
	}, nil
}

func (c *AuthUseCase) ForgotPassword(request *model.ForgotPasswordRequest) (*model.AuthResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid forgot password request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	// SECURITY: Blind Response — jangan bocorkan apakah email terdaftar atau tidak
	user, err := c.UserRepository.FindByEmail(request.Email)
	if err != nil {
		c.Log.Infof("forgot password requested for unregistered email: %s", request.Email)
		return nil, nil
	}

	if err := c.OTPRepository.DeleteByEmail(request.Email); err != nil {
		c.Log.Errorf("failed to cleanup old OTP: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process request")
	}

	code, err := otp.GenerateCode()
	if err != nil {
		c.Log.Errorf("failed to generate OTP code: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process request")
	}

	newOTP := &entity.OTP{
		Email:     request.Email,
		Code:      code,
		Purpose:   otp.PurposeResetPassword,
		ExpiresAt: timeNowAddMinutes(OTPValidDuration),
	}
	if err := c.OTPRepository.Create(newOTP); err != nil {
		c.Log.Errorf("failed to save OTP: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process request")
	}

	if err := c.EmailService.SendOTP(request.Email, code); err != nil {
		c.Log.Errorf("failed to send OTP email: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to send OTP email")
	}

	tempToken, err := c.JWTManager.GenerateTempToken(user.ID, request.Email, jwt.PurposeVerifyEmail, 5)
	if err != nil {
		c.Log.Errorf("failed to generate temp token: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process request")
	}

	return &model.AuthResponse{
		TempToken: tempToken,
		ExpiresIn: 5 * 60,
	}, nil
}

func (c *AuthUseCase) VerifyOTP(request *model.VerifyOTPRequest) (*model.AuthResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid verify otp request: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByID(request.UserID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	otpRecord, err := c.OTPRepository.FindValidByEmailAndPurpose(user.Email, otp.PurposeResetPassword)
	if err != nil {
		c.Log.Warnf("no valid OTP found for user %d", request.UserID)
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid or expires OTP")
	}
	// batas percobaan tercapai: hanguskan OTP, paksa minta kode baru
	if otpRecord.Attempts >= MaxOTPAttempts {
		c.Log.Warnf("OTP attempts limit exceeded for user %d, invalidating code", request.UserID)
		if err := c.OTPRepository.DeleteByEmail(user.Email); err != nil {
			c.Log.Errorf("failedd to invalidate OTP after attempts limmit: %+v", err)
		}
		return nil, fiber.NewError(fiber.StatusTooManyRequests, "to many failed attempts, please request a new OTP code")
	}
	// Perbandingan waktu-konstan supaya durasi respons tidak membocorkan
	// seberapa banyak digit awal yang sudah benar.
	if subtle.ConstantTimeCompare([]byte(otpRecord.Code), []byte(request.OTPCode)) != 1 {
		if err := c.OTPRepository.IncrementAttempts(otpRecord.ID); err != nil {
			c.Log.Errorf("failed to increment OTP attempts: %+v", err)
		}
		remaining := MaxOTPAttempts - (otpRecord.Attempts + 1)
		c.Log.Warnf("incorrect OTP for user %d, %d attempts remaining", request.UserID, remaining)

		if remaining <= 0 {
			return nil, fiber.NewError(fiber.StatusTooManyRequests, "to many failed attempts, please request a new OTP code")
		}
		return nil, fiber.NewErrorf(fiber.StatusBadRequest, 
			fmt.Sprintf("invalid OTP code, %d attempt(s) remaining", remaining))
	}
	if err := c.OTPRepository.DeleteByEmail(user.Email); err != nil {
		c.Log.Errorf("failed to cleanup OTP: %+v", err)
	}
	
	tempToken, err := c.JWTManager.GenerateTempToken(user.ID, user.Email, jwt.PurposeResetPassword, 5)
	if err != nil {
		c.Log.Errorf("failed to generate temp token: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to process request")
	}
	return &model.AuthResponse{
		TempToken: tempToken,
		ExpiresIn: 5 * 60,
	}, nil
}

func (c *AuthUseCase) ResetPassword(request *model.ResetPasswordRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid reset password request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByID(request.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	hashed, err := password.Hash(request.NewPassword)
	if err != nil {
		c.Log.Errorf("failed to hash password: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to reset password")
	}

	if err := c.UserRepository.UpdatePassword(request.UserID, hashed); err != nil {
		c.Log.Errorf("failed to update password: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to reset password")
	}

	_ = c.OTPRepository.DeleteByEmail(user.Email)

	return nil
}

func (c *AuthUseCase) ChangePassword(request *model.ChangePasswordRequest) error {
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("invalid change password request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByID(request.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	if !password.Check(request.OldPassword, user.Password) {
		c.Log.Warnf("change password failed: incorrect old password for user %d", request.UserID)
		return fiber.NewError(fiber.StatusBadRequest, "current password is incorrect")
	}

	hashed, err := password.Hash(request.NewPassword)
	if err != nil {
		c.Log.Errorf("failed to hash password: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "failed to change password")
	}

	return c.UserRepository.UpdatePassword(request.UserID, hashed)
}

func (c *AuthUseCase) GetMe(request *model.GetMeRequest) (*model.UserResponse, error) {
	if err := c.Validate.Struct(request); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	user, err := c.UserRepository.FindByID(request.UserID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
	}

	return converter.UserToResponse(user), nil
}

func timeNowAddMinutes(minutes int) (t time.Time) {
	return time.Now().Add(time.Duration(minutes) * time.Minute)
}