package kimsufi

import (
	"fmt"
	"net/url"
	"time"

	kimsufiorder "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/order"
)

// CreateCart creates a new cart which will expire at the given time.
func (s *Service) CreateCart(ovhSubsidiary string, expire time.Time) (*kimsufiorder.CartResponse, error) {
	u := "/order/cart"

	req := kimsufiorder.CartRequest{
		Description:   "kimsufi-notifier",
		Expire:        expire.Format(time.RFC3339),
		OvhSubsidiary: ovhSubsidiary,
	}
	s.logger.Debugf("CreateCart request: %+#v", req)

	var resp kimsufiorder.CartResponse
	err := s.client.PostUnAuth(u, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// AddEcoItem adds an OVH eco item to the cart with the given planCode, quantity and duration, mode from priceConfig.
func (s *Service) AddEcoItem(cartID, planCode string, quantity int, priceConfig kimsufiorder.EcoItemPriceConfig) (*kimsufiorder.EcoItemResponse, error) {
	u := fmt.Sprintf("/order/cart/%s/eco", cartID)

	req := kimsufiorder.EcoItemRequest{
		EcoItemPriceConfig: priceConfig,

		PlanCode: planCode,
		Quantity: quantity,
	}
	s.logger.Debugf("AddEcoItem request: %+#v", req)

	var resp kimsufiorder.EcoItemResponse
	err := s.client.PostUnAuth(u, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetEcoInfo returns information about an eco item in the cart.
func (s *Service) GetEcoInfo(cartID, planCode string) (kimsufiorder.EcoItemInfos, error) {
	u, err := url.Parse(fmt.Sprintf("/order/cart/%s/eco", cartID))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("planCode", planCode)
	u.RawQuery = q.Encode()

	var resp kimsufiorder.EcoItemInfos
	err = s.client.GetUnAuth(u.String(), &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetEcoOptions returns the options for an eco item in the cart.
func (s *Service) GetEcoOptions(cartID string, planCode string) ([]kimsufiorder.EcoItemOption, error) {
	u, err := url.Parse(fmt.Sprintf("/order/cart/%s/eco/options", cartID))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("planCode", planCode)
	u.RawQuery = q.Encode()

	var options []kimsufiorder.EcoItemOption
	err = s.client.GetUnAuth(u.String(), &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}

// ConfigureEcoItemOptions configures the item options in the cart.
// It finds the cheapest mandatory options and merges them into the user options.
func (s *Service) ConfigureEcoItemOptions(cartID string, itemID int, options kimsufiorder.EcoItemOptions, userOptions kimsufiorder.Options, priceConfig kimsufiorder.EcoItemPriceConfig) (kimsufiorder.Options, error) {
	u := fmt.Sprintf("/order/cart/%s/eco/options", cartID)

	mandatoryOptions := options.GetCheapestMandatoryOptions()

	mergedOptions := userOptions.Merge(mandatoryOptions.ToOptions())

	for _, option := range mergedOptions {
		req := kimsufiorder.EcoItemOptionRequest{
			EcoItemRequest: kimsufiorder.EcoItemRequest{
				EcoItemPriceConfig: priceConfig,
				PlanCode:           option.PlanCode,
				Quantity:           kimsufiorder.QuantityDefault,
			},
			ItemID: itemID,
		}

		s.logger.Debugf("ConfigureItemOptions request: %+#v", req)
		err := s.client.PostUnAuth(u, req, nil)
		if err != nil {
			return nil, err
		}
	}

	return mergedOptions, nil
}

// GetItemRequiredConfiguration returns the required configuration options for an item in the cart.
func (s *Service) GetItemRequiredConfiguration(cartID string, itemID int) ([]kimsufiorder.ItemConfiguration, error) {
	u := fmt.Sprintf("/order/cart/%s/item/%d/requiredConfiguration", cartID, itemID)

	var resp []kimsufiorder.ItemConfiguration
	err := s.client.GetUnAuth(u, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ConfigureItem configures an item in the cart with the given configurations.
func (s *Service) AddItemConfiguration(cartID string, itemID int, configuration kimsufiorder.ItemConfigurationRequest) (*kimsufiorder.ItemConfigurationResponse, error) {
	u := fmt.Sprintf("/order/cart/%s/item/%d/configuration", cartID, itemID)

	var resp kimsufiorder.ItemConfigurationResponse
	s.logger.Debugf("ConfigureItem request: %+#v", configuration)
	err := s.client.PostUnAuth(u, configuration, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ConfigureItem configures an item in the cart with the given configurations.
func (s *Service) RemoveItemConfiguration(cartID string, itemID, configurationID int) error {
	u := fmt.Sprintf("/order/cart/%s/item/%d/configuration/%d", cartID, itemID, configurationID)

	return s.client.DeleteUnAuth(u, nil)
}

// AssignCart assigns the cart to the user's account.
func (s *Service) AssignCart(cartID string) error {
	u := fmt.Sprintf("/order/cart/%s/assign", cartID)

	err := s.client.Post(u, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// CheckoutCart checks out the cart to place the order.
// If autoPay is true, the order will be paid automatically using the preferred payment method.
func (s *Service) CheckoutCart(cartID string, autoPay bool) (*kimsufiorder.CheckoutResponse, error) {
	u := fmt.Sprintf("/order/cart/%s/checkout", cartID)

	req := kimsufiorder.CheckoutRequest{
		AutoPayWithPreferredPaymentMethod: autoPay,
		WaiveRetractationPeriod:           false,
	}
	s.logger.Debugf("CheckoutCart request: %+#v", req)

	var resp kimsufiorder.CheckoutResponse
	err := s.client.Post(u, req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GenerateItemAutoConfigurations generates configurations for an item..
// It returns every required configuration with only one allowed value.
func (s *Service) GenerateItemAutoConfigurations(itemConfigurationOptions []kimsufiorder.ItemConfiguration) kimsufiorder.ItemConfigurationRequests {
	var autoConfigurations []kimsufiorder.ItemConfigurationRequest

	for _, option := range itemConfigurationOptions {
		if !option.Required {
			continue
		}

		if len(option.AllowedValues) == 1 {
			i := kimsufiorder.ItemConfigurationRequest{
				Label: option.Label,
				Value: option.AllowedValues[0],
			}
			autoConfigurations = append(autoConfigurations, i)
		}
	}

	return autoConfigurations
}
