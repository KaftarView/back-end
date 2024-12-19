package application

import (
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants      *bootstrap.Constants
	userRepository *repository_database.UserRepository
	otpService     *OTPService
}

func NewUserService(
	constants *bootstrap.Constants, userRepository *repository_database.UserRepository, otpService *OTPService,
) *UserService {
	return &UserService{
		constants:      constants,
		userRepository: userRepository,
		otpService:     otpService,
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
	var conflictError exceptions.ConflictError
	isRegError := false
	_, usernameExist := userService.userRepository.FindActiveOrVerifiedUserByUsername(username)
	if usernameExist {
		isRegError = true
		conflictError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
	}
	_, emailExist := userService.userRepository.FindActiveOrVerifiedUserByEmail(email)
	if emailExist {
		isRegError = true
		conflictError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.AlreadyExist)
	}
	if isRegError {
		panic(conflictError)
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
		role, _ := userService.userRepository.FindRoleByType(enums.User.String())
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
	userService.otpService.VerifyOTP(
		user, otp, userService.constants.ErrorField.OTP,
		userService.constants.ErrorTag.ExpiredToken,
		userService.constants.ErrorTag.InvalidToken)
	userService.userRepository.ActivateUserAccount(user)
}

func (userService *UserService) AuthenticateUser(username string, password string) (user entities.User) {
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

func (userService *UserService) UpdateUserOTPIfExists(email, otp string) {
	var registrationError exceptions.UserRegistrationError
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.EmailNotExist)
		panic(registrationError)
	}
	userService.userRepository.UpdateUserToken(user, otp)
}

func (userService *UserService) ValidateUserOTP(email, otp string) uint {
	var registrationError exceptions.UserRegistrationError
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.EmailNotExist)
		panic(registrationError)
	}
	userService.otpService.VerifyOTP(
		user, otp, userService.constants.ErrorField.OTP,
		userService.constants.ErrorTag.ExpiredToken,
		userService.constants.ErrorTag.InvalidToken)
	return user.ID
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

func (userService *UserService) CreateNewRole(name string) entities.Role {
	var registrationError exceptions.UserRegistrationError
	_, roleExist := userService.userRepository.FindRoleByType(name)
	if roleExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Role,
			userService.constants.ErrorTag.AlreadyExist)
		panic(registrationError)
	}
	role := userService.userRepository.CreateNewRole(name)
	return role
}

func (userService *UserService) AssignPermissionsToRole(roleID uint, permissions []string) {
	var notFoundError exceptions.NotFoundError
	role, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	permissionsMap := make(map[string]bool)
	for _, permission := range permissions {
		permissionsMap[permission] = true
	}
	existingPermissions := userService.userRepository.FindPermissionsByRole(roleID)
	for _, permission := range existingPermissions {
		permissionsMap[permission.String()] = false
	}
	permissionTypes := enums.GetAllPermissionTypes()
	for _, permission := range permissionTypes {
		if permissionsMap[permission.String()] {
			permission, _ := userService.userRepository.FindPermissionByType(permission)
			userService.userRepository.AssignPermissionToRole(role, permission)
		}
	}
}

func (userService *UserService) UpdateUserRoles(userID uint, roles []string) {
	var registrationError exceptions.UserRegistrationError
	user, verifiedUserExist := userService.userRepository.FindByUserID(userID)
	if !verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.EmailNotExist)
		panic(registrationError)
	}
	allowedRolesMap := make(map[string]bool)
	for _, role := range roles {
		allowedRolesMap[role] = true
	}
	existingRoles := userService.userRepository.FindUserRoleTypesByUserID(user.ID)
	for _, role := range existingRoles {
		allowedRolesMap[role.Type] = true
	}

	for roleType, ok := range allowedRolesMap {
		if ok {
			role, _ := userService.userRepository.FindRoleByType(roleType)
			userService.userRepository.AssignRoleToUser(user, role)
		}
	}
}

func (userService *UserService) FindUserRolesAndPermissions(userID uint) ([]string, []string) {
	var roleTypes []string
	var permissionTypes []string
	roles := userService.userRepository.FindUserRoleTypesByUserID(userID)
	for _, role := range roles {
		roleTypes = append(roleTypes, role.Type)
		permissions := userService.userRepository.FindPermissionsByRole(role.ID)
		for _, permission := range permissions {
			permissionTypes = append(permissionTypes, permission.String())
		}
	}
	return roleTypes, permissionTypes
}

func (userService *UserService) GetRolesList() []dto.RoleDetailsResponse {
	roles := userService.userRepository.FindAllRolesWithPermissions()
	rolesDetails := make([]dto.RoleDetailsResponse, len(roles))

	for i, role := range roles {
		permissions := make(map[uint]string)
		for _, permission := range role.Permissions {
			permissions[permission.ID] = permission.Type.String()
		}
		rolesDetails[i] = dto.RoleDetailsResponse{
			ID:          role.ID,
			Type:        role.Type,
			CreatedAt:   role.CreatedAt,
			Permissions: permissions,
		}
	}
	return rolesDetails
}

func (userService *UserService) GetRoleOwners(roleID uint) map[string]string {
	var notFoundError exceptions.NotFoundError
	_, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	users := userService.userRepository.FindUsersByRoleID(roleID)
	userDetails := make(map[string]string)
	for _, user := range users {
		userDetails[user.Email] = user.Name
	}
	return userDetails
}

func (userService *UserService) DeleteRole(roleID uint) {
	var notFoundError exceptions.NotFoundError
	_, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	userService.userRepository.DeleteRoleByRoleID(roleID)
}
