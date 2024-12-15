package multimon

// LogicalToPhysical converts logical coordinates to physical coordinates for a given monitor
func LogicalToPhysical(m Monitor, logical Rect) Rect {
	// Calculate the scaling factors
	scaleX := float64(m.PhysicalBounds.Right-m.PhysicalBounds.Left) / float64(m.LogicalBounds.Right-m.LogicalBounds.Left)
	scaleY := float64(m.PhysicalBounds.Bottom-m.PhysicalBounds.Top) / float64(m.LogicalBounds.Bottom-m.LogicalBounds.Top)

	// Apply scaling and offset
	return Rect{
		Left:   int(float64(logical.Left-m.LogicalBounds.Left)*scaleX) + m.PhysicalBounds.Left,
		Top:    int(float64(logical.Top-m.LogicalBounds.Top)*scaleY) + m.PhysicalBounds.Top,
		Right:  int(float64(logical.Right-m.LogicalBounds.Left)*scaleX) + m.PhysicalBounds.Left,
		Bottom: int(float64(logical.Bottom-m.LogicalBounds.Top)*scaleY) + m.PhysicalBounds.Top,
	}
}

// PhysicalToLogical converts physical coordinates to logical coordinates for a given monitor
func PhysicalToLogical(m Monitor, physical Rect) Rect {
	// Calculate the scaling factors
	scaleX := float64(m.LogicalBounds.Right-m.LogicalBounds.Left) / float64(m.PhysicalBounds.Right-m.PhysicalBounds.Left)
	scaleY := float64(m.LogicalBounds.Bottom-m.LogicalBounds.Top) / float64(m.PhysicalBounds.Bottom-m.PhysicalBounds.Top)

	// Apply scaling and offset
	return Rect{
		Left:   int(float64(physical.Left-m.PhysicalBounds.Left)*scaleX) + m.LogicalBounds.Left,
		Top:    int(float64(physical.Top-m.PhysicalBounds.Top)*scaleY) + m.LogicalBounds.Top,
		Right:  int(float64(physical.Right-m.PhysicalBounds.Left)*scaleX) + m.LogicalBounds.Left,
		Bottom: int(float64(physical.Bottom-m.PhysicalBounds.Top)*scaleY) + m.LogicalBounds.Top,
	}
}

// ContainsLogicalPoint checks if a logical point is within the monitor's logical bounds
func ContainsLogicalPoint(m Monitor, x, y int) bool {
	return x >= m.LogicalBounds.Left &&
		x < m.LogicalBounds.Right &&
		y >= m.LogicalBounds.Top &&
		y < m.LogicalBounds.Bottom
}

// ContainsPhysicalPoint checks if a physical point is within the monitor's physical bounds
func ContainsPhysicalPoint(m Monitor, x, y int) bool {
	return x >= m.PhysicalBounds.Left &&
		x < m.PhysicalBounds.Right &&
		y >= m.PhysicalBounds.Top &&
		y < m.PhysicalBounds.Bottom
}
