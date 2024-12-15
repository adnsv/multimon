package types

// Rect represents a rectangle with coordinates in screen space
type Rect struct {
	Left   int // X coordinate of the left edge
	Top    int // Y coordinate of the top edge
	Right  int // X coordinate of the right edge
	Bottom int // Y coordinate of the bottom edge
}

// Monitor represents a display monitor and its properties
type Monitor struct {
	LogicalBounds    Rect // Logical (scaled) bounds of the monitor in screen coordinates
	LogicalWorkArea  Rect // Logical (scaled) work area (excluding taskbar, etc.)
	PhysicalBounds   Rect // Physical (unscaled) bounds of the monitor in screen coordinates
	PhysicalWorkArea Rect // Physical (unscaled) work area (excluding taskbar, etc.)
}
