package multimon

import (
	"errors"
	"fmt"
)

// ErrNoMonitors is returned when no monitors are available for fitting
var ErrNoMonitors = errors.New("no monitors available")

// ErrInvalidDimensions is returned when a window or monitor has negative dimensions
var ErrInvalidDimensions = errors.New("invalid dimensions: negative width or height")

// FitMode specifies how to fit a window to a monitor
type FitMode int

const (
	// FitModeBounds fits to monitor's total bounds
	FitModeBounds FitMode = iota
	// FitModeWorkArea fits to monitor's work area (excluding taskbar, dock, etc.)
	FitModeWorkArea
)

// validateRect checks if a rectangle has valid dimensions
func validateRect(r Rect) error {
	if r.Right-r.Left <= 0 || r.Bottom-r.Top <= 0 {
		return fmt.Errorf("%w: width=%d, height=%d", ErrInvalidDimensions, r.Right-r.Left, r.Bottom-r.Top)
	}
	return nil
}

// validateMonitor checks if a monitor has valid dimensions
func validateMonitor(m Monitor) error {
	if err := validateRect(m.LogicalBounds); err != nil {
		return fmt.Errorf("invalid logical bounds: %w", err)
	}
	if err := validateRect(m.LogicalWorkArea); err != nil {
		return fmt.Errorf("invalid logical work area: %w", err)
	}
	return nil
}

// FitToMonitor fits a window to a specific monitor.
// All coordinates are in logical space.
// Returns error if window or monitor has negative dimensions.
func FitToMonitor(m Monitor, mode FitMode, window Rect) (Rect, error) {
	// Validate input dimensions
	if err := validateRect(window); err != nil {
		return window, fmt.Errorf("invalid window: %w", err)
	}
	if err := validateMonitor(m); err != nil {
		return window, err
	}

	// Get target bounds based on mode
	bounds := m.LogicalBounds
	if mode == FitModeWorkArea {
		bounds = m.LogicalWorkArea
	}

	// Calculate window dimensions
	width := window.Right - window.Left
	height := window.Bottom - window.Top

	// If window fits within bounds, just reposition if needed
	if width <= (bounds.Right-bounds.Left) && height <= (bounds.Bottom-bounds.Top) {
		newLeft := window.Left
		newTop := window.Top

		// If window is completely outside, position it relative to the monitor's edge
		if window.Left > bounds.Right || window.Right < bounds.Left ||
			window.Top > bounds.Bottom || window.Bottom < bounds.Top {
			newLeft = bounds.Left
			newTop = bounds.Top
		} else {
			// Normal positioning logic for partially overlapping windows
			if window.Right > bounds.Right {
				newLeft = bounds.Right - width
			}
			if newLeft < bounds.Left {
				newLeft = bounds.Left
			}
			if window.Bottom > bounds.Bottom {
				newTop = bounds.Bottom - height
			}
			if newTop < bounds.Top {
				newTop = bounds.Top
			}
		}

		return Rect{
			Left:   newLeft,
			Top:    newTop,
			Right:  newLeft + width,
			Bottom: newTop + height,
		}, nil
	}

	// Window needs to be resized
	newWidth := min(width, bounds.Right-bounds.Left)
	newHeight := min(height, bounds.Bottom-bounds.Top)
	newLeft := max(bounds.Left, min(window.Left, bounds.Right-newWidth))
	newTop := max(bounds.Top, min(window.Top, bounds.Bottom-newHeight))

	return Rect{
		Left:   newLeft,
		Top:    newTop,
		Right:  newLeft + newWidth,
		Bottom: newTop + newHeight,
	}, nil
}

// FitToNearestMonitor finds the most appropriate monitor and fits the window to it.
// If minWidth and minHeight are > 0, it will try to find a monitor that can fit these minimum dimensions.
// Returns error if window has negative dimensions or if no valid monitors are available.
func FitToNearestMonitor(monitors []Monitor, mode FitMode, window Rect, minWidth, minHeight int) (Rect, error) {
	// Validate window dimensions
	if err := validateRect(window); err != nil {
		return window, fmt.Errorf("invalid window: %w", err)
	}

	if len(monitors) == 0 {
		return window, ErrNoMonitors
	}

	// Filter out invalid monitors
	var validMonitors []Monitor
	for _, m := range monitors {
		if validateMonitor(m) == nil {
			validMonitors = append(validMonitors, m)
		}
	}

	if len(validMonitors) == 0 {
		return window, fmt.Errorf("%w: no valid monitors found", ErrNoMonitors)
	}

	type monitorScore struct {
		monitor Monitor
		overlap int
		dist    int
		canFit  bool
	}

	var candidates []monitorScore
	maxOverlap := 0

	// First pass: calculate overlap areas and check minimum size requirements
	for _, m := range validMonitors {
		bounds := m.LogicalBounds
		if mode == FitModeWorkArea {
			bounds = m.LogicalWorkArea
		}

		// Check if monitor can fit minimum dimensions
		canFit := true
		if minWidth > 0 && minHeight > 0 {
			width := bounds.Right - bounds.Left
			height := bounds.Bottom - bounds.Top
			if width < minWidth || height < minHeight {
				canFit = false
			}
		}

		// Calculate overlap area
		overlap := getOverlapArea(window, bounds)
		if overlap > maxOverlap {
			maxOverlap = overlap
		}

		// Calculate edge distance (used when no overlap)
		dist := getEdgeDistance(window, bounds)

		candidates = append(candidates, monitorScore{m, overlap, dist, canFit})
	}

	// Sort candidates by overlap area (primary) and edge distance (secondary)
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			shouldSwap := false

			// If both have overlap, prefer larger overlap
			if candidates[i].overlap > 0 && candidates[j].overlap > 0 {
				shouldSwap = candidates[j].overlap > candidates[i].overlap
			} else if candidates[i].overlap == 0 && candidates[j].overlap == 0 {
				// If neither has overlap, prefer closer distance
				shouldSwap = candidates[j].dist < candidates[i].dist
			} else {
				// Prefer any overlap over no overlap
				shouldSwap = candidates[j].overlap > candidates[i].overlap
			}

			if shouldSwap {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	// If we have minimum size requirements
	if minWidth > 0 && minHeight > 0 {
		// First try the monitor with the most overlap if it can fit minimum size
		bestCandidate := candidates[0]
		if bestCandidate.canFit {
			fitted, err := FitToMonitor(bestCandidate.monitor, mode, window)
			if err == nil {
				return fitted, nil
			}
		}

		// If best candidate can't fit, look for another monitor that can
		for _, c := range candidates[1:] { // Start from second candidate
			if c.canFit {
				if mode == FitModeWorkArea {
					return c.monitor.LogicalWorkArea, nil
				}
				return c.monitor.LogicalBounds, nil
			}
		}

		// If no monitor can fit, use the largest one
		largest := findLargestMonitor(validMonitors, mode)
		if largest != nil {
			if mode == FitModeWorkArea {
				return largest.LogicalWorkArea, nil
			}
			return largest.LogicalBounds, nil
		}
		return window, fmt.Errorf("%w: no monitor can fit minimum size", ErrNoMonitors)
	}

	// No minimum size requirements, use the best candidate
	return FitToMonitor(candidates[0].monitor, mode, window)
}

// findLargestMonitor returns the monitor with the largest area
func findLargestMonitor(monitors []Monitor, mode FitMode) *Monitor {
	var maxArea int
	var largest *Monitor

	for i := range monitors {
		m := &monitors[i]
		if validateMonitor(*m) != nil {
			continue
		}

		bounds := m.LogicalBounds
		if mode == FitModeWorkArea {
			bounds = m.LogicalWorkArea
		}

		area := (bounds.Right - bounds.Left) * (bounds.Bottom - bounds.Top)
		if area > maxArea {
			maxArea = area
			largest = m
		}
	}
	return largest
}
