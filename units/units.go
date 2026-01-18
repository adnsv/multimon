// Package units provides dimension types for specifying window sizes
// in pixels, em-units (relative to system font), or percentages (relative to work area).
package units

import "fmt"

// DimensionUnit represents the type of measurement unit
type DimensionUnit int

const (
	Pixel   DimensionUnit = iota // Absolute pixels
	Em                           // Relative to system em-height
	Percent                      // Relative to work area
)

// Dimension represents a size value with its unit type
type Dimension struct {
	Value float64
	Unit  DimensionUnit
}

// Pixels returns a pixel dimension
func Pixels(v int) Dimension {
	return Dimension{Value: float64(v), Unit: Pixel}
}

// Ems returns an em-unit dimension
func Ems(v float64) Dimension {
	return Dimension{Value: v, Unit: Em}
}

// Pct returns a percentage dimension
func Pct(v float64) Dimension {
	return Dimension{Value: v, Unit: Percent}
}

// IsZero returns true if the dimension has no value
func (d Dimension) IsZero() bool {
	return d.Value == 0
}

// String returns the string representation (e.g., "60em", "80%", "1024")
func (d Dimension) String() string {
	switch d.Unit {
	case Em:
		if d.Value == float64(int(d.Value)) {
			return fmt.Sprintf("%dem", int(d.Value))
		}
		return fmt.Sprintf("%.2gem", d.Value)
	case Percent:
		if d.Value == float64(int(d.Value)) {
			return fmt.Sprintf("%d%%", int(d.Value))
		}
		return fmt.Sprintf("%.2g%%", d.Value)
	default:
		return fmt.Sprintf("%d", int(d.Value))
	}
}

// WorkArea provides the dimensions needed for percentage calculations
type WorkArea struct {
	Width  int
	Height int
}

// ResolveContext contains all context needed to resolve dimensions to pixels
type ResolveContext struct {
	EmHeight int      // System em-height in pixels
	WorkArea WorkArea // Monitor work area dimensions
}
