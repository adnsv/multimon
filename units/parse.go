package units

import (
	"strconv"
	"strings"
)

// ParseDimension parses a dimension string like "1024", "1024px", "60em", or "80%".
// Returns a zero dimension if parsing fails.
func ParseDimension(s string) Dimension {
	s = strings.TrimSpace(s)
	if s == "" {
		return Dimension{}
	}

	// Check for em suffix
	if strings.HasSuffix(s, "em") {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(s, "em"), 64); err == nil {
			return Dimension{Value: v, Unit: Em}
		}
		return Dimension{}
	}

	// Check for percent suffix
	if strings.HasSuffix(s, "%") {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64); err == nil {
			return Dimension{Value: v, Unit: Percent}
		}
		return Dimension{}
	}

	// Check for px suffix (explicit pixels)
	if strings.HasSuffix(s, "px") {
		if v, err := strconv.ParseFloat(strings.TrimSuffix(s, "px"), 64); err == nil {
			return Dimension{Value: v, Unit: Pixel}
		}
		return Dimension{}
	}

	// Try as pixel value (plain number)
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return Dimension{Value: v, Unit: Pixel}
	}

	return Dimension{}
}

// ParseDimensionWithDefault parses a dimension string, returning defaultVal if parsing fails.
func ParseDimensionWithDefault(s string, defaultVal Dimension) Dimension {
	d := ParseDimension(s)
	if d.IsZero() && s != "" && s != "0" {
		return defaultVal
	}
	if d.IsZero() {
		return defaultVal
	}
	return d
}
