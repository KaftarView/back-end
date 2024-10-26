package controller_v1_general

import (
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/exceptions"
	"first-project/src/jwt"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	constants *bootstrap.Constants
}

func NewAuthController(constants *bootstrap.Constants) *AuthController {
	return &AuthController{
		constants: constants,
	}
}

func (authController *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(authController.constants.Context.RefreshToken)
	if err != nil {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}

	claims := jwt.VerifyToken(c, "./jwtKeys", authController.constants.Context.IsLoadedJWTPrivateKey, refreshToken)
	subject := claims["sub"].(string)

	accessToken, newRefreshToken := jwt.GenerateJWT(
		c, "./jwtKeys", authController.constants.Context.IsLoadedJWTPrivateKey, subject)

	jwt.SetAuthCookies(
		c, accessToken, newRefreshToken,
		authController.constants.Context.AccessToken,
		authController.constants.Context.RefreshToken)

	trans := controller.GetTranslator(c, authController.constants.Context.Translator)
	message, _ := trans.T("successMessage.refreshToken")
	controller.Response(c, 200, message, nil)
}
