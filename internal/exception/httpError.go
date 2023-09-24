package exception

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type HttpError struct {
	Description string `json:"description,omitempty"`
	StatusCode  int    `json:"statusCode"`
	Metadata    error  `json:"metadata"`
}

func NewHttpError(statusCode int, description string, metadata error) HttpError {
	return HttpError{
		Description: description,
		StatusCode:  statusCode,
		Metadata:    metadata,
	}
}

func ParseError(err error) HttpError {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewHttpError(http.StatusNotFound, "Not Found", err)
	case strings.Contains(err.Error(), "strconv."):
		return NewHttpError(http.StatusBadRequest, "Bad Request", err)
	case strings.Contains(strings.ToLower(err.Error()), "param"):
		return NewHttpError(http.StatusBadRequest, "Invalid param", err)
	case strings.Contains(strings.ToLower(err.Error()), "email_unique"):
		return NewHttpError(http.StatusBadRequest, "User with given email already exists", err)
	case strings.Contains(strings.ToLower(err.Error()), "otp"):
		return NewHttpError(http.StatusBadRequest, "User is already verified", err)
	case strings.Contains(strings.ToLower(err.Error()), "invalid.token"):
		return NewHttpError(http.StatusBadRequest, "Invalid Token", err)
	default:
		return NewHttpError(http.StatusInternalServerError, "Internal Server Error", err)
	}
}

func (e HttpError) Error() string {
	return fmt.Sprintf("description: %s", e.Description)
}

func (e HttpError) Status() int {
	return e.StatusCode
}

func ErrorResponse(err error) (int, HttpError) {
	httpError := ParseError(err)
	fmt.Println(err.Error())
	return httpError.Status(), httpError
}
