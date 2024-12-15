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
	// some code here
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

func (podcastController *PodcastController) GetPodcastDetails(c *gin.Context) {
	// some code here
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
	// some code here
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
	type subscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[subscribePodcastParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	podcastController.podcastService.UnSubscribePodcast(param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.unSubscribePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) GetEpisodesList(c *gin.Context) {
	// some code here
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
