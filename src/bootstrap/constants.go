package bootstrap

import (
	"fmt"
)

type Constants struct {
	Context    Context
	ErrorField ErrorField
	ErrorTag   ErrorTag
	Redis      Redis
}

type Context struct {
	Translator                    string
	IsLoadedValidationTranslator  string
	IsLoadedCustomValidationError string
}

type ErrorField struct {
	Username string
	Password string
	Email    string
	OTP      string
}

type ErrorTag struct {
	AlreadyExist            string
	MinimumLength           string
	ContainsLowercase       string
	ContainsUppercase       string
	ContainsNumber          string
	ContainsSpecialChar     string
	NotMatchConfirmPAssword string
	InvalidToken            string
	AlreadyVerified         string
	OTPExpired              string
	InvalidOTP              string
	LoginFailed             string
}

type Redis struct {
}

func NewConstants() *Constants {
	return &Constants{
		Context: Context{
			Translator:                    "translator",
			IsLoadedValidationTranslator:  "isLoadedValidationTranslator",
			IsLoadedCustomValidationError: "isLoadedCustomValidationError",
		},
		ErrorField: ErrorField{
			Username: "username",
			Password: "password",
			Email:    "email",
			OTP:      "OTP",
		},
		ErrorTag: ErrorTag{
			AlreadyExist:            "alreadyExist",
			ContainsLowercase:       "containsLowercase",
			MinimumLength:           "minimumLength",
			ContainsUppercase:       "containsUppercase",
			ContainsNumber:          "containsNumber",
			ContainsSpecialChar:     "containsSpecialChar",
			NotMatchConfirmPAssword: "notMatchConfirmPAssword",
			InvalidToken:            "invalidToken",
			AlreadyVerified:         "alreadyVerified",
			OTPExpired:              "expiredOTP",
			InvalidOTP:              "invalidOTP",
			LoginFailed:             "loginFailed",
		},
		Redis: Redis{},
	}
}

func (r *Redis) GetUserID(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}
