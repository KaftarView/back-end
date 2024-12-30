package controller_v1_podcast

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type GeneralPodcastController struct {
	constants      *bootstrap.Constants
	podcastService application_interfaces.PodcastService
}

func NewGeneralPodcastController(
	constants *bootstrap.Constants,
	podcastService application_interfaces.PodcastService,
) *GeneralPodcastController {
	return &GeneralPodcastController{
		constants:      constants,
		podcastService: podcastService,
	}
}

func (generalPodcastController *GeneralPodcastController) GetPodcastsList(c *gin.Context) {
	pagination := controller.GetPagination(c, &generalPodcastController.constants.Context)
	podcasts := generalPodcastController.podcastService.GetPodcastList(pagination.Page, pagination.PageSize)
	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) SearchPodcast(c *gin.Context) {
	type searchPodcastsParams struct {
		Query string `form:"query"`
	}
	param := controller.Validated[searchPodcastsParams](c, &generalPodcastController.constants.Context)
	pagination := controller.GetPagination(c, &generalPodcastController.constants.Context)
	podcasts := generalPodcastController.podcastService.SearchEvents(param.Query, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) FilterPodcastByCategory(c *gin.Context) {
	type filterPodcastsParams struct {
		Categories []string `form:"categories"`
	}
	param := controller.Validated[filterPodcastsParams](c, &generalPodcastController.constants.Context)
	pagination := controller.GetPagination(c, &generalPodcastController.constants.Context)
	podcasts := generalPodcastController.podcastService.FilterPodcastsByCategory(param.Categories, pagination.Page, pagination.PageSize)

	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) GetEpisodesList(c *gin.Context) {
	pagination := controller.GetPagination(c, &generalPodcastController.constants.Context)
	episodes := generalPodcastController.podcastService.GetEpisodesList(pagination.Page, pagination.PageSize)
	controller.Response(c, 200, "", episodes)
}

func (generalPodcastController *GeneralPodcastController) GetPodcastDetails(c *gin.Context) {
	type podcastDetailsParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[podcastDetailsParams](c, &generalPodcastController.constants.Context)
	podcast := generalPodcastController.podcastService.GetPodcastDetails(param.PodcastID)
	controller.Response(c, 200, "", podcast)
}

func (generalPodcastController *GeneralPodcastController) GetEpisodeDetails(c *gin.Context) {
	type episodeListParams struct {
		EpisodeID uint `uri:"episodeID" validate:"required"`
	}
	param := controller.Validated[episodeListParams](c, &generalPodcastController.constants.Context)
	episodes := generalPodcastController.podcastService.GetEpisodeDetails(param.EpisodeID)

	controller.Response(c, 200, "", episodes)
}
