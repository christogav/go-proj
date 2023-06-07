package utils

func StringWithDefault(value string, defaultValue string) string {
	if value != "" {
		return value
	}

	return defaultValue
}
