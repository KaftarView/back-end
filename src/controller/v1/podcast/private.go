package controller_v1_podcast

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type AdminPodcastController struct {
	constants      *bootstrap.Constants
	podcastService application_interfaces.PodcastService
}

func NewAdminPodcastController(
	constants *bootstrap.Constants,
	podcastService application_interfaces.PodcastService,
) *AdminPodcastController {
	return &AdminPodcastController{
		constants:      constants,
		podcastService: podcastService,
	}
}

func (adminPodcastController *AdminPodcastController) CreatePodcast(c *gin.Context) {
	type createPodcastParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  []string              `form:"categories"`
	}
	param := controller.Validated[createPodcastParams](c, &adminPodcastController.constants.Context)
	userID, _ := c.Get(adminPodcastController.constants.Context.UserID)
	podcast := adminPodcastController.podcastService.CreatePodcast(param.Name, param.Description, param.Categories, param.Banner, userID.(uint))

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcast")
	controller.Response(c, 200, message, podcast.ID)
}

func (adminPodcastController *AdminPodcastController) UpdatePodcast(c *gin.Context) {
	type updatePodcastParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Categories  *[]string             `form:"categories"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[updatePodcastParams](c, &adminPodcastController.constants.Context)
	adminPodcastController.podcastService.UpdatePodcast(param.PodcastID, param.Name, param.Description, param.Categories, param.Banner)

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updatePodcast")
	controller.Response(c, 200, message, nil)
}

func (adminPodcastController *AdminPodcastController) DeletePodcast(c *gin.Context) {
	type deletePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[deletePodcastParams](c, &adminPodcastController.constants.Context)
	adminPodcastController.podcastService.DeletePodcast(param.PodcastID)

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deletePodcast")
	controller.Response(c, 200, message, nil)
}

func (adminPodcastController *AdminPodcastController) CreateEpisode(c *gin.Context) {
	type createEpisodeParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Audio       *multipart.FileHeader `form:"audio" validate:"required"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[createEpisodeParams](c, &adminPodcastController.constants.Context)
	userID, _ := c.Get(adminPodcastController.constants.Context.UserID)
	episode := adminPodcastController.podcastService.CreateEpisode(param.Name, param.Description, param.Banner, param.Audio, param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcastEpisode")
	controller.Response(c, 200, message, episode.ID)
}

func (adminPodcastController *AdminPodcastController) UpdateEpisode(c *gin.Context) {
	type updateEpisodeParams struct {
		Name        *string               `form:"name" validate:"omitempty,max=50"`
		Description *string               `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Audio       *multipart.FileHeader `form:"audio"`
		EpisodeID   uint                  `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[updateEpisodeParams](c, &adminPodcastController.constants.Context)
	adminPodcastController.podcastService.UpdateEpisode(param.EpisodeID, param.Name, param.Description, param.Banner, param.Audio)

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updatePodcastEpisode")
	controller.Response(c, 200, message, nil)
}

func (adminPodcastController *AdminPodcastController) DeleteEpisode(c *gin.Context) {
	type deleteEpisodeParams struct {
		EpisodeID uint `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[deleteEpisodeParams](c, &adminPodcastController.constants.Context)
	adminPodcastController.podcastService.DeleteEpisode(param.EpisodeID)

	trans := controller.GetTranslator(c, adminPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.deletePodcastEpisode")
	controller.Response(c, 200, message, nil)
}
