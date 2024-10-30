package middleware_authentication

import (
	"first-project/src/bootstrap"
	"first-project/src/enums"
	"first-project/src/exceptions"
	"first-project/src/jwt"
	"first-project/src/repository"

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

func (authMiddleware *AuthMiddleware) AuthenticateMiddleware(c *gin.Context, allowedRules []enums.RoleType) {
	tokenString, err := c.Cookie(authMiddleware.constants.Context.AccessToken)
	if err != nil {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}
	claims := jwt.VerifyToken(c, "./jwtKeys", authMiddleware.constants.Context.IsLoadedJWTPrivateKey, tokenString)
	userID := claims["sub"].(uint)

	roles := authMiddleware.userRepository.FindUserRoleTypesByUserID(userID)

	if !isAllowRole(allowedRules, roles) {
		authError := exceptions.NewForbiddenError()
		panic(authError)
	}
	c.Next()
}

func isAllowRole(allowedRoles []enums.RoleType, userRoles []enums.RoleType) bool {
	allowedRolesMap := make(map[enums.RoleType]bool)
	for _, allowedRole := range allowedRoles {
		allowedRolesMap[allowedRole] = true
	}

	for _, userRole := range userRoles {
		if allowedRolesMap[userRole] {
			return true
		}
	}
	return false
}
