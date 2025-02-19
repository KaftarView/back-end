package repository_cache

import (
	"context"
	"encoding/json"
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/exceptions"
	repository_database_interfaces "first-project/src/repository/database/interfaces"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserCacheData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Roles []entities.Role
}

type UserCache struct {
	constants      *bootstrap.Constants
	rdb            *redis.Client
	userRepository repository_database_interfaces.UserRepository
	db             *gorm.DB
}

func NewUserCache(
	db *gorm.DB,
	constants *bootstrap.Constants,
	rdb *redis.Client,
	userRepository repository_database_interfaces.UserRepository,
) *UserCache {
	return &UserCache{
		constants:      constants,
		rdb:            rdb,
		userRepository: userRepository,
		db:             db,
	}
}

var ctx = context.Background()

func (userCache *UserCache) SetUser(userID uint, username, email string) {
	key := userCache.constants.Redis.GetUserID(int(userID))
	roles := userCache.userRepository.FindUserRoleTypesByUserID(userCache.db, userID)
	userData := UserCacheData{
		Name:  username,
		Email: email,
		Roles: roles,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		panic(err)
	}

	err = userCache.rdb.Set(ctx, key, userDataJSON, 3600).Err()
	if err != nil {
		panic(err)
	}
}

func (userCache *UserCache) GetUser(userID uint) UserCacheData {
	key := userCache.constants.Redis.GetUserID(int(userID))
	val, err := userCache.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		user, userExist := userCache.userRepository.FindByUserID(userCache.db, userID)
		if !userExist {
			unauthorizedError := exceptions.NewUnauthorizedError()
			panic(unauthorizedError)
		}
		userCache.SetUser(userID, user.Name, user.Email)
		val, _ = userCache.rdb.Get(ctx, key).Result()
	} else if err != nil {
		panic(err)
	}

	var userData UserCacheData
	err = json.Unmarshal([]byte(val), &userData)
	if err != nil {
		panic(err)
	}
	return userData
}
