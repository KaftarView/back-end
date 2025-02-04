package wire

import (
	"first-project/src/bootstrap"
	controller_v1_category "first-project/src/controller/v1/category"
	controller_v1_chat "first-project/src/controller/v1/chat"
	controller_v1_comment "first-project/src/controller/v1/comment"
	controller_v1_event "first-project/src/controller/v1/event"
	controller_v1_journal "first-project/src/controller/v1/journal"
	controller_v1_news "first-project/src/controller/v1/news"
	controller_v1_podcast "first-project/src/controller/v1/podcast"
	controller_v1_user "first-project/src/controller/v1/user"
	"first-project/src/websocket"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
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

type Application struct {
	Admin    *AdminControllers
	Customer *CustomerControllers
	General  *GeneralControllers
}

func InitializeApplication(container *bootstrap.Di, db *gorm.DB, rdb *redis.Client, hub *websocket.Hub) (*Application, error) {
	wire.Build(
		ProvideConstants,
		ProvideEmailInfo,
		ProvideStorage,
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
