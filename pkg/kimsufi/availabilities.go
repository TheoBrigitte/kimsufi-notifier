package kimsufi

type Result map[string]map[string][]string

func (list Availabilities) Format(firstKey, secondKey KeyFunc, formatDatacenters func([]Datacenter) []string) Result {
	result := make(Result)

	for _, a := range list {
		datacenters := formatDatacenters(a.Datacenters)
		if len(datacenters) > 0 {
			firstGroup := firstKey(a)
			secondGroup := secondKey(a)

			if result[firstGroup] == nil {
				result[firstGroup] = make(map[string][]string)
			}
			result[firstGroup][secondGroup] = append(result[firstGroup][secondGroup], datacenters...)
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
