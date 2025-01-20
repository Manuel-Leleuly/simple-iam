package validation

import (
	"github.com/Manuel-Leleuly/simple-iam/helpers"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator = helpers.GetTranslation()

func GetValidator() *validator.Validate {
	// create validation
	validate := validator.New(validator.WithRequiredStructEnabled())

	// register alias
	validate.RegisterAlias("username", "alphanum,min=5")
	validate.RegisterAlias("name", "alpha,min=2,max=10")
	validate.RegisterAlias("password", "required,min=8")

	// register translation
	en_translations.RegisterDefaultTranslations(validate, trans)

	RegisterTranslation(validate, "required", "{0} is required!")
	RegisterTranslation(validate, "email", "{0} incorrect email format")
	RegisterTranslation(validate, "name", "{0} must be alpha, min length 2, max length 10")
	RegisterTranslation(validate, "password", "{0} is required, min length 8")
	RegisterTranslation(validate, "username", "{0} must be alphanumeric, min length 5")

	return validate
}

func RegisterTranslation(validate *validator.Validate, key, message string) {
	validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
		return ut.Add(key, message, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(key, fe.Field())
		return t
	})
}

func TranslateValidationErrors(err error) []string {
	errs := err.(validator.ValidationErrors)

	translatedErrors := errs.Translate(trans)

	var errorStrings []string

	for k := range translatedErrors {
		errorStrings = append(errorStrings, translatedErrors[k])
	}

	return errorStrings
}
