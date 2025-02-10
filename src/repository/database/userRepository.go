package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindByUserID(db *gorm.DB, userID uint) (*entities.User, bool) {
	var user entities.User
	result := db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *UserRepository) FindActiveOrVerifiedUserByUsername(db *gorm.DB, username string) (*entities.User, bool) {
	var user entities.User
	result := db.Where("name = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return &user, true
	}
	db.Delete(&user)
	return nil, false
}

func (repo *UserRepository) FindActiveOrVerifiedUserByEmail(db *gorm.DB, email string) (*entities.User, bool) {
	var user entities.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return &user, true
	}
	db.Delete(&user)
	return nil, false
}

func (repo *UserRepository) FindByUsernameAndVerified(db *gorm.DB, username string, verified bool) (*entities.User, bool) {
	var user entities.User
	result := db.Where("name = ? AND verified = ?", username, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if time.Since(user.UpdatedAt) < 2*time.Minute {
				return &user, true
			}
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *UserRepository) FindByEmailAndVerified(db *gorm.DB, email string, verified bool) (*entities.User, bool) {
	var user entities.User
	result := db.Where("email = ? AND verified = ?", email, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *UserRepository) UpdateUser(db *gorm.DB, user *entities.User) error {
	return db.Save(user).Error
}

func (repo *UserRepository) UpdateUserToken(db *gorm.DB, user *entities.User, token string) {
	user.Token = token
	db.Save(user)
}

func (repo *UserRepository) CreateNewUser(db *gorm.DB, user *entities.User) error {
	return db.Create(&user).Error
}

func (repo *UserRepository) ActivateUserAccount(db *gorm.DB, user *entities.User) {
	user.Verified = true
	user.Token = ""
	if err := db.Save(user).Error; err != nil {
		panic(err)
	}
}

func (repo *UserRepository) UpdateUserPassword(db *gorm.DB, user *entities.User, password string) {
	user.Password = password
	user.Token = ""
	db.Save(user)
}

func (repo *UserRepository) FindUnverifiedUsersWeekAgo(db *gorm.DB, startOfWeekAgo, endOfWeekAgo time.Time) []*entities.User {
	var users []*entities.User
	err := db.Where(
		"verified = ? AND created_at >= ? AND created_at < ?",
		false, startOfWeekAgo, endOfWeekAgo).Find(&users).Error
	if err != nil {
		panic(err)
	}
	return users
}

func (repo *UserRepository) FindRoleByType(db *gorm.DB, roleType string) (*entities.Role, bool) {
	var role entities.Role
	result := db.Where("type = ?", roleType).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *UserRepository) FindRoleByID(db *gorm.DB, roleID uint) (*entities.Role, bool) {
	var role entities.Role
	result := db.First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *UserRepository) FindPermissionByID(db *gorm.DB, permissionID uint) (*entities.Permission, bool) {
	var permission entities.Permission
	result := db.First(&permission, permissionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *UserRepository) FindPermissionByType(db *gorm.DB, permissionType enums.PermissionType) (*entities.Permission, bool) {
	var permission entities.Permission
	result := db.Where("type = ?", permissionType).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *UserRepository) CreateNewRole(db *gorm.DB, roleType string) *entities.Role {
	role := entities.Role{
		Type: roleType,
	}
	result := db.Create(&role)
	if result.Error != nil {
		panic(result.Error)
	}
	return &role
}

func (repo *UserRepository) CreateNewPermission(db *gorm.DB, permissionType enums.PermissionType, description string) *entities.Permission {
	permission := entities.Permission{
		Type:        permissionType,
		Description: description,
	}
	result := db.Create(&permission)
	if result.Error != nil {
		panic(result.Error)
	}
	return &permission
}

func (repo *UserRepository) AssignRoleToUser(db *gorm.DB, user *entities.User, role *entities.Role) error {
	return db.Model(user).Association("Roles").Append(role)
}

func (repo *UserRepository) AssignPermissionToRole(db *gorm.DB, role *entities.Role, permission *entities.Permission) error {
	return db.Model(role).Association("Permissions").Append(permission)
}

func (repo *UserRepository) FindUserRoleTypesByUserID(db *gorm.DB, userID uint) []entities.Role {
	var user entities.User
	err := db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		panic(err)
	}
	return user.Roles
}

func (repo *UserRepository) FindPermissionsByRole(db *gorm.DB, roleID uint) []enums.PermissionType {
	var role entities.Role
	err := db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		panic(err)
	}
	permissionTypes := make([]enums.PermissionType, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissionTypes[i] = permission.Type
	}
	return permissionTypes
}

func (repo *UserRepository) FindAllRolesWithPermissions(db *gorm.DB) []*entities.Role {
	var roles []*entities.Role
	result := db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return roles
}

func (repo *UserRepository) FindUsersByRoleID(db *gorm.DB, roleID uint) []*entities.User {
	var users []*entities.User
	result := db.
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return users
}

func (repo *UserRepository) DeleteRoleByRoleID(db *gorm.DB, roleID uint) {
	db.Exec("DELETE FROM user_roles WHERE role_id = ?", roleID)
	db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID)

	if err := db.Delete(&entities.Role{}, roleID).Error; err != nil {
		panic(err)
	}
}

func (repo *UserRepository) DeleteRolePermission(db *gorm.DB, role *entities.Role, permission *entities.Permission) {
	if err := db.Model(role).Association("Permissions").Delete(permission); err != nil {
		panic(err)
	}
}

func (repo *UserRepository) DeleteUserRole(db *gorm.DB, user *entities.User, role *entities.Role) {
	if err := db.Model(user).Association("Roles").Delete(role); err != nil {
		panic(err)
	}
}

func (repo *UserRepository) FindAllPermissions(db *gorm.DB) []*entities.Permission {
	var permissions []*entities.Permission
	result := db.Find(&permissions)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return permissions
}

func (repo *UserRepository) FindCouncilorByID(db *gorm.DB, councilorID uint) (*entities.Councilor, bool) {
	var councilor entities.Councilor
	result := db.First(&councilor, councilorID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &councilor, true
}

func (repo *UserRepository) FindCouncilorByUserIDAndPromotedYear(db *gorm.DB, userID uint, promotedYear int) (*entities.Councilor, bool) {
	var councilor entities.Councilor
	result := db.Where("user_id = ? AND promoted_year = ?", userID, promotedYear).First(&councilor)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &councilor, true
}

func (repo *UserRepository) FindAllCouncilorsByPromotedYear(db *gorm.DB, promotedYear int) []*entities.Councilor {
	var councilors []*entities.Councilor
	result := OrderByCreatedAtDesc(db).Where("promoted_year  = ?", promotedYear).Find(&councilors)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return councilors
}

func (repo *UserRepository) CreateNewCouncilor(db *gorm.DB, councilor *entities.Councilor) error {
	return db.Create(&councilor).Error
}

func (repo *UserRepository) UpdateCouncilor(db *gorm.DB, councilor *entities.Councilor) error {
	return db.Save(councilor).Error
}

func (repo *UserRepository) DeleteCouncilor(db *gorm.DB, councilorID uint) {
	err := db.Unscoped().Delete(&entities.Councilor{}, councilorID).Error
	if err != nil {
		panic(err)
	}
}

func (repo *UserRepository) FindUsersByPermissions(db *gorm.DB, permissions []enums.PermissionType) []entities.User {
	var users []entities.User
	result := db.Joins("JOIN user_roles ur ON ur.user_id = users.id").
		Joins("JOIN role_permissions rp ON rp.role_id = ur.role_id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("p.type IN ?", permissions).
		Find(&users)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return users
}
