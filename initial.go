package multimon

// CalcPlacementSize calculates the window size in screen units, attempting to satisfy
// the desired size while fitting within monitor bounds:
// 1. Converts desired size to screen units using monitor's scale
// 2. Attempts to fit within work area minus margins
// 3. If needed, allows using margin area to satisfy minimum size
//
// Parameters:
// - desiredWidth, desiredHeight: preferred window size in logical units
// - minWidth, minHeight: minimum required window size in logical units
// - margin: minimum distance from work area edges in logical units
//
// Returns width and height in screen units.
func CalcPlacementSize(m *Monitor, desiredWidth, desiredHeight, minWidth, minHeight, margin int) (width, height int) {
	if m == nil {
		// Convert logical units to screen units assuming 1:1 scale
		width = max(minWidth, desiredWidth)
		height = max(minHeight, desiredHeight)
		return
	}

	// Convert input parameters to screen units
	screenMinWidth := int(float64(minWidth) * m.Scale)
	screenMinHeight := int(float64(minHeight) * m.Scale)
	screenDesiredWidth := int(float64(desiredWidth) * m.Scale)
	screenDesiredHeight := int(float64(desiredHeight) * m.Scale)
	screenMargin := int(float64(margin) * m.Scale)

	// First try to satisfy desired size with margins
	availWidth := m.WorkArea.Right - m.WorkArea.Left - 2*screenMargin
	availHeight := m.WorkArea.Bottom - m.WorkArea.Top - 2*screenMargin

	// If minimum size exceeds available space with margins, allow using margin area
	width = min(max(screenMinWidth, screenDesiredWidth), availWidth)
	if width < screenMinWidth {
		width = min(screenMinWidth, m.WorkArea.Right-m.WorkArea.Left)
	}

	height = min(max(screenMinHeight, screenDesiredHeight), availHeight)
	if height < screenMinHeight {
		height = min(screenMinHeight, m.WorkArea.Bottom-m.WorkArea.Top)
	}

	return
}

// InitialPlacement calculates the initial window placement centered on a primary
// monitor. Window size is determined by logic implemented in CalcPlacementSize.
//
// Parameters:
// - desiredWidth, desiredHeight: preferred window size in logical units
// - minWidth, minHeight: minimum required window size in logical units
// - margin: minimum distance from work area edges in logical units
//
// Returns a Rect with the calculated window position and size in screen units.
func InitialPlacement(desiredWidth, desiredHeight, minWidth, minHeight, margin int) Rect {
	monitors := GetMonitors()
	if len(monitors) == 0 {
		width, height := CalcPlacementSize(nil, desiredWidth, desiredHeight, minWidth, minHeight, margin)
		return Rect{
			Left:   0,
			Top:    0,
			Right:  width,
			Bottom: height,
		}
	}

	// Find default mon (containing 0,0 or first available)
	mon := FindPrimaryMonitor(monitors)

	width, height := CalcPlacementSize(mon, desiredWidth, desiredHeight, minWidth, minHeight, margin)

	// Center the window in work area
	centerX := (mon.WorkArea.Left + mon.WorkArea.Right) / 2
	centerY := (mon.WorkArea.Top + mon.WorkArea.Bottom) / 2

	return Rect{
		Left:   centerX - width/2,
		Top:    centerY - height/2,
		Right:  centerX + (width+1)/2,
		Bottom: centerY + (height+1)/2,
	}
}
