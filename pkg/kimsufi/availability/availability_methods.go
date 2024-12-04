package availability

import (
	"slices"
	"strings"
)

func (a Availabilities) GetByPlanCode(planCode string) Availabilities {
	var availabilities Availabilities

	for _, availability := range a {
		if availability.PlanCode == planCode {
			availabilities = append(availabilities, availability)
		}
	}

	return availabilities
}

// GetAvailableDatacenters returns the list of available datacenters.
func (a Availability) GetAvailableDatacenters() Datacenters {
	var datacenters []Datacenter

	for _, datacenter := range a.Datacenters {
		if datacenter.IsAvailable() && !slices.Contains(datacenters, datacenter) {
			datacenters = append(datacenters, datacenter)
		}
	}

	return datacenters
}

// GetPlanCodeAvailableDatacenters returns the list of available datacenters for a given plan code.
func (a Availabilities) GetAvailableDatacenters() Datacenters {
	var datacenters []Datacenter

	for _, availability := range a {
		datacenters = append(datacenters, availability.GetAvailableDatacenters()...)
	}

	slices.SortFunc(datacenters, func(i, j Datacenter) int {
		return strings.Compare(i.Datacenter, j.Datacenter)
	})

	uniqDatacenters := slices.CompactFunc(datacenters, func(i, j Datacenter) bool {
		return i.Datacenter == j.Datacenter
	})

	return uniqDatacenters
}
