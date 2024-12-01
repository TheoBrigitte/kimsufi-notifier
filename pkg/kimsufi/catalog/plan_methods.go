package catalog

import (
	"math"
	"slices"
)

const (
	AddonMemory    = "memory"
	AddonStorage   = "storage"
	AddonBandwidth = "bandwidth"

	PriceModeDefault = "default"
)

var (
	priceDecimals = 8
	priceDivider  = math.Pow10(priceDecimals)
)

// GetAddons returns addons that match the provided names.
func (p *Plan) GetAddons(names ...string) []PlanAddonFamily {
	addons := make([]PlanAddonFamily, 0)

	for _, addon := range p.AddonFamilies {
		if slices.Contains(names, addon.Name) {
			addons = append(addons, addon)
		}
	}

	return addons
}

// GetAddon returns the first addon that matches the provided name.
func (p *Plan) GetAddon(name string) *PlanAddonFamily {
	for _, addon := range p.AddonFamilies {
		if addon.Name == name {
			return &addon
		}
	}

	return nil
}

// GetConfiguration returns the first configuration that matches the provided name.
func (p Plan) GetConfiguration(name string) *PlanConfiguration {
	for _, config := range p.Configurations {
		if config.Name == name {
			return &config
		}
	}

	return nil
}

// GetPrices returns the prices that have an interval greater than the provided minInterval.
// This allow to filter out prices like installation fees which correspond to 0 minInterval.
func (p *Plan) GetPrices(minInterval int) []PlanPricing {
	prices := make([]PlanPricing, 0)

	for _, price := range p.Pricings {
		if price.Interval <= minInterval {
			continue
		}

		prices = append(prices, price)
	}

	return prices
}

// GetPriceOrFirst returns the price that matches the provided PlanPricing
// or the first price if the provided PlanPricing is nil.
func (p Plan) GetPriceOrFirst(needle *PlanPricing) PlanPricing {
	if needle != nil {
		price := p.FindPrice(*needle)
		if price != nil {
			return *price
		}
	}

	return p.GetFirstPrice()
}

// GetFirstPrice does best effort to return the first price of the plan.
// If plan holds no prices, it returns an empty PlanPricing.
// Otherwise, it tries to find a price matching the following criteria:
// - IntervalUnit: "month"
// - Interval: 1
// - Commitement: 0
// - Phase: 1
// - Mode: PriceModeDefault
// - Type: "rental"
// - Strategy: "tiered"
// - Capacities: "renew"
// If no price matches the criteria, it returns the first price.
func (p Plan) GetFirstPrice() PlanPricing {
	if len(p.Pricings) == 0 {
		return PlanPricing{}
	}

	priceMatcher := PlanPricing{
		IntervalUnit: "month",
		Interval:     1,
		Commitement:  0,
		Phase:        1,
		Mode:         PriceModeDefault,
		Type:         "rental",
		Strategy:     "tiered",
		Capacities:   []string{"renew"},
	}

	price := p.FindPrice(priceMatcher)
	if price != nil {
		return *price
	}

	return p.Pricings[0]
}

// FindPrice returns the price that matches the provided PlanPricing.
func (p Plan) FindPrice(needle PlanPricing) *PlanPricing {
	for _, price := range p.Pricings {
		if price.Equals(needle) {
			return &price
		}
	}

	return nil
}

// Equals returns true if the provided PlanPricing is equal to the current PlanPricing.
// It compares the following fields:
// - Phase
// - Interval
// - IntervalUnit
// - Mode
// - Type
// - Commitement
// - Strategy
// - Capacities
func (price PlanPricing) Equals(other PlanPricing) bool {
	if price.Phase == other.Phase &&
		price.Interval == other.Interval &&
		price.IntervalUnit == other.IntervalUnit &&
		price.Mode == other.Mode &&
		price.Type == other.Type &&
		price.Commitement == other.Commitement &&
		price.Strategy == other.Strategy {
		for _, capacity := range other.Capacities {
			if !slices.Contains(price.Capacities, capacity) {
				return false
			}
		}

		return true
	}

	return false
}

// GetPrice returns the human readable price as a float64.
func (price PlanPricing) GetPrice() float64 {
	return float64(price.Price) / priceDivider
}
