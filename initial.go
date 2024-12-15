package multimon

// InitialPlacement calculates the initial window placement on a default monitor.
// The default monitor is either:
// - The monitor containing point (0,0) in logical coordinates
// - If no monitor contains (0,0), the monitor with largest logical area
//
// All calculations are performed in logical (scaled) coordinates, accounting for
// the system's DPI settings and scaling factors.
//
// Parameters:
// - minWidth, minHeight: minimum required window size in logical pixels
// - desiredWidth, desiredHeight: preferred window size in logical pixels (will be clamped if exceeds work area)
// - margin: minimum distance from work area edges in logical pixels
//
// Returns a Rect with the calculated window position and size in logical coordinates.
func InitialPlacement(minWidth, minHeight, desiredWidth, desiredHeight, margin int) Rect {
	monitors := GetMonitors()
	if len(monitors) == 0 {
		return Rect{
			Left:   0,
			Top:    0,
			Right:  max(minWidth, desiredWidth),
			Bottom: max(minHeight, desiredHeight),
		}
	}

	// Find default monitor (containing 0,0 or largest)
	var defaultMonitor Monitor
	var maxArea int

	for _, m := range monitors {
		// Check if monitor contains (0,0)
		if m.LogicalBounds.Left <= 0 && m.LogicalBounds.Right > 0 &&
			m.LogicalBounds.Top <= 0 && m.LogicalBounds.Bottom > 0 {
			defaultMonitor = m
			break
		}

		// Track largest monitor as fallback
		area := getRectArea(m.LogicalBounds)
		if area > maxArea {
			maxArea = area
			defaultMonitor = m
		}
	}

	workArea := defaultMonitor.LogicalWorkArea
	workWidth := getRectWidth(workArea)
	workHeight := getRectHeight(workArea)

	// First try to satisfy desired size with margins
	availWidth := workWidth - 2*margin
	availHeight := workHeight - 2*margin

	// If minimum size exceeds available space with margins, allow using margin area
	width := min(max(minWidth, desiredWidth), availWidth)
	if width < minWidth {
		width = min(minWidth, workWidth)
	}

	height := min(max(minHeight, desiredHeight), availHeight)
	if height < minHeight {
		height = min(minHeight, workHeight)
	}

	// Center the window in work area
	centerX := (workArea.Left + workArea.Right) / 2
	centerY := (workArea.Top + workArea.Bottom) / 2

	return Rect{
		Left:   centerX - width/2,
		Top:    centerY - height/2,
		Right:  centerX + (width+1)/2,
		Bottom: centerY + (height+1)/2,
	}
}
