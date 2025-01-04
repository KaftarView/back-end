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

func (customerUserController *CustomerUserController) ChangeUsername(c *gin.Context) {
	type changeUsernameParams struct {
		Username string `json:"username" validate:"required"`
	}
	param := controller.Validated[changeUsernameParams](c, &customerUserController.constants.Context)
	userID, _ := c.Get(customerUserController.constants.Context.UserID)
	customerUserController.userService.UpdateUser(userID.(uint), param.Username)

	trans := controller.GetTranslator(c, customerUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.changeUsername")
	controller.Response(c, 200, message, nil)
}

func (customerUserController *CustomerUserController) ResetPassword(c *gin.Context) {
	type resetPasswordParams struct {
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[resetPasswordParams](c, &customerUserController.constants.Context)
	userID, _ := c.Get(customerUserController.constants.Context.UserID)
	customerUserController.userService.ResetPasswordService(userID.(uint), param.Password, param.ConfirmPassword)

	trans := controller.GetTranslator(c, customerUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.resetPassword")
	controller.Response(c, 200, message, nil)
}
