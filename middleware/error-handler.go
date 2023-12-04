package middleware

import (
	"log"
	"net/http"

	"test_iam/apierrors"
	"test_iam/generated/swagger/models"
)

// ServeError is a custom error handler for the API
func ServeError() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var apiErr *apierrors.APIError
		apiErr, ok := apierrors.ParseAPIError(err)
		if !ok {
			apiErr, ok = apierrors.ParseSwaggerAPIError(err)
			if !ok {
				apiErr = apierrors.NewAPIError(http.StatusInternalServerError, err.Error())
			}
		}
		errorResponse := models.Error{
			Message: apiErr.Message,
		}

		raw, marshalErr := errorResponse.MarshalBinary()
		if marshalErr != nil {
			http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.Code)

		_, writeErr := w.Write(raw)
		if writeErr != nil {
			log.Print("writing error response", "error", writeErr)
		}
	}
}
