package localization

import (
	"fmt"

	"github.com/go-playground/locales/en_US"
	"github.com/go-playground/locales/fa_IR"
	ut "github.com/go-playground/universal-translator"
)

var translationMap = make(map[string]map[string]string)

func GetTranslator(locale string) ut.Translator {
	universalTranslator := createUniversalTranslator()
	loadAndAddTranslations(universalTranslator)

	translator, found := universalTranslator.GetTranslator(locale)
	if !found {
		translator, _ = universalTranslator.GetTranslator("fa_IR")
	}

	return translator
}

func createUniversalTranslator() *ut.UniversalTranslator {
	en := en_US.New()
	fa := fa_IR.New()
	return ut.New(en, en, fa)
}

func loadAndAddTranslations(universalTranslator *ut.UniversalTranslator) {
	addTranslations("fa_IR", Persian, universalTranslator)
	addTranslations("en_US", English, universalTranslator)
}

func addTranslations(locale string, translations map[string]interface{}, universalTranslator *ut.UniversalTranslator) {
	translator, found := universalTranslator.GetTranslator(locale)
	if !found {
		panic(fmt.Sprintf("translator for locale %s not found", locale))
	}

	flattenedTranslations := loadTranslations(locale, translations)

	for key, translation := range flattenedTranslations {
		translator.Add(key, translation, true)
	}
}

func loadTranslations(locale string, translations map[string]interface{}) map[string]string {
	if translations, ok := translationMap[locale]; ok {
		return translations
	}

	flattenedTranslations := make(map[string]string)
	flattenMap("", translations, flattenedTranslations)

	translationMap[locale] = flattenedTranslations

	return flattenedTranslations
}

func flattenMap(prefix string, input map[string]interface{}, output map[string]string) {
	for k, v := range input {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}
		switch value := v.(type) {
		case map[string]interface{}:
			flattenMap(fullKey, value, output)
		case string:
			output[fullKey] = value
		default:
			// Handle other types as needed, e.g., numbers, booleans, etc.
		}
	}
}
