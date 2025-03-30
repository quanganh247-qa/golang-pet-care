package test

// Helper functions to get category details
func getCategoryName(categoryID string) string {
	// Ideally, fetch this from database
	categoryNames := map[string]string{
		"blood": "Blood Tests",
		// Add other categories
	}

	if name, exists := categoryNames[categoryID]; exists {
		return name
	}
	return categoryID // Fallback to ID if name not found
}

func getCategoryIcon(categoryID string) string {
	// Return icon identifier that frontend can interpret
	categoryIcons := map[string]string{
		"blood": "beaker",
		// Add other categories
	}

	if icon, exists := categoryIcons[categoryID]; exists {
		return icon
	}
	return "default-icon"
}

func getCategoryDescription(categoryID string) string {
	// Ideally, fetch this from database
	categoryDescriptions := map[string]string{
		"blood": "Check blood count, liver and kidney enzymes",
		// Add other categories
	}

	if desc, exists := categoryDescriptions[categoryID]; exists {
		return desc
	}
	return ""
}
