package controller_v1_user

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_interfaces "first-project/src/application/interfaces"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	jwt_keys "first-project/src/jwtKeys"
	repository_cache "first-project/src/repository/redis"

	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	constants    *bootstrap.Constants
	userService  application_interfaces.UserService
	emailService *application_communication.EmailService
	userCache    *repository_cache.UserCache
	otpService   *application.OTPService
	jwtService   *application_jwt.JWTToken
}

func NewGeneralUserController(
	constants *bootstrap.Constants,
	userService application_interfaces.UserService,
	emailService *application_communication.EmailService,
	userCache *repository_cache.UserCache,
	otpService *application.OTPService,
	jwtService *application_jwt.JWTToken,
) *GeneralUserController {
	return &GeneralUserController{
		constants:    constants,
		userService:  userService,
		emailService: emailService,
		userCache:    userCache,
		otpService:   otpService,
		jwtService:   jwtService,
	}
}

func (generalUserController *GeneralUserController) Register(c *gin.Context) {
	type registerParams struct {
		Username        string `json:"username" validate:"required,gt=2,lt=20"`
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[registerParams](c, &generalUserController.constants.Context)
	generalUserController.userService.ValidateUserRegistrationDetails(param.Username, param.Email, param.Password, param.ConfirmPassword)
	otp := generalUserController.otpService.GenerateOTP()
	generalUserController.userService.UpdateOrCreateUser(param.Username, param.Email, param.Password, otp)

	emailTemplateData := struct {
		Username string
		OTP      string
	}{
		Username: param.Username,
		OTP:      otp,
	}
	templatePath := controller.GetTemplatePath(c, generalUserController.constants.Context.Translator)
	generalUserController.emailService.SendEmail(
		param.Email, "Activate account", "activateAccount/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.userRegistration")
	controller.Response(c, 201, message, nil)
}

func (generalUserController *GeneralUserController) VerifyEmail(c *gin.Context) {
	type verifyEmailParams struct {
		OTP   string `json:"otp" validate:"required"`
		Email string `json:"email" validate:"required"`
	}
	param := controller.Validated[verifyEmailParams](c, &generalUserController.constants.Context)
	generalUserController.userService.ActivateUser(param.Email, param.OTP)

	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.emailVerification")
	controller.Response(c, 200, message, nil)
}

func (generalUserController *GeneralUserController) Login(c *gin.Context) {
	type loginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	param := controller.Validated[loginParams](c, &generalUserController.constants.Context)
	user := generalUserController.userService.AuthenticateUser(param.Username, param.Password)
	jwt_keys.SetupJWTKeys(c,
		generalUserController.constants.Context.IsLoadedJWTKeys,
		generalUserController.constants.JWTKeysPath)
	accessToken, refreshToken := generalUserController.jwtService.GenerateJWT(user.ID)
	generalUserController.userCache.SetUser(user.ID, user.Name, user.Email)
	roles, permissions := generalUserController.userService.FindUserRolesAndPermissions(user.ID)
	userDataResponse := struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		ID           uint     `json:"id"`
		Name         string   `json:"username"`
		Email        string   `json:"email"`
		Roles        []string `json:"roles"`
		Permissions  []string `json:"permissions"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Roles:        roles,
		Permissions:  permissions,
	}
	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.login")
	controller.Response(c, 200, message, userDataResponse)
}

func (generalUserController *GeneralUserController) ForgotPassword(c *gin.Context) {
	type forgotPasswordParams struct {
		Email string `json:"email" validate:"required,email"`
	}
	param := controller.Validated[forgotPasswordParams](c, &generalUserController.constants.Context)
	otp := generalUserController.otpService.GenerateOTP()
	generalUserController.userService.UpdateUserOTPIfExists(param.Email, otp)
	emailTemplateData := struct{ OTP string }{OTP: otp}
	templatePath := controller.GetTemplatePath(c, generalUserController.constants.Context.Translator)
	generalUserController.emailService.SendEmail(
		param.Email, "Forgot Password", "forgotPassword/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.forgotPassword")
	controller.Response(c, 200, message, nil)
}

func (generalUserController *GeneralUserController) ConfirmOTP(c *gin.Context) {
	type confirmOTPParams struct {
		Email string `json:"email" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}
	param := controller.Validated[confirmOTPParams](c, &generalUserController.constants.Context)
	userID := generalUserController.userService.ValidateUserOTP(param.Email, param.OTP)
	jwt_keys.SetupJWTKeys(c,
		generalUserController.constants.Context.IsLoadedJWTKeys,
		generalUserController.constants.JWTKeysPath)
	accessToken, refreshToken := generalUserController.jwtService.GenerateJWT(userID)
	type tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	controller.Response(c, 200, "", tokens{AccessToken: accessToken, RefreshToken: refreshToken})
}

func (generalUserController *GeneralUserController) ResetPassword(c *gin.Context) {
	type resetPasswordParams struct {
		Email           string `json:"email" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}
	param := controller.Validated[resetPasswordParams](c, &generalUserController.constants.Context)
	generalUserController.userService.ResetPasswordService(param.Email, param.Password, param.ConfirmPassword)

	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.resetPassword")
	controller.Response(c, 200, message, nil)
}

func (generalUserController *GeneralUserController) RefreshToken(c *gin.Context) {
	type refreshTokenParams struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	param := controller.Validated[refreshTokenParams](c, &generalUserController.constants.Context)

	jwt_keys.SetupJWTKeys(c,
		generalUserController.constants.Context.IsLoadedJWTKeys,
		generalUserController.constants.JWTKeysPath)
	claims := generalUserController.jwtService.VerifyToken(param.RefreshToken)
	userID := uint(claims["sub"].(float64))
	accessToken, _ := generalUserController.jwtService.GenerateJWT(userID)

	trans := controller.GetTranslator(c, generalUserController.constants.Context.Translator)
	message, _ := trans.T("successMessage.refreshToken")
	controller.Response(c, 200, message, accessToken)
}
