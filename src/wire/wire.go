//go:build wireinject
// +build wireinject

package wire

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_cron "first-project/src/application/cron"
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	controller_v1_category "first-project/src/controller/v1/category"
	controller_v1_chat "first-project/src/controller/v1/chat"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_podcast "first-project/src/controller/v1/podcast"
	controller_v1_user "first-project/src/controller/v1/user"
	middleware_authentication "first-project/src/middleware/Authentication"
	middleware_exceptions "first-project/src/middleware/exceptions"
	middleware_i18n "first-project/src/middleware/i18n"
	middleware_rate_limit "first-project/src/middleware/rateLimit"
	middleware_websocket "first-project/src/middleware/websocket"
	repository_database "first-project/src/repository/database"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	repository_cache "first-project/src/repository/redis"
	"first-project/src/seed"
	"first-project/src/websocket"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var DatabaseProviderSet = wire.NewSet(
	repository_database.NewCategoryRepository,
	repository_database.NewChatRepository,
	repository_database.NewCommentRepository,
	repository_database.NewEventRepository,
	repository_database.NewJournalRepository,
	repository_database.NewNewsRepository,
	repository_database.NewPodcastRepository,
	repository_database.NewPurchaseRepository,
	repository_database.NewUserRepository,
	wire.Bind(new(repository_database_interfaces.CategoryRepository), new(*repository_database.CategoryRepository)),
	wire.Bind(new(repository_database_interfaces.ChatRepository), new(*repository_database.ChatRepository)),
	wire.Bind(new(repository_database_interfaces.CommentRepository), new(*repository_database.CommentRepository)),
	wire.Bind(new(repository_database_interfaces.EventRepository), new(*repository_database.EventRepository)),
	wire.Bind(new(repository_database_interfaces.JournalRepository), new(*repository_database.JournalRepository)),
	wire.Bind(new(repository_database_interfaces.NewsRepository), new(*repository_database.NewsRepository)),
	wire.Bind(new(repository_database_interfaces.PodcastRepository), new(*repository_database.PodcastRepository)),
	wire.Bind(new(repository_database_interfaces.PurchaseRepository), new(*repository_database.PurchaseRepository)),
	wire.Bind(new(repository_database_interfaces.UserRepository), new(*repository_database.UserRepository)),
)

var RedisProviderSet = wire.NewSet(
	repository_cache.NewUserCache,
)

// no cron here!
var ServiceProviderSet = wire.NewSet(
	application.NewS3Service,
	application.NewCategoryService,
	application.NewChatService,
	application.NewCommentService,
	application.NewEventService,
	application.NewJournalService,
	application.NewJWTToken,
	application.NewNewsService,
	application.NewOTPService,
	application.NewPodcastService,
	application.NewUserService,
	application_communication.NewEmailService,
	wire.Bind(new(application_interfaces.S3Service), new(*application.S3Service)),
	wire.Bind(new(application_interfaces.CategoryService), new(*application.CategoryService)),
	wire.Bind(new(application_interfaces.ChatService), new(*application.ChatService)),
	wire.Bind(new(application_interfaces.CommentService), new(*application.CommentService)),
	wire.Bind(new(application_interfaces.EventService), new(*application.EventService)),
	wire.Bind(new(application_interfaces.JournalService), new(*application.JournalService)),
	wire.Bind(new(application_interfaces.JWTToken), new(*application.JWTToken)),
	wire.Bind(new(application_interfaces.NewsService), new(*application.NewsService)),
	wire.Bind(new(application_interfaces.OTPService), new(*application.OTPService)),
	wire.Bind(new(application_interfaces.PodcastService), new(*application.PodcastService)),
	wire.Bind(new(application_interfaces.UserService), new(*application.UserService)),
	wire.Bind(new(application_interfaces.EmailService), new(*application_communication.EmailService)),
)

var AdminControllerProviderSet = wire.NewSet(
	controller_v1_comment.NewAdminCommentController,
	controller_v1_event.NewAdminEventController,
	controller_v1_journal.NewAdminJournalController,
	controller_v1_news.NewAdminNewsController,
	controller_v1_podcast.NewAdminPodcastController,
	controller_v1_user.NewAdminUserController,
	wire.Struct(new(AdminControllers), "*"),
)

var CustomerControllerProviderSet = wire.NewSet(
	controller_v1_chat.NewCustomerChatController,
	controller_v1_comment.NewCustomerCommentController,
	controller_v1_event.NewCustomerEventController,
	controller_v1_podcast.NewCustomerPodcastController,
	controller_v1_user.NewCustomerUserController,
	wire.Struct(new(CustomerControllers), "*"),
)

var GeneralControllerProviderSet = wire.NewSet(
	controller_v1_category.NewGeneralCategoryController,
	controller_v1_comment.NewGeneralCommentController,
	controller_v1_event.NewGeneralEventController,
	controller_v1_journal.NewGeneralJournalController,
	controller_v1_news.NewGeneralNewsController,
	controller_v1_podcast.NewGeneralPodcastController,
	controller_v1_user.NewGeneralUserController,
	wire.Struct(new(GeneralControllers), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	middleware_authentication.NewAuthMiddleware,
	middleware_exceptions.NewRecovery,
	middleware_i18n.NewLocalization,
	middleware_rate_limit.NewRateLimit,
	middleware_websocket.NewWebsocketMiddleware,
	wire.Struct(new(Middlewares), "*"),
)

var SeederCronJobSet = wire.NewSet(
	application_cron.NewCronJob,
	wire.Struct(new(CronJobs), "*"),
)

var SeederProviderSet = wire.NewSet(
	seed.NewRoleSeeder,
	wire.Struct(new(Seeders), "*"),
)

func ProvideConstants(container *bootstrap.Di) *bootstrap.Constants {
	return container.Constants
}

func ProvideEmailInfo(container *bootstrap.Di) *bootstrap.EmailInfo {
	return &container.Env.Email
}

func ProvideStorage(container *bootstrap.Di) *bootstrap.S3 {
	return &container.Env.Storage
}

func ProvideSuperAdminInfo(container *bootstrap.Di) *bootstrap.AdminCredentials {
	return &container.Env.SuperAdmin
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	RedisProviderSet,
	ServiceProviderSet,
	AdminControllerProviderSet,
	CustomerControllerProviderSet,
	GeneralControllerProviderSet,
	MiddlewareProviderSet,
	SeederCronJobSet,
	SeederProviderSet,
)

type AdminControllers struct {
	CommentController *controller_v1_comment.AdminCommentController
	EventController   *controller_v1_event.AdminEventController
	JournalController *controller_v1_journal.AdminJournalController
	NewsController    *controller_v1_news.AdminNewsController
	PodcastController *controller_v1_podcast.AdminPodcastController
	UserController    *controller_v1_user.AdminUserController
}

type CustomerControllers struct {
	ChatController    *controller_v1_chat.CustomerChatController
	CommentController *controller_v1_comment.CustomerCommentController
	EventController   *controller_v1_event.CustomerEventController
	PodcastController *controller_v1_podcast.CustomerPodcastController
	UserController    *controller_v1_user.CustomerUserController
}

type GeneralControllers struct {
	CategoryController *controller_v1_category.GeneralCategoryController
	CommentController  *controller_v1_comment.GeneralCommentController
	EventController    *controller_v1_event.GeneralEventController
	JournalController  *controller_v1_journal.GeneralJournalController
	NewsController     *controller_v1_news.GeneralNewsController
	PodcastController  *controller_v1_podcast.GeneralPodcastController
	UserController     *controller_v1_user.GeneralUserController
}

type Middlewares struct {
	Auth         *middleware_authentication.AuthMiddleware
	Recovery     *middleware_exceptions.RecoveryMiddleware
	Localization *middleware_i18n.LocalizationMiddleware
	RateLimit    *middleware_rate_limit.RateLimitMiddleware
	Websocket    *middleware_websocket.WebsocketMiddleware
}

type CronJobs struct {
	CronJob *application_cron.CronJob
}

type Seeders struct {
	RoleSeeder *seed.RoleSeeder
}

type Application struct {
	AdminControllers    *AdminControllers
	CustomerControllers *CustomerControllers
	GeneralControllers  *GeneralControllers
	Middlewares         *Middlewares
	CronJobs            *CronJobs
	Seeders             *Seeders
}

func InitializeApplication(container *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) (*Application, error) {
	wire.Build(
		ProvideConstants,
		ProvideEmailInfo,
		ProvideStorage,
		ProvideSuperAdminInfo,
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
