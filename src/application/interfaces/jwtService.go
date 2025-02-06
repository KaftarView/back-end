package application_interfaces

import "github.com/golang-jwt/jwt/v5"

type JWTToken interface {
	GenerateJWT(userID uint) (string, string)
	VerifyToken(tokenString string) jwt.MapClaims
}
