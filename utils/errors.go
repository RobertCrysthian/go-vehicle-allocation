package utils

import "net/http"

func BadRequestError(writer http.ResponseWriter, message string) {
	http.Error(writer, message, http.StatusBadRequest)
}

func InternalServerError(writer http.ResponseWriter, message string) {
	http.Error(writer, message, http.StatusBadRequest)
}

func ConflitctError(writer http.ResponseWriter, message string) {
	http.Error(writer, message, http.StatusConflict)
}