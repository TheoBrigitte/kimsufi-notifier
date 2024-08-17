package list

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"text/tabwriter"

	"github.com/ovh/go-ovh/ovh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
)

var (
	Cmd = &cobra.Command{
		Use:   "list",
		Short: "list available servers",
		RunE:  runner,
	}

	datacenters   []string
	logLevel      string
	planCode      string
	ovhSubsidiary string
)

const (
	kimsufiAPI = ovh.OvhEU
)

func init() {
	Cmd.PersistentFlags().StringVarP(&ovhSubsidiary, "country", "c", "FR", "country code to filter entries (allowed values: CZ, DE, ES, FI, FR, GB, IE, IT, LT, MA, NL, PL, PT, SN, TN)")
	Cmd.PersistentFlags().StringSliceVarP(&datacenters, "datacenters", "d", nil, "comma separated list of datacenters to check (allowed values: bhs, ca, de, fr, fra, gb, gra, lon, pl, rbx, sbg, waw)")
	Cmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "error", "log level (allowed values: debug, info, warn, error, fatal, panic)")
	Cmd.PersistentFlags().StringVarP(&planCode, "plan-code", "p", "", "plan code name (e.g. 22sk011)")
}

func runner(cmd *cobra.Command, args []string) error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatalf("failed to parse log-level: %v\n", err)
	}
	log.SetLevel(level)

	d := kimsufi.Config{
		URL:    kimsufiAPI,
		Logger: log.StandardLogger(),
	}
	k, err := kimsufi.NewService(d)
	if err != nil {
		log.Fatalf("failed to initialize kimsufi service: %v\n", err)
	}

	c, err := k.ListServers(ovhSubsidiary)
	if err != nil {
		log.Fatalf("failed to list servers: %v\n", err)
	}

	a, err := k.GetAvailabilities(datacenters, planCode)
	if err != nil {
		log.Fatalf("failed to list availabilities: %v\n", err)
	}

	log.Infof("Found %d plans\n", len(c.Plans))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "planCode\tcategory\tname\tprice\tstatus\tdatacenters")
	fmt.Fprintln(w, "--------\t--------\t----\t-----\t------\t-----------")

	sort.Slice(c.Plans, func(i, j int) bool {
		return c.Plans[i].FirstPrice().Price < c.Plans[j].FirstPrice().Price
	})

	var nothingAvailable bool = true
	for _, category := range kimsufi.PlanCategories {
		for _, plan := range c.Plans {
			if planCode != "" && plan.PlanCode != planCode {
				continue
			}
			if plan.Blobs.Commercial.Range != category {
				continue
			}

			var price float64
			planPrice := plan.FirstPrice()
			if !reflect.DeepEqual(planPrice, kimsufi.Pricing{}) {
				price = float64(planPrice.Price) / kimsufi.PriceDivider
			}

			var status string
			datacenters := a.GetPlanCodeAvailableDatacenters(plan.PlanCode)
			if len(datacenters) == 0 {
				status = kimsufi.StatusUnavailable
			} else {
				nothingAvailable = false
				status = kimsufi.StatusAvailable
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%.2f\t%s\t%s\n", plan.PlanCode, category, plan.InvoiceName, price, status, strings.Join(datacenters, ", "))
		}
	}
	w.Flush()

	if nothingAvailable {
		log.Warnf("no server available\n")
		os.Exit(1)
	}

	return nil
}
