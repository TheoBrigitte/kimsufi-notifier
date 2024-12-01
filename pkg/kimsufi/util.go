package kimsufi

import (
	"fmt"
	"strings"
)

// AddonGenericName returns the generic name of an addon.
// It strips out the last part of the name after the last dash.
// e.g. ram-64g-ecc-2400-24sk50 -> ram-64g-ecc-2400
func AddonGenericName(name string) string {
	index := strings.LastIndex(name, "-")
	if index > 0 {
		return name[:index]
	}

	return name
}

// IntervalToDuration converts an interval and a unit to a duration string.
// examples:
// - 1  year   -> P1Y
// - 2  years  -> P2Y
// - 1  month  -> P1M
// - 2  months -> P2M
// - 12 months -> P1Y
// - 0  days   -> P0D
// - 1  day    -> P1D
// - 2  days   -> P2D
// - 31 days   -> P31D
func IntervalToDuration(interval int, unit string) string {
	switch unit {
	case "year":
		return fmt.Sprintf("P%dY", interval)
	case "month":
		if interval%12 == 0 {
			return fmt.Sprintf("P%dY", interval/12)
		}

		return fmt.Sprintf("P%dM", interval)
	case "day":
		return fmt.Sprintf("P%dD", interval)
	}

	return "P0D"
}
