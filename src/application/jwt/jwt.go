package application_jwt

import (
	"first-project/src/exceptions"
	jwt_keys "first-project/src/jwtKeys"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct{}

func NewJWTToken() *JWTToken {
	return &JWTToken{}
}

func (jt *JWTToken) GenerateJWT(userID uint) (string, string) {
	jwtKeys := jwt_keys.GetJWTKeys()
	accessTokenClaims := jwt.MapClaims{
		"iss": "test",
		"sub": userID,
		// "exp": time.Now().Add(time.Hour * 1).Unix(),
		"exp": time.Now().Add(time.Minute * 1).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtKeys.PrivateKey)
	if err != nil {
		panic(err)
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": "test",
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKeys.PrivateKey)
	if err != nil {
		panic(err)
	}

	return accessTokenString, refreshTokenString
}

func (jt *JWTToken) VerifyToken(tokenString string) jwt.MapClaims {
	jwtKeys := jwt_keys.GetJWTKeys()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return jwtKeys.PublicKey, nil
	})
	if err != nil {
		panic(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	unauthorizedError := exceptions.NewUnauthorizedError()
	panic(unauthorizedError)
}
