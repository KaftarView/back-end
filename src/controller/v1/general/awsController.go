package controller_v1_general

import (
	application_aws "first-project/src/application/aws"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/exceptions"

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

func (ac *AWSController) CreateBucketController(c *gin.Context) {
	ac.awsService.CreateBucket()
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.createBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) UploadObjectController(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		bindingError := exceptions.BindingError{Err: err}
		panic(bindingError)
	}
	ac.awsService.UploadObject(file)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.uploadObjectToBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) DeleteObjectController(c *gin.Context) {
	type deletedObjectParams struct {
		ObjectName string `json:"otp" validate:"required"`
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
