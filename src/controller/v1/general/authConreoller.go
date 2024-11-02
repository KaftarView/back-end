package controller_v1_general

import (
	"first-project/src/bootstrap"
	"first-project/src/controller"
	"first-project/src/exceptions"
	application_jwt "first-project/src/jwt"

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
	claims := authController.jwtService.VerifyToken(c, "./jwtKeys", refreshToken)
	userID := claims["sub"].(uint)
	accessToken, newRefreshToken := authController.jwtService.GenerateJWT(c, "./jwtKeys", userID)
	authController.jwtService.SetAuthCookies(c, accessToken, newRefreshToken)

	trans := controller.GetTranslator(c, authController.constants.Context.Translator)
	message, _ := trans.T("successMessage.refreshToken")
	controller.Response(c, 200, message, nil)
}
