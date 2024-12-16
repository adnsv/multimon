package multimon

import (
	"errors"
	"fmt"
	"math"
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

// FitToMonitor fits a window to a specific monitor.
// Input window coordinates are in screen units.
// windowScale specifies what scale factor the window was designed for:
// - If 0.0: keep window as is, no rescaling needed
// - If > 0.0: rescale window from windowScale to monitor's scale
// Returns error if window or monitor has negative dimensions.
// Returns the fitted rect and the monitor's scale factor.
// If monitor is nil, returns windowScale if non-zero, otherwise 1.0.
func FitToMonitor(m *Monitor, mode FitMode, window Rect, windowScale float64) (Rect, float64, error) {
	// Validate input dimensions
	if err := validateRect(window); err != nil {
		return window, windowScale, fmt.Errorf("invalid window: %w", err)
	}
	if m == nil {
		// When no monitor is available, use windowScale if non-zero, otherwise 1.0
		outputScale := 1.0
		if windowScale > 0.0 {
			outputScale = windowScale
		}
		return window, outputScale, fmt.Errorf("invalid monitor: nil")
	}
	if err := validateMonitor(*m); err != nil {
		return window, windowScale, err
	}

	// Get target bounds based on mode
	bounds := m.Bounds
	if mode == FitModeWorkArea {
		bounds = m.WorkArea
	}

	// Scale window if needed
	var targetWindow Rect
	if windowScale == 0.0 {
		targetWindow = window
	} else {
		// Scale window dimensions relative to top-left corner
		scaledWidth := int(float64(window.Right-window.Left) * (m.Scale / windowScale))
		scaledHeight := int(float64(window.Bottom-window.Top) * (m.Scale / windowScale))
		targetWindow = Rect{
			Left:   window.Left,
			Top:    window.Top,
			Right:  window.Left + scaledWidth,
			Bottom: window.Top + scaledHeight,
		}
	}

	// Fit dimensions and positions within bounds
	newLeft, newWidth := fitRectDimension(targetWindow.Left, targetWindow.Right-targetWindow.Left, bounds.Left, bounds.Right)
	newTop, newHeight := fitRectDimension(targetWindow.Top, targetWindow.Bottom-targetWindow.Top, bounds.Top, bounds.Bottom)

	return Rect{
		Left:   newLeft,
		Top:    newTop,
		Right:  newLeft + newWidth,
		Bottom: newTop + newHeight,
	}, m.Scale, nil
}

// validateRect checks if a rectangle has valid dimensions
func validateRect(r Rect) error {
	width := r.Right - r.Left
	height := r.Bottom - r.Top
	if width <= 0 || height <= 0 {
		return fmt.Errorf("%w: width=%d, height=%d", ErrInvalidDimensions, width, height)
	}
	return nil
}

// validateMonitor checks if a monitor has valid dimensions
func validateMonitor(m Monitor) error {
	if err := validateRect(m.Bounds); err != nil {
		return fmt.Errorf("invalid bounds: %w", err)
	}
	if err := validateRect(m.WorkArea); err != nil {
		return fmt.Errorf("invalid work area: %w", err)
	}
	if m.Scale <= 0.0 {
		return fmt.Errorf("invalid scale: %v (must be positive non-zero)", m.Scale)
	}
	// Check that work area is contained within bounds
	if m.WorkArea.Left < m.Bounds.Left || m.WorkArea.Right > m.Bounds.Right ||
		m.WorkArea.Top < m.Bounds.Top || m.WorkArea.Bottom > m.Bounds.Bottom {
		return fmt.Errorf("work area outside bounds: work_area=%v bounds=%v", m.WorkArea, m.Bounds)
	}
	return nil
}

// FitToNearestMonitor finds the most appropriate monitor and fits the window to it.
// Input window coordinates are in screen units.
// windowScale specifies what scale factor the window was designed for:
// - If 0.0: keep window as is, no rescaling needed
// - If > 0.0: rescale window from windowScale to monitor's scale
// minWidth and minHeight specify the minimum dimensions the window should have (in logical units).
// Returns error if window has negative dimensions, if no valid monitors are available,
// or if no monitor can fit the minimum size requirements.
// If no monitors are available, returns windowScale if non-zero, otherwise 1.0.
func FitToNearestMonitor(monitors []Monitor, mode FitMode, window Rect, windowScale float64, minWidth, minHeight int) (Rect, float64, error) {
	// Validate window dimensions
	if err := validateRect(window); err != nil {
		return window, windowScale, fmt.Errorf("invalid window: %w", err)
	}

	if len(monitors) == 0 {
		// When no monitors are available, use windowScale if non-zero, otherwise 1.0
		outputScale := 1.0
		if windowScale > 0.0 {
			outputScale = windowScale
		}
		return window, outputScale, ErrNoMonitors
	}

	// Filter out invalid monitors
	var validMonitors []Monitor
	for _, m := range monitors {
		if validateMonitor(m) == nil {
			validMonitors = append(validMonitors, m)
		}
	}

	if len(validMonitors) == 0 {
		// When no valid monitors are available, use windowScale if non-zero, otherwise 1.0
		outputScale := 1.0
		if windowScale > 0.0 {
			outputScale = windowScale
		}
		return window, outputScale, fmt.Errorf("%w: no valid monitors found", ErrNoMonitors)
	}

	// Collect monitors that can fit minimum size with fallback to remaining monitors
	var suitableMonitors []Monitor
	if minWidth <= 0 || minHeight <= 0 {
		suitableMonitors = validMonitors
	} else {
		for _, m := range validMonitors {
			bounds := m.Bounds
			if mode == FitModeWorkArea {
				bounds = m.WorkArea
			}

			factor := 1.0
			if windowScale > 0.0 {
				factor = m.Scale / windowScale
			}

			// Check if monitor can fit minimum dimensions
			screenMinWidth := int(float64(minWidth) * factor)
			screenMinHeight := int(float64(minHeight) * factor)
			width := bounds.Right - bounds.Left
			height := bounds.Bottom - bounds.Top

			if width >= screenMinWidth && height >= screenMinHeight {
				suitableMonitors = append(suitableMonitors, m)
			}
		}
		if len(suitableMonitors) == 0 && mode == FitModeWorkArea {
			// If nothing fits within work area, as fallback try to fit to total bounds
			for _, m := range validMonitors {
				// Check if monitor can fit minimum dimensions
				screenMinWidth := int(float64(minWidth) * m.Scale)
				screenMinHeight := int(float64(minHeight) * m.Scale)
				width := m.Bounds.Right - m.Bounds.Left
				height := m.Bounds.Bottom - m.Bounds.Top
				if width >= screenMinWidth && height >= screenMinHeight {
					suitableMonitors = append(suitableMonitors, m)
				}
			}
		}
		if len(suitableMonitors) == 0 {
			// If nothing fits so far, try all the monitors, effectively
			// ignoring minWidth and minHeight.
			suitableMonitors = validMonitors
		}
	}

	// Find monitor with best overlap
	var bestMonitor *Monitor
	maxOverlap := 0
	for i := range suitableMonitors {
		m := &suitableMonitors[i]
		overlap := getOverlapArea(window, m.Bounds)
		if overlap > maxOverlap {
			maxOverlap = overlap
			bestMonitor = m
		}
	}

	// If nothing overlaps, use the nearest monitor
	if bestMonitor == nil {
		minDistance := math.MaxInt
		for i := range suitableMonitors {
			m := &suitableMonitors[i]
			dist := getEdgeDistance(window, m.Bounds)
			if dist < minDistance {
				minDistance = dist
				bestMonitor = m
			}
		}
	}

	if bestMonitor == nil {
		bestMonitor = &suitableMonitors[0]
	}

	return FitToMonitor(bestMonitor, mode, window, windowScale)
}
