package middleware_i18n

import (
	"first-project/src/bootstrap"
	"first-project/src/localization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LocalizationMiddleware struct {
	constants *bootstrap.Constants
}

func NewLocalization(constants *bootstrap.Constants) *LocalizationMiddleware {
	return &LocalizationMiddleware{
		constants: constants,
	}
}

func (lm LocalizationMiddleware) Localization(c *gin.Context) {
	locale := getLocale(c.Request)
	c.Set(lm.constants.Context.Translator, localization.GetTranslator(locale))

	c.Next()
}

func getLocale(request *http.Request) string {
	return request.Header.Get("Accept-Language")
}
