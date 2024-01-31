package helpers

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func ValidateInput(input interface{}) []string {
	validate := validator.New()
	enligh := en.New()
	universalTranslator := ut.New(enligh, enligh)
	trans, _ := universalTranslator.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(input)
	errorMessages := translateErrors(err, trans)
	if errorMessages != nil {
		return errorMessages
	}

	return nil
}

func translateErrors(err error, trans ut.Translator) (errs []string) {
	if err == nil {
		return nil
	}

	errors := err.(validator.ValidationErrors)
	for _, e := range errors {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr.Error())
	}

	return errs
}
