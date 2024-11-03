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

func (repo *UserRepository) Test() []string {
	var tables []string
	repo.db.Raw("SHOW TABLES").Scan(&tables)

	return tables
}

func (repo *UserRepository) Test2() []entities.Test {
	var results []entities.Test
	repo.db.Where("name = ?", "ali").Find(&results)

	return results
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

func (repo *UserRepository) FindByUsernameAndVerified(username string, verified bool) (entities.User, bool) {
	var user entities.User
	result := repo.db.Where("name = ? AND verified = ?", username, verified).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
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

func (repo *UserRepository) FindRoleByType(roleType enums.RoleType) (entities.Role, bool) {
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

func (repo *UserRepository) CreateNewRole(roleType enums.RoleType) entities.Role {
	role := entities.Role{
		Type: roleType,
	}
	result := repo.db.Create(&role)
	if result.Error != nil {
		panic(result.Error)
	}
	return role
}

func (repo *UserRepository) AssignRoleToUser(user entities.User, role entities.Role) {
	err := repo.db.Model(&user).Association("Roles").Append(&role)
	if err != nil {
		panic(err)
	}
}

func (repo *UserRepository) FindUserRoleTypesByUserID(userID uint) []enums.RoleType {
	var user entities.User
	err := repo.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		panic(err)
	}
	roleTypes := make([]enums.RoleType, len(user.Roles))
	for i, role := range user.Roles {
		roleTypes[i] = role.Type
	}
	return roleTypes
}
