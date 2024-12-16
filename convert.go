package multimon

// Point represents a point in 2D space
type Point struct {
	X, Y int
}

// LogicalToScreenRect converts logical coordinates to screen units for a given monitor
func LogicalToScreenRect(m Monitor, logical Rect) Rect {
	// Convert from logical units to screen units by multiplying by scale factor
	return Rect{
		Left:   int(float64(logical.Left) * m.Scale),
		Top:    int(float64(logical.Top) * m.Scale),
		Right:  int(float64(logical.Right) * m.Scale),
		Bottom: int(float64(logical.Bottom) * m.Scale),
	}
}

// ScreenToLogicalRect converts screen coordinates to logical units for a given monitor
func ScreenToLogicalRect(m Monitor, screen Rect) Rect {
	// Convert from screen units to logical units by dividing by scale factor
	return Rect{
		Left:   int(float64(screen.Left) / m.Scale),
		Top:    int(float64(screen.Top) / m.Scale),
		Right:  int(float64(screen.Right) / m.Scale),
		Bottom: int(float64(screen.Bottom) / m.Scale),
	}
}

// LogicalToScreenPoint converts a logical point to screen coordinates for a given monitor
func LogicalToScreenPoint(m Monitor, x, y int) Point {
	return Point{
		X: int(float64(x) * m.Scale),
		Y: int(float64(y) * m.Scale),
	}
}

// ScreenToLogicalPoint converts a screen point to logical coordinates for a given monitor
func ScreenToLogicalPoint(m Monitor, x, y int) Point {
	return Point{
		X: int(float64(x) / m.Scale),
		Y: int(float64(y) / m.Scale),
	}
}
