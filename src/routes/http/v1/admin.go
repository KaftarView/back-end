package routes_http_v1

import (
	"first-project/src/enums"
	"first-project/src/wire"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	events := routerGroup.Group("/events")
	{
		readGroup := events.Group("")
		readGroup.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.ManageEvent, enums.ReviewEvent}))
		{
			readGroup.GET("", app.AdminControllers.EventController.GetEventsList)
			readGroup.GET("/search", app.AdminControllers.EventController.SearchEvents)
			readGroup.GET("/filter", app.AdminControllers.EventController.FilterEvents)

			readGroup.GET("ticket/:ticketID", app.AdminControllers.EventController.GetTicketDetails)
			readGroup.GET("discount/:discountID", app.AdminControllers.EventController.GetDiscountDetails)
			readGroup.GET("media/:mediaID", app.AdminControllers.EventController.GetMediaDetails)

			readSingleEventGroup := readGroup.Group("/:eventID")
			{
				readSingleEventGroup.GET("", app.AdminControllers.EventController.GetEventDetails)
				readSingleEventGroup.GET("/tickets", app.AdminControllers.EventController.GetAllTicketDetails)
				readSingleEventGroup.GET("/discounts", app.AdminControllers.EventController.GetAllDiscountDetails)
				readSingleEventGroup.GET("/media", app.AdminControllers.EventController.GetEventMedia)
				readSingleEventGroup.GET("/attendees", app.AdminControllers.EventController.GetEventAttendees)
			}
		}

		createGroup := events.Group("")
		createGroup.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.CreateEvent}))
		{
			createGroup.POST("/create", app.AdminControllers.EventController.CreateEvent)
			createGroup.POST("/add-ticket/:eventID", app.AdminControllers.EventController.AddEventTicket)
			createGroup.POST("/add-discount/:eventID", app.AdminControllers.EventController.AddEventDiscount)
			createGroup.POST("/add-organizer/:eventID", app.AdminControllers.EventController.AddEventOrganizer)
		}

		manageEventsGroup := events.Group("")
		manageEventsGroup.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.CreateEvent, enums.ManageEvent}))
		{
			eventSubGroup := events.Group("/:eventID")
			{
				eventSubGroup.PUT("", app.AdminControllers.EventController.UpdateEvent)
				eventSubGroup.POST("/publish", app.AdminControllers.EventController.PublishEvent)
				eventSubGroup.POST("/unpublish", app.AdminControllers.EventController.UnpublishEvent)
				eventSubGroup.DELETE("", app.AdminControllers.EventController.DeleteEvent)
				eventSubGroup.POST("/media", app.AdminControllers.EventController.UploadEventMedia)
			}

			ticketSubGroup := manageEventsGroup.Group("/ticket/:ticketID")
			{
				ticketSubGroup.PUT("", app.AdminControllers.EventController.UpdateEventTicket)
				ticketSubGroup.DELETE("", app.AdminControllers.EventController.DeleteTicket)
			}

			discountSubGroup := manageEventsGroup.Group("/discount/:discountID")
			{
				discountSubGroup.PUT("", app.AdminControllers.EventController.UpdateEventDiscount)
				discountSubGroup.DELETE("", app.AdminControllers.EventController.DeleteDiscount)
			}

			mediaSubGroup := manageEventsGroup.Group("/media/:mediaID")
			{
				mediaSubGroup.PUT("", app.AdminControllers.EventController.UpdateEventMedia)
				mediaSubGroup.DELETE("", app.AdminControllers.EventController.DeleteEventMedia)
			}

			manageEventsGroup.DELETE("/organizer/:organizerID", app.AdminControllers.EventController.DeleteOrganizer)
		}
	}

	comments := routerGroup.Group("/comments")
	comments.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ModerateComments}))
	{
		comments.DELETE("/:commentID", app.AdminControllers.CommentController.DeleteComment)
	}

	podcasts := routerGroup.Group("/podcasts")
	podcasts.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcasts.POST("", app.AdminControllers.PodcastController.CreatePodcast)
		podcastSubGroup := podcasts.Group("/:podcastID")
		{
			podcastSubGroup.PUT("", app.AdminControllers.PodcastController.UpdatePodcast)
			podcastSubGroup.DELETE("", app.AdminControllers.PodcastController.DeletePodcast)

			podcastEpisodesSubRouter := podcastSubGroup.Group("/episodes")
			{
				podcastEpisodesSubRouter.POST("", app.AdminControllers.PodcastController.CreateEpisode)
			}
		}
	}

	podcastEpisodes := routerGroup.Group("/episodes")
	podcastEpisodes.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManagePodcasts}))
	{
		podcastEpisodes.PUT("/:episodeID", app.AdminControllers.PodcastController.UpdateEpisode)
		podcastEpisodes.DELETE("/:episodeID", app.AdminControllers.PodcastController.DeleteEpisode)
	}

	accessManagement := routerGroup.Group("")
	accessManagement.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManageUsers, enums.ManageRoles}))
	{
		roles := accessManagement.Group("/roles")
		{
			roles.GET("", app.AdminControllers.UserController.GetRolesList)
			roles.POST("", app.AdminControllers.UserController.CreateRole)

			roleSubGroup := roles.Group("/:roleID")
			{
				roleSubGroup.GET("/owners", app.AdminControllers.UserController.GetRoleOwners)
				roleSubGroup.DELETE("", app.AdminControllers.UserController.DeleteRole)

				rolePermissions := roleSubGroup.Group("/permissions")
				{
					rolePermissions.PUT("", app.AdminControllers.UserController.UpdateRole)
					rolePermissions.DELETE("/:permissionID", app.AdminControllers.UserController.DeleteRolePermission)
				}
			}
		}

		accessManagement.GET("/permissions", app.AdminControllers.UserController.GetPermissionsList)

		userRoles := accessManagement.Group("/users/roles")
		{
			userRoles.PUT("", app.AdminControllers.UserController.UpdateUserRoles)
			userRoles.DELETE("/:roleID", app.AdminControllers.UserController.DeleteUserRole)
		}

		councilors := accessManagement.Group("/councilors")
		{
			councilors.POST("", app.AdminControllers.UserController.CreateCouncilor)
			councilors.DELETE("/:councilorID", app.AdminControllers.UserController.DeleteCouncilor)
		}
	}

	news := routerGroup.Group("/news")
	news.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManageNews}))
	{
		news.POST("", app.AdminControllers.NewsController.CreateNews)
		newsSubGroup := news.Group("/:newsID")
		{
			newsSubGroup.PUT("", app.AdminControllers.NewsController.UpdateNews)
			newsSubGroup.DELETE("", app.AdminControllers.NewsController.DeleteNews)
		}
	}

	journals := routerGroup.Group("/journal")
	journals.Use(app.Middlewares.Auth.RequirePermission([]enums.PermissionType{enums.ManageJournal}))
	{
		journals.POST("", app.AdminControllers.JournalController.CreateJournal)
		journalSubGroup := journals.Group("/:journalID")
		{
			journalSubGroup.PUT("", app.AdminControllers.JournalController.UpdateJournal)
			journalSubGroup.DELETE("", app.AdminControllers.JournalController.DeleteJournal)
		}
	}
}
