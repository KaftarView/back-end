package controller

import (
	"first-project/src/bootstrap"

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

func GetTemplatePath(c *gin.Context, transKey string) string {
	trans := GetTranslator(c, transKey)
	if trans.Locale() == "fa_IR" {
		return "fa.html"
	}
	return "en.html"
}

type PaginationParams struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`
}

func GetPagination(c *gin.Context, constants *bootstrap.Context) PaginationParams {
	param := Validated[PaginationParams](c, constants)

	if param.Page == 0 {
		param.Page = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	return param
}
