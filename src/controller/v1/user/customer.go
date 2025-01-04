package controller_v1_user

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type CustomerUserController struct {
	constants   *bootstrap.Constants
	userService application_interfaces.UserService
}

func NewCustomerUserController(
	constants *bootstrap.Constants,
	userService application_interfaces.UserService,
) *CustomerUserController {
	return &CustomerUserController{
		constants:   constants,
		userService: userService,
	}
}

func (customerUserController *CustomerUserController) ResetPassword(c *gin.Context) {
	type resetPasswordParams struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[resetPasswordParams](c, &customerUserController.constants.Context)
	customerUserController.userService.ResetPasswordService(param.Email, param.Password, param.ConfirmPassword)

	trans := controller.GetTranslator(c, customerUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.resetPassword")
	controller.Response(c, 200, message, nil)
}
