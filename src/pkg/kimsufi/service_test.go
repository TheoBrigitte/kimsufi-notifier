package kimsufi

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestNewMultiService(t *testing.T) {
	testCases := []struct {
		name     string
		logger   *log.Logger
		expected int
	}{
		{
			name:   "nil logger",
			logger: nil,
		},
		{
			name:   "standard logger",
			logger: log.StandardLogger(),
		},
		{
			name:   "dummy logger",
			logger: &log.Logger{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewMultiService(tc.logger, nil)
			if err != nil {
				t.Errorf("NewMultiService failed: %v", err)
			}
		})
	}
}

func TestMultiServiceEndpoint(t *testing.T) {
	m, err := NewMultiService(nil, nil)
	if err != nil {
		t.Errorf("NewMultiService failed: %v", err)
	}

	endpoints := []string{
		"ovh-eu",
		"ovh-ca",
		"ovh-us",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			s := m.Endpoint(endpoint)
			if s == nil {
				t.Errorf("expected endpoint %s to be found", endpoint)
			}
		})
	}
}

func TestNewService(t *testing.T) {
	testCases := []struct {
		name          string
		endpoint      string
		expectedError bool
	}{
		{
			name:          "valid endpoint",
			endpoint:      "ovh-eu",
			expectedError: false,
		},
		{
			name:          "invalid endpoint",
			endpoint:      "invalid",
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewService(tc.endpoint, nil, nil)
			if tc.expectedError && err == nil {
				t.Errorf("expected error for %s endpoint", tc.endpoint)
			}
		})
	}
}

func TestServiceLogger(t *testing.T) {
	s, err := NewService("ovh-eu", nil, nil)
	if err != nil {
		t.Errorf("NewService failed: %v", err)
	}

	if s.logger == nil {
		t.Error("expected logger to be set")
	}
}
