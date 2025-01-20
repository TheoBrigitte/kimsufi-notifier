package catalog

// Catalog represents the OVH Eco catalog.
// Definition can be found at https://eu.api.ovh.com/console/?section=%2Forder&branch=v1#get-/order/catalog/public/eco
type Catalog struct {
	Addons    []Addon   `json:"addons"`
	CatalogID int       `json:"catalogId"`
	Locale    Locale    `json:"locale"`
	Plans     []Plan    `json:"plans"`
	Products  []Product `json:"products"`
}

type Addon struct {
	InvoiceName string        `json:"invoiceName"`
	PlanCode    string        `json:"planCode"`
	PricingType string        `json:"pricingType"`
	Pricings    []PlanPricing `json:"pricings"`
	Product     string        `json:"product"`
}

type Locale struct {
	CurrencyCode string `json:"currencyCode"`
	Subsidiary   string `json:"subsidiary"`
	TaxRate      int    `json:"taxRate"`
}
