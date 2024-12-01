package availability

// GetDatacentersKnownCodes returns the list of known datacenter codes.
func GetDatacentersKnownCodes() []string {
	var result []string

	for _, dc := range DatacentersKnown {
		result = append(result, dc.Code)
	}

	return result
}

// GetDatacenterInfoByCode returns the datacenter info.
// for a given datacenter code.
func GetDatacenterInfoByCode(code string) *DatacenterInfo {
	for _, dc := range DatacentersKnown {
		if dc.Code == code {
			return &dc
		}
	}

	return nil
}

// IsAvailable returns true if the datacenter is available.
func (d Datacenter) IsAvailable() bool {
	return d.Availability != StatusUnavailable
}

// Codes returns a list of datacenter codes.
func (d Datacenters) Codes() []string {
	var names []string

	for _, dc := range d {
		names = append(names, dc.Datacenter)
	}

	return names
}

// Status returns the status of the datacenters.
// Either StatusAvailable or StatusUnavailable.
func (d Datacenters) Status() string {
	if len(d) == 0 {
		return StatusUnavailable
	}

	return StatusAvailable
}

// ToFullNamesOrCodes returns the list of datacenter names.
// Either the full name if found or the code.
func (d Datacenters) ToFullNamesOrCodes() []string {
	var result []string

	for _, dc := range d {
		fullName := dc.GetFullName()
		if fullName != nil {
			result = append(result, *fullName)
		} else {
			result = append(result, dc.Datacenter)
		}
	}

	return result
}

// GetFullName returns the full name of the datacenter.
func (d Datacenter) GetFullName() *string {
	di := GetDatacenterInfoByCode(d.Datacenter)
	if di != nil {
		return &di.Name
	}

	return nil
}
