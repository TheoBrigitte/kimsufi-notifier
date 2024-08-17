package check

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ovh/go-ovh/ovh"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/sms"
)

// Cmd represents the check command
var (
	Cmd = &cobra.Command{
		Use:   "check",
		Short: "check availability",
		RunE:  runner,
	}

	datacenters []string
	planCode    string

	country   string
	hardware  string
	checkOnly bool
)

const (
	kimsufiAPI = ovh.OvhEU
	smsAPI     = "https://smsapi.free-mobile.fr/sendmsg"
)

func init() {
	Cmd.PersistentFlags().StringSliceVarP(&datacenters, "datacenters", "d", []string{"fr"}, "datacenter comma separated list")
	Cmd.PersistentFlags().StringVarP(&planCode, "planCode", "p", "", "plan code name (e.g. 1801sk143)")

	Cmd.PersistentFlags().StringVarP(&country, "country", "c", "", "country code (e.g. fr)")
	Cmd.PersistentFlags().StringVarP(&hardware, "hardware", "w", "", "harware code name (e.g. 1801sk143)")
	Cmd.PersistentFlags().BoolVarP(&checkOnly, "check-only", "o", true, "run check only (no sms)")
}

func runner(cmd *cobra.Command, args []string) error {
	d := kimsufi.Config{
		URL:    kimsufiAPI,
		Logger: log.StandardLogger(),
	}
	k, err := kimsufi.NewService(d)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	a, err := k.GetAvailabilities(datacenters, planCode)
	if kimsufi.IsNotAvailableError(err) {
		log.Printf("%s is not available in %s\n", hardware, country)
		os.Exit(0)
	} else if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	formatter := kimsufi.DatacenterFormatter(kimsufi.IsDatacenterAvailable, kimsufi.DatacenterKey)
	result := a.Format(kimsufi.PlanCode, formatter)
	//data, err := json.Marshal(result)
	//log.Printf("%s\n", data)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "planCode\tstatus\tdatacenters")
	fmt.Fprintln(w, "--------\t------\t-----------")

	for k, v := range result {
		status := "available"
		if len(v) == 0 {
			status = "unavailable"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", k, status, strings.Join(v, ", "))
	}

	w.Flush()
	//for _, availability := range *a {
	//	if availability.IsAvailable() {
	//		planDatacenters := formatter(availability.Datacenters)
	//		fmt.Printf("%s is available in following datacenters %v\n", availability.PlanCode, planDatacenters)
	//	} else {
	//		fmt.Printf("%s is not available\n", availability.PlanCode)
	//	}
	//}

	if checkOnly {
		os.Exit(0)
	}

	c := sms.Config{
		URL:    smsAPI,
		Logger: log.StandardLogger(),
		User:   "",
		Pass:   "",
	}

	s, err := sms.NewService(c)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	err = s.SendMessage("")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	log.Printf("message sent\n")

	return nil
}
