package resterr

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4"
)

type RestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *RestError) Error() string {
	if r == nil {
		return ""
	}
	return fmt.Sprintf("Code: %d, Message: %s", r.Code, r.Message)
}

func (r *RestError) Respond(c echo.Context) error {
	return c.JSON(r.Code, *r)
}

func NewInfraError(err error) *RestError {
	return &RestError{
		Code:    http.StatusGatewayTimeout,
		Message: err.Error(),
	}
}

func NewInternalServerError(err error) *RestError {
	fmt.Println("Error: ", err)
	return &RestError{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
}

func NewNotFoundError(entity string) *RestError {
	return &RestError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", entity),
	}
}

func NewBadRequestError(message string) *RestError {
	return &RestError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

var InvitationRequired = &RestError{
	Code:    http.StatusForbidden,
	Message: "Invitation required",
}

var InternalServerError = &RestError{
	Code:    http.StatusInternalServerError,
	Message: "Internal server error",
}

var BadRequest = &RestError{
	Code:    http.StatusBadRequest,
	Message: "Bad request",
}

var NotFound = &RestError{
	Code:    http.StatusNotFound,
	Message: "Not found",
}

var Unauthorized = &RestError{
	Code:    http.StatusUnauthorized,
	Message: "Unauthorized",
}
