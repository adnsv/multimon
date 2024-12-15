package multimon

// FindMonitorFromLogicalPoint finds the monitor that contains the given logical point
func FindMonitorFromLogicalPoint(monitors []Monitor, x, y int) *Monitor {
	for i := range monitors {
		if ContainsLogicalPoint(monitors[i], x, y) {
			return &monitors[i]
		}
	}
	return nil
}

// FindMonitorFromPhysicalPoint finds the monitor that contains the given physical point
func FindMonitorFromPhysicalPoint(monitors []Monitor, x, y int) *Monitor {
	for i := range monitors {
		if ContainsPhysicalPoint(monitors[i], x, y) {
			return &monitors[i]
		}
	}
	return nil
}

// FindMonitorFromLogicalRect finds a monitor with the largest overlap with the given rect
func FindMonitorFromLogicalRect(monitors []Monitor, rect Rect) *Monitor {
	var bestMonitor *Monitor
	maxArea := 0

	for i := range monitors {
		m := &monitors[i]
		area := getOverlapArea(rect, m.LogicalBounds)
		if area > maxArea {
			maxArea = area
			bestMonitor = m
		}
	}

	return bestMonitor
}
