package kimsufi

type Availabilities []Availability

type Availability struct {
	FQN         string       `json:"fqn"`
	Memory      string       `json:"memory"`
	PlanCode    string       `json:"planCode"`
	Server      string       `json:"server"`
	Storage     string       `json:"storage"`
	Datacenters []Datacenter `json:"datacenters"`
}

type Datacenter struct {
	Datacenter   string `json:"datacenter"`
	Availability string `json:"availability"`
}
