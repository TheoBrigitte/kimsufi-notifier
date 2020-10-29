package kimsufi

func DatacenterFormatter(filter func(Datacenter) bool, format func(Datacenter) string) func([]Datacenter) []string {
	f := func(datacenters []Datacenter) []string {
		var result []string

		for _, d := range datacenters {
			if filter(d) {
				result = append(result, format(d))
			}
		}

		return result
	}

	return f
}

func IsDatacenterAvailable(d Datacenter) bool {
	return d.Availability != "unavailable"
}

func DatacenterKey(d Datacenter) string {
	return d.Datacenter
}
