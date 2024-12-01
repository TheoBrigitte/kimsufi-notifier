package order

const (
	// PriceDuration P1M represents a duration of 1 month.
	// There are other durations like P0D for installation price only, P1Y for 1 year, etc...
	PriceDuration = "P1M"
	// PricingMode default means no commitment.
	// there are other price modes like degressivityX for degressiv prices over X months, upfrontX, etc...
	PricingMode = "default"
	PricingType = "rental"
	// PricingCapacityRenew renew means renewals are possible
	PricingCapacityRenew = "renew"

	QuantityDefault = 1

	ConfigurationLabelDatacenter = "dedicated_datacenter"
	ConfigurationLabelRegion     = "region"
)

// CartRequest represents the request to create a cart.
type CartRequest struct {
	Description   string `json:"description"`
	Expire        string `json:"expire"`
	OvhSubsidiary string `json:"ovhSubsidiary"`
}

// CartResponse represents the response of a cart creation.
type CartResponse struct {
	CartID string `json:"cartId"`
}

// CheckoutRequest represents the request to checkout.
type CheckoutRequest struct {
	AutoPayWithPreferredPaymentMethod bool `json:"autoPayWithPreferredPaymentMethod,omitempty"`
	WaiveRetractationPeriod           bool `json:"waiveRetractationPeriod,omitempty"`
}

// CheckoutResponse represents the response of a checkout.
type CheckoutResponse struct {
	OrderID   int                `json:"orderId,omitempty"`
	URL       string             `json:"url,omitempty"`
	Contracts []CheckoutContract `json:"contracts,omitempty"`
	Prices    CheckoutPrices     `json:"prices,omitempty"`
}

type CheckoutContract struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

type CheckoutPrices struct {
	OriginalWithoutTax Price `json:"originalWithoutTax,omitempty"`
	Reduction          Price `json:"reduction,omitempty"`
	Tax                Price `json:"tax,omitempty"`
	WithTax            Price `json:"withTax,omitempty"`
	WithoutTax         Price `json:"withoutTax,omitempty"`
}

type Price struct {
	CurrencyCode  string  `json:"currencyCode"`
	PriceInUcents int     `json:"priceInUcents"`
	Text          string  `json:"text"`
	Value         float64 `json:"value"`
}
