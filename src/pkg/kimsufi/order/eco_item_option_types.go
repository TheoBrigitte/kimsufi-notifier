package order

type EcoItemOptions []EcoItemOption

// EcoItemOptions represents the options for an eco item.
type EcoItemOption struct {
	Option `json:",inline"`

	Mandatory   bool                 `json:"mandatory"`
	Prices      []EcoItemOptionPrice `json:"prices"`
	ProductName string               `json:"productName"`
}

type Options []Option

type Option struct {
	Family   string `json:"family"`
	PlanCode string `json:"planCode"`
}

type EcoItemOptionPrice struct {
	Duration      string `json:"duration"`
	PricingMode   string `json:"pricingMode"`
	PriceInUcents int    `json:"priceInUcents"`
	Price         Price  `json:"price"`
}

// EcoItemOptionRequest represents the request to add an option
// to an eco item in the cart.
type EcoItemOptionRequest struct {
	EcoItemRequest `json:",inline"`

	ItemID int `json:"itemId"`
}
