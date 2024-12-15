package multimon

import "testing"

func TestCoordinateConversion(t *testing.T) {
	// Test monitor with different logical and physical coordinates
	m := Monitor{
		LogicalBounds: Rect{
			Left:   0,
			Top:    0,
			Right:  1920,
			Bottom: 1080,
		},
		PhysicalBounds: Rect{
			Left:   0,
			Top:    0,
			Right:  3840,
			Bottom: 2160,
		},
	}

	tests := []struct {
		name     string
		logical  Rect
		physical Rect
	}{
		{
			name: "origin",
			logical: Rect{
				Left: 0, Top: 0, Right: 100, Bottom: 100,
			},
			physical: Rect{
				Left: 0, Top: 0, Right: 200, Bottom: 200,
			},
		},
		{
			name: "center",
			logical: Rect{
				Left: 960, Top: 540, Right: 1060, Bottom: 640,
			},
			physical: Rect{
				Left: 1920, Top: 1080, Right: 2120, Bottom: 1280,
			},
		},
		{
			name: "bottom-right",
			logical: Rect{
				Left: 1820, Top: 980, Right: 1920, Bottom: 1080,
			},
			physical: Rect{
				Left: 3640, Top: 1960, Right: 3840, Bottom: 2160,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"-LogicalToPhysical", func(t *testing.T) {
			got := LogicalToPhysical(m, tt.logical)
			if got != tt.physical {
				t.Errorf("LogicalToPhysical() = %+v, want %+v", got, tt.physical)
			}
		})

		t.Run(tt.name+"-PhysicalToLogical", func(t *testing.T) {
			got := PhysicalToLogical(m, tt.physical)
			if got != tt.logical {
				t.Errorf("PhysicalToLogical() = %+v, want %+v", got, tt.logical)
			}
		})
	}
}

func TestContainsPoint(t *testing.T) {
	m := Monitor{
		LogicalBounds: Rect{
			Left:   0,
			Top:    0,
			Right:  1920,
			Bottom: 1080,
		},
		PhysicalBounds: Rect{
			Left:   0,
			Top:    0,
			Right:  3840,
			Bottom: 2160,
		},
	}

	tests := []struct {
		name         string
		x, y         int
		wantLogical  bool
		wantPhysical bool
		isPhysical   bool
	}{
		{
			name: "inside-logical",
			x:    960, y: 540,
			wantLogical: true,
			isPhysical:  false,
		},
		{
			name: "outside-logical",
			x:    2000, y: 540,
			wantLogical: false,
			isPhysical:  false,
		},
		{
			name: "inside-physical",
			x:    1920, y: 1080,
			wantPhysical: true,
			isPhysical:   true,
		},
		{
			name: "outside-physical",
			x:    4000, y: 1080,
			wantPhysical: false,
			isPhysical:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.isPhysical {
				got := ContainsLogicalPoint(m, tt.x, tt.y)
				if got != tt.wantLogical {
					t.Errorf("ContainsLogicalPoint(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.wantLogical)
				}
			} else {
				got := ContainsPhysicalPoint(m, tt.x, tt.y)
				if got != tt.wantPhysical {
					t.Errorf("ContainsPhysicalPoint(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.wantPhysical)
				}
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	m := Monitor{
		LogicalBounds: Rect{
			Left:   100,
			Top:    200,
			Right:  2020,
			Bottom: 1280,
		},
		PhysicalBounds: Rect{
			Left:   200,
			Top:    400,
			Right:  4040,
			Bottom: 2560,
		},
	}

	// Test various points to ensure they convert back correctly
	testPoints := []Rect{
		{Left: 150, Top: 250, Right: 1000, Bottom: 800},
		{Left: 1500, Top: 1000, Right: 1900, Bottom: 1200},
		{Left: 500, Top: 600, Right: 1500, Bottom: 900},
	}

	for i, point := range testPoints {
		physical := LogicalToPhysical(m, point)
		logical := PhysicalToLogical(m, physical)
		if logical != point {
			t.Errorf("Point %d: Round trip conversion failed. Started with %+v, got back %+v", i, point, logical)
		}
	}
}
