package controller_v1_podcast

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type CustomerPodcastController struct {
	constants      *bootstrap.Constants
	podcastService application_interfaces.PodcastService
}

func NewCustomerPodcastController(
	constants *bootstrap.Constants,
	podcastService application_interfaces.PodcastService,
) *CustomerPodcastController {
	return &CustomerPodcastController{
		constants:      constants,
		podcastService: podcastService,
	}
}

func (customerPodcastController *CustomerPodcastController) SubscribePodcast(c *gin.Context) {
	type subscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[subscribePodcastParams](c, &customerPodcastController.constants.Context)
	userID, _ := c.Get(customerPodcastController.constants.Context.UserID)
	customerPodcastController.podcastService.SubscribePodcast(param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, customerPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.subscribePodcast")
	controller.Response(c, 200, message, nil)
}

func (customerPodcastController *CustomerPodcastController) UnSubscribePodcast(c *gin.Context) {
	type unSubscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[unSubscribePodcastParams](c, &customerPodcastController.constants.Context)
	userID, _ := c.Get(customerPodcastController.constants.Context.UserID)
	customerPodcastController.podcastService.UnSubscribePodcast(param.PodcastID, userID.(uint))

	trans := controller.GetTranslator(c, customerPodcastController.constants.Context.Translator)
	message, _ := trans.T("successMessage.unSubscribePodcast")
	controller.Response(c, 200, message, nil)
}

func (customerPodcastController *CustomerPodcastController) SubscribeStatus(c *gin.Context) {
	type unSubscribePodcastParams struct {
		PodcastID uint `uri:"podcastID" validate:"required"`
	}
	param := controller.Validated[unSubscribePodcastParams](c, &customerPodcastController.constants.Context)
	userID, _ := c.Get(customerPodcastController.constants.Context.UserID)
	status := customerPodcastController.podcastService.IsUserSubscribedPodcast(param.PodcastID, userID.(uint))

	controller.Response(c, 200, "", status)
}
