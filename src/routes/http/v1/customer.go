package routes_http_v1

import (
	"first-project/src/wire"

	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	event := routerGroup.Group("/events/:eventID")
	{
		event.GET("/attendance", app.CustomerControllers.EventController.IsUserAttended)
		event.GET("/tickets", app.CustomerControllers.EventController.GetAvailableEventTicketsList)
		event.GET("/media", app.CustomerControllers.EventController.GetEventMedia)
		event.POST("/reserve", app.CustomerControllers.EventController.ReserveTickets)
		event.POST("/purchase/:reservationID", app.CustomerControllers.EventController.PurchaseTickets)
	}

	comments := routerGroup.Group("/comments")
	{
		comments.POST("/post/:postID", app.CustomerControllers.CommentController.CreateComment)

		commentSubGroup := comments.Group("/:commentID")
		{
			commentSubGroup.PUT("", app.CustomerControllers.CommentController.EditComment)
			commentSubGroup.DELETE("", app.CustomerControllers.CommentController.DeleteComment)
		}
	}

	podcast := routerGroup.Group("/podcasts/:podcastID/subscribe")
	{
		podcast.POST("", app.CustomerControllers.PodcastController.SubscribePodcast)
		podcast.DELETE("", app.CustomerControllers.PodcastController.UnSubscribePodcast)
		podcast.GET("/status", app.CustomerControllers.PodcastController.SubscribeStatus)
	}

	profile := routerGroup.Group("/profile")
	{
		profile.PUT("/username", app.CustomerControllers.UserController.ChangeUsername)
		profile.PUT("/reset-password", app.CustomerControllers.UserController.ResetPassword)
		profile.GET("/events", app.CustomerControllers.EventController.GetAllUserJoinedEvents)
	}

	chat := routerGroup.Group("/chat")
	{
		chat.POST("/room", app.CustomerControllers.ChatController.CreateOrGetRoom)
		chat.GET("/room/:roomID/messages", app.CustomerControllers.ChatController.GetMessages)
	}
}
