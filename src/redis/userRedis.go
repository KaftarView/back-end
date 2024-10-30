package cache

import (
	"context"
	"encoding/json"
	"first-project/src/enums"
	"first-project/src/exceptions"
	"first-project/src/repository"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type UserCacheData struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Roles []enums.RoleType
}

type UserCache struct {
	rdb            *redis.Client
	userRepository *repository.UserRepository
}

func NewUserCache(rdb *redis.Client, userRepository *repository.UserRepository) *UserCache {
	return &UserCache{
		rdb:            rdb,
		userRepository: userRepository,
	}
}

var ctx = context.Background()

func (userRedis *UserCache) SetUser(userID uint, username, email string) {
	key := "user:" + strconv.Itoa(int(userID))
	roles := userRedis.userRepository.FindUserRoleTypesByUserID(userID)
	userData := UserCacheData{
		Name:  username,
		Email: email,
		Roles: roles,
	}

	userDataJSON, err := json.Marshal(userData)
	if err != nil {
		panic(err)
	}

	// TODO set label for 3600 ttl of redis -> sync it with jwt
	err = userRedis.rdb.Set(ctx, key, userDataJSON, 3600).Err()
	if err != nil {
		panic(err)
	}
}

func (userRedis *UserCache) GetUser(userID uint) UserCacheData {
	key := "user:" + strconv.Itoa(int(userID))
	val, err := userRedis.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		user, userExist := userRedis.userRepository.FindByUserID(userID)
		if !userExist {
			unauthorizedError := exceptions.NewUnauthorizedError()
			panic(unauthorizedError)
		}
		userRedis.SetUser(userID, user.Name, user.Email)
		val, _ = userRedis.rdb.Get(ctx, key).Result()
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
