package controller_v1_general

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/jwt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	constants    *bootstrap.Constants
	userService  *application.UserService
	emailService *application_communication.EmailService
}

func NewUserController(
	constants *bootstrap.Constants,
	userService *application.UserService,
	emailService *application_communication.EmailService,
) *UserController {
	return &UserController{
		constants:    constants,
		userService:  userService,
		emailService: emailService,
	}
}

func getTemplatePath(c *gin.Context, transKey string) string {
	trans := controller.GetTranslator(c, transKey)
	if trans.Locale() == "fa_IR" {
		return "fa.html"
	}
	return "en.html"
}

func (userController *UserController) Register(c *gin.Context) {
	type registerParams struct {
		Username        string `json:"username" validate:"required,gt=2,lt=20"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[registerParams](c, &userController.constants.Context)
	userController.userService.ValidateUserRegistrationDetails(param.Username, param.Email, param.Password, param.ConfirmPassword)
	otp := application.GenerateOTP()
	userController.userService.UpdateOrCreateUser(param.Username, param.Email, param.Password, otp)

	emailTemplateData := struct {
		Username string
		OTP      string
	}{
		Username: param.Username,
		OTP:      otp,
	}
	templatePath := getTemplatePath(c, userController.constants.Context.Translator)
	userController.emailService.SendEmail(
		param.Email, "Activate account", "activateAccount/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.userRegistration")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) VerifyEmail(c *gin.Context) {
	type verifyEmailParams struct {
		OTP   string `json:"otp" validate:"required"`
		Email string `json:"email" validate:"required"`
	}
	param := controller.Validated[verifyEmailParams](c, &userController.constants.Context)
	userController.userService.ActivateUser(param.Email, param.OTP)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.emailVerification")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) Login(c *gin.Context) {
	type loginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	param := controller.Validated[loginParams](c, &userController.constants.Context)
	userController.userService.VerifyLogin(param.Username, param.Password)
	jwtString := jwt.GenerateJWT(c, "./jwtKeys", userController.constants.Context.IsLoadedJWTPrivateKey, param.Username)
	c.SetCookie(userController.constants.Context.Token, jwtString, 3600, "/", "localhost", true, true)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.login")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) ForgotPassword(c *gin.Context) {
	type forgotPasswordParams struct {
		Email string `json:"email" validate:"required,email"`
	}
	param := controller.Validated[forgotPasswordParams](c, &userController.constants.Context)
	userController.userService.VerifyUserActivated(param.Email)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.forgotPassword")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) ResetPassword(c *gin.Context) {
	type resetPasswordParams struct {
		Email           string `json:"email" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[resetPasswordParams](c, &userController.constants.Context)
	userController.userService.ResetPasswordService(param.Email, param.Password, param.ConfirmPassword)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.resetPassword")
	controller.Response(c, 200, message, nil)
}
