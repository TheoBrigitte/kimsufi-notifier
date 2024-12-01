package category

var (
	// List of known plan categories, including an empty string for uncategorized plans
	Categories = []category{
		{"kimsufi", "Kimsufi"},
		{"soyoustart", "So you Start"},
		{"rise", "Rise"},
		{"", "uncategorized"},
	}
)

type category struct {
	Name        string
	DisplayName string
}

func Contains(name string) bool {
	for _, category := range Categories {
		if category.Name == name {
			return true
		}
	}

	return false
}

func Names() []string {
	var values []string
	for _, category := range Categories {
		values = append(values, category.Name)
	}

	return values
}
