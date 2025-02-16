package kimsufi

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ovh/go-ovh/ovh"
)

// IsAvailabilityNotFoundError checks if the error is an ovh.APIError
// which contains an availability not found error message.
func IsAvailabilityNotFoundError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusNotFound && strings.HasPrefix(ovhAPIError.Message, "No availabilities found")
}

// IsNotAvailableError checks if the error is an ovh.APIError
// which contains an not available error message.
func IsNotAvailableError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusBadRequest && strings.Contains(ovhAPIError.Message, "is not available in")
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

func IsPreferredPaymentMethodNotSetError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusBadRequest && strings.Contains(ovhAPIError.Message, "You do not have preferred payment method")
}

func IsPreferredPaymentMethodInvalidError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusBadRequest && strings.Contains(ovhAPIError.Message, "Your preferred payment method is not valid")
}

func IsPlanNotFoundError(err error) bool {
	var ovhAPIError *ovh.APIError
	if !errors.As(err, &ovhAPIError) {
		return false
	}

	return ovhAPIError.Code == http.StatusBadRequest && strings.HasPrefix(ovhAPIError.Message, "Plan code not found")
}
