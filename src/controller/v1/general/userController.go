package controller_v1_general

import (
	"first-project/src/application"
	application_communication "first-project/src/application/communication/emailService"
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	jwt_keys "first-project/src/jwtKeys"
	repository_cache "first-project/src/repository/redis"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	constants    *bootstrap.Constants
	userService  *application.UserService
	emailService *application_communication.EmailService
	userCache    *repository_cache.UserCache
	otpService   *application.OTPService
	jwtService   *application_jwt.JWTToken
}

func NewUserController(
	constants *bootstrap.Constants,
	userService *application.UserService,
	emailService *application_communication.EmailService,
	userCache *repository_cache.UserCache,
	otpService *application.OTPService,
	jwtService *application_jwt.JWTToken,
) *UserController {
	return &UserController{
		constants:    constants,
		userService:  userService,
		emailService: emailService,
		userCache:    userCache,
		otpService:   otpService,
		jwtService:   jwtService,
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
	otp := userController.otpService.GenerateOTP()
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
	controller.Response(c, 201, message, nil)
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
	user := userController.userService.AuthenticateUser(param.Username, param.Password)
	jwt_keys.SetupJWTKeys(c, userController.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	accessToken, refreshToken := userController.jwtService.GenerateJWT(user.ID)
	userController.userCache.SetUser(user.ID, user.Name, user.Email)
	roles, permissions := userController.userService.FindUserRolesAndPermissions(user.ID)
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
	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.login")
	controller.Response(c, 200, message, userDataResponse)
}

func (userController *UserController) ForgotPassword(c *gin.Context) {
	type forgotPasswordParams struct {
		Email string `json:"email" validate:"required,email"`
	}
	param := controller.Validated[forgotPasswordParams](c, &userController.constants.Context)
	otp := userController.otpService.GenerateOTP()
	userController.userService.UpdateUserOTPIfExists(param.Email, otp)
	emailTemplateData := struct{ OTP string }{OTP: otp}
	templatePath := getTemplatePath(c, userController.constants.Context.Translator)
	userController.emailService.SendEmail(
		param.Email, "Forgot Password", "forgotPassword/"+templatePath, emailTemplateData)

	trans := controller.GetTranslator(c, userController.constants.Context.Translator)
	message, _ := trans.T("successMessage.forgotPassword")
	controller.Response(c, 200, message, nil)
}

func (userController *UserController) ConfirmOTP(c *gin.Context) {
	type confirmOTPParams struct {
		Email string `json:"email" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}
	param := controller.Validated[confirmOTPParams](c, &userController.constants.Context)
	userID := userController.userService.ValidateUserOTP(param.Email, param.OTP)
	jwt_keys.SetupJWTKeys(c, userController.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	accessToken, refreshToken := userController.jwtService.GenerateJWT(userID)
	type tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	controller.Response(c, 200, "", tokens{AccessToken: accessToken, RefreshToken: refreshToken})
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
