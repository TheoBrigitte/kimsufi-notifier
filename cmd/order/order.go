package order

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/TheoBrigitte/kimsufi-notifier/cmd/flag"
	"github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi"
	kimsufiavailability "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/availability"
	kimsufiorder "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/order"
	kimsufiregion "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/region"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:   "order",
		Short: "Place an order",
		Long:  "Place an order for a servers from OVH Eco (including Kimsufi) catalog",
		RunE:  runner,
	}

	// Flags variables
	autoPay    bool
	datacenter string
	planCode   string
	quantity   int

	itemUserConfigurations map[string]string
	itemUserOptions        map[string]string

	listConfigurations bool
	listOptions        bool
	listPrices         bool

	priceDuration string
	priceMode     string

	ovhAppKeyEnvVarName      string
	ovhAppSecretEnvVarName   string
	ovhConsumerKeyEnvVarName string

	dryRun bool
)

func init() {
	flag.BindPlanCodeFlag(Cmd, &planCode)

	Cmd.PersistentFlags().BoolVar(&autoPay, "auto-pay", false, "automatically pay the order")
	Cmd.PersistentFlags().StringVarP(&datacenter, "datacenter", "d", "", fmt.Sprintf("datacenter (known values: %s)", strings.Join(kimsufiavailability.GetDatacentersKnownCodes(), ", ")))
	Cmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", kimsufiorder.QuantityDefault, "item quantity")

	Cmd.PersistentFlags().StringToStringVarP(&itemUserConfigurations, "item-configuration", "i", nil, "item configuration, see --list-configurations for available values (e.g. region=europe)")
	Cmd.PersistentFlags().StringToStringVarP(&itemUserOptions, "item-option", "o", nil, "item option, see --list-options for available values (e.g. memory=ram-64g-noecc-2133-24ska01)")

	Cmd.PersistentFlags().BoolVar(&listConfigurations, "list-configurations", false, "list available item configurations")
	Cmd.PersistentFlags().BoolVar(&listOptions, "list-options", false, "list available item options")
	Cmd.PersistentFlags().BoolVar(&listPrices, "list-prices", false, "list available prices")

	Cmd.PersistentFlags().StringVar(&priceMode, "price-mode", kimsufiorder.PricingMode, "price mode, see --list-prices for available values")
	Cmd.PersistentFlags().StringVar(&priceDuration, "price-duration", kimsufiorder.PriceDuration, "price duration, see --list-prices for available values")

	Cmd.PersistentFlags().StringVar(&ovhAppKeyEnvVarName, "ovh-app-key", "OVH_APP_KEY", "environement variable name for OVH API application key")
	Cmd.PersistentFlags().StringVar(&ovhAppSecretEnvVarName, "ovh-app-secret", "OVH_APP_SECRET", "environement variable name for OVH API application secret")
	Cmd.PersistentFlags().StringVar(&ovhConsumerKeyEnvVarName, "ovh-consumer-key", "OVH_CONSUMER_KEY", "environement variable name for OVH API consumer key")

	Cmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "only create a cart and do not submit the order")
}

func runner(cmd *cobra.Command, args []string) error {
	ovhSubsidiary := cmd.Flag(flag.CountryFlagName).Value.String()

	// Validate command arguments
	if planCode == "" {
		return fmt.Errorf("--plan-code is required")
	}
	if ovhSubsidiary == "" {
		return fmt.Errorf("--country is required")
	}

	// Initialize kimsufi service
	endpoint := cmd.Flag(flag.OVHAPIEndpointFlagName).Value.String()
	k, err := kimsufi.NewService(endpoint, log.StandardLogger(), nil)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Create cart
	expire := time.Now().AddDate(0, 0, 1)
	cart, err := k.CreateCart(ovhSubsidiary, expire)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Printf("> cart created id=%s\n", cart.CartID)

	// Retrieve item options
	ecoOptions, err := k.GetEcoOptions(cart.CartID, planCode)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	priceConfig := kimsufiorder.EcoItemPriceConfig{
		Duration:    priceDuration,
		PricingMode: priceMode,
	}

	if listOptions {
		printItemOptions(ecoOptions, priceConfig)
		return nil
	}

	// Retrieve item informations
	ecoInfo, err := k.GetEcoInfo(cart.CartID, planCode)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if listPrices {
		return printPrices(ecoInfo, planCode)
	}

	// Ensure price config is valid, otherwise use default
	priceConfig = ecoInfo.GetPriceConfigOrDefault(planCode, priceConfig)

	// Add plan to cart
	item, err := k.AddEcoItem(cart.CartID, planCode, quantity, priceConfig)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Printf("> cart item added id=%d\n", item.ItemID)

	requiredConfigurations, err := k.GetItemRequiredConfiguration(cart.CartID, item.ItemID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	if listConfigurations {
		printConfigurations(requiredConfigurations)
		return nil
	}

	if datacenter == "" {
		return fmt.Errorf("--datacenter is required")
	}

	// Prepare item configurations
	itemConfigurations := kimsufiorder.NewItemConfigurationsFromMap(itemUserConfigurations)
	itemConfigurations.Add(kimsufiorder.ConfigurationLabelDatacenter, datacenter)
	r := kimsufiregion.AllowedRegions.GetRegionFromEndpoint(endpoint)
	if r != nil {
		itemConfigurations.Add(kimsufiorder.ConfigurationLabelRegion, r.Region)
	}

	autoConfigs := k.GenerateItemAutoConfigurations(requiredConfigurations)
	userConfigs := autoConfigs.Merge(itemConfigurations)
	manualConfigs, err := generateItemManualConfiguration(userConfigs, requiredConfigurations)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	configurations := userConfigs.Merge(manualConfigs)

	// Configure item
	resp, err := k.ConfigureItem(cart.CartID, item.ItemID, configurations)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	for _, r := range resp {
		fmt.Printf("> cart item configured: %s=%s\n", r.Label, r.Value)
	}

	// Prepare item options
	userOptions := kimsufiorder.NewOptionsFromMap(itemUserOptions)

	// Configure item options
	options, err := k.ConfigureEcoItemOptions(cart.CartID, item.ItemID, ecoOptions, userOptions, priceConfig)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	for _, o := range options {
		fmt.Printf("> cart option set: %s=%s\n", o.Family, o.PlanCode)
	}

	// Stop on dry-run
	if dryRun {
		fmt.Println("> dry-run enabled, skipping order submission")
		return nil
	}

	// Read OVH API credentials from environment
	appKey := os.Getenv(cmd.Flag("ovh-app-key").Value.String())
	if appKey == "" {
		return fmt.Errorf("%s env var is required", cmd.Flag("ovh-app-key").Value.String())
	}
	appSecret := os.Getenv(cmd.Flag("ovh-app-secret").Value.String())
	if appSecret == "" {
		return fmt.Errorf("%s env var is required", cmd.Flag("ovh-app-secret").Value.String())
	}
	consumerKey := os.Getenv(cmd.Flag("ovh-consumer-key").Value.String())
	if consumerKey == "" {
		return fmt.Errorf("%s env var is required", cmd.Flag("ovh-consumer-key").Value.String())
	}

	// Authenticate
	k, err = k.WithAuth(appKey, appSecret, consumerKey)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Assign cart to user account
	err = k.AssignCart(cart.CartID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Println("> cart assigned")

	// Checkout and complete the order
	checkoutResp, err := k.CheckoutCart(cart.CartID, autoPay)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}
	fmt.Printf("> cart checked out %s\n", checkoutResp.URL)

	return nil
}

func printItemOptions(options []kimsufiorder.EcoItemOption, priceConfig kimsufiorder.EcoItemPriceConfig) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "item-option\tname\tprice")
	fmt.Fprintln(w, "-----------\t----\t-----")
	for _, o := range options {
		if !o.Mandatory {
			continue
		}

		price := ""
		p := o.GetPriceByConfig(priceConfig)
		if p != nil {
			price = p.Price.Text
		}

		fmt.Fprintf(w, "%s=%s\t%s\t%s\n", o.Family, o.PlanCode, o.ProductName, price)
	}
	w.Flush()
}

func printPrices(ecoInfos kimsufiorder.EcoItemInfos, planCode string) error {
	planInfo := ecoInfos.GetByPlanCode(planCode)
	if planInfo == nil {
		return fmt.Errorf("plan %s not found", planCode)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "price-duration\tprice-mode\tprice\tdescription")
	fmt.Fprintln(w, "--------------\t----------\t-----\t-----------")

	slices.SortFunc(planInfo.Prices, func(i, j kimsufiorder.EcoItemInfoPrice) int {
		return strings.Compare(i.Duration+i.PricingMode, j.Duration+j.PricingMode)
	})

	for _, p := range planInfo.Prices {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.Duration, p.PricingMode, p.Price.Text, p.Description)
	}
	w.Flush()

	return nil
}

func printConfigurations(configurations []kimsufiorder.ItemConfiguration) {
	fmt.Println("item-configuration")
	fmt.Println("------------------")

	for _, config := range configurations {
		if !config.Required || len(config.AllowedValues) <= 1 {
			continue
		}

		for _, value := range config.AllowedValues {
			fmt.Printf("%s=%s\n", config.Label, value)
		}
	}
}
