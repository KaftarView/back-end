package repository_database_interfaces

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	ActivateUserAccount(db *gorm.DB, user *entities.User)
	AssignPermissionToRole(db *gorm.DB, role *entities.Role, permission *entities.Permission) error
	AssignRoleToUser(db *gorm.DB, user *entities.User, role *entities.Role) error
	CreateNewCouncilor(db *gorm.DB, councilor *entities.Councilor) error
	CreateNewPermission(db *gorm.DB, permissionType enums.PermissionType, description string) *entities.Permission
	CreateNewRole(db *gorm.DB, roleType string) *entities.Role
	CreateNewUser(db *gorm.DB, use *entities.User) error
	DeleteCouncilor(db *gorm.DB, councilorID uint)
	DeleteRoleByRoleID(db *gorm.DB, roleID uint)
	DeleteRolePermission(db *gorm.DB, role *entities.Role, permission *entities.Permission)
	DeleteUserRole(db *gorm.DB, user *entities.User, role *entities.Role)
	FindActiveOrVerifiedUserByEmail(db *gorm.DB, email string) (*entities.User, bool)
	FindActiveOrVerifiedUserByUsername(db *gorm.DB, username string) (*entities.User, bool)
	FindAllCouncilorsByPromotedYear(db *gorm.DB, promotedYear int) []*entities.Councilor
	FindAllPermissions(db *gorm.DB) []*entities.Permission
	FindAllRolesWithPermissions(db *gorm.DB) []*entities.Role
	FindByEmailAndVerified(db *gorm.DB, email string, verified bool) (*entities.User, bool)
	FindByUserID(db *gorm.DB, userID uint) (*entities.User, bool)
	FindByUsernameAndVerified(db *gorm.DB, username string, verified bool) (*entities.User, bool)
	FindCouncilorByID(db *gorm.DB, councilorID uint) (*entities.Councilor, bool)
	FindCouncilorByUserIDAndPromotedYear(db *gorm.DB, userID uint, promotedYear int) (*entities.Councilor, bool)
	FindPermissionByID(db *gorm.DB, permissionID uint) (*entities.Permission, bool)
	FindPermissionByType(db *gorm.DB, permissionType enums.PermissionType) (*entities.Permission, bool)
	FindPermissionsByRole(db *gorm.DB, roleID uint) []enums.PermissionType
	FindRoleByID(db *gorm.DB, roleID uint) (*entities.Role, bool)
	FindRoleByType(db *gorm.DB, roleType string) (*entities.Role, bool)
	FindUnverifiedUsersWeekAgo(db *gorm.DB, startOfWeekAgo time.Time, endOfWeekAgo time.Time) []*entities.User
	FindUserRoleTypesByUserID(db *gorm.DB, userID uint) []entities.Role
	FindUsersByRoleID(db *gorm.DB, roleID uint) []*entities.User
	UpdateCouncilor(db *gorm.DB, councilor *entities.Councilor) error
	UpdateUserPassword(db *gorm.DB, user *entities.User, password string)
	UpdateUserToken(db *gorm.DB, user *entities.User, token string)
}
