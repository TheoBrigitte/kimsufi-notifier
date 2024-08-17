package kimsufi

func (a Availabilities) GetPlanCodeAvailableDatacenters(planCode string) []string {
	var datacenters []string
	for _, availability := range a {
		if availability.PlanCode == planCode {
			for _, datacenter := range availability.Datacenters {
				if IsDatacenterAvailable(datacenter) && !inList(datacenters, datacenter.Datacenter) {
					datacenters = append(datacenters, datacenter.Datacenter)
				}
			}
		}
	}
	return datacenters
}

func (a Availability) IsAvailable() bool {
	for _, d := range a.Datacenters {
		if IsDatacenterAvailable(d) {
			return true
		}
	}

	return false
}
