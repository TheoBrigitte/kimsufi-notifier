package order

import (
	"slices"
)

// Get returns the option that matches the provided family.
func (i EcoItemOptions) Get(family string) *EcoItemOption {
	for index := range i {
		o := &i[index]
		if o.Family == family {
			return o
		}
	}

	return nil
}

// GetCheapestMandatoryOptions returns the cheapest mandatory options.
// It compares prices matching PriceDuration and PricingMode.
func (i EcoItemOptions) GetCheapestMandatoryOptions() EcoItemOptions {
	defautPriceConfig := EcoItemPriceConfig{
		Duration:    PriceDuration,
		PricingMode: PricingMode,
	}

	var options = EcoItemOptions{}

	for _, option := range i {
		if !option.Mandatory {
			continue
		}

		current := options.Get(option.Family)
		if current == nil {
			options = append(options, option)
			continue
		}

		newPrice := option.GetPriceByConfig(defautPriceConfig)
		if newPrice == nil {
			continue
		}

		currentPrice := current.GetPriceByConfig(defautPriceConfig)
		if currentPrice == nil {
			continue
		}

		if newPrice.PriceInUcents < currentPrice.PriceInUcents {
			options = append(options, option)
		}
	}

	return options
}

// GetPriceByConfig returns the first price that matches the provided EcoItemPriceConfig.
func (i EcoItemOption) GetPriceByConfig(priceConfig EcoItemPriceConfig) *EcoItemOptionPrice {
	for _, price := range i.Prices {
		if price.Duration == priceConfig.Duration && price.PricingMode == priceConfig.PricingMode {
			return &price
		}
	}

	return nil
}

// ToOptions converts EcoItemOptions to Options.
func (i EcoItemOptions) ToOptions() Options {
	var options []Option

	for _, o := range i {
		options = append(options, o.Option)
	}

	return options
}

// NewOptionsFromMap creates a new Options from a map.
// Using the map keys as family and the values as planCode.
func NewOptionsFromMap(optionsMap map[string]string) []Option {
	var options []Option

	for family, planCode := range optionsMap {
		o := Option{
			Family:   family,
			PlanCode: planCode,
		}
		options = append(options, o)
	}

	return options
}

// Merge merges the provided Options with the current Options.
// It does not overwrite existing options.
func (opts Options) Merge(other []Option) Options {
	families := opts.families()

	for _, o := range other {
		if !slices.Contains(families, o.Family) {
			opts = append(opts, o)
		}
	}

	return opts
}

func (opts Options) families() []string {
	var families []string

	for _, o := range opts {
		families = append(families, o.Family)
	}

	return families
}
