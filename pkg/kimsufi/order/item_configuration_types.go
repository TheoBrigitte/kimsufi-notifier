package order

// ItemConfiguration represents an available option for an item configuration.
type ItemConfiguration struct {
	AllowedValues []string `json:"allowedValues"`
	Fields        []string `json:"fields"`
	Label         string   `json:"label"`
	Required      bool     `json:"required"`
	Type          string   `json:"type"`
}

type ItemConfigurationRequests []ItemConfigurationRequest

// ItemConfigurationRequest represents the request
// to add a configuration to an item in the cart.
type ItemConfigurationRequest struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ItemConfigurationResponse represents the response
// of an item configuration added to the cart.
type ItemConfigurationResponse struct {
	ItemConfigurationRequest `json:",inline"`

	ID int `json:"id"`
}
