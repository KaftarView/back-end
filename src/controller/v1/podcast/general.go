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
	type podcastListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[podcastListParams](c, &generalPodcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	podcasts := generalPodcastController.podcastService.GetPodcastList(param.Page, param.PageSize)
	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) SearchPodcast(c *gin.Context) {
	type searchPodcastsParams struct {
		Query    string `form:"query"`
		Page     int    `form:"page"`
		PageSize int    `form:"pageSize"`
	}
	param := controller.Validated[searchPodcastsParams](c, &generalPodcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	podcasts := generalPodcastController.podcastService.SearchEvents(param.Query, param.Page, param.PageSize)

	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) FilterPodcastByCategory(c *gin.Context) {
	type filterPodcastsParams struct {
		Categories []string `form:"categories"`
		Page       int      `form:"page"`
		PageSize   int      `form:"pageSize"`
	}
	param := controller.Validated[filterPodcastsParams](c, &generalPodcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}

	podcasts := generalPodcastController.podcastService.FilterPodcastsByCategory(param.Categories, param.Page, param.PageSize)

	controller.Response(c, 200, "", podcasts)
}

func (generalPodcastController *GeneralPodcastController) GetEpisodesList(c *gin.Context) {
	type episodeListParams struct {
		Page     int `form:"page"`
		PageSize int `form:"pageSize"`
	}
	param := controller.Validated[episodeListParams](c, &generalPodcastController.constants.Context)
	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	episodes := generalPodcastController.podcastService.GetEpisodesList(param.Page, param.PageSize)
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
