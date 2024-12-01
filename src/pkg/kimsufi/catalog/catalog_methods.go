package catalog

// GetPlan returns the plan with the given plan code.
func (c Catalog) GetPlan(planCode string) *Plan {
	for _, plan := range c.Plans {
		if plan.PlanCode == planCode {
			return &plan
		}
	}

	return nil
}

// GetProduct returns the product with the given product name.
func (c Catalog) GetProduct(productName string) *Product {
	for _, product := range c.Products {
		if product.Name == productName {
			return &product
		}
	}

	return nil
}
