package application_jwt

import (
	"crypto/rsa"
	"first-project/src/bootstrap"
	"first-project/src/exceptions"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	constants *bootstrap.Constants
}

func NewJWTToken(constants *bootstrap.Constants) *JWTToken {
	return &JWTToken{
		constants: constants,
	}
}

type JWTKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var jwtKeys = &JWTKeys{}

func loadPrivateKey(jwtKeys *JWTKeys, privateKeyPath string) {
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes) // Use the instance
	if err != nil {
		panic(err)
	}
}

func loadPublicKey(jwtKeys *JWTKeys, publicKeyPath string) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
}

func (jt *JWTToken) setupJWTKeys(c *gin.Context, jwtKeysPath string) {
	_, exists := c.Get(jt.constants.Context.IsLoadedJWTKeys)
	if !exists {
		privateKeyPath := jwtKeysPath + "/privateKey.pem"
		loadPrivateKey(jwtKeys, privateKeyPath)
		publicKeyPath := jwtKeysPath + "/publicKey.pem"
		loadPublicKey(jwtKeys, publicKeyPath)
		c.Set(jt.constants.Context.IsLoadedJWTKeys, true)
	}
}

func (jt *JWTToken) GenerateJWT(c *gin.Context, jwtKeysPath string, userID uint) (string, string) {
	jt.setupJWTKeys(c, jwtKeysPath)
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

func (jt *JWTToken) SetAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	c.SetCookie(jt.constants.Context.AccessToken, accessToken, 60, "/", "localhost", true, true)
	c.SetCookie(jt.constants.Context.RefreshToken, refreshToken, 3600*24*7, "/", "localhost", true, true)
}

func (jt *JWTToken) VerifyToken(c *gin.Context, jwtKeysPath, tokenString string) jwt.MapClaims {
	jt.setupJWTKeys(c, jwtKeysPath)
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
