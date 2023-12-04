package apierrors

import (
	"errors"

	openApiErrors "github.com/go-openapi/errors"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{code, message}
}

func ParseAPIError(err error) (*APIError, bool) {
	var apiError *APIError
	if errors.As(err, &apiError) {
		var unwrappedErr = err
		var resultErr *APIError
		ok := false
		for !ok {
			resultErr, ok = unwrappedErr.(*APIError) //nolint: errorlint
			unwrappedErr = errors.Unwrap(unwrappedErr)
		}
		return resultErr, true
	}
	return nil, false
}

func ParseSwaggerAPIError(err error) (*APIError, bool) {
	var apiError openApiErrors.Error
	if errors.As(err, &apiError) {
		var unwrappedErr = err
		var resultErr openApiErrors.Error
		ok := false
		for !ok {
			resultErr, ok = unwrappedErr.(openApiErrors.Error) //nolint: errorlint
			unwrappedErr = errors.Unwrap(unwrappedErr)
		}
		return &APIError{
			Code:    int(resultErr.Code()),
			Message: resultErr.Error(),
		}, true
	}
	return nil, false
}
