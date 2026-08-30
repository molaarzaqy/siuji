package email

import "log"

type Service interface {
	SendOTP(email, otpCode string) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) SendOTP(email, otpCode string) error {
	log.Printf("OTP EMAIL SIMULATION")
	log.Printf("To: %s", email)
	log.Printf("OTP Code: %s", otpCode)
	return nil
}