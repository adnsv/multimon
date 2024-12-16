package multimon

import "testing"

func TestCalcPlacementSize(t *testing.T) {
	monitor := &Monitor{
		Bounds:   Rect{0, 0, 1920, 1080},
		WorkArea: Rect{0, 40, 1920, 1040},
		Scale:    1.5,
	}

	tests := []struct {
		name          string
		monitor       *Monitor
		desiredWidth  int
		desiredHeight int
		minWidth      int
		minHeight     int
		margin        int
		wantWidth     int
		wantHeight    int
	}{
		// No monitor tests (1:1 scale)
		{
			name:          "nil monitor uses desired size",
			monitor:       nil,
			desiredWidth:  800,
			desiredHeight: 600,
			minWidth:      400,
			minHeight:     300,
			margin:        10,
			wantWidth:     800,
			wantHeight:    600,
		},
		{
			name:          "nil monitor enforces minimum size",
			monitor:       nil,
			desiredWidth:  300,
			desiredHeight: 200,
			minWidth:      400,
			minHeight:     300,
			margin:        10,
			wantWidth:     400,
			wantHeight:    300,
		},

		// Basic sizing tests
		{
			name:          "fits desired size with margins",
			monitor:       monitor,
			desiredWidth:  800,
			desiredHeight: 600,
			minWidth:      400,
			minHeight:     300,
			margin:        20,
			wantWidth:     1200, // 800 * 1.5
			wantHeight:    900,  // 600 * 1.5
		},
		{
			name:          "clamps to available space with margins",
			monitor:       monitor,
			desiredWidth:  1200,
			desiredHeight: 900,
			minWidth:      400,
			minHeight:     300,
			margin:        20,
			wantWidth:     1800, // 1920 - (20 * 1.5 * 2)
			wantHeight:    940,  // 1040 - 40 - (20 * 1.5 * 2)
		},

		// Minimum size tests
		{
			name:          "enforces minimum size within margins",
			monitor:       monitor,
			desiredWidth:  300,
			desiredHeight: 200,
			minWidth:      400,
			minHeight:     300,
			margin:        20,
			wantWidth:     600, // 400 * 1.5
			wantHeight:    450, // 300 * 1.5
		},
		{
			name:          "allows using margin area for minimum size",
			monitor:       monitor,
			desiredWidth:  300,
			desiredHeight: 200,
			minWidth:      1300,
			minHeight:     700,
			margin:        20,
			wantWidth:     1920, // full width
			wantHeight:    1000, // work area height
		},

		// Edge cases
		{
			name:          "zero margin",
			monitor:       monitor,
			desiredWidth:  800,
			desiredHeight: 600,
			minWidth:      400,
			minHeight:     300,
			margin:        0,
			wantWidth:     1200, // 800 * 1.5
			wantHeight:    900,  // 600 * 1.5
		},
		{
			name:          "zero desired size uses minimum",
			monitor:       monitor,
			desiredWidth:  0,
			desiredHeight: 0,
			minWidth:      400,
			minHeight:     300,
			margin:        20,
			wantWidth:     600, // 400 * 1.5
			wantHeight:    450, // 300 * 1.5
		},
		{
			name:          "zero minimum size allows any size",
			monitor:       monitor,
			desiredWidth:  800,
			desiredHeight: 600,
			minWidth:      0,
			minHeight:     0,
			margin:        20,
			wantWidth:     1200, // 800 * 1.5
			wantHeight:    900,  // 600 * 1.5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWidth, gotHeight := CalcPlacementSize(tt.monitor, tt.desiredWidth, tt.desiredHeight, tt.minWidth, tt.minHeight, tt.margin)
			if gotWidth != tt.wantWidth || gotHeight != tt.wantHeight {
				t.Errorf("CalcPlacementSize() = (%v, %v), want (%v, %v)", gotWidth, gotHeight, tt.wantWidth, tt.wantHeight)
			}
		})
	}
}
