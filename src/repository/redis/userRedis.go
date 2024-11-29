package repository_cache

import (
	"context"
	"encoding/json"
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/exceptions"
	repository_database "first-project/src/repository/database"

	"github.com/redis/go-redis/v9"
)

type UserCacheData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Roles []entities.Role
}

type UserCache struct {
	constants      *bootstrap.Constants
	rdb            *redis.Client
	userRepository *repository_database.UserRepository
}

func NewUserCache(
	constants *bootstrap.Constants,
	rdb *redis.Client,
	userRepository *repository_database.UserRepository,
) *UserCache {
	return &UserCache{
		constants:      constants,
		rdb:            rdb,
		userRepository: userRepository,
	}
}

var ctx = context.Background()

func (userCache *UserCache) SetUser(userID uint, username, email string) {
	key := userCache.constants.Redis.GetUserID(int(userID))
	roles := userCache.userRepository.FindUserRoleTypesByUserID(userID)
	userData := UserCacheData{
		Name:  username,
		Email: email,
		Roles: roles,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		panic(err)
	}

	// TODO set label for 3600 ttl of redis
	err = userCache.rdb.Set(ctx, key, userDataJSON, 3600).Err()
	if err != nil {
		panic(err)
	}
}

func (userCache *UserCache) GetUser(userID uint) UserCacheData {
	key := userCache.constants.Redis.GetUserID(int(userID))
	val, err := userCache.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		user, userExist := userCache.userRepository.FindByUserID(userID)
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
