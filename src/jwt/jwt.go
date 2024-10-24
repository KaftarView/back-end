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

func setupJWTKeys(c *gin.Context, privateKeyPath, publicKeyPath, contextJWTKey string) {
	_, exists := c.Get(contextJWTKey)
	if !exists {
		loadPrivateKey(jwtToken, privateKeyPath)
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

func GenerateJWT(c *gin.Context, jwtKeysPath, contextJWTKey, username string) string {
	privateKeyPath := jwtKeysPath + "/privateKey.pem"
	publicKeyPath := jwtKeysPath + "/publicKey.pem"
	setupJWTKeys(c, privateKeyPath, publicKeyPath, contextJWTKey)
	claims := jwt.MapClaims{
		"iss": "test",
		"sub": username,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(jwtToken.PrivateKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func VerifyToken(tokenString string) jwt.MapClaims {
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
