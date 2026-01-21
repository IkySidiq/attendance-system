package utils

import "strings"

func ParseBooleanField(value interface{}) bool {
	// Handle nil
	if value == nil {
		return false
	}

	switch v := value.(type) {

	// If already boolean
	case bool:
		return v

	// If number (Excel often uses 1/0)
	case int:
		return v == 1
	case int8:
		return v == 1
	case int16:
		return v == 1
	case int32:
		return v == 1
	case int64:
		return v == 1
	case float32:
		return v == 1
	case float64:
		return v == 1

	// If string
	case string:
		str := strings.ToLower(strings.TrimSpace(v))
		if str == "" {
			return false
		}
		return str == "yes" || str == "true" || str == "1" || str == "y"

	default:
		return false
	}
}
