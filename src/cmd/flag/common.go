package flag

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/pkg/category"
	kimsufiavailability "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/availability"
)

const (
	CategoryFlagName = "category"

	DatacentersFlagName      = "datacenters"
	DatacentersFlagShortName = "d"

	OptionsFlagName      = "options"
	OptionsFlagShortName = "o"

	HumanFlagName      = "human"
	HumanFlagShortName = "h"

	PlanCodeFlagName      = "plan-code"
	PlanCodeFlagShortName = "p"
	PlanCodeExample       = "24ska01"
)

// BindCountryFlag binds the country flag to the provided cmd and value.
func BindCategoryFlag(cmd *cobra.Command, value *string) {
	categories := slices.DeleteFunc(category.Names(), func(s string) bool {
		return s == ""
	})
	allowedValues := strings.Join(categories, ", ")
	cmd.PersistentFlags().StringVar(value, CategoryFlagName, "", fmt.Sprintf("category to filter on (allowed values: %s)", allowedValues))
}

// BindDatacentersFlag binds the datacenters flag to the provided cmd and value.
func BindDatacentersFlag(cmd *cobra.Command, value *[]string) {
	cmd.PersistentFlags().StringSliceVarP(value, DatacentersFlagName, DatacentersFlagShortName, nil, fmt.Sprintf("datacenter(s) to filter on, comma separated list (known values: %s)", strings.Join(kimsufiavailability.GetDatacentersKnownCodes(), ", ")))
}

// BindOptionsFlag binds the datacenters flag to the provided cmd and value.
func BindOptionsFlag(cmd *cobra.Command, value *map[string]string) {
	cmd.PersistentFlags().StringToStringVarP(value, OptionsFlagName, OptionsFlagShortName, nil, "options in key=value format, comma separated list, keys are column names (e.g. memory=ram-64g-noecc-2133)")
}

// BindHumanFlag binds the verbose flag to the provided cmd and value.
// Warning: this redefine the help flag to only be a long --help flag.
func BindHumanFlag(cmd *cobra.Command, value *int) {
	cmd.PersistentFlags().CountVarP(value, HumanFlagName, HumanFlagShortName, "human output, more h makes it better (e.g. -h, -hh)")
	cmd.Flags().Bool("help", false, "help for "+cmd.Name())
}

// BindPlanCodeFlag binds the plan code flag to the provided cmd and value.
func BindPlanCodeFlag(cmd *cobra.Command, value *string) {
	cmd.PersistentFlags().StringVarP(value, PlanCodeFlagName, PlanCodeFlagShortName, "", fmt.Sprintf("plan code name (e.g. %s)", PlanCodeExample))
}
