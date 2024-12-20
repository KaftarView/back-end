package controller_v1_general

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/enums"

	"github.com/gin-gonic/gin"
)

type AWSController struct {
	constants  *bootstrap.Constants
	awsService *application_aws.S3service
}

func NewAWSController(constants *bootstrap.Constants, awsService *application_aws.S3service) *AWSController {
	return &AWSController{
		constants:  constants,
		awsService: awsService,
	}
}

func (ac *AWSController) DeleteObjectController(c *gin.Context) {
	type deletedObjectParams struct {
		ObjectName string `json:"obj_name" validate:"required"`
	}
	param := controller.Validated[deletedObjectParams](c, &ac.constants.Context)
	ac.awsService.DeleteObject(enums.BannersBucket, param.ObjectName)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteObjectFromBucket")
	controller.Response(c, 200, message, nil)
}
