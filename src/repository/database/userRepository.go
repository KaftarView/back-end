package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (repo *userRepository) FindByUserID(userID uint) (*entities.User, bool) {
	var user entities.User
	result := repo.db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *userRepository) FindActiveOrVerifiedUserByUsername(username string) (*entities.User, bool) {
	var user entities.User
	result := repo.db.Where("name = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return &user, true
	}
	repo.db.Delete(&user)
	return nil, false
}

func (repo *userRepository) FindActiveOrVerifiedUserByEmail(email string) (*entities.User, bool) {
	var user entities.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return &user, true
	}
	repo.db.Delete(&user)
	return nil, false
}

func (repo *userRepository) FindByUsernameAndVerified(username string, verified bool) (*entities.User, bool) {
	var user entities.User
	result := repo.db.Where("name = ? AND verified = ?", username, verified).First(&user)
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

func (repo *userRepository) FindByEmailAndVerified(email string, verified bool) (*entities.User, bool) {
	var user entities.User
	result := repo.db.Where("email = ? AND verified = ?", email, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &user, true
}

func (repo *userRepository) UpdateUserToken(user *entities.User, token string) {
	user.Token = token
	repo.db.Save(user)
}

func (repo *userRepository) CreateNewUser(
	username string, email string, password string, token string, verified bool) *entities.User {
	user := entities.User{
		Name:     username,
		Email:    email,
		Password: password,
		Token:    token,
		Verified: verified,
	}
	result := repo.db.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return &user
}

func (repo *userRepository) ActivateUserAccount(user *entities.User) {
	user.Verified = true
	user.Token = ""
	if err := repo.db.Save(user).Error; err != nil {
		panic(err)
	}
}

func (repo *userRepository) UpdateUserPassword(user *entities.User, password string) {
	user.Password = password
	user.Token = ""
	repo.db.Save(user)
}

func (repo *userRepository) FindUnverifiedUsersWeekAgo(startOfWeekAgo, endOfWeekAgo time.Time) []*entities.User {
	var users []*entities.User
	err := repo.db.Where(
		"verified = ? AND created_at >= ? AND created_at < ?",
		false, startOfWeekAgo, endOfWeekAgo).Find(&users).Error
	if err != nil {
		panic(err)
	}
	return users
}

func (repo *userRepository) FindRoleByType(roleType string) (*entities.Role, bool) {
	var role entities.Role
	result := repo.db.Where("type = ?", roleType).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *userRepository) FindRoleByID(roleID uint) (*entities.Role, bool) {
	var role entities.Role
	result := repo.db.First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &role, true
}

func (repo *userRepository) FindPermissionByID(permissionID uint) (*entities.Permission, bool) {
	var permission entities.Permission
	result := repo.db.First(&permission, permissionID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *userRepository) FindPermissionByType(permissionType enums.PermissionType) (*entities.Permission, bool) {
	var permission entities.Permission
	result := repo.db.Where("type = ?", permissionType).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &permission, true
}

func (repo *userRepository) CreateNewRole(roleType string) *entities.Role {
	role := entities.Role{
		Type: roleType,
	}
	result := repo.db.Create(&role)
	if result.Error != nil {
		panic(result.Error)
	}
	return &role
}

func (repo *userRepository) CreateNewPermission(permissionType enums.PermissionType) *entities.Permission {
	permission := entities.Permission{
		Type: permissionType,
	}
	result := repo.db.Create(&permission)
	if result.Error != nil {
		panic(result.Error)
	}
	return &permission
}

func (repo *userRepository) AssignRoleToUser(user *entities.User, role *entities.Role) {
	if err := repo.db.Model(user).Association("Roles").Append(role); err != nil {
		panic(err)
	}
}

func (repo *userRepository) AssignPermissionToRole(role *entities.Role, permission *entities.Permission) {
	if err := repo.db.Model(role).Association("Permissions").Append(permission); err != nil {
		panic(err)
	}
}

func (repo *userRepository) FindUserRoleTypesByUserID(userID uint) []entities.Role {
	var user entities.User
	err := repo.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		panic(err)
	}
	return user.Roles
}

func (repo *userRepository) FindPermissionsByRole(roleID uint) []enums.PermissionType {
	var role entities.Role
	err := repo.db.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		panic(err)
	}
	permissionTypes := make([]enums.PermissionType, len(role.Permissions))
	for i, permission := range role.Permissions {
		permissionTypes[i] = permission.Type
	}
	return permissionTypes
}

func (repo *userRepository) FindAllRolesWithPermissions() []*entities.Role {
	var roles []*entities.Role
	result := repo.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return roles
}

func (repo *userRepository) FindUsersByRoleID(roleID uint) []*entities.User {
	var users []*entities.User
	result := repo.db.
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

func (repo *userRepository) DeleteRoleByRoleID(roleID uint) {
	repo.db.Exec("DELETE FROM user_roles WHERE role_id = ?", roleID)
	repo.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID)

	if err := repo.db.Delete(&entities.Role{}, roleID).Error; err != nil {
		panic(err)
	}
}

func (repo *userRepository) DeleteRolePermission(role *entities.Role, permission *entities.Permission) {
	if err := repo.db.Model(role).Association("Permissions").Delete(permission); err != nil {
		panic(err)
	}
}

func (repo *userRepository) DeleteUserRole(user *entities.User, role *entities.Role) {
	if err := repo.db.Model(user).Association("Roles").Delete(role); err != nil {
		panic(err)
	}
}

func (repo *userRepository) FindAllPermissions() []*entities.Permission {
	var permissions []*entities.Permission
	result := repo.db.Find(&permissions)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		}
		panic(result.Error)
	}
	return permissions
}

func (repo *userRepository) FindCouncilorByID(councilorID uint) (*entities.Councilor, bool) {
	var councilor entities.Councilor
	result := repo.db.First(&councilor, councilorID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &councilor, true
}

func (repo *userRepository) FindCouncilorByUserIDAndPromoteDate(userID uint, promotedDate time.Time) (*entities.Councilor, bool) {
	var councilor entities.Councilor
	result := repo.db.Where("user_id = ? AND promoted_date = ?", userID, promotedDate).First(&councilor)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, false
		}
		panic(result.Error)
	}
	return &councilor, true
}

func (repo *userRepository) CreateNewCouncilor(councilor *entities.Councilor) {
	result := repo.db.Create(&councilor)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (repo *userRepository) UpdateCouncilor(councilor *entities.Councilor) {
	if err := repo.db.Save(councilor).Error; err != nil {
		panic(err)
	}
}

func (repo *userRepository) DeleteCouncilor(councilorID uint) {
	err := repo.db.Unscoped().Delete(&entities.Councilor{}, councilorID).Error
	if err != nil {
		panic(err)
	}
}
