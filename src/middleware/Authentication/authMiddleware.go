package middleware_authentication

import (
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/exceptions"
	"first-project/src/jwt"
	"first-project/src/repository"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	userRepository *repository.UserRepository
}

func NewAuthMiddleware(constants *bootstrap.Constants, userRepository *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		constants:      constants,
		userRepository: userRepository,
	}
}

func (authMiddleware *AuthMiddleware) AuthenticateMiddleware(c *gin.Context, allowedRules []string) {
	tokenString, err := c.Cookie(authMiddleware.constants.Context.Token)
	if err != nil {
		panic(err)
	}
	claims := jwt.VerifyToken(tokenString)
	subject := claims["sub"].(string)
	user, userExist := authMiddleware.userRepository.FindByUsernameAndVerified(subject, true)
	if !userExist {
		panic(fmt.Errorf("no user found for this jwt token subject"))
	}

	roles := authMiddleware.userRepository.FindUserRolesByUserID(user.ID)

	if !isAllowRole(allowedRules, roles) {
		authError := exceptions.NewAuthError()
		panic(authError)
	}
	c.Next()
}

func isAllowRole(allowedRoles []string, userRoles []entities.Role) bool {
	for _, allowedRole := range allowedRoles {
		for _, userRole := range userRoles {
			if allowedRole == userRole.Name {
				return true
			}
		}
	}
	return false
}
