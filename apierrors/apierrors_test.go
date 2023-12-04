package apierrors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-openapi/errors"
	"github.com/stretchr/testify/require"
)

func TestParseAPIError_Valid(t *testing.T) {
	code := http.StatusInternalServerError
	message := "test"
	err := NewAPIError(code, message)
	apiErr, ok := ParseAPIError(err)
	require.True(t, ok)
	require.NotNil(t, apiErr)
	require.Equal(t, code, apiErr.Code)
	require.Equal(t, err.Message, apiErr.Message)
}

func TestParseAPIError_Wrap(t *testing.T) {
	code := http.StatusInternalServerError
	message := "test"
	err := NewAPIError(code, message)
	wrappedErr := fmt.Errorf("wraped error: %w", err)
	apiErr, ok := ParseAPIError(wrappedErr)
	require.True(t, ok)
	require.NotNil(t, apiErr)
	require.Equal(t, code, apiErr.Code)
	require.Equal(t, err.Message, apiErr.Message)
}

func TestParseAPIError_Invalid(t *testing.T) {
	err := fmt.Errorf("test")
	apiErr, ok := ParseAPIError(err)
	require.False(t, ok)
	require.Nil(t, apiErr)
}

func TestParseSwaggerAPIError_Valid(t *testing.T) {
	code := http.StatusInternalServerError
	message := "test"
	err := errors.New(int32(code), message)
	apiErr, ok := ParseSwaggerAPIError(err)
	require.True(t, ok)
	require.NotNil(t, apiErr)
	require.Equal(t, code, apiErr.Code)
	require.Equal(t, message, apiErr.Message)
}

func TestParseSwaggerAPIError_Wrap(t *testing.T) {
	code := http.StatusInternalServerError
	message := "test"
	err := errors.New(int32(code), message)
	wrappedErr := fmt.Errorf("wraped error: %w", err)
	apiErr, ok := ParseSwaggerAPIError(wrappedErr)
	require.True(t, ok)
	require.NotNil(t, apiErr)
	require.Equal(t, code, apiErr.Code)
	require.Equal(t, message, apiErr.Message)
}

func TestParseSwaggerAPIError_Invalid(t *testing.T) {
	err := fmt.Errorf("test")
	apiErr, ok := ParseSwaggerAPIError(err)
	require.False(t, ok)
	require.Nil(t, apiErr)
}
