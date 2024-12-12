package controller_v1_private

import (
	"first-project/src/application"
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"

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
	// some code here
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
