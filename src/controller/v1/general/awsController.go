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
	type createBucketParams struct {
		BucketName string `json:"bucket_name" validate:"required"`
	}
	param := controller.Validated[createBucketParams](c, &ac.constants.Context)
	ac.awsService.CreateBucket(param.BucketName)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.createBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) UploadObjectController(c *gin.Context) {
	type uploadObjectParams struct {
		BucketName string `json:"bucket_name" validate:"required"`
	}
	param := controller.Validated[uploadObjectParams](c, &ac.constants.Context)
	file, err := c.FormFile("file")
	if err != nil {
		bindingError := exceptions.BindingError{Err: err}
		panic(bindingError)
	}
	ac.awsService.UploadObject(file, param.BucketName)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.uploadObjectToBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) DeleteObjectController(c *gin.Context) {
	type deletedObjectParams struct {
		BucketName string `json:"bucket_name" validate:"required"`
		ObjectName string `json:"obj_name" validate:"required"`
	}
	param := controller.Validated[deletedObjectParams](c, &ac.constants.Context)
	ac.awsService.DeleteObject(param.ObjectName, param.BucketName)
	trans := controller.GetTranslator(c, ac.constants.Context.Translator)
	message, _ := trans.T("successMessage.deleteObjectFromBucket")
	controller.Response(c, 200, message, nil)
}

func (ac *AWSController) GetListOfObjectsController(c *gin.Context) {
	type listBucketParams struct {
		BucketName string `json:"bucket_name" validate:"required"`
	}
	param := controller.Validated[listBucketParams](c, &ac.constants.Context)
	objects := ac.awsService.ListObjects(param.BucketName)
	controller.Response(c, 200, "", objects)
}
