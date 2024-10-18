package repository

import (
	"first-project/src/entities"
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

func (repo *UserRepository) CheckUsernameExists(username string) bool {
	var user entities.User
	result := repo.db.Where("name = ? AND verified = ?", username, true).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(result.Error)
	}

	return true
}

func (repo *UserRepository) CheckEmailExists(email string) bool {
	var user entities.User
	result := repo.db.Where("email = ? AND verified = ?", email, true).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false
		}
		panic(result.Error)
	}

	return true
}

func (repo *UserRepository) registerNewUser(username string, email string, password string, otp string) {
	user := entities.User{
		Name:     username,
		Email:    email,
		Password: password,
		Token:    otp,
		Verified: false,
	}
	result := repo.db.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
}

func (repo *UserRepository) updateUserToken(user entities.User, token string) {
	user.Token = token
	repo.db.Save(&user)
}

func (repo *UserRepository) ForgotPassword(email string, token string) bool {
	var user entities.User

	result := repo.db.Where("email = ? AND verified = ?", email, true).First(&user)
	if result.Error == nil {
		repo.updateUserToken(user, token)
	} else if result.Error == gorm.ErrRecordNotFound {
		return false
	} else {
		panic(result.Error)
	}
	return true
}

func (repo *UserRepository) RegisterUser(username string, email string, password string, otp string) {
	var user entities.User

	result := repo.db.Where("email = ? AND verified = ?", email, false).First(&user)
	if result.Error == nil {
		repo.updateUserToken(user, otp)
	} else if result.Error == gorm.ErrRecordNotFound {
		repo.registerNewUser(username, email, password, otp)
	} else {
		panic(result.Error)
	}
}

func (repo *UserRepository) GetOTPByEmail(email string) (string, time.Time) {
	var user entities.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	return user.Token, user.UpdatedAt
}

func (repo *UserRepository) VerifyEmail(email string) {
	var user entities.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	user.Verified = true
	user.Token = ""
	if err := repo.db.Save(&user).Error; err != nil {
		panic(err)
	}
}

func (repo *UserRepository) GetPasswordByVerifiedUsername(username string) (string, error) {
	var user entities.User
	result := repo.db.Where("name = ? AND verified = ?", username, true).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", result.Error
		}
		panic(result.Error)
	}

	return user.Password, nil
}

func (repo *UserRepository) UpdatePasswordByEmail(email, password string) {
	var user entities.User
	result := repo.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	user.Password = password
	user.Token = ""
	repo.db.Save(&user)
}
