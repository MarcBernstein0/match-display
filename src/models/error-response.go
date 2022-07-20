package models

type ErrorResponse struct {
	Message      string `json:"message"`
	ErrorMessage string `json:"error_message"`
}
