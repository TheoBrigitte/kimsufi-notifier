package availability

const (
	StatusAvailable   = "available"
	StatusUnavailable = "unavailable"
)

var (
	// DatacentersKnown is the non exhaustive list of known datacenters.
	// https://www.ovhcloud.com/en/about-us/global-infrastructure/expansion-regions-az/
	DatacentersKnown = DatacentersInfo{
		{Code: "bhs", Name: "Beauharnois (Canada)"},
		{Code: "fra", Name: "Frankfurt (Germany)"},
		{Code: "gra", Name: "Gravelines (France)"},
		{Code: "hil", Name: "Hillsboro (United States)"},
		{Code: "lon", Name: "London (United Kingdom)"},
		{Code: "par", Name: "Paris (France)"},
		{Code: "rbx", Name: "Roubaix (France)"},
		{Code: "sbg", Name: "Strasbourg (France)"},
		{Code: "sgp", Name: "Singapore"},
		{Code: "syd", Name: "Sydney (Australia)"},
		{Code: "vin", Name: "Vint Hill (United States)"},
		{Code: "waw", Name: "Warsaw (Poland)"},
		{Code: "ynm", Name: "Mumbai (India)"},
		{Code: "yyz", Name: "Toronto (Canada)"},
	}
)

type Datacenters []Datacenter

type Datacenter struct {
	Datacenter   string `json:"datacenter"`
	Availability string `json:"availability"`
}

type DatacentersInfo []DatacenterInfo

type DatacenterInfo struct {
	Code string
	Name string
}
