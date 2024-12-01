package order

import "slices"

// GetByPlanCode returns the EcoItemInfo with the given plan code.
func (e EcoItemInfos) GetByPlanCode(planCode string) *EcoItemInfo {
	for _, i := range e {
		if i.PlanCode == planCode {
			return &i
		}
	}

	return nil
}

// GetPriceConfigOrDefault does best effort to return a price config for
// the given plan code, and priceConfig. It tries to find a matching
// priceConfig, and if not it returns a default price config.
func (e EcoItemInfos) GetPriceConfigOrDefault(planCode string, priceConfig EcoItemPriceConfig) EcoItemPriceConfig {
	defaultPriceConfig := EcoItemPriceConfig{
		PricingMode: PricingMode,
		Duration:    PriceDuration,
	}

	planInfo := e.GetByPlanCode(planCode)
	if planInfo == nil {
		return defaultPriceConfig
	}

	exists := planInfo.priceConfigExists(priceConfig)
	if exists {
		return priceConfig
	}

	p := planInfo.firstPriceConfig()
	if p != nil {
		return *p
	}

	return defaultPriceConfig
}

func (e EcoItemInfo) priceConfigExists(price EcoItemPriceConfig) bool {
	for _, p := range e.Prices {
		if p.Duration == price.Duration && p.PricingMode == price.PricingMode {
			return true
		}
	}

	return false
}

func (e EcoItemInfo) firstPriceConfig() *EcoItemPriceConfig {
	for _, price := range e.Prices {
		if price.Interval == 1 &&
			price.PricingMode == PricingMode &&
			price.PricingType == PricingType &&
			slices.Contains(price.Capacities, PricingCapacityRenew) {
			return &EcoItemPriceConfig{
				Duration:    price.Duration,
				PricingMode: price.PricingMode,
			}
		}
	}

	return nil
}
