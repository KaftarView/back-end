package controller_v1_private

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type PodcastController struct {
	constants      *bootstrap.Constants
	podcastService *application.PodcastService
}

func NewPodcastController(
	constants *bootstrap.Constants,
	podcastService *application.PodcastService,
) *PodcastController {
	return &PodcastController{
		constants:      constants,
		podcastService: podcastService,
	}
}

func (podcastController *PodcastController) GetPodcastsList(c *gin.Context) {
	type podcastListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[podcastListParams](c, &podcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	podcasts := podcastController.podcastService.GetPodcastList(param.Page, param.PageSize)
	controller.Response(c, 200, "", podcasts)
}

func (podcastController *PodcastController) GetPodcastDetails(c *gin.Context) {
	type podcastDetailsParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[podcastDetailsParams](c, &podcastController.constants.Context)
	podcast := podcastController.podcastService.GetPodcastDetails(param.PodcastID)
	controller.Response(c, 200, "", podcast)
}

func (podcastController *PodcastController) CreatePodcast(c *gin.Context) {
	type createPodcastParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  []string              `form:"categories"`
	}
	param := controller.Validated[createPodcastParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	podcast := podcastController.podcastService.CreatePodcast(param.Name, param.Description, param.Categories, param.Banner, userID.(uint))

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcast")
	controller.Response(c, 200, message, podcast.ID)
}

func (podcastController *PodcastController) UpdatePodcast(c *gin.Context) {
	type updatePodcastParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  *[]string             `form:"categories"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[updatePodcastParams](c, &podcastController.constants.Context)
	podcastController.podcastService.UpdatePodcast(param.PodcastID, param.Name, param.Description, param.Categories, param.Banner)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updatePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) DeletePodcast(c *gin.Context) {
	type deletePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[deletePodcastParams](c, &podcastController.constants.Context)
	podcastController.podcastService.DeletePodcast(param.PodcastID)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deletePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) SubscribePodcast(c *gin.Context) {
	type subscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[subscribePodcastParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	podcastController.podcastService.SubscribePodcast(param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.subscribePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) UnSubscribePodcast(c *gin.Context) {
	type unSubscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[unSubscribePodcastParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	podcastController.podcastService.UnSubscribePodcast(param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.unSubscribePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) GetEpisodesList(c *gin.Context) {
	type episodeListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[episodeListParams](c, &podcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	episodes := podcastController.podcastService.GetEpisodesList(param.Page, param.PageSize)
	controller.Response(c, 200, "", episodes)
}

func (podcastController *PodcastController) GetEpisodeDetails(c *gin.Context) {
	type episodeListParams struct {
		EpisodeID uint `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[episodeListParams](c, &podcastController.constants.Context)
	episodes := podcastController.podcastService.GetEpisodeDetails(param.EpisodeID)

	controller.Response(c, 200, "", episodes)
}

func (podcastController *PodcastController) CreateEpisode(c *gin.Context) {
	type createEpisodeParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Audio       *multipart.FileHeader `form:"audio" validate:"required"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[createEpisodeParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	episode := podcastController.podcastService.CreateEpisode(param.Name, param.Description, param.Banner, param.Audio, param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcastEpisode")
	controller.Response(c, 200, message, episode.ID)
}

func (podcastController *PodcastController) UpdateEpisode(c *gin.Context) {
	type updateEpisodeParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Audio       *multipart.FileHeader `form:"audio"`
		EpisodeID   uint                  `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[updateEpisodeParams](c, &podcastController.constants.Context)
	podcastController.podcastService.UpdateEpisode(param.EpisodeID, param.Name, param.Description, param.Banner, param.Audio)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updatePodcastEpisode")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) DeleteEpisode(c *gin.Context) {
	type deleteEpisodeParams struct {
		EpisodeID uint `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[deleteEpisodeParams](c, &podcastController.constants.Context)
	podcastController.podcastService.DeleteEpisode(param.EpisodeID)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deletePodcastEpisode")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) SearchPodcast(c *gin.Context) {
	type searchPodcastsParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchPodcastsParams](c, &podcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	podcasts := podcastController.podcastService.SearchEvents(param.Query, param.Page, param.PageSize)

	controller.Response(c, 200, "", podcasts)
}

func (podcastController *PodcastController) FilterPodcastByCategory(c *gin.Context) {
	type filterPodcastsParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterPodcastsParams](c, &podcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}

	podcasts := podcastController.podcastService.FilterPodcastsByCategory(param.Categories, param.Page, param.PageSize)

	controller.Response(c, 200, "", podcasts)
}
