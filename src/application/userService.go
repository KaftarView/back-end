package application

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"mime/multipart"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userService struct {
	constants      *bootstrap.Constants
	userRepository repository_database_interfaces.UserRepository
	otpService     *OTPService
	awsS3Service   *application_aws.S3service
	db             *gorm.DB
}

func NewUserService(
	constants *bootstrap.Constants,
	userRepository repository_database_interfaces.UserRepository,
	otpService *OTPService,
	awsS3Service *application_aws.S3service,
	db *gorm.DB,
) *userService {
	return &userService{
		constants:      constants,
		userRepository: userRepository,
		otpService:     otpService,
		awsS3Service:   awsS3Service,
		db:             db,
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

func (userService *userService) FindByUserID(id uint) (*entities.User, bool) {
	return userService.userRepository.FindByUserID(userService.db, id)
}

func (userService *userService) ValidateUserRegistrationDetails(
	username string, email string, password string, confirmPassword string) {
	var registrationError exceptions.UserRegistrationError
	var conflictError exceptions.ConflictError
	isRegError := false
	_, usernameExist := userService.userRepository.FindActiveOrVerifiedUserByUsername(userService.db, username)
	if usernameExist {
		isRegError = true
		conflictError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
	}
	_, emailExist := userService.userRepository.FindActiveOrVerifiedUserByEmail(userService.db, email)
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
	user, notVerifiedUserExist := userService.userRepository.FindByUsernameAndVerified(userService.db, username, false)
	if notVerifiedUserExist {
		userService.userRepository.UpdateUserToken(userService.db, user, otp)
	} else {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			panic(err)
		}

		user := &entities.User{
			Name:     username,
			Email:    email,
			Password: hashedPassword,
			Token:    otp,
			Verified: false,
		}
		err = repository_database.ExecuteInTransaction(userService.db, func(tx *gorm.DB) error {
			userService.userRepository.CreateNewUser(tx, user)
			role, _ := userService.userRepository.FindRoleByType(userService.db, enums.User.String())
			userService.userRepository.AssignRoleToUser(tx, user, role)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
}

func (userService *userService) ActivateUser(email, otp string) {
	var registrationError exceptions.UserRegistrationError
	_, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
	if verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.AlreadyVerified)
		panic(registrationError)
	}

	user, _ := userService.userRepository.FindByEmailAndVerified(userService.db, email, false)
	userService.otpService.VerifyOTP(
		user, otp, userService.constants.ErrorField.OTP,
		userService.constants.ErrorTag.ExpiredToken,
		userService.constants.ErrorTag.InvalidToken)
	userService.userRepository.ActivateUserAccount(userService.db, user)
}

func (userService *userService) AuthenticateUser(username string, password string) (user *entities.User) {
	user, verifiedUserExist := userService.userRepository.FindByUsernameAndVerified(userService.db, username, true)
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
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
	if !verifiedUserExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Email,
			userService.constants.ErrorTag.EmailNotExist)
		panic(registrationError)
	}
	userService.userRepository.UpdateUserToken(userService.db, user, otp)
}

func (userService *userService) ValidateUserOTP(email, otp string) uint {
	var registrationError exceptions.UserRegistrationError
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
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

func (userService *userService) UpdateUser(userID uint, username string) {
	var conflictError exceptions.ConflictError
	var notFoundError exceptions.NotFoundError
	user, userExist := userService.FindByUserID(userID)
	if !userExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	_, usernameExist := userService.userRepository.FindActiveOrVerifiedUserByUsername(userService.db, username)
	if usernameExist {
		conflictError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
	}
	user.Name = username
	if err := userService.userRepository.UpdateUser(userService.db, user); err != nil {
		panic(err)
	}
}

func (userService *userService) ResetPasswordService(userID uint, password, confirmPassword string) {
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

	user, _ := userService.FindByUserID(userID)
	userService.userRepository.UpdateUserPassword(userService.db, user, hashedPassword)
}

func (userService *userService) CreateNewRole(name string) *entities.Role {
	var registrationError exceptions.UserRegistrationError
	_, roleExist := userService.userRepository.FindRoleByType(userService.db, name)
	if roleExist {
		registrationError.AppendError(
			userService.constants.ErrorField.Role,
			userService.constants.ErrorTag.AlreadyExist)
		panic(registrationError)
	}
	role := userService.userRepository.CreateNewRole(userService.db, name)
	return role
}

func (userService *userService) AssignPermissionsToRole(roleID uint, permissions []string) {
	var notFoundError exceptions.NotFoundError
	role, roleExist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	permissionsMap := make(map[string]bool)
	for _, permission := range permissions {
		permissionsMap[permission] = true
	}
	existingPermissions := userService.userRepository.FindPermissionsByRole(userService.db, roleID)
	for _, permission := range existingPermissions {
		permissionsMap[permission.String()] = false
	}
	permissionTypes := enums.GetAllPermissionTypes()

	err := repository_database.ExecuteInTransaction(userService.db, func(tx *gorm.DB) error {
		for _, permission := range permissionTypes {
			if !permissionsMap[permission.String()] {
				continue
			}
			permission, _ := userService.userRepository.FindPermissionByType(userService.db, permission)
			if err := userService.userRepository.AssignPermissionToRole(tx, role, permission); err != nil {
				panic(err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (userService *userService) UpdateUserRoles(email string, roles []string) {
	var notFoundError exceptions.NotFoundError
	user, verifiedUserExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
	if !verifiedUserExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	allowedRolesMap := make(map[string]bool)
	for _, role := range roles {
		allowedRolesMap[role] = true
	}
	existingRoles := userService.userRepository.FindUserRoleTypesByUserID(userService.db, user.ID)
	for _, role := range existingRoles {
		allowedRolesMap[role.Type] = false
	}

	err := repository_database.ExecuteInTransaction(userService.db, func(tx *gorm.DB) error {
		for roleType, ok := range allowedRolesMap {
			if !ok {
				continue
			}
			role, _ := userService.userRepository.FindRoleByType(userService.db, roleType)
			if err := userService.userRepository.AssignRoleToUser(tx, user, role); err != nil {
				panic(err)
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (userService *userService) FindUserRolesAndPermissions(userID uint) ([]string, []string) {
	var roleTypes []string
	var permissionTypes []string
	roles := userService.userRepository.FindUserRoleTypesByUserID(userService.db, userID)
	for _, role := range roles {
		roleTypes = append(roleTypes, role.Type)
		permissions := userService.userRepository.FindPermissionsByRole(userService.db, role.ID)
		for _, permission := range permissions {
			permissionTypes = append(permissionTypes, permission.String())
		}
	}
	return roleTypes, permissionTypes
}

func (userService *userService) GetRolesList() []dto.RoleDetailsResponse {
	roles := userService.userRepository.FindAllRolesWithPermissions(userService.db)
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
	_, roleExist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	users := userService.userRepository.FindUsersByRoleID(userService.db, roleID)
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
	_, roleExist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	userService.userRepository.DeleteRoleByRoleID(userService.db, roleID)
}

func (userService *userService) DeleteRolePermission(roleID, permissionID uint) {
	var notFoundError exceptions.NotFoundError
	role, roleExist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	permission, permissionExist := userService.userRepository.FindPermissionByID(userService.db, permissionID)
	if !permissionExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Permission
		panic(notFoundError)
	}
	userService.userRepository.DeleteRolePermission(userService.db, role, permission)
}

func (userService *userService) GetPermissionsList() []dto.PermissionDetailsResponse {
	permissions := userService.userRepository.FindAllPermissions(userService.db)
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
	user, userExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
	if !userExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	role, roleExist := userService.userRepository.FindRoleByID(userService.db, roleID)
	if !roleExist {
		notFoundError.ErrorField = userService.constants.ErrorField.Role
		panic(notFoundError)
	}
	userService.userRepository.DeleteUserRole(userService.db, user, role)
}

func (userService *userService) GetCouncilorsList(promotedYear int) []dto.CouncilorsDetailsResponse {
	councilors := userService.userRepository.FindAllCouncilorsByPromotedYear(userService.db, promotedYear)
	councilorsDetails := make([]dto.CouncilorsDetailsResponse, len(councilors))
	for i, councilor := range councilors {
		profile := ""
		if councilor.ProfilePath != "" {
			profile = userService.awsS3Service.GetPresignedURL(enums.ProfilesBucket, councilor.ProfilePath, 8*time.Hour)
		}
		user, _ := userService.FindByUserID(councilor.UserID)
		councilorsDetails[i] = dto.CouncilorsDetailsResponse{
			ID:           councilor.ID,
			FirstName:    councilor.FirstName,
			LastName:     councilor.LastName,
			Email:        user.Email,
			EnteringYear: councilor.EnteringYear,
			Description:  councilor.Description,
			Profile:      profile,
		}
	}
	return councilorsDetails
}

func (userService *userService) CreateCouncilor(email, firstName, lastName, description string, promotedYear int, enteringYear int, profile *multipart.FileHeader) {
	var notFoundError exceptions.NotFoundError
	var conflictError exceptions.ConflictError
	user, userExist := userService.userRepository.FindByEmailAndVerified(userService.db, email, true)
	if !userExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	_, councilorExist := userService.userRepository.FindCouncilorByUserIDAndPromotedYear(userService.db, user.ID, promotedYear)
	if councilorExist {
		conflictError.AppendError(
			userService.constants.ErrorField.Username,
			userService.constants.ErrorTag.AlreadyExist)
		panic(conflictError)
	}

	councilor := &entities.Councilor{
		FirstName:    firstName,
		LastName:     lastName,
		EnteringYear: enteringYear,
		Description:  description,
		PromotedYear: promotedYear,
		UserID:       user.ID,
	}
	err := repository_database.ExecuteInTransaction(userService.db, func(tx *gorm.DB) error {
		if err := userService.userRepository.CreateNewCouncilor(tx, councilor); err != nil {
			panic(err)
		}

		profilePath := userService.constants.S3Service.GetCouncilorProfileKey(councilor.ID, profile.Filename)
		userService.awsS3Service.UploadObject(enums.ProfilesBucket, profilePath, profile)
		councilor.ProfilePath = profilePath

		if err := userService.userRepository.UpdateCouncilor(tx, councilor); err != nil {
			panic(err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func (userService *userService) DeleteCouncilor(councilorID uint) {
	var notFoundError exceptions.NotFoundError
	councilor, councilorExist := userService.userRepository.FindCouncilorByID(userService.db, councilorID)
	if !councilorExist {
		notFoundError.ErrorField = userService.constants.ErrorField.User
		panic(notFoundError)
	}
	userService.userRepository.DeleteCouncilor(userService.db, councilorID)
	if councilor.ProfilePath != "" {
		userService.awsS3Service.DeleteObject(enums.ProfilesBucket, councilor.ProfilePath)
	}
}

func (userService *userService) GetUsersByPermissions(permissions []enums.PermissionType) []entities.User {
	users := userService.userRepository.FindUsersByPermissions(userService.db, permissions)
	if len(users) == 0 {
		return nil
	}
	return users
}
