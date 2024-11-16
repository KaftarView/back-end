package jwt_keys

import (
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var jwtKeys = &JWTKeys{}

func loadPrivateKey(privateKeyPath string) {
	privKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBytes) // Use the instance
	if err != nil {
		panic(err)
	}
}

func loadPublicKey(publicKeyPath string) {
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
}

func SetupJWTKeys(c *gin.Context, isLoadedKeys, jwtKeysPath string) {
	_, exists := c.Get(isLoadedKeys)
	if !exists {
		privateKeyPath := jwtKeysPath + "/privateKey.pem"
		loadPrivateKey(privateKeyPath)
		publicKeyPath := jwtKeysPath + "/publicKey.pem"
		loadPublicKey(publicKeyPath)
		c.Set(isLoadedKeys, true)
	}
}

func GetJWTKeys() *JWTKeys {
	if jwtKeys.PrivateKey == nil || jwtKeys.PublicKey == nil {
		panic(fmt.Errorf("private key or public key are nil"))
	}
	return jwtKeys
}
