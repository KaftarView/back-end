package repository_database_interfaces

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"
)

type UserRepository interface {
	ActivateUserAccount(user *entities.User)
	AssignPermissionToRole(role *entities.Role, permission *entities.Permission)
	AssignRoleToUser(user *entities.User, role *entities.Role)
	CreateNewPermission(permissionType enums.PermissionType) *entities.Permission
	CreateNewRole(roleType string) *entities.Role
	CreateNewUser(username string, email string, password string, token string, verified bool) *entities.User
	DeleteRoleByRoleID(roleID uint)
	DeleteRolePermission(role *entities.Role, permission *entities.Permission)
	DeleteUserRole(user *entities.User, role *entities.Role)
	FindActiveOrVerifiedUserByEmail(email string) (*entities.User, bool)
	FindActiveOrVerifiedUserByUsername(username string) (*entities.User, bool)
	FindAllPermissions() []*entities.Permission
	FindAllRolesWithPermissions() []*entities.Role
	FindByEmailAndVerified(email string, verified bool) (*entities.User, bool)
	FindByUserID(userID uint) (*entities.User, bool)
	FindByUsernameAndVerified(username string, verified bool) (*entities.User, bool)
	FindPermissionByID(permissionID uint) (*entities.Permission, bool)
	FindPermissionByType(permissionType enums.PermissionType) (*entities.Permission, bool)
	FindPermissionsByRole(roleID uint) []enums.PermissionType
	FindRoleByID(roleID uint) (*entities.Role, bool)
	FindRoleByType(roleType string) (*entities.Role, bool)
	FindUnverifiedUsersWeekAgo(startOfWeekAgo time.Time, endOfWeekAgo time.Time) []*entities.User
	FindUserRoleTypesByUserID(userID uint) []entities.Role
	FindUsersByRoleID(roleID uint) []*entities.User
	UpdateUserPassword(user *entities.User, password string)
	UpdateUserToken(user *entities.User, token string)
}
