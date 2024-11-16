package application

import (
	"crypto/rand"
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	"first-project/src/repository"
	"io"
	"time"
)

type OTPService struct {
	constants      *bootstrap.Constants
	userRepository *repository.UserRepository
}

func NewOTPService(constants *bootstrap.Constants, userRepository *repository.UserRepository) *OTPService {
	return &OTPService{
		constants:      constants,
		userRepository: userRepository,
	}
}

const otpLength = 6

var table = []byte("123456789")

func GenerateOTP() string {
	otp := make([]byte, otpLength)
	n, err := io.ReadAtLeast(rand.Reader, otp, otpLength)
	if n != otpLength {
		panic(err)
	}
	for i := 0; i < len(otp); i++ {
		otp[i] = table[int(otp[i])%len(table)]
	}
	return string(otp)
}

func (otpService *OTPService) VerifyOTP(inputOTP string, email string) {
	var registrationError exceptions.UserRegistrationError
	otp, lastSentOtp := otpService.userRepository.GetOTPByEmail(email)

	if time.Since(lastSentOtp) > 5*time.Minute {
		registrationError.AppendError(
			otpService.constants.ErrorField.OTP,
			otpService.constants.ErrorTag.OTPExpired)
		panic(registrationError)
	}
	if inputOTP != otp {
		registrationError.AppendError(
			otpService.constants.ErrorField.OTP,
			otpService.constants.ErrorTag.InvalidOTP)
		panic(registrationError)
	}
}
