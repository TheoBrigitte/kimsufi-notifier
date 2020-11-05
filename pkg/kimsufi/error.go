package kimsufi

import (
	"net/http"
	"strings"

	"github.com/ovh/go-ovh/ovh"
)

func IsNotAvailableError(err error) bool {
	e, ok := err.(*ovh.APIError)
	if !ok {
		return false
	}

	return e.Code == http.StatusNotFound && strings.HasPrefix(e.Message, "No availabilities found")
}
