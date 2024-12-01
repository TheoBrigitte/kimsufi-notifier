package kimsufi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ovh/go-ovh/ovh"
)

// IsNotAvailableError checks if the error is an ovh.APIError
// which contains an availability error message.
func IsNotAvailableError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusNotFound && strings.HasPrefix(ovhAPIError.Message, "No availabilities found")
}

// IsForbiddenError checks if the error is an ovh.APIError
// with http.StatusForbidden code.
func IsForbiddenError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusForbidden
}
