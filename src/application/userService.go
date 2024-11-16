package application

import (
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	"first-project/src/repository"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants      *bootstrap.Constants
	userRepository *repository.UserRepository
}

func NewUserService(constants *bootstrap.Constants, userRepository *repository.UserRepository) *UserService {
	return &UserService{
		constants:      constants,
		userRepository: userRepository,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func validatePasswordTests(errors *[]string, test string, password string, tag string) {
	matched, _ := regexp.MatchString(test, password)
	if !matched {
		*errors = append(*errors, tag)
	}
}

func (userService *UserService) passwordValidation(password string) []string {
	var errors []string

	validatePasswordTests(&errors, ".{8,}", password, userService.constants.ErrorTag.MinimumLength)
	validatePasswordTests(&errors, "[a-z]", password, userService.constants.ErrorTag.ContainsLowercase)
	validatePasswordTests(&errors, "[A-Z]", password, userService.constants.ErrorTag.ContainsUppercase)
	validatePasswordTests(&errors, "[0-9]", password, userService.constants.ErrorTag.ContainsNumber)
	validatePasswordTests(&errors, "[^\\d\\w]", password, userService.constants.ErrorTag.ContainsSpecialChar)

	return errors
}

func (userService *UserService) VerifyUserRegistration(username string, email string, password string, confirmPassword string) {
	var registrationError exceptions.UserRegistrationError
	isRegError := false
	usernameExist := userService.userRepository.CheckUsernameExists(username)
	if usernameExist {
		isRegError = true
		// maybe adding translates here but I was not agreed
		registrationError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
	}
	emailExist := userService.userRepository.CheckEmailExists(email)
	if emailExist {
		isRegError = true
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.AlreadyExist)
	}
	passwordErrorTags := userService.passwordValidation(password)
	if len(passwordErrorTags) > 0 {
		isRegError = true
		for _, v := range passwordErrorTags {
			registrationError.AppendError(userService.constants.ErrorField.Password, v)
		}
	}

	if confirmPassword != password {
		isRegError = true
		registrationError.AppendError(
			userService.constants.ErrorField.Password,
			userService.constants.ErrorTag.NotMatchConfirmPAssword)
	}

	if isRegError {
		panic(registrationError)
	}
}

func (userService *UserService) RegisterUser(username string, email string, password string, otp string) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		panic(err)
	}
	userService.userRepository.RegisterUser(username, email, hashedPassword, otp)
}

func (userService *UserService) CheckUserAlreadyVerifiedByEmail(email string) {
	var registrationError exceptions.UserRegistrationError
	alreadyVerified := userService.userRepository.CheckEmailExists(email)
	if alreadyVerified {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.AlreadyVerified)
		panic(registrationError)
	}
}

func (userService *UserService) VerifyEmail(email string) {
	userService.userRepository.VerifyEmail(email)
}

func (userService *UserService) LoginService(username string, password string) {
	hashedPassword, err := userService.userRepository.GetPasswordByUsername(username)
	if err != nil {
		loginError := exceptions.NewLoginError()
		panic(loginError)
	}
	passwordMatch := verifyPassword(hashedPassword, password)
	if !passwordMatch {
		loginError := exceptions.NewLoginError()
		panic(loginError)
	}
}
