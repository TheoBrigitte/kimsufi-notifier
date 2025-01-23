package category

var (
	// Categories is a list of known plan categories, including an empty string for uncategorized plans
	Categories = []category{
		{"kimsufi", "Kimsufi", "sk"},
		{"soyoustart", "So you Start", "sys"},
		{"rise", "Rise", "rise"},
		{"", "uncategorized", ""},
	}
)

type category struct {
	Name        string
	DisplayName string
	ShortCode   string
}

func GetDisplayName(name string) string {
	for _, category := range Categories {
		if category.Name == name {
			return category.DisplayName
		}
	}

	return ""
}

// Contains checks if a category name is in the list of known categories.
func Contains(name string) bool {
	for _, category := range Categories {
		if category.Name == name {
			return true
		}
	}

	return false
}

// Names returns a list of known category names.
func Names() []string {
	var values []string
	for _, category := range Categories {
		values = append(values, category.Name)
	}

	return values
}
