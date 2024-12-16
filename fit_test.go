package multimon

import (
	"strings"
	"testing"
)

func TestFitToMonitor(t *testing.T) {
	monitor := &Monitor{
		Bounds:   Rect{0, 0, 1920, 1080},
		WorkArea: Rect{0, 40, 1920, 1040},
		Scale:    1.5,
	}

	tests := []struct {
		name        string
		monitor     *Monitor
		mode        FitMode
		window      Rect
		scale       float64
		want        Rect
		wantScale   float64
		wantErr     bool
		errContains string
	}{
		// Basic fitting tests
		{
			name:      "window in screen units",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 300, 300},
			scale:     0.0,
			want:      Rect{100, 100, 300, 300},
			wantScale: 1.5,
		},
		{
			name:      "window with custom scale 1.25",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 300, 300},
			scale:     1.25,
			want:      Rect{100, 100, 340, 340},
			wantScale: 1.5,
		},
		{
			name:      "window with custom scale 2.0",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 300, 300},
			scale:     2.0,
			want:      Rect{100, 100, 250, 250},
			wantScale: 1.5,
		},

		// Test rescaling that causes window to protrude
		{
			name:      "window near edge with scale up",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{1800, 900, 1900, 1000},
			scale:     0.5,                         // Scaling from 0.5 to 1.5 means 3x size increase
			want:      Rect{1620, 780, 1920, 1080}, // Window grows and gets moved to stay within bounds
			wantScale: 1.5,
		},

		// Mode tests
		{
			name:      "fit to work area",
			monitor:   monitor,
			mode:      FitModeWorkArea,
			window:    Rect{100, 20, 300, 300},
			scale:     0.0,
			want:      Rect{100, 40, 300, 320},
			wantScale: 1.5,
		},

		// Boundary tests
		{
			name:      "window too large",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 2100, 1200},
			scale:     0.0,
			want:      Rect{0, 0, 1920, 1080},
			wantScale: 1.5,
		},
		{
			name:      "window outside bounds left",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{-100, 100, 100, 300},
			scale:     0.0,
			want:      Rect{0, 100, 200, 300},
			wantScale: 1.5,
		},
		{
			name:      "window outside bounds right",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{1900, 100, 2100, 300},
			scale:     0.0,
			want:      Rect{1720, 100, 1920, 300},
			wantScale: 1.5,
		},
		{
			name:      "window outside bounds top",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, -100, 300, 100},
			scale:     0.0,
			want:      Rect{100, 0, 300, 200},
			wantScale: 1.5,
		},
		{
			name:      "window outside bounds bottom",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{100, 1000, 300, 1200},
			scale:     0.0,
			want:      Rect{100, 880, 300, 1080},
			wantScale: 1.5,
		},
		{
			name:      "window outside bounds all sides",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{-100, -100, 2100, 1200},
			scale:     0.0,
			want:      Rect{0, 0, 1920, 1080},
			wantScale: 1.5,
		},

		// Scale edge cases
		{
			name:      "large window with 2x scale",
			monitor:   monitor,
			mode:      FitModeBounds,
			window:    Rect{0, 0, 2000, 2000},
			scale:     2.0,
			want:      Rect{0, 0, 1500, 1080}, // First scaled (2000 * 1.5/2.0 = 1500), then height clamped to 1080
			wantScale: 1.5,
		},

		// Work area edge cases
		{
			name: "work area equals bounds",
			monitor: &Monitor{
				Bounds:   Rect{0, 0, 1920, 1080},
				WorkArea: Rect{0, 0, 1920, 1080},
				Scale:    1.5,
			},
			mode:      FitModeWorkArea,
			window:    Rect{100, 100, 300, 300},
			scale:     0.0,
			want:      Rect{100, 100, 300, 300},
			wantScale: 1.5,
		},
		{
			name: "small work area",
			monitor: &Monitor{
				Bounds:   Rect{0, 0, 1920, 1080},
				WorkArea: Rect{100, 100, 200, 200},
				Scale:    1.5,
			},
			mode:      FitModeWorkArea,
			window:    Rect{0, 0, 200, 200},
			scale:     0.0,
			want:      Rect{100, 100, 200, 200},
			wantScale: 1.5,
		},

		// Error cases
		{
			name:        "nil monitor with zero scale",
			monitor:     nil,
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.0,
			wantErr:     true,
			errContains: "nil",
		},
		{
			name:        "nil monitor with custom scale",
			monitor:     nil,
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       1.25,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.25,
			wantErr:     true,
			errContains: "nil",
		},
		{
			name:        "invalid window dimensions",
			monitor:     monitor,
			mode:        FitModeBounds,
			window:      Rect{300, 300, 100, 100},
			scale:       0.0,
			wantErr:     true,
			errContains: "invalid dimensions",
		},
		{
			name: "invalid monitor scale",
			monitor: &Monitor{
				Bounds:   Rect{0, 0, 1920, 1080},
				WorkArea: Rect{0, 40, 1920, 1040},
				Scale:    0.0,
			},
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			wantErr:     true,
			errContains: "invalid scale",
		},
		{
			name: "invalid monitor bounds",
			monitor: &Monitor{
				Bounds:   Rect{100, 0, 0, 1080},
				WorkArea: Rect{0, 40, 1920, 1040},
				Scale:    1.5,
			},
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			wantErr:     true,
			errContains: "invalid bounds",
		},
		{
			name: "invalid monitor work area",
			monitor: &Monitor{
				Bounds:   Rect{0, 0, 1920, 1080},
				WorkArea: Rect{0, 1040, 1920, 40},
				Scale:    1.5,
			},
			mode:        FitModeWorkArea,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			wantErr:     true,
			errContains: "invalid work area",
		},
		{
			name: "work area outside bounds",
			monitor: &Monitor{
				Bounds:   Rect{0, 0, 1920, 1080},
				WorkArea: Rect{-100, -100, 2000, 2000},
				Scale:    1.5,
			},
			mode:        FitModeWorkArea,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			wantErr:     true,
			errContains: "work area outside bounds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotScale, err := FitToMonitor(tt.monitor, tt.mode, tt.window, tt.scale)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got rect %v, want %v", got, tt.want)
			}
			if gotScale != tt.wantScale {
				t.Errorf("got scale %v, want %v", gotScale, tt.wantScale)
			}
		})
	}
}

func TestFitToNearestMonitor(t *testing.T) {
	// Common test monitors
	monitors := []Monitor{
		{
			// Primary monitor: 1920x1080, taskbar at bottom
			Bounds:   Rect{0, 0, 1920, 1080},
			WorkArea: Rect{0, 0, 1920, 1040},
			Scale:    1.0,
		},
		{
			// Secondary monitor: 1920x1080, taskbar at top
			Bounds:   Rect{1920, 0, 3840, 1080},
			WorkArea: Rect{1920, 40, 3840, 1080},
			Scale:    1.5,
		},
		{
			// Small monitor: 1280x720, no taskbar
			Bounds:   Rect{0, 1080, 1280, 1800},
			WorkArea: Rect{0, 1080, 1280, 1800},
			Scale:    1.25,
		},
	}

	tests := []struct {
		name        string
		monitors    []Monitor
		mode        FitMode
		window      Rect
		scale       float64
		minWidth    int
		minHeight   int
		want        Rect
		wantScale   float64
		wantErr     bool
		errContains string
	}{
		// Monitor filtering tests
		{
			name:        "no monitors with zero scale",
			monitors:    nil,
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.0, // default to 1.0 when no monitors and zero scale
			wantErr:     true,
			errContains: "no monitors available",
		},
		{
			name:        "no monitors with custom scale",
			monitors:    nil,
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       1.25,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.25, // use provided scale when no monitors
			wantErr:     true,
			errContains: "no monitors available",
		},
		{
			name: "all monitors invalid with zero scale",
			monitors: []Monitor{
				{Bounds: Rect{0, 0, -100, 100}, Scale: 1.0},
				{Bounds: Rect{0, 0, 100, -100}, Scale: 1.0},
			},
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       0.0,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.0, // default to 1.0 when no valid monitors and zero scale
			wantErr:     true,
			errContains: "no valid monitors found",
		},
		{
			name: "all monitors invalid with custom scale",
			monitors: []Monitor{
				{Bounds: Rect{0, 0, -100, 100}, Scale: 1.0},
				{Bounds: Rect{0, 0, 100, -100}, Scale: 1.0},
			},
			mode:        FitModeBounds,
			window:      Rect{100, 100, 300, 300},
			scale:       1.25,
			want:        Rect{100, 100, 300, 300},
			wantScale:   1.25, // use provided scale when no valid monitors
			wantErr:     true,
			errContains: "no valid monitors found",
		},

		// Minimum size tests
		{
			name:      "fits work area minimum size",
			monitors:  monitors,
			mode:      FitModeWorkArea,
			window:    Rect{100, 100, 300, 300},
			scale:     0.0,
			minWidth:  1000,
			minHeight: 900,
			want:      Rect{100, 100, 300, 300},
			wantScale: 1.0, // primary monitor's scale
		},
		{
			name:      "fits bounds but not work area",
			monitors:  monitors,
			mode:      FitModeWorkArea,
			window:    Rect{100, 100, 300, 300},
			scale:     0.0,
			minWidth:  1000,
			minHeight: 1000,
			want:      Rect{100, 100, 300, 300},
			wantScale: 1.0, // primary monitor's scale
		},
		{
			name:      "ignores minimum size if no monitor fits",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{100, 1080, 300, 1280},
			scale:     0.0,
			minWidth:  2000,
			minHeight: 2000,
			want:      Rect{100, 1080, 300, 1280},
			wantScale: 1.25, // small monitor's scale
		},

		// Overlap tests
		{
			name:      "window fully inside primary monitor",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 300, 300},
			scale:     0.0,
			want:      Rect{100, 100, 300, 300},
			wantScale: 1.0, // primary monitor's scale
		},
		{
			name:      "window overlaps two monitors",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{1800, 100, 2100, 300},
			scale:     0.0,
			want:      Rect{1920, 100, 2220, 300}, // Window moved to fit within monitor bounds
			wantScale: 1.5,                        // secondary monitor's scale (more overlap in screen units)
		},
		{
			name:      "window mostly in secondary monitor",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{1900, 100, 2100, 300},
			scale:     0.0,
			want:      Rect{1920, 100, 2120, 300}, // Window moved to fit within monitor bounds
			wantScale: 1.5,                        // secondary monitor's scale
		},

		// Distance tests
		{
			name:      "window outside all monitors",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{4000, 100, 4200, 300},
			scale:     0.0,
			want:      Rect{3640, 100, 3840, 300},
			wantScale: 1.5, // secondary monitor's scale (nearest)
		},
		{
			name:      "window closest to small monitor",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{100, 1900, 300, 2100},
			scale:     0.0,
			want:      Rect{100, 1600, 300, 1800},
			wantScale: 1.25, // small monitor's scale
		},

		// Scale tests
		{
			name:      "window with custom scale",
			monitors:  monitors,
			mode:      FitModeBounds,
			window:    Rect{100, 100, 300, 300},
			scale:     2.0,
			want:      Rect{100, 100, 200, 200}, // Shrinks from top-left: width/height * (1.0/2.0) = 200 * 0.5 = 100
			wantScale: 1.0,                      // primary monitor's scale
		},

		// Window scale tests
		{
			name:     "window designed for 1.0 scale on 1.5 monitor",
			monitors: monitors[1:2], // Only secondary monitor
			mode:     FitModeBounds,
			window:   Rect{2000, 100, 2200, 300},
			scale:    1.0,
			want: Rect{
				Left:   2000,
				Top:    100,
				Right:  2300, // Width increased by 1.5x
				Bottom: 400,  // Height increased by 1.5x
			},
			wantScale: 1.5,
		},
		{
			name:     "window designed for 2.0 scale on 1.0 monitor",
			monitors: monitors[0:1], // Only primary monitor
			mode:     FitModeBounds,
			window:   Rect{100, 100, 500, 300},
			scale:    2.0,
			want: Rect{
				Left:   100,
				Top:    100,
				Right:  300, // Width reduced by 0.5x
				Bottom: 200, // Height reduced by 0.5x
			},
			wantScale: 1.0,
		},
		{
			name:     "window designed for 1.25 scale on 1.5 monitor",
			monitors: monitors[1:2], // Only secondary monitor
			mode:     FitModeBounds,
			window:   Rect{2000, 100, 2250, 350},
			scale:    1.25,
			want: Rect{
				Left:   2000,
				Top:    100,
				Right:  2300, // Width increased by 1.5/1.25x
				Bottom: 400,  // Height increased by 1.5/1.25x
			},
			wantScale: 1.5,
		},
		{
			name:     "window designed for 1.5 scale on 1.25 monitor",
			monitors: monitors[2:3], // Only small monitor
			mode:     FitModeBounds,
			window:   Rect{100, 1100, 400, 1400},
			scale:    1.5,
			want: Rect{
				Left:   100,
				Top:    1100,
				Right:  350,  // Width reduced by 1.25/1.5x
				Bottom: 1350, // Height reduced by 1.25/1.5x
			},
			wantScale: 1.25,
		},
		{
			name:     "window with zero scale on 1.5 monitor",
			monitors: monitors[1:2], // Only secondary monitor
			mode:     FitModeBounds,
			window:   Rect{2000, 100, 2200, 300},
			scale:    0.0,
			want: Rect{
				Left:   2000,
				Top:    100,
				Right:  2200,
				Bottom: 300,
			},
			wantScale: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotScale, err := FitToNearestMonitor(tt.monitors, tt.mode, tt.window, tt.scale, tt.minWidth, tt.minHeight)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got rect %v, want %v", got, tt.want)
			}
			if gotScale != tt.wantScale {
				t.Errorf("got scale %v, want %v", gotScale, tt.wantScale)
			}
		})
	}
}
