package application

import (
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: improve and use this logic in login
type JwtService struct {
	secretKey string
	constants *bootstrap.Constants
}

func NewJwtService(secretKey string, constants *bootstrap.Constants) *JwtService {
	return &JwtService{
		secretKey: secretKey,
		constants: constants,
	}
}

func (jwtService *JwtService) CreateToken(email string) string {
	expirationTime := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"exp": expirationTime,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(jwtService.secretKey)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func (jwtService *JwtService) VerifyToken(tokenString string) string {
	var registrationError exceptions.UserRegistrationError

	jwtSecret := []byte(jwtService.secretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email := claims["email"].(string)
		return email
	}
	registrationError.AppendError(
		jwtService.constants.ErrorField.Email,
		jwtService.constants.ErrorTag.InvalidToken)
	panic(registrationError)
}
