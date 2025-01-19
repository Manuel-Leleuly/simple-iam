package models

type ErrorMessage struct {
	Message string `json:"message"`
}

type ValidationErrorMessage struct {
	Message []string `json:"message"`
}
