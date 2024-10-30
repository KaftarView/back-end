package jwt

import (
	"crypto/rsa"
	"first-project/src/exceptions"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTToken struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var jwtToken = &JWTToken{}

func setupJWTKeys(c *gin.Context, jwtKeysPath, contextJWTKey string) {
	_, exists := c.Get(contextJWTKey)
	if !exists {
		privateKeyPath := jwtKeysPath + "/privateKey.pem"
		loadPrivateKey(jwtToken, privateKeyPath)
		publicKeyPath := jwtKeysPath + "/publicKey.pem"
		loadPublicKey(jwtToken, publicKeyPath)
		c.Set(contextJWTKey, jwtToken)
	}
}

func loadPrivateKey(jwtToken *JWTToken, privateKeyPath string) {
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	jwtToken.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes) // Use the instance
	if err != nil {
		panic(err)
	}
}

func loadPublicKey(jwtToken *JWTToken, publicKeyPath string) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	jwtToken.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
}

func GenerateJWT(c *gin.Context, jwtKeysPath, contextJWTKey string, userID uint) (string, string) {
	setupJWTKeys(c, jwtKeysPath, contextJWTKey)
	accessTokenClaims := jwt.MapClaims{
		"iss": "test",
		"sub": userID,
		// "exp": time.Now().Add(time.Hour * 1).Unix(),
		"exp": time.Now().Add(time.Minute * 1).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtToken.PrivateKey)
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
	refreshTokenString, err := refreshToken.SignedString(jwtToken.PrivateKey)
	if err != nil {
		panic(err)
	}

	return accessTokenString, refreshTokenString
}

func SetAuthCookies(c *gin.Context, accessToken, refreshToken, accessTokenKey, refreshTokenKey string) {
	c.SetCookie(accessTokenKey, accessToken, 60, "/", "localhost", true, true)
	c.SetCookie(refreshTokenKey, refreshToken, 3600*24*7, "/", "localhost", true, true)
}

func VerifyToken(c *gin.Context, jwtKeysPath, contextJWTKey, tokenString string) jwt.MapClaims {
	setupJWTKeys(c, jwtKeysPath, contextJWTKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return jwtToken.PublicKey, nil
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
