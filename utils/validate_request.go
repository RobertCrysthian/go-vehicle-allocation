package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ValidateRequest(request *http.Request, into any) error {
	if err := json.NewDecoder(request.Body).Decode(into); err != nil {
		return errors.New("A resposta deve ser um JSON válido")
	}
	if err := Validate.Struct(into); err != nil {
		errs := TranslateError(err)
		return errors.New(FormatErrorListToString(errs))
	}
	return nil
}