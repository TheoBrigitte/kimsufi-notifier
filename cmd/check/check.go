package check

import (
	"fmt"
	"os"

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

	country  string
	hardware string
	smsUser  string
	smsPass  string
)

const (
	kimsufiAPI = ovh.OvhEU
	smsAPI     = "https://smsapi.free-mobile.fr/sendmsg"
)

func init() {
	Cmd.PersistentFlags().StringVarP(&country, "country", "c", "", "country code (e.g. fr)")
	Cmd.PersistentFlags().StringVarP(&hardware, "hardware", "w", "", "harware code name (e.g. 1801sk143)")
	Cmd.PersistentFlags().StringVarP(&smsUser, "user", "u", "", "sms api username")
	Cmd.PersistentFlags().StringVarP(&smsPass, "pass", "p", "", "sms api password")
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

	_, err = k.GetAvailabilities(country, hardware)
	if kimsufi.IsNotAvailableError(err) {
		log.Printf("%s is not available in %s\n", hardware, country)
		os.Exit(0)
	} else if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	//formatter := kimsufi.DatacenterFormatter(kimsufi.IsDatacenterAvailable, kimsufi.DatacenterKey)
	//result := a.Format(kimsufi.HardwareKey, kimsufi.RegionKey, formatter)
	//data, err := json.Marshal(result)
	//fmt.Printf("%s\n", data)

	message := fmt.Sprintf("%s is available", hardware)
	fmt.Println(message)

	c := sms.Config{
		URL:    smsAPI,
		Logger: log.StandardLogger(),
		User:   smsUser,
		Pass:   smsPass,
	}

	s, err := sms.NewService(c)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	err = s.SendMessage(message)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	log.Printf("message sent\n")

	return nil
}
