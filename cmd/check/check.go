package check

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
	kimsufiavailability "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/availability"
	kimsuficatalog "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/catalog"
)

var (
	Cmd = &cobra.Command{
		Use:   "check",
		Short: "Check server availability",
		Long:  "Check OVH Eco (including Kimsufi) server availability\n\ndatacenters are the available datacenters for this plan",
		Example: `  kimsufi-notifier check --plan-code 24ska01
  kimsufi-notifier check --plan-code 24ska01 --datacenters gra,rbx`,
		RunE: runner,
	}

	// Flags variables
	datacenters []string
	options     map[string]string
	planCode    string
	humanLevel  int

	listOptions bool
)

// init registers all flags
func init() {
	flag.BindPlanCodeFlag(Cmd, &planCode)
	flag.BindDatacentersFlag(Cmd, &datacenters)

	Cmd.PersistentFlags().CountVarP(&humanLevel, "human", "h", "Human output, more h makes it better (e.g. -h, -hh)")
	// Redefine help flag to only be a long --help flag, to avoid conflict with human flag
	Cmd.Flags().Bool("help", false, "help for "+Cmd.Name())

	Cmd.PersistentFlags().BoolVar(&listOptions, "list-options", false, "list available item options")
	Cmd.PersistentFlags().StringToStringVarP(&options, "option", "o", nil, "options to filter on, comma separated list of key=value, see --list-options for available options (e.g. memory=ram-64g-noecc-2133)")
}

// runner is the main function for the check command
func runner(cmd *cobra.Command, args []string) error {
	// Initialize kimsufi service
	endpoint := cmd.Flag(flag.OVHAPIEndpointFlagName).Value.String()
	k, err := kimsufi.NewService(endpoint, log.StandardLogger(), nil)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Flag validation
	if planCode == "" {
		return fmt.Errorf("--%s is required", flag.PlanCodeFlagName)
	}

	var catalog *kimsuficatalog.Catalog
	if humanLevel > 0 || listOptions {
		// Get the catalog to display human readable information.
		catalog, err = k.ListServers(cmd.Flag(flag.CountryFlagName).Value.String())
		if err != nil {
			return fmt.Errorf("failed to list servers: %w", err)
		}
	}

	if listOptions {
		return printItemOptions(catalog, planCode)
	}

	// Check availability
	availabilities, err := k.GetAvailabilities(datacenters, planCode, options)
	if err != nil {
		if kimsufi.IsNotAvailableError(err) {
			message := datacenterAvailableMessageFormatter(datacenters)
			log.Printf("%s is not available in %s\n", planCode, message)
			return nil
		}

		return fmt.Errorf("error: %w", err)
	}

	// Display the server availabilities for each options.
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "planCode\tmemory\tstorage\tstatus\tdatacenters")
	fmt.Fprintln(w, "--------\t------\t-------\t------\t-----------")

	nothingAvailable := true
	for _, v := range *availabilities {
		var (
			name        = v.PlanCode
			memory      = v.Memory
			storage     = v.Storage
			datacenters = v.GetAvailableDatacenters()
		)

		if humanLevel > 0 {
			plan := catalog.GetPlan(v.PlanCode)
			if plan != nil {
				names := strings.Split(plan.InvoiceName, " | ")
				name = names[0]
			}

			memoryProduct := catalog.GetProduct(memory)
			if memoryProduct != nil {
				memory = memoryProduct.Description
			}

			storageProduct := catalog.GetProduct(storage)
			if storageProduct != nil {
				storage = storageProduct.Description
			}
		}

		var datacenterNames []string
		if humanLevel > 1 {
			datacenterNames = datacenters.ToFullNamesOrCodes()
		} else {
			datacenterNames = datacenters.Codes()
		}

		var status = datacenters.Status()
		if status == kimsufiavailability.StatusAvailable {
			nothingAvailable = false
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", name, memory, storage, status, strings.Join(datacenterNames, ", "))
	}
	w.Flush()

	if nothingAvailable {
		os.Exit(1)
	}

	return nil
}

func datacenterAvailableMessageFormatter(datacenters []string) string {
	var message string

	switch len(datacenters) {
	case 0:
		message = "any datacenter"
	case 1:
		message = datacenters[0] + " datacenter"
	default:
		message = strings.Join(datacenters, ", ") + " datacenters"
	}

	return message
}

func printItemOptions(catalog *kimsuficatalog.Catalog, planCode string) error {
	plan := catalog.GetPlan(planCode)
	if plan == nil {
		return fmt.Errorf("plan %s not found", planCode)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "option\tdescription\tdefault")
	fmt.Fprintln(w, "-----------\t-----------\t-------")

	for _, addon := range plan.AddonFamilies {
		if !addon.Mandatory {
			continue
		}

		for _, value := range addon.Addons {
			isDefault := value == addon.Default
			description := ""
			genericValue := kimsufi.AddonGenericName(value)
			product := catalog.GetProduct(genericValue)
			if product != nil {
				description = product.Description
			}

			fmt.Fprintf(w, "%s=%s\t%s\t%t\n", addon.Name, genericValue, description, isDefault)
		}
	}
	w.Flush()

	return nil
}
