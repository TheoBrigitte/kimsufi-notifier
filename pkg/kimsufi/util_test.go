package kimsufi

import "testing"

func TestIntervalToDuration(t *testing.T) {
	testCases := []struct {
		interval int
		unit     string
		expected string
	}{
		{interval: 1, unit: "year", expected: "P1Y"},
		{interval: 2, unit: "year", expected: "P2Y"},
		{interval: 1, unit: "month", expected: "P1M"},
		{interval: 2, unit: "month", expected: "P2M"},
		{interval: 12, unit: "month", expected: "P1Y"},
		{interval: 13, unit: "month", expected: "P13M"},
		{interval: 1, unit: "day", expected: "P1D"},
		{interval: 2, unit: "day", expected: "P2D"},
		{expected: "P0D"},
	}

	for _, tc := range testCases {
		actual := IntervalToDuration(tc.interval, tc.unit)
		if actual != tc.expected {
			t.Errorf("expected %s, got %s", tc.expected, actual)
		}
	}
}
