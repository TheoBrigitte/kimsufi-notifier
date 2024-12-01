package flag

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
	kimsufiregion "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/region"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/logger"
)

const (
	LogLevelFlagName      = "log-level"
	LogLevelFlagShortName = "l"

	OVHAPIEndpointFlagName      = "endpoint"
	OVHAPIEndpointFlagShortName = "e"
	OVHAPIEndpointDefault       = "ovh-eu"

	CountryFlagName      = "country"
	CountryFlagShortName = "c"
	CountryDefault       = "FR"
)

// Bind binds the global flags to the provided cmd.
func Bind(cmd *cobra.Command) {
	// Redefine help flag to only be a long --help flag
	cmd.PersistentFlags().Bool("help", false, "help for "+cmd.Name())

	// Log level
	cmd.PersistentFlags().StringP(LogLevelFlagName, LogLevelFlagShortName, log.ErrorLevel.String(), fmt.Sprintf("log level (allowed values: %s)", strings.Join(logger.AllLevelsString(), ", ")))

	// OVH API Endpoint
	cmd.PersistentFlags().StringP(OVHAPIEndpointFlagName, OVHAPIEndpointFlagShortName, OVHAPIEndpointDefault, fmt.Sprintf("OVH API Endpoint (allowed values: %s)", strings.Join(kimsufi.GetOVHEndpoints(), ", ")))

	// Country
	// Display all countries per endpoint
	var output = &bytes.Buffer{}
	for _, region := range kimsufiregion.AllowedRegions {
		fmt.Fprintf(output, "  %s: ", region.Endpoint)

		countries := []string{}
		for _, c := range region.Countries {
			countries = append(countries, c.Code)
		}
		fmt.Fprintf(output, "%s\n", strings.Join(countries, ", "))
	}

	cmd.PersistentFlags().StringP(CountryFlagName, CountryFlagShortName, CountryDefault, fmt.Sprintf("country code, known values per endpoints:\n%s", output.String()))
}
