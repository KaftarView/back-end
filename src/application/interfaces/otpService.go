package application_interfaces

import "first-project/src/entities"

type OTPService interface {
	GenerateOTP() string
	VerifyOTP(user *entities.User, inputOTP string, otpFieldError string, expiredTokenTagError string, invalidTokenTagError string)
}
