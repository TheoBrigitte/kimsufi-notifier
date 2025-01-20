package region

import (
	"strings"
)

// GetRegionFromCountry returns the region for a given country code.
func GetRegionFromCountry(country string) *Region {
	country = strings.ToUpper(country)

	for _, region := range AllowedRegions {
		for _, c := range region.Countries {
			if c.Code == country {
				return &region
			}
		}
	}

	return nil
}

// GetRegionFromEndpoint returns the region for a given endpoint.
func GetRegionFromEndpoint(endpoint string) *Region {
	for _, region := range AllowedRegions {
		if region.Endpoint == endpoint {
			return &region
		}
	}

	return nil
}
