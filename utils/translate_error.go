package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ptBRTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var Validate *validator.Validate
var Trans ut.Translator

func init () {
	Validate = validator.New()
	translateToPortuguese(Validate)
}

func TranslateError(err error) (errs []error) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(Trans))
		errs = append(errs, translatedErr)
	}

	return errs
}

func FormatErrorListToString (errs []error) (string) {
	var errorMessages []string
	for _, err := range errs {
		if err != nil {
			errorMessages = append(errorMessages, err.Error())
		}
	}
	return strings.Join(errorMessages, ", ")
}

func translateToPortuguese(validator *validator.Validate) {
	portuguese := pt_BR.New()
	uni := ut.New(portuguese, portuguese)
	trans, _ := uni.GetTranslator("pt_BR")
	Trans = trans
	_ = ptBRTranslations.RegisterDefaultTranslations(validator, Trans)
}