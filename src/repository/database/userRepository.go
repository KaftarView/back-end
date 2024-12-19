package repository_database

import (
	"first-project/src/entities"
	"first-project/src/enums"
	"time"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) FindByUserID(userID uint) (entities.User, bool) {
	var user entities.User
	result := repo.db.First(&user, userID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, false
		}
		panic(result.Error)
	}
	return user, true
}

func (repo *UserRepository) FindActiveOrVerifiedUserByUsername(username string) (entities.User, bool) {
	var user entities.User
	result := repo.db.Where("name = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return user, true
	}
	repo.db.Delete(&user)
	return user, false
}

func (repo *UserRepository) FindActiveOrVerifiedUserByEmail(email string) (entities.User, bool) {
	var user entities.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, false
		}
		panic(result.Error)
	}
	if user.Verified || time.Since(user.UpdatedAt) < 2*time.Minute {
		return user, true
	}
	repo.db.Delete(&user)
	return user, false
}

func (repo *UserRepository) FindByUsernameAndVerified(username string, verified bool) (entities.User, bool) {
	var user entities.User
	result := repo.db.Where("name = ? AND verified = ?", username, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			if time.Since(user.UpdatedAt) < 2*time.Minute {
				return user, true
			}
			return user, false
		}
		panic(result.Error)
	}
	return user, true
}

func (repo *UserRepository) FindByEmailAndVerified(email string, verified bool) (entities.User, bool) {
	var user entities.User
	result := repo.db.Where("email = ? AND verified = ?", email, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return user, false
		}
		panic(result.Error)
	}
	return user, true
}

func (repo *UserRepository) UpdateUserToken(user entities.User, token string) {
	user.Token = token
	repo.db.Save(&user)
}

func (repo *UserRepository) CreateNewUser(
	username string, email string, password string, token string, verified bool) entities.User {
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
	return user
}

func (repo *UserRepository) ActivateUserAccount(user entities.User) {
	user.Verified = true
	user.Token = ""
	if err := repo.db.Save(&user).Error; err != nil {
		panic(err)
	}
}

func (repo *UserRepository) UpdateUserPassword(user entities.User, password string) {
	user.Password = password
	user.Token = ""
	repo.db.Save(&user)
}

func (repo *UserRepository) FindUnverifiedUsersWeekAgo(startOfWeekAgo, endOfWeekAgo time.Time) []entities.User {
	var users []entities.User
	err := repo.db.Where(
		"verified = ? AND created_at >= ? AND created_at < ?",
		false, startOfWeekAgo, endOfWeekAgo).Find(&users).Error
	if err != nil {
		panic(err)
	}
	return users
}

func (repo *UserRepository) FindRoleByType(roleType string) (entities.Role, bool) {
	var role entities.Role
	result := repo.db.Where("type = ?", roleType).First(&role)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return role, false
		}
		panic(result.Error)
	}
	return role, true
}

func (repo *UserRepository) FindRoleByID(roleID uint) (entities.Role, bool) {
	var role entities.Role
	result := repo.db.First(&role, roleID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return role, false
		}
		panic(result.Error)
	}
	return role, true
}

func (repo *UserRepository) FindPermissionByType(permissionType enums.PermissionType) (entities.Permission, bool) {
	var permission entities.Permission
	result := repo.db.Where("type = ?", permissionType).First(&permission)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return permission, false
		}
		panic(result.Error)
	}
	return permission, true
}

func (repo *UserRepository) CreateNewRole(roleType string) entities.Role {
	role := entities.Role{
		Type: roleType,
	}
	result := repo.db.Create(&role)
	if result.Error != nil {
		panic(result.Error)
	}
	return role
}

func (repo *UserRepository) CreateNewPermission(permissionType enums.PermissionType) entities.Permission {
	permission := entities.Permission{
		Type: permissionType,
	}
	result := repo.db.Create(&permission)
	if result.Error != nil {
		panic(result.Error)
	}
	return permission
}

// TODO: it is wrong i think!!
func (repo *UserRepository) AssignRoleToUser(user entities.User, role entities.Role) {
	exists := repo.db.Model(&user).
		Where("id = ?", role.ID).
		Association("Roles").
		Count() > 0
	if exists {
		return
	}
	if err := repo.db.Model(&user).Association("Roles").Append(&role); err != nil {
		panic(err)
	}
}

func (repo *UserRepository) AssignPermissionToRole(role entities.Role, permission entities.Permission) {
	err := repo.db.Model(&role).Association("Permissions").Append(&permission)
	if err != nil {
		panic(err)
	}
}

func (repo *UserRepository) FindUserRoleTypesByUserID(userID uint) []entities.Role {
	var user entities.User
	err := repo.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		panic(err)
	}
	return user.Roles
}

func (repo *UserRepository) FindPermissionsByRole(roleID uint) []enums.PermissionType {
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

func (repo *UserRepository) FindAllRolesWithPermissions() []entities.Role {
	var roles []entities.Role
	result := repo.db.Preload("Permissions").Find(&roles)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []entities.Role{}
		}
		panic(result.Error)
	}
	return roles
}

func (repo *UserRepository) FindUsersByRoleID(roleID uint) []entities.User {
	var users []entities.User
	result := repo.db.
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return []entities.User{}
		}
		panic(result.Error)
	}
	return users
}

func (repo *UserRepository) DeleteRoleByRoleID(roleID uint) {
	repo.db.Exec("DELETE FROM user_roles WHERE role_id = ?", roleID)
	repo.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID)

	if err := repo.db.Delete(&entities.Role{}, roleID).Error; err != nil {
		panic(err)
	}
}
