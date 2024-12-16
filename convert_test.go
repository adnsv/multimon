package multimon

import "testing"

func TestCoordinateConversion(t *testing.T) {
	// Test monitor with 200% scaling
	m := Monitor{
		Bounds: Rect{
			Left:   0,
			Top:    0,
			Right:  3840,
			Bottom: 2160,
		},
		Scale: 2.0,
	}

	tests := []struct {
		name    string
		logical Rect
		screen  Rect
	}{
		{
			name: "origin",
			logical: Rect{
				Left: 0, Top: 0, Right: 100, Bottom: 100,
			},
			screen: Rect{
				Left: 0, Top: 0, Right: 200, Bottom: 200,
			},
		},
		{
			name: "center",
			logical: Rect{
				Left: 960, Top: 540, Right: 1060, Bottom: 640,
			},
			screen: Rect{
				Left: 1920, Top: 1080, Right: 2120, Bottom: 1280,
			},
		},
		{
			name: "bottom-right",
			logical: Rect{
				Left: 1820, Top: 980, Right: 1920, Bottom: 1080,
			},
			screen: Rect{
				Left: 3640, Top: 1960, Right: 3840, Bottom: 2160,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"-LogicalToScreenRect", func(t *testing.T) {
			got := LogicalToScreenRect(m, tt.logical)
			if got != tt.screen {
				t.Errorf("LogicalToScreenRect() = %+v, want %+v", got, tt.screen)
			}
		})

		t.Run(tt.name+"-ScreenToLogicalRect", func(t *testing.T) {
			got := ScreenToLogicalRect(m, tt.screen)
			if got != tt.logical {
				t.Errorf("ScreenToLogicalRect() = %+v, want %+v", got, tt.logical)
			}
		})
	}
}

func TestPointConversion(t *testing.T) {
	// Test monitor with 200% scaling
	m := Monitor{
		Bounds: Rect{
			Left:   0,
			Top:    0,
			Right:  3840,
			Bottom: 2160,
		},
		Scale: 2.0,
	}

	tests := []struct {
		name     string
		logicalX int
		logicalY int
		screenX  int
		screenY  int
	}{
		{
			name:     "origin",
			logicalX: 0,
			logicalY: 0,
			screenX:  0,
			screenY:  0,
		},
		{
			name:     "center",
			logicalX: 960,
			logicalY: 540,
			screenX:  1920,
			screenY:  1080,
		},
		{
			name:     "bottom-right",
			logicalX: 1920,
			logicalY: 1080,
			screenX:  3840,
			screenY:  2160,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"-LogicalToScreenPoint", func(t *testing.T) {
			got := LogicalToScreenPoint(m, tt.logicalX, tt.logicalY)
			want := Point{X: tt.screenX, Y: tt.screenY}
			if got != want {
				t.Errorf("LogicalToScreenPoint(%d, %d) = %+v, want %+v", tt.logicalX, tt.logicalY, got, want)
			}
		})

		t.Run(tt.name+"-ScreenToLogicalPoint", func(t *testing.T) {
			got := ScreenToLogicalPoint(m, tt.screenX, tt.screenY)
			want := Point{X: tt.logicalX, Y: tt.logicalY}
			if got != want {
				t.Errorf("ScreenToLogicalPoint(%d, %d) = %+v, want %+v", tt.screenX, tt.screenY, got, want)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	m := Monitor{
		Bounds: Rect{
			Left:   0,
			Top:    0,
			Right:  3840,
			Bottom: 2160,
		},
		Scale: 2.0,
	}

	// Test various points to ensure they convert back correctly
	testPoints := []Rect{
		{Left: 150, Top: 250, Right: 1000, Bottom: 800},
		{Left: 1500, Top: 1000, Right: 1900, Bottom: 1200},
		{Left: 500, Top: 600, Right: 1500, Bottom: 900},
	}

	for i, point := range testPoints {
		screen := LogicalToScreenRect(m, point)
		logical := ScreenToLogicalRect(m, screen)
		if logical != point {
			t.Errorf("Point %d: Round trip conversion failed. Started with %+v, got back %+v", i, point, logical)
		}
	}

	// Test point round trips
	testCoords := []struct{ x, y int }{
		{150, 250},
		{1500, 1000},
		{500, 600},
	}

	for i, coord := range testCoords {
		screen := LogicalToScreenPoint(m, coord.x, coord.y)
		logical := ScreenToLogicalPoint(m, screen.X, screen.Y)
		if logical.X != coord.x || logical.Y != coord.y {
			t.Errorf("Coord %d: Round trip conversion failed. Started with (%d,%d), got back (%d,%d)",
				i, coord.x, coord.y, logical.X, logical.Y)
		}
	}
}
