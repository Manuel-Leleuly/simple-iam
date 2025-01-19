package helpers

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

func GetTranslation() ut.Translator {
	enUs := en.New()
	uni := ut.New(enUs, enUs)

	trans, found := uni.GetTranslator("en")
	if !found {
		panic("Failed to get translation for validation")
	}

	return trans
}
