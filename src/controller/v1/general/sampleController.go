package controller_v1_general

import (
	application_math "first-project/src/application/math"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SampleController struct {
	constants  *bootstrap.Constants
	addService *application_math.AddService
}

func NewSampleController(constants *bootstrap.Constants, addService *application_math.AddService) *SampleController {
	return &SampleController{
		constants:  constants,
		addService: addService,
	}
}

func (sample *SampleController) Add(c *gin.Context) {
	type addParams struct {
		Num1 int `uri:"num1" validate:"required,number,gt=10"`
		Num2 int `uri:"num2" validate:"required,number,gt=11"`
	}
	param := controller.Validated[addParams](c, &sample.constants.Context)
	sum := sample.addService.Add(param.Num1, param.Num2)
	c.String(http.StatusOK, strconv.Itoa(sum))
}
