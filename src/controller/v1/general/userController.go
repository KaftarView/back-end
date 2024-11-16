package controller_v1_general

import (
	"first-project/src/application"
	"first-project/src/bootstrap"
	"first-project/src/controller"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	constants    *bootstrap.Constants
	userService  *application.UserService
	otpService   *application.OTPService
	emailService *application.EmailService
}

func NewUserController(
	constants *bootstrap.Constants,
	userService *application.UserService,
	otpService *application.OTPService,
	emailService *application.EmailService,
) *UserController {
	return &UserController{
		constants:    constants,
		userService:  userService,
		otpService:   otpService,
		emailService: emailService,
	}
}

func userRegistrationResponse(c *gin.Context, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	message, _ := trans.T("successMessage.userRegistration")
	controller.Response(c, 200, message, nil)
}

func emailVerificationResponse(c *gin.Context, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	message, _ := trans.T("successMessage.emailVerification")
	controller.Response(c, 200, message, nil)
}

func loginResponse(c *gin.Context, transKey string) {
	trans := controller.GetTranslator(c, transKey)
	message, _ := trans.T("successMessage.login")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) Register(c *gin.Context) {
	type registerParams struct {
		Username        string `json:"username" validate:"required,gt=2,lt=20"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[registerParams](c, &userController.constants.Context)
	userController.userService.VerifyUserRegistration(param.Username, param.Email, param.Password, param.ConfirmPassword)
	otp := application.GenerateOTP()
	userController.emailService.SendVerificationEmail(param.Username, param.Email, otp)
	userController.userService.RegisterUser(param.Username, param.Email, param.Password, otp)
	userRegistrationResponse(c, userController.constants.Context.Translator)
}

func (userController *UserController) VerifyEmail(c *gin.Context) {
	type verifyEmailParams struct {
		OTP   string `json:"otp" validate:"required"`
		Email string `json:"email" validate:"required"`
	}
	param := controller.Validated[verifyEmailParams](c, &userController.constants.Context)
	userController.userService.CheckUserAlreadyVerifiedByEmail(param.Email)
	userController.otpService.VerifyOTP(param.OTP, param.Email)
	userController.userService.VerifyEmail(param.Email)
	emailVerificationResponse(c, userController.constants.Context.Translator)
}

func (userController *UserController) Login(c *gin.Context) {
	type LoginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	param := controller.Validated[LoginParams](c, &userController.constants.Context)
	userController.userService.LoginService(param.Username, param.Password)
	loginResponse(c, userController.constants.Context.Translator)
}
