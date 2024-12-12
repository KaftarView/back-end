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
	podcastID := podcastController.podcastService.CreatePodcast(param.Name, param.Description, param.Categories, userID.(uint))
	objectPath := fmt.Sprintf("banners/podcasts/%d/images/%s", podcastID, param.Banner.Filename)
	podcastController.awsService.UploadObject(enums.BannersBucket, objectPath, param.Banner)
	podcastController.podcastService.SetPodcastBannerPath(objectPath, podcastID)

	trans := controller.GetTranslator(c, podcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.createPodcast")
	controller.Response(c, 200, message, podcastID)
}

func (podcastController *PodcastController) GetPodcastDetails(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) UpdatePodcast(c *gin.Context) {
	// some code here
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
	// some code here
}

func (podcastController *PodcastController) UpdateEpisode(c *gin.Context) {
	// some code here
}

func (podcastController *PodcastController) DeleteEpisode(c *gin.Context) {
	// some code here
}
