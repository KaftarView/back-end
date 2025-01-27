package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
	"first-project/src/enums"
	"mime/multipart"
)

type UserService interface {
	ActivateUser(email string, otp string)
	AssignPermissionsToRole(roleID uint, permissions []string)
	AuthenticateUser(username string, password string) (user *entities.User)
	CreateCouncilor(email, firstName, lastName, description string, promotedYear int, enteringYear int, profile *multipart.FileHeader)
	CreateNewRole(name string) *entities.Role
	DeleteCouncilor(councilorID uint)
	DeleteRole(roleID uint)
	DeleteRolePermission(roleID uint, permissionID uint)
	DeleteUserRole(email string, roleID uint)
	FindByUserID(id uint) (*entities.User, bool)
	FindUserRolesAndPermissions(userID uint) ([]string, []string)
	GetCouncilorsList(promotedYear int) []dto.CouncilorsDetailsResponse
	GetPermissionsList() []dto.PermissionDetailsResponse
	GetRoleOwners(roleID uint) []dto.UserDetailsResponse
	GetRolesList() []dto.RoleDetailsResponse
	GetUsersByPermissions(permissions []enums.PermissionType) []entities.User
	ResetPasswordService(userID uint, password string, confirmPassword string)
	UpdateUser(userID uint, username string)
	UpdateOrCreateUser(username string, email string, password string, otp string)
	UpdateUserOTPIfExists(email string, otp string)
	UpdateUserRoles(email string, roles []string)
	ValidateUserOTP(email string, otp string) uint
	ValidateUserRegistrationDetails(username string, email string, password string, confirmPassword string)
}
