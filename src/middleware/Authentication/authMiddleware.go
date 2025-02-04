package middleware_authentication

import (
	application_interfaces "first-project/src/application/interfaces"
	"first-project/src/bootstrap"
	"first-project/src/entities"
	"first-project/src/enums"
	"first-project/src/exceptions"
	jwt_keys "first-project/src/jwtKeys"
	repository_database_interfaces "first-project/src/repository/database/interfaces"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	userRepository repository_database_interfaces.UserRepository
	jwtService     application_interfaces.JWTToken
	db             *gorm.DB
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	userRepository repository_database_interfaces.UserRepository,
	jwtService application_interfaces.JWTToken,
	db *gorm.DB,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:      constants,
		userRepository: userRepository,
		jwtService:     jwtService,
		db:             db,
	}
}

func (am *AuthMiddleware) AuthRequired(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}

	tokenString := parts[1]
	if tokenString == "" {
		unauthorizedError := exceptions.NewUnauthorizedError()
		panic(unauthorizedError)
	}
	jwt_keys.SetupJWTKeys(c, am.constants.Context.IsLoadedJWTKeys, am.constants.JWTKeysPath)
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

		roles := am.userRepository.FindUserRoleTypesByUserID(am.db, userID.(uint))
		allowedPermissions = append(allowedPermissions, enums.All)
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
		permissions := am.userRepository.FindPermissionsByRole(am.db, userRole.ID)
		for _, permission := range permissions {
			if allowedPermissionMap[permission] {
				return true
			}
		}
	}
	return false
}
