package controller

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
)

func GetTranslator(c *gin.Context, key string) ut.Translator {
	translator, exists := c.Get(key)
	if !exists {
		panic("translator not registered!")
	}

	return translator.(ut.Translator)
}

func SetAuthCookies(c *gin.Context, accessToken, refreshToken, accessTokenKey, refreshTokenKey string) {
	c.SetCookie(accessTokenKey, accessToken, 60*15, "/", c.Request.Host, false, true)
	c.SetCookie(refreshTokenKey, refreshToken, 3600*24*7, "/", c.Request.Host, false, true)
}
