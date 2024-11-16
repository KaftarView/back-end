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
