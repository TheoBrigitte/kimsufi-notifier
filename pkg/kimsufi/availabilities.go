package kimsufi

type Result map[string][]string

func (list Availabilities) Format(firstKey KeyFunc, formatDatacenters func([]Datacenter) []string) Result {
	result := make(Result)

	for _, a := range list {
		firstGroup := firstKey(a)
		if result[firstGroup] == nil {
			result[firstGroup] = []string{}
		}

		datacenters := formatDatacenters(a.Datacenters)
		if len(datacenters) > 0 {
			for _, d := range datacenters {
				if !inList(result[firstGroup], d) {
					result[firstGroup] = append(result[firstGroup], d)
				}
			}
		}
	}

	return result
}

func (list Availabilities) IsAvailable() bool {
	for _, a := range list {
		for _, d := range a.Datacenters {
			if IsDatacenterAvailable(d) {
				return true
			}
		}
	}

	return false
}

func inList(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}
