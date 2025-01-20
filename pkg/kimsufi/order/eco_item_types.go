package order

// EcoItemRequest represents the request to add an eco item to the cart
type EcoItemRequest struct {
	EcoItemPriceConfig `json:",inline"`

	PlanCode string `json:"planCode"`
	Quantity int    `json:"quantity"`
}

type EcoItemPriceConfig struct {
	Duration    string `json:"duration"`
	PricingMode string `json:"pricingMode"`
}

// EcoItemResponse represents the response of an eco item added to the cart
type EcoItemResponse struct {
	CartID string `json:"cartId"`
	ItemID int    `json:"itemId"`
}

type EcoItemInfos []EcoItemInfo

type EcoItemInfo struct {
	PlanCode    string             `json:"planCode"`
	Prices      []EcoItemInfoPrice `json:"prices"`
	ProductName string             `json:"productName"`
	ProductType string             `json:"productType"`
}

type EcoItemInfoPrice struct {
	Capacities      []string `json:"capacities"`
	Description     string   `json:"description"`
	Duration        string   `json:"duration"`
	Interval        int      `json:"interval"`
	MaximumQuantity int      `json:"maximumQuantity"`
	MaximumRepeat   int      `json:"maximumRepeat"`
	MinimumQuantity int      `json:"minimumQuantity"`
	MinimumRepeat   int      `json:"minimumRepeat"`
	Price           Price    `json:"price"`
	PriceInUcents   int      `json:"priceInUcents"`
	PricingMode     string   `json:"pricingMode"`
	PricingType     string   `json:"pricingType"`
}
