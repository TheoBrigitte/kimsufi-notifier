package kimsufi

import "math"

var (
	PriceDecimals = 8
	PriceDivider  = math.Pow10(PriceDecimals)

	PlanCategories = []string{"kimsufi", "soyoustart", "rise"}

	StatusAvailable   = "available"
	StatusUnavailable = "unavailable"
)

func (p Plan) FirstPrice() Pricing {
	if len(p.Pricings) == 0 {
		return Pricing{}
	}

	for _, price := range p.Pricings {
		if price.Phase == 1 && price.Mode == "default" {
			return price
		}
	}

	return p.Pricings[0]
}
