package multimon

// getOverlapArea calculates the overlapping area between two rectangles.
// Returns 0 if rectangles don't overlap.
func getOverlapArea(r1, r2 Rect) int {
	// Calculate intersection
	left := max(r1.Left, r2.Left)
	right := min(r1.Right, r2.Right)
	top := max(r1.Top, r2.Top)
	bottom := min(r1.Bottom, r2.Bottom)

	// Check if there is an overlap
	if left < right && top < bottom {
		return (right - left) * (bottom - top)
	}
	return 0
}

// getEdgeDistance calculates the Manhattan distance from window center to nearest monitor edge.
// Returns 0 if window center is inside monitor bounds, otherwise returns the sum of
// horizontal and vertical distances to the nearest edges.
func getEdgeDistance(window, bounds Rect) int {
	windowCenterX := (window.Left + window.Right) / 2
	windowCenterY := (window.Top + window.Bottom) / 2

	// Calculate distance to nearest edge or 0 if inside bounds
	var dx, dy int
	if windowCenterX < bounds.Left {
		dx = bounds.Left - windowCenterX
	} else if windowCenterX > bounds.Right {
		dx = windowCenterX - bounds.Right
	}
	if windowCenterY < bounds.Top {
		dy = bounds.Top - windowCenterY
	} else if windowCenterY > bounds.Bottom {
		dy = windowCenterY - bounds.Bottom
	}

	return dx + dy
}

// fitRectDimension fits a dimension within bounds and returns the fitted size and position.
// If size exceeds available space, it is clamped to bounds.
// Position is adjusted to keep the dimension within bounds while maintaining size.
// Returns (pos, 0) if size is negative.
// Returns (boundsMin, 0) if min > max.
func fitRectDimension(pos, size, boundsMin, boundsMax int) (newPos, newSize int) {
	// Handle negative size
	if size < 0 {
		return pos, 0
	}

	// Handle invalid bounds
	if boundsMin > boundsMax {
		return boundsMin, 0
	}

	// First fit the size
	if size > boundsMax-boundsMin {
		size = boundsMax - boundsMin
	}

	// Then fit the position
	if pos+size > boundsMax {
		pos = boundsMax - size
	}
	if pos < boundsMin {
		pos = boundsMin
	}

	return pos, size
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
