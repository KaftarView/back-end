package redis

import (
	"context"
	"encoding/json"
	"first-project/src/entities"
	"first-project/src/repository"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type UserRedis struct {
	rdb            *redis.Client
	userRepository *repository.UserRepository
}

func NewUserRedis(rdb *redis.Client, userRepository *repository.UserRepository) *UserRedis {
	return &UserRedis{
		rdb:            rdb,
		userRepository: userRepository,
	}
}

var ctx = context.Background()

func (userRedis *UserRedis) SetUser(user *entities.User) {
	key := "user:" + strconv.Itoa(int(user.ID))
	userData, err := json.Marshal(user)
	// TODO
	if err != nil {
		panic(err)
	}

	// TODO
	err = userRedis.rdb.Set(ctx, key, userData, 3600).Err()
	if err != nil {
		panic(err)
	}
}

func (userRedis *UserRedis) GetUser(userID uint) (entities.User, error) {
	key := "user:" + strconv.Itoa(int(userID))
	val, err := userRedis.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		user, userExist := userRedis.userRepository.FindByUserID(userID)
		if !userExist {
			// TODO
			// unauthorized error
			return user, err
		}
		userRedis.SetUser(&user)

		return user, nil

		// TODO
	} else if err != nil {
		panic(err)
	}

	// If the key exists, unmarshal the value
	var user entities.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		panic(err)
	}

	return user, nil
}
