package controller

import (
	"first-project/src/exceptions"
	middleware_i18n "first-project/src/middleware/i18n"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/fa"
)

var validate *validator.Validate = validator.New()

func setupTranslation(trans ut.Translator) {
	if trans.Locale() == "fa_IR" {
		fa.RegisterDefaultTranslations(validate, trans)
	} else {
		en.RegisterDefaultTranslations(validate, trans)
	}
}

func Validated[T any](c *gin.Context) T {
	translator := middleware_i18n.GetTranslator(c)
	setupTranslation(translator)

	var params T
	if err := c.ShouldBindUri(&params); err != nil {
		bindingError := exceptions.BindingError{Err: err}
		panic(bindingError)
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		bindingError := exceptions.BindingError{Err: err}
		panic(bindingError)
	}

	if err := c.ShouldBind(&params); err != nil {
		bindingError := exceptions.BindingError{Err: err}
		panic(bindingError)
	}

	if err := validate.Struct(params); err != nil {
		panic(err)
	}

	return params
}
