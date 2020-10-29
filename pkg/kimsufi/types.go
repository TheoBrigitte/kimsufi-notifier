package kimsufi

type Availabilities []Availability

type Availability struct {
	Hardware    string       `json:"hardware"`
	Region      string       `json:"region"`
	Datacenters []Datacenter `json:"datacenters"`
}

type Datacenter struct {
	Datacenter   string `json:"datacenter"`
	Availability string `json:"availability"`
}
