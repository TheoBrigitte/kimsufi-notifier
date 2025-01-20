package region

import (
	"fmt"
	"testing"
)

func TestGetRegionFromCountry(t *testing.T) {
	testCases := []struct {
		country  string
		expected *string
	}{
		{
			country:  "FR",
			expected: strPtr("Europe"),
		},
		{
			country:  "PT",
			expected: strPtr("Europe"),
		},
		{
			country:  "CZ",
			expected: strPtr("Europe"),
		},
		{
			country:  "AU",
			expected: strPtr("Other"),
		},
		{
			country:  "WS",
			expected: strPtr("Other"),
		},
		{
			country:  "US",
			expected: strPtr("US"),
		},
		{
			country:  "XX",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s -> %v", tc.country, tc.expected), func(t *testing.T) {
			actual := GetRegionFromCountry(tc.country)
			if tc.expected == nil && actual != nil {
				t.Errorf("expected nil, got %v", actual.DisplayName)
			} else if tc.expected != nil {
				if actual == nil {
					t.Errorf("expected %v, got nil", *tc.expected)
				} else if actual != nil && *tc.expected != actual.DisplayName {
					t.Errorf("expected %v, got %v", *tc.expected, actual.DisplayName)
				}
			}
		})
	}
}

func TestGetRegionFromEndpoint(t *testing.T) {
	testCases := []struct {
		endpoint string
		expected *string
	}{
		{
			endpoint: "ovh-eu",
			expected: strPtr("Europe"),
		},
		{
			endpoint: "ovh-ca",
			expected: strPtr("Other"),
		},
		{
			endpoint: "ovh-us",
			expected: strPtr("US"),
		},
		{
			endpoint: "ovh-xx",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s -> %v", tc.endpoint, tc.expected), func(t *testing.T) {
			actual := GetRegionFromEndpoint(tc.endpoint)
			if tc.expected == nil && actual != nil {
				t.Errorf("expected nil, got %v", actual.DisplayName)
			} else if tc.expected != nil {
				if actual == nil {
					t.Errorf("expected %v, got nil", *tc.expected)
				} else if actual != nil && *tc.expected != actual.DisplayName {
					t.Errorf("expected %v, got %v", *tc.expected, actual.DisplayName)
				}
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
