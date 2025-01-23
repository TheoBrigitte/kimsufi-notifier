package list

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	pkgcategory "github.com/TheoBrigitte/kimsufi-notifier/pkg/category"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
	kimsufiavailability "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/availability"
	kimsuficatalog "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/catalog"
)

var (
	Cmd = &cobra.Command{
		Use:   "list",
		Short: "List available servers",
		Long:  "List servers from OVH Eco (including Kimsufi) catalog",
		Example: `  kimsufi-notifier list --category kimsufi
  kimsufi-notifier list --country US --endpoint ovh-us`,
		RunE: runner,
	}

	// Flags variables
	category    string
	datacenters []string
	humanLevel  int
	planCode    string
)

// init registers all flags
func init() {
	flag.BindCategoryFlag(Cmd, &category)
	flag.BindDatacentersFlag(Cmd, &datacenters)
	flag.BindHumanFlag(Cmd, &humanLevel)

	Cmd.PersistentFlags().StringVarP(&planCode, flag.PlanCodeFlagName, flag.PlanCodeFlagShortName, "", fmt.Sprintf("plan code to filter on (e.g. %s)", flag.PlanCodeExample))
}

// runner is the main function for the list command
func runner(cmd *cobra.Command, args []string) error {
	// Initialize kimsufi service
	endpoint := cmd.Flag(flag.OVHAPIEndpointFlagName).Value.String()
	k, err := kimsufi.NewService(endpoint, log.StandardLogger(), nil)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// List servers
	catalog, err := k.ListServers(cmd.Flag(flag.CountryFlagName).Value.String())
	if err != nil {
		return fmt.Errorf("failed to list servers: %w", err)
	}

	// List availabilities
	availabilities, err := k.GetAvailabilities(datacenters, planCode, nil)
	if err != nil {
		return fmt.Errorf("failed to list availabilities: %w", err)
	}

	// Display servers availabilities
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "planCode\tcategory\tname\tprice\tstatus\tdatacenters")
	fmt.Fprintln(w, "--------\t--------\t----\t-----\t------\t-----------")

	// Sort plans by category and price
	sort.Slice(catalog.Plans, func(i, j int) bool {
		planCategoryI := catalog.Plans[i].GetCategory()
		planCategoryJ := catalog.Plans[j].GetCategory()
		if planCategoryI != planCategoryJ {
			// Group plans by category first
			a := pkgcategory.GetDisplayName(planCategoryI)
			if a == "" {
				a = planCategoryI
			}

			b := pkgcategory.GetDisplayName(planCategoryJ)
			if b == "" {
				b = planCategoryJ
			}

			return a < b
		}

		// Then sort by price
		return catalog.Plans[i].GetFirstPrice().Price < catalog.Plans[j].GetFirstPrice().Price
	})

	// Display servers plans
	nothingAvailable := true
	for _, plan := range catalog.Plans {
		// Filter plans by plan code code
		if planCode != "" && plan.PlanCode != planCode {
			continue
		}

		planCategory := plan.GetCategory()

		// Filter plans by category
		if category != "" && category != planCategory {
			continue
		}

		// Format price
		var price float64
		planPrice := plan.GetFirstPrice()
		if !reflect.DeepEqual(planPrice, kimsuficatalog.PlanPricing{}) {
			price = planPrice.GetPrice()
		}

		// Format availability status
		datacenters := availabilities.GetByPlanCode(plan.PlanCode).GetAvailableDatacenters()

		var datacenterNames []string
		if humanLevel > 0 {
			datacenterNames = datacenters.ToFullNamesOrCodes()
		} else {
			datacenterNames = datacenters.Codes()
		}

		status := datacenters.Status()
		if status == kimsufiavailability.StatusAvailable {
			nothingAvailable = false
		}

		categoryDisplay := pkgcategory.GetDisplayName(planCategory)
		if categoryDisplay == "" {
			categoryDisplay = planCategory
		}

		// Display plan
		fmt.Fprintf(w, "%s\t%s\t%s\t%.2f %s\t%s\t%s\n", plan.PlanCode, categoryDisplay, plan.InvoiceName, price, catalog.Locale.CurrencyCode, status, strings.Join(datacenterNames, ", "))
	}
	w.Flush()

	if nothingAvailable {
		os.Exit(1)
	}

	return nil
}
