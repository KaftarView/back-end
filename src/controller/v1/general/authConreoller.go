package controller_v1_general

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	jwt_keys "first-project/src/jwtKeys"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	constants  *bootstrap.Constants
	jwtService *application_jwt.JWTToken
}

func NewAuthController(constants *bootstrap.Constants, jwtService *application_jwt.JWTToken) *AuthController {
	return &AuthController{
		constants:  constants,
		jwtService: jwtService,
	}
}

func (authController *AuthController) RefreshToken(c *gin.Context) {
	type refreshTokenParams struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	param := controller.Validated[refreshTokenParams](c, &authController.constants.Context)

	jwt_keys.SetupJWTKeys(c, authController.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	claims, err := authController.jwtService.VerifyToken(param.RefreshToken)
	if err != nil {
		panic(err)
	}
	userID := uint(claims["sub"].(float64))
	accessToken, _ := authController.jwtService.GenerateJWT(userID)

	trans := controller.GetTranslator(c, authController.constants.Context.Translator)
	message, _ := trans.T("successMessage.refreshToken")
	controller.Response(c, 200, message, accessToken)
}
