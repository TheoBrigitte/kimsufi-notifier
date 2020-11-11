package order

import (
	"context"
	"strings"

	"github.com/ovh/go-ovh/ovh"

	"github.com/chromedp/chromedp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Cmd represents the check command
var (
	Cmd = &cobra.Command{
		Use:   "order",
		Short: "order it",
		RunE:  runner,
	}

	kimsufiUser string
	kimsufiPass string
	kimsufiKey  string
	country     string
	hardware    string
)

const (
	kimsufiAPI = ovh.OvhEU
	smsAPI     = "https://smsapi.free-mobile.fr/sendmsg"
)

func init() {
	Cmd.PersistentFlags().StringVarP(&country, "country", "c", "", "country code (e.g. fr)")
	Cmd.PersistentFlags().StringVarP(&hardware, "hardware", "w", "", "harware code name (e.g. 1801sk143)")
	Cmd.PersistentFlags().StringVarP(&kimsufiUser, "kimsufi-user", "u", "", "kimsufi api username")
	Cmd.PersistentFlags().StringVarP(&kimsufiPass, "kimsufi-pass", "p", "", "kimsufi api password")
}

func runner(cmd *cobra.Command, args []string) error {
	log.Printf("hello\n")

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		chromedp.Text(`#pkg-overview`, &res, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(strings.TrimSpace(res))

	return nil
}
