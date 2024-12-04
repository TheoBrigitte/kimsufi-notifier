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

const (
	anyOption = "any"
)

var (
	Cmd = &cobra.Command{
		Use:   "order",
		Short: "Place an order",
		Long:  "Place an order for a servers from OVH Eco (including Kimsufi) catalog",
		Example: `  kimsufi-notifier order --plan-code 24ska01 --datacenter rbx --dry-run
  kimsufi-notifier order --plan-code 25skle01 --datacenter bhs --item-option memory=ram-32g-noecc-1333-25skle01,storage=softraid-3x2000sa-25skle01`,
		RunE: runner,
	}

	// Flags variables
	autoPay     bool
	datacenters []string
	planCode    string
	quantity    int

	itemUserConfigurations map[string]string
	itemUserOptions        []string

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
	Cmd.PersistentFlags().StringSliceVarP(&datacenters, "datacenters", "d", nil, fmt.Sprintf(`datacenters, comma separated list, %q to try all datacenters (known values: %s)`, anyOption, strings.Join(kimsufiavailability.GetDatacentersKnownCodes(), ", ")))
	Cmd.PersistentFlags().IntVarP(&quantity, "quantity", "q", kimsufiorder.QuantityDefault, "item quantity")

	Cmd.PersistentFlags().StringToStringVarP(&itemUserConfigurations, "item-configuration", "i", nil, "item configuration, comma separated list, see --list-configurations for available values (e.g. region=europe)")
	Cmd.PersistentFlags().StringSliceVarP(&itemUserOptions, "item-option", "o", nil, fmt.Sprintf("item option, comma separated list, use any to include all options, see --list-options for available values (e.g. memory=ram-64g-noecc-2133-24ska01, memory=%[1]s, %[1]s)", anyOption))

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

	if len(datacenters) == 0 {
		return fmt.Errorf("--datacenter is required")
	} else if slices.Contains(datacenters, anyOption) {
		catalog, err := k.ListServers(cmd.Flag(flag.CountryFlagName).Value.String())
		if err != nil {
			return fmt.Errorf("failed to list servers: %w", err)
		}

		plan := catalog.GetPlan(planCode)
		if plan == nil {
			return fmt.Errorf("plan %s not found", planCode)
		}

		datacenterConfiguration := plan.GetConfiguration(kimsufiorder.ConfigurationLabelDatacenter)
		if datacenterConfiguration == nil {
			return fmt.Errorf("datacenter configuration not found")
		}

		datacenters = datacenterConfiguration.Values
	}

	// Prepare item configurations
	itemConfigurations := kimsufiorder.NewItemConfigurationsFromMap(itemUserConfigurations)
	r := kimsufiregion.GetRegionFromEndpoint(endpoint)
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
	for _, configuration := range configurations {
		resp, err := k.AddItemConfiguration(cart.CartID, item.ItemID, configuration)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		fmt.Printf("> cart item configured: %s=%s\n", resp.Label, resp.Value)
	}

	// Prepare item options
	var optionsCombinations []kimsufiorder.Options
	var mergedOptions kimsufiorder.Options
	allOptions := slices.Contains(itemUserOptions, anyOption)
	if allOptions {
		// Get all mandatory options
		mergedOptions = ecoOptions.GetMandatoryOptions(nil).ToOptions()
		optionsCombinations = kimsufiorder.NewOptionsCombinationsFromSlice(mergedOptions)
	} else {
		userOptions, err := kimsufiorder.NewOptionsFromSlice(itemUserOptions)
		if err != nil {
			return fmt.Errorf("error: %w", err)
		}
		anyOptions, userOptions := userOptions.SplitByPlanCode(anyOption)
		anyFamilies := anyOptions.Families()

		optionFilter := func(opts kimsufiorder.EcoItemOptions, o kimsufiorder.EcoItemOption) bool {
			defautPriceConfig := kimsufiorder.EcoItemPriceConfig{
				Duration:    kimsufiorder.PriceDuration,
				PricingMode: kimsufiorder.PricingMode,
			}

			// Inclue option if it is marked as any
			if slices.Contains(anyFamilies, o.Family) {
				return true
			}

			// Include option if its family is not already included
			current := opts.Get(o.Family)
			if current == nil {
				return true
			}

			newPrice := o.GetPriceByConfig(defautPriceConfig)
			if newPrice == nil {
				return false
			}

			currentPrice := current.GetPriceByConfig(defautPriceConfig)
			if currentPrice == nil {
				return false
			}

			// Include option if its price is lower than the current one
			return newPrice.PriceInUcents < currentPrice.PriceInUcents
		}

		mandatoryOptions := ecoOptions.GetMandatoryOptions(optionFilter)
		mergedOptions = userOptions.Merge(mandatoryOptions.ToOptions())
		optionsCombinations = kimsufiorder.NewOptionsCombinationsFromSlice(mergedOptions)
	}

	fmt.Printf("> item options: %d %v\n", len(mergedOptions), mergedOptions.PlanCodes())
	fmt.Printf("> datacenter(s): %d\n", len(datacenters))
	fmt.Printf("> combinations: %d\n", len(optionsCombinations)*len(datacenters))

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

	// Try all options combinations
	for _, options := range optionsCombinations {
		// Configure item options
		for _, option := range options {
			err = k.ConfigureEcoItemOption(cart.CartID, item.ItemID, option, priceConfig)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
			fmt.Printf("> cart option set: %s=%s\n", option.Family, option.PlanCode)
		}

		// Try all datacenters
		for _, datacenter := range datacenters {
			datacenterConfiguration := kimsufiorder.ItemConfigurationRequest{
				Label: kimsufiorder.ConfigurationLabelDatacenter,
				Value: datacenter,
			}

			resp, err := k.AddItemConfiguration(cart.CartID, item.ItemID, datacenterConfiguration)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
			fmt.Printf("> datacenter %s configured\n", resp.Value)

			// Checkout and complete the order
			checkoutResp, err := k.CheckoutCart(cart.CartID, autoPay)
			if err == nil {
				fmt.Printf("> order completed: %s\n", checkoutResp.URL)
				return nil
			}

			if kimsufi.IsNotAvailableError(err) {
				fmt.Printf("> datacenter %s not available\n", datacenter)
			} else {
				fmt.Printf("> error: %v\n", err)
			}

			err = k.RemoveItemConfiguration(cart.CartID, item.ItemID, resp.ID)
			if err != nil {
				return fmt.Errorf("error: %w", err)
			}
		}
	}

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
