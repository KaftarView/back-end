package middleware_authentication

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/enums"
	"first-project/src/exceptions"
	jwt_keys "first-project/src/jwtKeys"
	repository_database "first-project/src/repository/database"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	userRepository *repository_database.UserRepository
	jwtService     *application_jwt.JWTToken
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	userRepository *repository_database.UserRepository,
	jwtService *application_jwt.JWTToken,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:      constants,
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (authMiddleware *AuthMiddleware) AuthenticateMiddleware(c *gin.Context, allowedRules []enums.RoleType) {
	tokenString, err := c.Cookie(authMiddleware.constants.Context.AccessToken)
	if err != nil {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}
	jwt_keys.SetupJWTKeys(c, authMiddleware.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	claims := authMiddleware.jwtService.VerifyToken(tokenString)
	userID := uint(claims["sub"].(float64))

	roles := authMiddleware.userRepository.FindUserRoleTypesByUserID(uint(userID))

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
