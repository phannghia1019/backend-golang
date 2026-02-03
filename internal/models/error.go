package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError represents a structured application error
type AppError struct {
	Where       string `json:"where,omitempty"`
	Message     string `json:"message"`
	Detail      string `json:"detail,omitempty"`
	StatusCode  int    `json:"statusCode"`
	RedirectURL string `json:"redirectUrl,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(where, message, detail string, statusCode int) *AppError {
	return &AppError{
		Where:      where,
		Message:    message,
		Detail:     detail,
		StatusCode: statusCode,
	}
}

// Unauthorized creates an unauthorized error
func Unauthorized(where, message, detail string) *AppError {
	return NewAppError(where, message, detail, http.StatusUnauthorized)
}

// Forbidden creates a forbidden error
func Forbidden(where, message, detail string) *AppError {
	return NewAppError(where, message, detail, http.StatusForbidden)
}

// BadRequest creates a bad request error
func BadRequest(where, message, detail string) *AppError {
	return NewAppError(where, message, detail, http.StatusBadRequest)
}

// InternalError creates an internal server error
func InternalError(where, message, detail string) *AppError {
	return NewAppError(where, message, detail, http.StatusInternalServerError)
}

// HandleError sends error response to client
func HandleError(c *gin.Context, err *AppError) {
	c.JSON(err.StatusCode, err)
}
