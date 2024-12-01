package order

import (
	"fmt"

	kimsufiorder "github.com/TheoBrigitte/kimsufi-notifier/pkg/kimsufi/order"
)

const (
	maxInputRetries = 3
)

// generateItemManualConfiguration asks the user to select a value for each required configuration not already set in mergedConfigs
func generateItemManualConfiguration(mergedConfigs kimsufiorder.ItemConfigurationRequests, requiredConfigs []kimsufiorder.ItemConfiguration) (kimsufiorder.ItemConfigurationRequests, error) {
	var manualConfigs kimsufiorder.ItemConfigurationRequests
	for _, option := range requiredConfigs {
		if !option.Required {
			continue
		}
		i := mergedConfigs.GetByLabel(option.Label)
		if i != nil {
			continue
		}

		fmt.Printf("> cart item manual configuration, select a value for %s\n", option.Label)
		for index, value := range option.AllowedValues {
			fmt.Printf("  %d. %s\n", index, value)
		}
		var choice int
		var err error
		for i := 0; i < maxInputRetries; i++ {
			fmt.Printf("> Choice: ")
			_, err = fmt.Scan(&choice)
			if err != nil {
				fmt.Printf("  invalid choice: %v\n", err)
			} else if choice < 0 || choice >= len(option.AllowedValues) {
				fmt.Printf("  invalid choice: %d\n", choice)
			} else {
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("too many invalid choices: %w", err)
		}
		if choice < 0 || choice >= len(option.AllowedValues) {
			return nil, fmt.Errorf("invalid choice: %d", choice)
		}

		m := kimsufiorder.ItemConfigurationRequest{
			Label: option.Label,
			Value: option.AllowedValues[choice],
		}
		manualConfigs = append(manualConfigs, m)
	}

	return manualConfigs, nil
}
