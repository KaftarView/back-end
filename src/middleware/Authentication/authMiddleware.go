package middleware_authentication

import (
	application_jwt "first-project/src/application/jwt"
	"first-project/src/bootstrap"
	"first-project/src/entities"
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

func (am *AuthMiddleware) Authentication(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString != "" {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}
	jwt_keys.SetupJWTKeys(c, am.constants.Context.IsLoadedJWTKeys, "./src/jwtKeys")
	claims := am.jwtService.VerifyToken(tokenString)

	c.Set(am.constants.Context.UserID, uint(claims["sub"].(float64)))

	c.Next()
}

func (am *AuthMiddleware) RequirePermission(allowedPermissions []enums.PermissionType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exist := c.Get(am.constants.Context.UserID)
		if !exist {
			unauthorizedError := exceptions.NewUnauthorizedError()
			panic(unauthorizedError)
		}

		roles := am.userRepository.FindUserRoleTypesByUserID(userID.(uint))

		if !am.isAllowRole(allowedPermissions, roles) {
			authError := exceptions.NewForbiddenError()
			panic(authError)
		}

		c.Next()
	}
}

func (am *AuthMiddleware) isAllowRole(allowedPermissions []enums.PermissionType, userRoles []entities.Role) bool {
	allowedPermissionMap := make(map[enums.PermissionType]bool)
	for _, permission := range allowedPermissions {
		allowedPermissionMap[permission] = true
	}

	for _, userRole := range userRoles {
		permissions := am.userRepository.FindPermissionsByRole(userRole.ID)
		for _, permission := range permissions {
			if allowedPermissionMap[permission] {
				return true
			}
		}
	}
	return false
}
