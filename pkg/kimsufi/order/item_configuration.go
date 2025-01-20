package order

// NewItemConfigurationsFromMap creates a new ItemConfigurations.
// It uses the map keys as labels and the map values as values.
func NewItemConfigurationsFromMap(values map[string]string) ItemConfigurationRequests {
	var itemConfigurations ItemConfigurationRequests

	for label, value := range values {
		i := ItemConfigurationRequest{
			Label: label,
			Value: value,
		}
		itemConfigurations = append(itemConfigurations, i)
	}

	return itemConfigurations
}

// Merge merges other ItemConfigurations into the current one.
// It does not overwrite existing configurations.
func (c ItemConfigurationRequests) Merge(other ItemConfigurationRequests) ItemConfigurationRequests {
	for _, o := range other {
		i := c.GetByLabel(o.Label)
		if i == nil {
			c = append(c, o)
		}
	}

	return c
}

// GetByLabel returns the first configuration with the given label.
func (c ItemConfigurationRequests) GetByLabel(label string) *ItemConfigurationRequest {
	for index := range c {
		i := &c[index]
		if i.Label == label {
			return i
		}
	}

	return nil
}

// Add adds a configuration if it does not exist.
func (c *ItemConfigurationRequests) Add(label, value string) {
	if c.GetByLabel(label) == nil {
		i := ItemConfigurationRequest{
			Label: label,
			Value: value,
		}
		*c = append(*c, i)
	}
}
