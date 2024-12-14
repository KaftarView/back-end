package controller_v1_private

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"
	"fmt"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type PodcastController struct {
	constants      *bootstrap.Constants
	podcastService *application.PodcastService
	awsService     *application_aws.S3service
}

func NewPodcastController(
	constants *bootstrap.Constants,
	podcastService *application.PodcastService,
	awsService *application_aws.S3service,
) *PodcastController {
	return &PodcastController{
		constants:      constants,
		podcastService: podcastService,
		awsService:     awsService,
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
		Categories  []string              `form:"category"`
	}
	param := controller.Validated[createPodcastParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	podcast := podcastController.podcastService.CreatePodcast(param.Name, param.Description, param.Categories, userID.(uint))
	objectPath := fmt.Sprintf("banners/podcasts/%d/images/%s", podcast.ID, param.Banner.Filename)
	podcastController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	podcastController.podcastService.SetPodcastBannerPath(objectPath, podcast)

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
		Categories  *[]string             `form:"category"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[updatePodcastParams](c, &podcastController.constants.Context)
	podcast := podcastController.podcastService.UpdatePodcast(param.PodcastID, param.Name, param.Description, param.Categories)
	if param.Banner != nil {
		podcastController.awsService.DeleteObject(enums.BannersBucket, podcast.BannerPath)
		objectPath := fmt.Sprintf("banners/podcasts/%d/images/%s", param.PodcastID, param.Banner.Filename)
		podcastController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
		podcastController.podcastService.SetPodcastBannerPath(objectPath, podcast)
	}

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.updatePodcast")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) DeletePodcast(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) SubscribePodcast(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) UnSubscribePodcast(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) GetEpisodesList(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) CreateEpisode(c *gin.Context) {
	type createEpisodeParams struct {
		Name        string                `form:"name" validate:"required,max=50"`
		Description string                `form:"description"`
		Banner      *multipart.FileHeader `form:"banner"`
		Audio       *multipart.FileHeader `form:"audio"`
		PodcastID   uint                  `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[createEpisodeParams](c, &podcastController.constants.Context)
	userID, _ := c.Get(podcastController.constants.Context.UserID)
	episode := podcastController.podcastService.CreateEpisode(param.Name, param.Description, param.PodcastID, userID.(uint))

	bannerPath := fmt.Sprintf("banners/podcasts/%d/episodes/%d/images/%s", param.PodcastID, episode.ID, param.Banner.Filename)
	podcastController.awsService.UploadObject(enums.BannersBucket, bannerPath, param.Banner)
	podcastController.podcastService.SetEpisodeBannerPath(bannerPath, episode)

	audioPath := fmt.Sprintf("media/podcasts/%d/episodes/%d/audio/%s", param.PodcastID, episode.ID, param.Audio.Filename)
	podcastController.awsService.UploadObject(enums.PodcastsBucket, audioPath, param.Banner)
	podcastController.podcastService.SetEpisodeAudioPath(audioPath, episode)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcastEpisode")
	controller.Response(c, 200, message, nil)
}

func (podcastController *PodcastController) UpdateEpisode(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) DeleteEpisode(c *gin.Context) {
	// some code here
}
