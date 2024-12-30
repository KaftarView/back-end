package application_interfaces

import (
	"first-project/src/dto"
	"first-project/src/entities"
	"mime/multipart"
)

type UserService interface {
	ActivateUser(email string, otp string)
	AssignPermissionsToRole(roleID uint, permissions []string)
	AuthenticateUser(username string, password string) (user *entities.User)
	CreateCouncilor(email, firstName, lastName, description string, promotedYear, semester int, profile *multipart.FileHeader)
	CreateNewRole(name string) *entities.Role
	DeleteRole(roleID uint)
	DeleteRolePermission(roleID uint, permissionID uint)
	DeleteUserRole(email string, roleID uint)
	FindUserRolesAndPermissions(userID uint) ([]string, []string)
	GetPermissionsList() []dto.PermissionDetailsResponse
	GetRoleOwners(roleID uint) []dto.UserDetailsResponse
	GetRolesList() []dto.RoleDetailsResponse
	ResetPasswordService(email string, password string, confirmPassword string)
	UpdateOrCreateUser(username string, email string, password string, otp string)
	UpdateUserOTPIfExists(email string, otp string)
	UpdateUserRoles(email string, roles []string)
	ValidateUserOTP(email string, otp string) uint
	ValidateUserRegistrationDetails(username string, email string, password string, confirmPassword string)
}
