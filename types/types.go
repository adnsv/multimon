package types

// Rect represents a rectangle with coordinates in screen units
type Rect struct {
	Left   int // X coordinate of the left edge
	Top    int // Y coordinate of the top edge
	Right  int // X coordinate of the right edge
	Bottom int // Y coordinate of the bottom edge
}

// Monitor represents a display monitor and its properties.
// All coordinates are in screen units (physical pixels on Windows/Linux, points on macOS).
// Scale factor is used to convert between screen units and logical units.
type Monitor struct {
	Bounds   Rect    // Monitor bounds in screen units
	WorkArea Rect    // Work area (excluding taskbar, etc.) in screen units
	Scale    float64 // Scale factor (1.0 = 100%, 2.0 = 200%, etc.)
}
