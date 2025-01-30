package wire

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	controller_v1_user "first-project/src/controller/v1/user"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	repository_cache "first-project/src/repository/redis"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var UserRepositorySet = wire.NewSet(
	repository_database.NewUserRepository,
	wire.Bind(new(repository_database_interfaces.UserRepository), new(*repository_database.UserRepository)),
)

var UserServiceSet = wire.NewSet(
	application.NewUserService,
	wire.Bind(new(application_interfaces.UserService), new(*application.UserService)),
)

func provideConstants(container *bootstrap.Di) *bootstrap.Constants {
	return container.Constants
}

func provideEmailInfo(container *bootstrap.Di) *bootstrap.EmailInfo {
	return &container.Env.Email
}

func provideStorage(container *bootstrap.Di) *bootstrap.S3 {
	return &container.Env.Storage
}

func InitializeUserRouterImpl(container *bootstrap.Di, db *gorm.DB, rdb *redis.Client) *controller_v1_user.GeneralUserController {
	wire.Build(
		provideConstants,
		provideEmailInfo,
		provideStorage,
		UserRepositorySet,
		UserServiceSet,
		application_aws.NewS3Service,
		application_communication.NewEmailService,
		repository_cache.NewUserCache,
		application.NewOTPService,
		application_jwt.NewJWTToken,
		controller_v1_user.NewGeneralUserController,
	)

	return nil
}
