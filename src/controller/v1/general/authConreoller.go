package controller_v1_general

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/exceptions"
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
	refreshToken, err := c.Cookie(authController.constants.Context.RefreshToken)
	if err != nil {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}
	claims := authController.jwtService.VerifyToken(refreshToken)
	userID := claims["sub"].(uint)
	jwt_keys.SetupJWTKeys(c, authController.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	accessToken, newRefreshToken := authController.jwtService.GenerateJWT(userID)
	controller.SetAuthCookies(
		c, accessToken, newRefreshToken,
		authController.constants.Context.AccessToken,
		authController.constants.Context.RefreshToken,
	)

	trans := controller.GetTranslator(c, authController.constants.Context.Translator)
	message, _ := trans.T("successMessage.refreshToken")
	controller.Response(c, 200, message, nil)
}
