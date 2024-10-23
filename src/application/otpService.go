package application

import (
	"crypto/rand"
	"first-project/src/entities"
	"first-project/src/exceptions"
	"io"
	"time"
)

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

func VerifyOTP(
	user entities.User, email, inputOTP, otpFieldError, expiredTokenTagError, invalidTokenTagError string) {
	var registrationError exceptions.UserRegistrationError

	if time.Since(user.UpdatedAt) > 20*time.Minute {
		registrationError.AppendError(
			otpFieldError,
			expiredTokenTagError)
		panic(registrationError)
	}
	if inputOTP != user.Token {
		registrationError.AppendError(
			otpFieldError,
			invalidTokenTagError)
		panic(registrationError)
	}
}
