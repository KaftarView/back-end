package application

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/enums"
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

func (userService *UserService) ValidateUserRegistrationDetails(
	username string, email string, password string, confirmPassword string) {
	var registrationError exceptions.UserRegistrationError
	isRegError := false
	_, usernameExist := userService.userRepository.FindByUsernameAndVerified(username, true)
	if usernameExist {
		isRegError = true
		registrationError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
	}
	_, emailExist := userService.userRepository.FindByEmailAndVerified(email, true)
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

func (userService *UserService) UpdateOrCreateUser(username string, email string, password string, otp string) {
	user, notVerifiedUserExist := userService.userRepository.FindByUsernameAndVerified(username, false)
	if notVerifiedUserExist {
		userService.userRepository.UpdateUserToken(user, otp)
	} else {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			panic(err)
		}
		user := userService.userRepository.CreateNewUser(username, email, hashedPassword, otp, false)
		role, _ := userService.userRepository.FindRoleByType(enums.User)
		userService.userRepository.AssignRoleToUser(user, role)
	}
}

func (userService *UserService) ActivateUser(email, otp string) {
	var registrationError exceptions.UserRegistrationError
	_, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.AlreadyVerified)
		panic(registrationError)
	}

	user, _ := userService.userRepository.FindByEmailAndVerified(email, false)
	VerifyOTP(
		user, email, otp,
		userService.constants.ErrorField.OTP,
		userService.constants.ErrorTag.ExpiredToken,
		userService.constants.ErrorTag.InvalidToken)

	userService.userRepository.ActivateUserAccount(user)
}

func (userService *UserService) VerifyLogin(username string, password string) (user entities.User) {
	user, verifiedUserExist := userService.userRepository.FindByUsernameAndVerified(username, true)
	if !verifiedUserExist {
		loginError := exceptions.NewLoginError()
		panic(loginError)
	}
	passwordMatch := verifyPassword(user.Password, password)
	if !passwordMatch {
		loginError := exceptions.NewLoginError()
		panic(loginError)
	}
	return user
}

func (userService *UserService) VerifyUserActivated(email string) {
	var registrationError exceptions.UserRegistrationError
	_, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.EmailNotExist)
		panic(registrationError)
	}
}

func (userService *UserService) ResetPasswordService(email, password, confirmPassword string) {
	var registrationError exceptions.UserRegistrationError
	passwordErrorTags := userService.passwordValidation(password)
	if len(passwordErrorTags) > 0 {
		for _, v := range passwordErrorTags {
			registrationError.AppendError(userService.constants.ErrorField.Password, v)
		}
		panic(registrationError)
	}
	if confirmPassword != password {
		registrationError.AppendError(
			userService.constants.ErrorField.Password,
			userService.constants.ErrorTag.NotMatchConfirmPAssword)
		panic(registrationError)
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		panic(err)
	}

	user, _ := userService.userRepository.FindByEmailAndVerified(email, true)
	userService.userRepository.UpdateUserPassword(user, hashedPassword)
}
