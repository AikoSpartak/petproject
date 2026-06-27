package handler

import (
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Errors  string `json:"errors,omitempty"`
}

func ResponseHandler(msg, errors string) *Response {

	return &Response{
		Message: msg,
		Errors:  errors,
	}
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
}
