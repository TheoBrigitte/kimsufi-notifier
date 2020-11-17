package order

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ovh/go-ovh/ovh"

	"github.com/TheoBrigitte/kimsufi-notifier/pkg/sms"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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
	smsUser     string
	smsPass     string

	country          string
	hardware         string
	quantity         int
	paymentMethod    int
	paymentFrequency string

	timeout    time.Duration
	screenshot string
	dryRun     bool
)

const (
	kimsufiAPI = ovh.OvhEU
	smsAPI     = "https://smsapi.free-mobile.fr/sendmsg"
)

func init() {
	Cmd.PersistentFlags().StringVarP(&country, "country", "c", "fr", "country code")
	Cmd.PersistentFlags().StringVarP(&hardware, "hardware", "w", "", "hardware code name (e.g. 1801sk143)")
	Cmd.PersistentFlags().StringVarP(&kimsufiUser, "kimsufi-user", "u", "", "kimsufi api username")
	Cmd.PersistentFlags().StringVarP(&kimsufiPass, "kimsufi-pass", "p", "", "kimsufi api password")
	Cmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 30*time.Second, "command timeout")
	Cmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", 1, "quantity of hardware to order")
	Cmd.PersistentFlags().IntVarP(&paymentMethod, "payment-method", "m", 1, "payment method index")
	Cmd.PersistentFlags().StringVarP(&paymentFrequency, "frequency", "f", "Mensuelle", "payement frequency (Mensuelle, Trimestrielle, Semestrielle, or Annuelle)")
	Cmd.PersistentFlags().StringVarP(&screenshot, "screenshot", "s", "kimsufi-order.png", "screenshot filename")
	Cmd.PersistentFlags().StringVar(&smsUser, "sms-user", "", "sms api username")
	Cmd.PersistentFlags().StringVar(&smsPass, "sms-pass", "", "sms api password")
	Cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "do not submit order")
}

func runner(cmd *cobra.Command, args []string) error {
	ordered, err := alreadyOrdered(screenshot)
	if err != nil {
		return err
	}

	if ordered {
		log.Printf("screenshot %s already exists.", screenshot)
		log.Println("stopping, not to order multiple times.\n")
		return nil
	}

	u := fmt.Sprintf("https://www.kimsufi.com/fr/commande/kimsufi.xml?reference=%s", hardware)

	// create context
	execOptions := []chromedp.ExecAllocatorOption(chromedp.DefaultExecAllocatorOptions[:])
	execOptions = append(execOptions, chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:82.0) Gecko/20100101 Firefox/82.0"))
	ctx := context.Background()
	ctx, _ = chromedp.NewExecAllocator(ctx, execOptions...)
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	log.Printf("open: %s\n", u)
	err = chromedp.Run(ctx)
	if err != nil {
		return err
	}

	{
		// run task list
		ctx, _ := context.WithTimeout(ctx, timeout)
		ctx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		//defer chromedp.Run(ctx, )
		err = chromedp.Run(ctx,
			chromedp.Navigate(u),
			chromedp.Click("button#header_tc_privacy_button", chromedp.NodeVisible),
			isAvailable(),
			setQuantity(quantity),
			setPaymentFrequency(paymentFrequency),
			login(kimsufiUser, kimsufiPass),
			selectPayement(paymentMethod),
			chromedp.Click("#contracts-validation"),
			chromedp.Click("#customConractAccepted"),
			confirm(),
			waitNextPage(5*time.Second),
			fullScreenshot(90, screenshot),
		)
		if err != nil {
			return err
		}
	}

	message := fmt.Sprintf("%s ordered\ncheck your mail", hardware)
	log.Println(message)

	c := sms.Config{
		URL:    smsAPI,
		Logger: log.StandardLogger(),
		User:   smsUser,
		Pass:   smsPass,
	}

	s, err := sms.NewService(c)
	if err != nil {
		return err
	}

	err = s.SendMessage(message)
	if err != nil {
		return err
	}
	log.Printf("message sent\n")

	return nil
}

func alreadyOrdered(screenshot string) (bool, error) {
	// Only order once.
	_, err := os.Stat(screenshot)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func isAvailable() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var page string
		f := chromedp.Text("#main", &page, chromedp.NodeVisible, chromedp.ByID)
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		ok := strings.Contains(page, "Récapitulatif de votre commande")
		if !ok {
			return fmt.Errorf("order not available")
		}

		log.Println("order available")

		return nil
	})
}

func setQuantity(desired int) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		//var debug string
		//f := chromedp.TextContent("tbody.configuration tr.editable:nth-child(2) > td:nth-child(3)", &debug)
		//err := f.Do(ctx)
		//if err != nil {
		//	return err
		//}
		//fmt.Printf("debug\n%s\n", debug)

		var current string
		f := chromedp.TextContent("#main tbody.configuration tr.editable:nth-child(2) > td:nth-child(3)", &current)
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		c, err := strconv.Atoi(current)
		if err != nil {
			return err
		}

		log.Printf("quantity: current=%d desired=%d\n", c, desired)

		if c == desired {
			return nil
		}

		f = chromedp.Click(fmt.Sprintf("#main tbody.configuration tr.editable:nth-child(2) > td:nth-child(2) > ul:nth-child(1) > li:nth-child(%d) > label:nth-child(2)", desired))
		err = f.Do(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func setPaymentFrequency(desired string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		//var debug string
		//f := chromedp.TextContent("tbody.configuration tr.editable:nth-child(2) > td:nth-child(3)", &debug)
		//err := f.Do(ctx)
		//if err != nil {
		//	return err
		//}
		//fmt.Printf("debug\n%s\n", debug)

		var current string

		f := chromedp.TextContent("#main tbody.configuration tr.editable:nth-child(3) > td:nth-child(3)", &current)
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		log.Printf("frequency: current=%s desired=%s\n", current, desired)

		if current == desired {
			return nil
		}

		f = chromedp.Click(fmt.Sprintf("//label[contains(text(),'%s')]", desired), chromedp.BySearch)
		err = f.Do(ctx)
		if err != nil {
			return err
		}

		return nil
	})
}

func login(user, pass string) chromedp.Tasks {
	messageStart := chromedp.ActionFunc(func(ctx context.Context) error {
		log.Printf("login user=%s pass=%s\n", user, strings.Repeat("x", len(pass)))
		return nil
	})
	selectLogin := chromedp.Click("#main div.customer div.you-are span.existing label")
	inputeUser := chromedp.SendKeys("#existing-customer-login", user, chromedp.ByID)
	inputePassword := chromedp.SendKeys("#existing-customer-password", pass, chromedp.ByID)
	submitLogin := chromedp.Click("div.customer-existing form span.ec-button span.middle button span", chromedp.NodeVisible)
	wait := chromedp.WaitEnabled("#contracts-validation")
	isErr := chromedp.ActionFunc(func(ctx context.Context) error {
		var page string
		f := chromedp.Text("#main", &page, chromedp.NodeVisible, chromedp.ByID)
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		ok := strings.Contains(page, "Mauvais identifiant ou mot-de-passe")
		if ok {
			return fmt.Errorf("wrong crendetials")
		}

		return nil
	})
	messageEnd := chromedp.ActionFunc(func(ctx context.Context) error {
		log.Println("logged in")
		return nil
	})

	return []chromedp.Action{
		messageStart,
		selectLogin,
		inputeUser,
		inputePassword,
		submitLogin,
		isErr,
		wait,
		messageEnd,
	}
}

func selectPayement(index int) chromedp.Tasks {
	// offset by one, to skip header.
	i := index + 1
	wait := chromedp.WaitVisible(".payment-means form", chromedp.ByQuery)
	click := chromedp.Click(fmt.Sprintf(".payment-means form span:nth-child(%d) > .first > label", i), chromedp.NodeVisible, chromedp.ByQuery)
	debug := chromedp.ActionFunc(func(ctx context.Context) error {
		var payment string
		f := chromedp.Text(".payment-means form .selected > .type > label", &payment, chromedp.NodeVisible, chromedp.ByQuery)
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		log.Printf("payement=%s\n", strings.ReplaceAll(payment, "\n", " - "))

		return nil
	})

	return []chromedp.Action{
		wait,
		click,
		debug,
	}
}

func confirm() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		f := chromedp.Click(".dedicated-contracts div.center:nth-child(2) button")
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		log.Println("confirm")

		return nil
	})
}

func waitNextPage(duration time.Duration) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		f := chromedp.WaitNotVisible(".dedicated-contracts div.center:nth-child(2) button")
		err := f.Do(ctx)
		if err != nil {
			return err
		}

		log.Printf("sleeping=%v\n", duration)

		time.Sleep(duration)

		return nil
	})
}

// fullScreenshot takes a screenshot of the entire browser viewport.
//
// Liberally copied from puppeteer's source.
//
// Note: this will override the viewport emulation settings.
func fullScreenshot(quality int64, filename string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// get layout metrics
		_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// capture screenshot
		buf, err := page.CaptureScreenshot().
			WithQuality(quality).
			WithClip(&page.Viewport{
				X:      contentSize.X,
				Y:      contentSize.Y,
				Width:  contentSize.Width,
				Height: contentSize.Height,
				Scale:  1,
			}).Do(ctx)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
			return err
		}

		log.Printf("took screenshots: %s\n", filename)

		return nil
	})
}
