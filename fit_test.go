package multimon

import (
	"errors"
	"testing"
)

func TestFitToMonitor_InputValidation(t *testing.T) {
	monitor := Monitor{
		LogicalBounds: Rect{
			Left: 0, Top: 0, Right: 1920, Bottom: 1080,
		},
		LogicalWorkArea: Rect{
			Left: 0, Top: 40, Right: 1920, Bottom: 1040,
		},
	}

	invalidMonitor := Monitor{
		LogicalBounds: Rect{
			Left: 0, Top: 0, Right: -1920, Bottom: 1080,
		},
		LogicalWorkArea: Rect{
			Left: 0, Top: 40, Right: -1920, Bottom: 1040,
		},
	}

	tests := []struct {
		name      string
		monitor   Monitor
		window    Rect
		wantError bool
	}{
		{
			name:      "window with negative width",
			monitor:   monitor,
			window:    Rect{Left: 500, Top: 100, Right: 100, Bottom: 400},
			wantError: true,
		},
		{
			name:      "window with negative height",
			monitor:   monitor,
			window:    Rect{Left: 100, Top: 500, Right: 400, Bottom: 100},
			wantError: true,
		},
		{
			name:      "monitor with negative width",
			monitor:   invalidMonitor,
			window:    Rect{Left: 100, Top: 100, Right: 400, Bottom: 400},
			wantError: true,
		},
		{
			name:    "valid dimensions",
			monitor: monitor,
			window:  Rect{Left: 100, Top: 100, Right: 400, Bottom: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FitToMonitor(tt.monitor, FitModeBounds, tt.window)
			if tt.wantError {
				if err == nil {
					t.Error("FitToMonitor() error = nil, want error")
				}
				if !errors.Is(err, ErrInvalidDimensions) {
					t.Errorf("FitToMonitor() error = %v, want %v", err, ErrInvalidDimensions)
				}
				return
			}
			if err != nil {
				t.Errorf("FitToMonitor() unexpected error = %v", err)
			}
		})
	}
}

func TestFitToMonitor_BasicFitting(t *testing.T) {
	monitor := Monitor{
		LogicalBounds: Rect{
			Left: 0, Top: 0, Right: 1920, Bottom: 1080,
		},
		LogicalWorkArea: Rect{
			Left: 0, Top: 40, Right: 1920, Bottom: 1040,
		},
	}

	tests := []struct {
		name   string
		mode   FitMode
		window Rect
		want   Rect
	}{
		{
			name:   "window fits within bounds",
			mode:   FitModeBounds,
			window: Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			want:   Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
		},
		{
			name:   "window outside right edge",
			mode:   FitModeBounds,
			window: Rect{Left: 1720, Top: 100, Right: 2120, Bottom: 400},
			want:   Rect{Left: 1520, Top: 100, Right: 1920, Bottom: 400},
		},
		{
			name:   "window outside bottom edge of work area",
			mode:   FitModeWorkArea,
			window: Rect{Left: 100, Top: 900, Right: 500, Bottom: 1200},
			want:   Rect{Left: 100, Top: 740, Right: 500, Bottom: 1040},
		},
		{
			name:   "window too large for bounds",
			mode:   FitModeBounds,
			window: Rect{Left: -100, Top: -100, Right: 2020, Bottom: 1180},
			want:   Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080},
		},
		{
			name:   "window too large for work area",
			mode:   FitModeWorkArea,
			window: Rect{Left: -100, Top: -100, Right: 2020, Bottom: 1180},
			want:   Rect{Left: 0, Top: 40, Right: 1920, Bottom: 1040},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FitToMonitor(monitor, tt.mode, tt.window)
			if err != nil {
				t.Errorf("FitToMonitor() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FitToMonitor() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFitToNearestMonitor_OverlapBased(t *testing.T) {
	monitors := []Monitor{
		{
			// First monitor: 1920x1080
			LogicalBounds: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 0, Top: 40, Right: 1920, Bottom: 1040,
			},
		},
		{
			// Second monitor: 1920x1080, to the right
			LogicalBounds: Rect{
				Left: 1920, Top: 0, Right: 3840, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 1920, Top: 40, Right: 3840, Bottom: 1040,
			},
		},
	}

	tests := []struct {
		name   string
		window Rect
		want   Rect
	}{
		{
			name:   "window fully within first monitor",
			window: Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			want:   Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
		},
		{
			name:   "window overlaps more with second monitor",
			window: Rect{Left: 1800, Top: 100, Right: 2200, Bottom: 400},
			want:   Rect{Left: 1920, Top: 100, Right: 2320, Bottom: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FitToNearestMonitor(monitors, FitModeBounds, tt.window, 0, 0)
			if err != nil {
				t.Errorf("FitToNearestMonitor() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FitToNearestMonitor() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFitToNearestMonitor_MinimumSize(t *testing.T) {
	monitors := []Monitor{
		{
			// First monitor: 1920x1080
			LogicalBounds: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
		},
		{
			// Second monitor: 2560x1440
			LogicalBounds: Rect{
				Left: 1920, Top: 0, Right: 4480, Bottom: 1440,
			},
			LogicalWorkArea: Rect{
				Left: 1920, Top: 0, Right: 4480, Bottom: 1440,
			},
		},
	}

	tests := []struct {
		name      string
		window    Rect
		minWidth  int
		minHeight int
		want      Rect
	}{
		{
			name:      "window fits min size on current monitor",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			minWidth:  800,
			minHeight: 600,
			want:      Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
		},
		{
			name:      "window needs larger monitor for min size",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			minWidth:  2000,
			minHeight: 1200,
			want:      Rect{Left: 1920, Top: 0, Right: 4480, Bottom: 1440},
		},
		{
			name:      "no monitor can fit min size, use largest",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			minWidth:  3000,
			minHeight: 2000,
			want:      Rect{Left: 1920, Top: 0, Right: 4480, Bottom: 1440},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FitToNearestMonitor(monitors, FitModeBounds, tt.window, tt.minWidth, tt.minHeight)
			if err != nil {
				t.Errorf("FitToNearestMonitor() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FitToNearestMonitor() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFitToNearestMonitor_EdgeDistance(t *testing.T) {
	monitors := []Monitor{
		{
			// First monitor: 1920x1080
			LogicalBounds: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
		},
		{
			// Second monitor: 1920x1080, to the right
			LogicalBounds: Rect{
				Left: 1920, Top: 0, Right: 3840, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 1920, Top: 0, Right: 3840, Bottom: 1080,
			},
		},
		{
			// Third monitor: 1920x1080, below
			LogicalBounds: Rect{
				Left: 0, Top: 1080, Right: 1920, Bottom: 2160,
			},
			LogicalWorkArea: Rect{
				Left: 0, Top: 1080, Right: 1920, Bottom: 2160,
			},
		},
	}

	tests := []struct {
		name   string
		window Rect
		want   Rect
	}{
		{
			name:   "window closer to first monitor",
			window: Rect{Left: -200, Top: 100, Right: 100, Bottom: 400},
			want:   Rect{Left: 0, Top: 100, Right: 300, Bottom: 400},
		},
		{
			name:   "window closer to second monitor",
			window: Rect{Left: 2000, Top: 100, Right: 2300, Bottom: 400},
			want:   Rect{Left: 2000, Top: 100, Right: 2300, Bottom: 400},
		},
		{
			name:   "window closer to third monitor",
			window: Rect{Left: 100, Top: 1000, Right: 400, Bottom: 1300},
			want:   Rect{Left: 100, Top: 1080, Right: 400, Bottom: 1380},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FitToNearestMonitor(monitors, FitModeBounds, tt.window, 0, 0)
			if err != nil {
				t.Errorf("FitToNearestMonitor() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FitToNearestMonitor() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFitToNearestMonitor_EdgeCases(t *testing.T) {
	monitors := []Monitor{
		{
			// First monitor: 1920x1080
			LogicalBounds: Rect{
				Left: 0, Top: 0, Right: 1920, Bottom: 1080,
			},
			LogicalWorkArea: Rect{
				Left: 0, Top: 40, Right: 1920, Bottom: 1040,
			},
		},
		{
			// Second monitor: 2560x1440, to the right
			LogicalBounds: Rect{
				Left: 1920, Top: 0, Right: 4480, Bottom: 1440,
			},
			LogicalWorkArea: Rect{
				Left: 1920, Top: 40, Right: 4480, Bottom: 1400,
			},
		},
	}

	tests := []struct {
		name      string
		window    Rect
		mode      FitMode
		minWidth  int
		minHeight int
		want      Rect
		wantErr   error
	}{
		{
			name:    "zero-size window",
			window:  Rect{Left: 100, Top: 100, Right: 100, Bottom: 100},
			mode:    FitModeBounds,
			want:    Rect{Left: 100, Top: 100, Right: 100, Bottom: 100},
			wantErr: ErrInvalidDimensions,
		},
		{
			name:    "negative-size window",
			window:  Rect{Left: 200, Top: 200, Right: 100, Bottom: 100},
			mode:    FitModeBounds,
			want:    Rect{Left: 200, Top: 200, Right: 100, Bottom: 100},
			wantErr: ErrInvalidDimensions,
		},
		{
			name:      "minimum size larger than any monitor",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			mode:      FitModeBounds,
			minWidth:  3000,
			minHeight: 2000,
			want:      Rect{Left: 1920, Top: 0, Right: 4480, Bottom: 1440}, // Should return largest monitor
		},
		{
			name:      "minimum size exactly matches monitor",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			mode:      FitModeBounds,
			minWidth:  1920,
			minHeight: 1080,
			want:      Rect{Left: 100, Top: 100, Right: 500, Bottom: 400}, // Should keep original position
		},
		{
			name:   "window spans monitor boundary exactly",
			window: Rect{Left: 1919, Top: 100, Right: 1920, Bottom: 400},
			mode:   FitModeBounds,
			want:   Rect{Left: 1919, Top: 100, Right: 1920, Bottom: 400},
		},
		{
			name:   "window entirely outside all monitors",
			window: Rect{Left: 5000, Top: 5000, Right: 5500, Bottom: 5500},
			mode:   FitModeBounds,
			want:   Rect{Left: 1920, Top: 0, Right: 2420, Bottom: 500}, // Should move to nearest monitor
		},
		{
			name:   "work area mode with taskbar",
			window: Rect{Left: 100, Top: 0, Right: 500, Bottom: 100},
			mode:   FitModeWorkArea,
			want:   Rect{Left: 100, Top: 40, Right: 500, Bottom: 140}, // Should adjust for taskbar
		},
		{
			name:   "window larger than work area but smaller than bounds",
			window: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080},
			mode:   FitModeWorkArea,
			want:   Rect{Left: 0, Top: 40, Right: 1920, Bottom: 1040}, // Should fit to work area
		},
		{
			name:      "minimum size fits work area but not taskbar area",
			window:    Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			mode:      FitModeWorkArea,
			minWidth:  1920,
			minHeight: 1000, // Just fits in work area (1040-40=1000)
			want:      Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FitToNearestMonitor(monitors, tt.mode, tt.window, tt.minWidth, tt.minHeight)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("FitToNearestMonitor() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Errorf("FitToNearestMonitor() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("FitToNearestMonitor() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
