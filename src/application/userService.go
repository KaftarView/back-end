package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"mime/multipart"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	constants      *bootstrap.Constants
	userRepository repository_database_interfaces.UserRepository
	otpService     *OTPService
	awsS3Service   *application_aws.S3service
}

func NewUserService(
	constants *bootstrap.Constants,
	userRepository repository_database_interfaces.UserRepository,
	otpService *OTPService,
	awsS3Service *application_aws.S3service,
) *userService {
	return &userService{
		constants:      constants,
		userRepository: userRepository,
		otpService:     otpService,
		awsS3Service:   awsS3Service,
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

func (userService *userService) passwordValidation(password string) []string {
	var errors []string

	validatePasswordTests(&errors, ".{8,}", password, userService.constants.ErrorTag.MinimumLength)
	validatePasswordTests(&errors, "[a-z]", password, userService.constants.ErrorTag.ContainsLowercase)
	validatePasswordTests(&errors, "[A-Z]", password, userService.constants.ErrorTag.ContainsUppercase)
	validatePasswordTests(&errors, "[0-9]", password, userService.constants.ErrorTag.ContainsNumber)
	validatePasswordTests(&errors, "[^\\d\\w]", password, userService.constants.ErrorTag.ContainsSpecialChar)

	return errors
}

func (userService *userService) ValidateUserRegistrationDetails(
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

func (userService *userService) UpdateOrCreateUser(username string, email string, password string, otp string) {
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

func (userService *userService) ActivateUser(email, otp string) {
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

func (userService *userService) AuthenticateUser(username string, password string) (user *entities.User) {
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

func (userService *userService) UpdateUserOTPIfExists(email, otp string) {
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

func (userService *userService) ValidateUserOTP(email, otp string) uint {
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

func (userService *userService) ResetPasswordService(email, password, confirmPassword string) {
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

func (userService *userService) CreateNewRole(name string) *entities.Role {
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

func (userService *userService) AssignPermissionsToRole(roleID uint, permissions []string) {
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

func (userService *userService) UpdateUserRoles(email string, roles []string) {
	var notFoundError exceptions.NotFoundError
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !verifiedUserExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	allowedRolesMap := make(map[string]bool)
	for _, role := range roles {
		allowedRolesMap[role] = true
	}
	existingRoles := userService.userRepository.FindUserRoleTypesByUserID(user.ID)
	for _, role := range existingRoles {
		allowedRolesMap[role.Type] = false
	}

	for roleType, ok := range allowedRolesMap {
		if ok {
			role, _ := userService.userRepository.FindRoleByType(roleType)
			userService.userRepository.AssignRoleToUser(user, role)
		}
	}
}

func (userService *userService) FindUserRolesAndPermissions(userID uint) ([]string, []string) {
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

func (userService *userService) GetRolesList() []dto.RoleDetailsResponse {
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

func (userService *userService) GetRoleOwners(roleID uint) []dto.UserDetailsResponse {
	var notFoundError exceptions.NotFoundError
	_, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	users := userService.userRepository.FindUsersByRoleID(roleID)
	userDetails := make([]dto.UserDetailsResponse, len(users))
	for i, user := range users {
		userDetails[i] = dto.UserDetailsResponse{
			Name:  user.Name,
			Email: user.Email,
		}
	}
	return userDetails
}

func (userService *userService) DeleteRole(roleID uint) {
	var notFoundError exceptions.NotFoundError
	_, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	userService.userRepository.DeleteRoleByRoleID(roleID)
}

func (userService *userService) DeleteRolePermission(roleID, permissionID uint) {
	var notFoundError exceptions.NotFoundError
	role, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	permission, permissionExist := userService.userRepository.FindPermissionByID(permissionID)
	if !permissionExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Permission
		panic(notFoundError)
	}
	userService.userRepository.DeleteRolePermission(role, permission)
}

func (userService *userService) GetPermissionsList() []dto.PermissionDetailsResponse {
	permissions := userService.userRepository.FindAllPermissions()
	permissionsDetails := make([]dto.PermissionDetailsResponse, len(permissions))
	for i, permission := range permissions {
		permissionsDetails[i] = dto.PermissionDetailsResponse{
			ID:          permission.ID,
			Name:        permission.Type.String(),
			Description: permission.Description,
		}
	}
	return permissionsDetails
}

func (userService *userService) DeleteUserRole(email string, roleID uint) {
	var notFoundError exceptions.NotFoundError
	user, userExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !userExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	role, roleExist := userService.userRepository.FindRoleByID(roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	userService.userRepository.DeleteUserRole(user, role)
}

func (userService *userService) CreateCouncilor(email, firstName, lastName, description string, promotedDate time.Time, semester int, profile *multipart.FileHeader) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	user, userExist := userService.userRepository.FindByEmailAndVerified(email, true)
	if !userExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	_, councilorExist := userService.userRepository.FindCouncilorByUserIDAndPromoteDate(user.ID, promotedDate)
	if councilorExist {
		conflictError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	councilor := &entities.Councilor{
		FirstName:    firstName,
		LastName:     lastName,
		Semester:     semester,
		Description:  description,
		PromotedDate: promotedDate,
		UserID:       user.ID,
	}
	userService.userRepository.CreateNewCouncilor(councilor)

	profilePath := userService.constants.S3Service.GetCouncilorProfileKey(councilor.ID, profile.Filename)
	userService.awsS3Service.UploadObject(enums.ProfilesBucket, profilePath, profile)
	councilor.ProfilePath = profilePath
	userService.userRepository.UpdateCouncilor(councilor)
}

func (userService *userService) DeleteCouncilor(councilorID uint) {
	var notFoundError exceptions.NotFoundError
	councilor, councilorExist := userService.userRepository.FindCouncilorByID(councilorID)
	if !councilorExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	userService.userRepository.DeleteCouncilor(councilorID)
	if councilor.ProfilePath != "" {
		userService.awsS3Service.DeleteObject(enums.ProfilesBucket, councilor.ProfilePath)
	}
}
