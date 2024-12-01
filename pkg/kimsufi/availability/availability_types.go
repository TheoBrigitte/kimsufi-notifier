package availability

type Availabilities []Availability

// Availability represents the availability of a server.
// Definition can be found at https://eu.api.ovh.com/console/?section=%2Fdedicated%2Fserver&branch=v1#get-/dedicated/server/datacenter/availabilities
type Availability struct {
	FQN         string       `json:"fqn"`
	Memory      string       `json:"memory"`
	PlanCode    string       `json:"planCode"`
	Server      string       `json:"server"`
	Storage     string       `json:"storage"`
	Datacenters []Datacenter `json:"datacenters"`
}
