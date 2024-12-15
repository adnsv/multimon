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

// getEdgeDistance calculates the Manhattan distance from window center to monitor bounds.
// Returns 0 if window center is inside bounds.
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

// getRectArea returns the area of a rectangle
func getRectArea(r Rect) int {
	width := r.Right - r.Left
	height := r.Bottom - r.Top
	if width <= 0 || height <= 0 {
		return 0
	}
	return width * height
}

// getRectWidth returns the width of a rectangle
func getRectWidth(r Rect) int {
	return r.Right - r.Left
}

// getRectHeight returns the height of a rectangle
func getRectHeight(r Rect) int {
	return r.Bottom - r.Top
}

// getRectCenter returns the center point of a rectangle
func getRectCenter(r Rect) (x, y int) {
	return (r.Left + r.Right) / 2, (r.Top + r.Bottom) / 2
}
