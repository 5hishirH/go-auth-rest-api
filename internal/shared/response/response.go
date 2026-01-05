package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status bool
	Error  string
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err string) *ErrorResponse {
	return &ErrorResponse{
		Status: false,
		Error:  err,
	}
}

func HandleValidationErrors(w http.ResponseWriter, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMsgs []string

		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errMsgs = append(errMsgs, fmt.Sprintf("%s is required", e.Field()))
			case "email":
				errMsgs = append(errMsgs, fmt.Sprintf("%s is not a valid email", e.Field()))
			case "min":
				errMsgs = append(errMsgs, fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param()))
			default:
				errMsgs = append(errMsgs, fmt.Sprintf("%s is invalid", e.Field()))
			}
		}

		WriteJSON(w, http.StatusBadRequest, ErrorResponse{
			Status: false,
			Error:  strings.Join(errMsgs, ", "),
		})

		return
	}

	// Fallback for other errors
	HandleInternalError(w, "Internal Server Error")
}

func HandleInternalError(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusInternalServerError, *GeneralError(err))
}

func HandleBadRequest(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusBadRequest, *GeneralError(err))
}

func HandleConflict(w http.ResponseWriter, err string) {
	WriteJSON(w, http.StatusConflict, *GeneralError(err))
}

type ResponseWrapper struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // omitempty hides it if null
}

func CreatedOne(w http.ResponseWriter, fieldName string, data any) error {
	return WriteJSON(w, http.StatusCreated, ResponseWrapper{
		Success: true,
		Message: fmt.Sprintf("The %s is created successfully", fieldName),
		Data:    data,
	})
}
