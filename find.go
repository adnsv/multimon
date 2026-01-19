package multimon

import "math"

// DefaultMonitorMode specifies how to select a monitor when no exact match is found
type DefaultMonitorMode int

const (
	// DefaultMonitorNull returns nil if no monitor matches the criteria
	DefaultMonitorNull DefaultMonitorMode = iota
	// DefaultMonitorPrimary returns the monitor containing (0,0) in screen coordinates,
	// or the first available monitor if no monitor contains (0,0)
	DefaultMonitorPrimary
	// DefaultMonitorNearest returns the monitor with smallest edge distance to the target
	DefaultMonitorNearest
)

// FindPrimaryMonitor returns the monitor containing (0,0) in screen coordinates,
// or the first available monitor if no monitor contains (0,0).
// Returns nil if no monitors are available.
func FindPrimaryMonitor(monitors []Monitor) *Monitor {
	if len(monitors) == 0 {
		return nil
	}

	// First try to find monitor containing (0,0)
	for i := range monitors {
		m := &monitors[i]
		if m.Bounds.Left <= 0 && m.Bounds.Right > 0 &&
			m.Bounds.Top <= 0 && m.Bounds.Bottom > 0 {
			return m
		}
	}
	// If no monitor contains (0,0), return first available monitor
	return &monitors[0]
}

// FindMonitorFromScreenRect finds a monitor with the largest overlap with the given rect in screen coordinates.
// If no monitor has overlap:
// - DefaultMonitorNearest: returns the monitor with smallest edge distance
// - DefaultMonitorPrimary: returns the primary monitor
// - DefaultMonitorNull: returns nil
func FindMonitorFromScreenRect(monitors []Monitor, rect Rect, defaultTo DefaultMonitorMode) *Monitor {
	if len(monitors) == 0 {
		return nil
	}

	// First try to find monitor with largest overlap
	var bestMonitor *Monitor
	maxArea := 0

	for i := range monitors {
		m := &monitors[i]
		area := getOverlapArea(rect, m.Bounds)
		if area > maxArea {
			maxArea = area
			bestMonitor = m
		}
	}

	// If we found a monitor with overlap, return it
	if maxArea > 0 {
		return bestMonitor
	}

	// No overlap found, handle according to defaultTo mode
	switch defaultTo {
	case DefaultMonitorNull:
		return nil

	case DefaultMonitorPrimary:
		return FindPrimaryMonitor(monitors)

	case DefaultMonitorNearest:
		// Find monitor with smallest edge distance
		minDist := math.MaxInt
		var nearest *Monitor

		for i := range monitors {
			m := &monitors[i]
			dist := getEdgeDistance(rect, m.Bounds)
			if nearest == nil || dist < minDist {
				minDist = dist
				nearest = m
			}
		}
		return nearest

	default:
		return nil
	}
}

// FindMonitorFromScreenPoint finds the monitor that contains the given screen point.
// If no monitor contains the point:
// - DefaultMonitorNearest: returns the nearest monitor
// - DefaultMonitorPrimary: returns the primary monitor
// - DefaultMonitorNull: returns nil
func FindMonitorFromScreenPoint(monitors []Monitor, x, y int, defaultTo DefaultMonitorMode) *Monitor {
	// Create a 1x1 rect around the point
	rect := Rect{
		Left:   x,
		Top:    y,
		Right:  x + 1,
		Bottom: y + 1,
	}
	return FindMonitorFromScreenRect(monitors, rect, defaultTo)
}

// GetWorkAreaForRect returns the work area of the monitor containing the given rect.
// Returns an empty Rect if no monitors are available.
func GetWorkAreaForRect(monitors []Monitor, rect Rect) Rect {
	m := FindMonitorFromScreenRect(monitors, rect, DefaultMonitorNearest)
	if m == nil {
		return Rect{}
	}
	return m.WorkArea
}
