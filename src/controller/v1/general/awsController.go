package controller_v1_general

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type AWSController struct {
	constants  *bootstrap.Constants
	awsService *application_aws.AWSS3
}

func NewAWSController(constants *bootstrap.Constants, awsService *application_aws.AWSS3) *AWSController {
	return &AWSController{
		constants:  constants,
		awsService: awsService,
	}
}

func (ac *AWSController) CreateBucket(c *gin.Context) {
	ac.awsService.CreateBucket()
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.createBucket")
	controller.Response(c, 200, message, nil)
}
