package routes_http_v1

import (
	"github.com/gin-gonic/gin"

	"first-project/src/wire"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	const (
		searchEndpoint = "/search"
		filterEndpoint = "/filter"
	)

	public := routerGroup.Group("/public")
	{
		categories := public.Group("/categories")
		{
			categories.GET("", app.GeneralControllers.CategoryController.GetListCategoryNames)
		}

		events := public.Group("/events")
		{
			events.GET("/published", app.GeneralControllers.EventController.ListEvents)
			events.GET(searchEndpoint, app.GeneralControllers.EventController.SearchEvents)
			events.GET(filterEndpoint, app.GeneralControllers.EventController.FilterEvents)
			events.GET("/:eventID", app.GeneralControllers.EventController.GetEventDetails)
			events.GET("/:eventID/organizers", app.GeneralControllers.EventController.GetEventOrganizers)
		}

		podcasts := public.Group("/podcasts")
		{
			podcasts.GET("", app.GeneralControllers.PodcastController.GetPodcastsList)
			podcasts.GET(searchEndpoint, app.GeneralControllers.PodcastController.SearchPodcast)
			podcasts.GET(filterEndpoint, app.GeneralControllers.PodcastController.FilterPodcastByCategory)
			podcasts.GET("/:podcastID/episodes", app.GeneralControllers.PodcastController.GetEpisodesList)
			podcasts.GET("/:podcastID", app.GeneralControllers.PodcastController.GetPodcastDetails)
		}

		episodes := public.Group("/episodes")
		{
			episodes.GET("/:episodeID", app.GeneralControllers.PodcastController.GetEpisodeDetails)
		}

		comments := public.Group("/comments/:postID")
		{
			comments.GET("", app.GeneralControllers.CommentController.GetComments)
		}

		news := public.Group("/news")
		{
			news.GET("", app.GeneralControllers.NewsController.GetNewsList)
			news.GET("/:newsID", app.GeneralControllers.NewsController.GetNewsDetails)
			news.GET(searchEndpoint, app.GeneralControllers.NewsController.SearchNews)
			news.GET(filterEndpoint, app.GeneralControllers.NewsController.FilterNewsByCategory)
		}

		journals := public.Group("/journals")
		{
			journals.GET("", app.GeneralControllers.JournalController.GetJournalsList)
			journals.GET(searchEndpoint, app.GeneralControllers.JournalController.SearchJournals)
		}

		councilors := public.Group("/councilors")
		{
			councilors.GET("", app.GeneralControllers.UserController.GetCouncilors)
		}
	}

	auth := routerGroup.Group("/auth")
	{
		auth.POST("/register", app.GeneralControllers.UserController.Register)
		auth.POST("/register/verify", app.GeneralControllers.UserController.VerifyEmail)
		auth.POST("/login", app.GeneralControllers.UserController.Login)
		auth.POST("/forgot-password", app.GeneralControllers.UserController.ForgotPassword)
		auth.POST("/confirm-otp", app.GeneralControllers.UserController.ConfirmOTP)
		auth.POST("/refresh-token", app.GeneralControllers.UserController.RefreshToken)
	}
}
