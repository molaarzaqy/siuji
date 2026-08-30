package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTPRequest struct {
	UserID  uint   `json:"-" validate:"required"`
	OTPCode string `json:"otp_code" validate:"required,len=6"`
}

type ResetPasswordRequest struct {
	UserID             uint   `json:"-" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}

type ChangePasswordRequest struct {
	UserID             uint   `json:"-" validate:"required"`
	OldPassword        string `json:"old_password" validate:"required"`
	NewPassword        string `json:"new_password" validate:"required,min=8"`
	ConfirmNewPassword string `json:"confirm_new_password" validate:"required,eqfield=NewPassword"`
}

type GetMeRequest struct {
	UserID uint `json:"-" validate:"required"`
}

type AuthResponse struct {
	TempToken    string        `json:"temp_token,omitempty"`
	AccessToken  string        `json:"access_token,omitempty"`
	RefreshToken string        `json:"refresh_token,omitempty"`
	ExpiresIn    int           `json:"expires_in,omitempty"`
	User         *UserResponse `json:"user,omitempty"`
}