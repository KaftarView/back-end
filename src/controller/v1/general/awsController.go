package controller_v1_general

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"mime/multipart"

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

func (ac *AWSController) UploadObjectController(c *gin.Context) {
	type uploadObjectParams struct {
		File *multipart.FileHeader `form:"fileSample" validate:"required"`
	}
	param := controller.Validated[uploadObjectParams](c, &ac.constants.Context)
	// TODO: should be based on request -> different for each req type
	ac.awsService.UploadObject(param.File, "session", 123)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.uploadObjectToBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) DeleteObjectController(c *gin.Context) {
	type deletedObjectParams struct {
		ObjectName string `json:"obj_name" validate:"required"`
	}
	param := controller.Validated[deletedObjectParams](c, &ac.constants.Context)
	ac.awsService.DeleteObject(param.ObjectName)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteObjectFromBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) GetListOfObjectsController(c *gin.Context) {
	objects := ac.awsService.ListObjects()
	controller.Response(c, 200, "", objects)
}

func (ac *AWSController) GetUserObjects(c *gin.Context) {
	// TODO: should be based on request
	objects := ac.awsService.GetSessionVideoURLs(123)
	controller.Response(c, 200, "", objects)
}
