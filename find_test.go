package multimon

import "testing"

func createTestMonitors() []Monitor {
	return []Monitor{
		{
			LogicalBounds:  Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080},
			PhysicalBounds: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080},
		},
		{
			LogicalBounds:  Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080},
			PhysicalBounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080},
		},
		{
			LogicalBounds:  Rect{Left: 0, Top: 1080, Right: 1920, Bottom: 2160},
			PhysicalBounds: Rect{Left: 0, Top: 1080, Right: 1920, Bottom: 2160},
		},
	}
}

func TestFindMonitorFromLogicalPoint(t *testing.T) {
	monitors := createTestMonitors()
	tests := []struct {
		name      string
		x, y      int
		wantIndex int // -1 means no monitor should be found
	}{
		{name: "first monitor center", x: 960, y: 540, wantIndex: 0},
		{name: "second monitor center", x: 2880, y: 540, wantIndex: 1},
		{name: "third monitor center", x: 960, y: 1620, wantIndex: 2},
		{name: "first monitor edge", x: 1919, y: 1079, wantIndex: 0},
		{name: "second monitor edge", x: 1920, y: 0, wantIndex: 1},
		{name: "outside all monitors", x: 4000, y: 4000, wantIndex: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromLogicalPoint(monitors, tt.x, tt.y)
			if tt.wantIndex == -1 {
				if got != nil {
					t.Errorf("FindMonitorFromLogicalPoint(%d, %d) = %+v, want nil", tt.x, tt.y, got)
				}
			} else {
				if got == nil {
					t.Errorf("FindMonitorFromLogicalPoint(%d, %d) = nil, want monitor at index %d", tt.x, tt.y, tt.wantIndex)
					return
				}
				if got != &monitors[tt.wantIndex] {
					t.Errorf("FindMonitorFromLogicalPoint(%d, %d) got monitor at wrong index", tt.x, tt.y)
				}
			}
		})
	}
}

func TestFindMonitorFromPhysicalPoint(t *testing.T) {
	monitors := createTestMonitors()
	tests := []struct {
		name      string
		x, y      int
		wantIndex int // -1 means no monitor should be found
	}{
		{name: "first monitor center", x: 960, y: 540, wantIndex: 0},
		{name: "second monitor center", x: 2880, y: 540, wantIndex: 1},
		{name: "third monitor center", x: 960, y: 1620, wantIndex: 2},
		{name: "first monitor edge", x: 1919, y: 1079, wantIndex: 0},
		{name: "second monitor edge", x: 1920, y: 0, wantIndex: 1},
		{name: "outside all monitors", x: 4000, y: 4000, wantIndex: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromPhysicalPoint(monitors, tt.x, tt.y)
			if tt.wantIndex == -1 {
				if got != nil {
					t.Errorf("FindMonitorFromPhysicalPoint(%d, %d) = %+v, want nil", tt.x, tt.y, got)
				}
			} else {
				if got == nil {
					t.Errorf("FindMonitorFromPhysicalPoint(%d, %d) = nil, want monitor at index %d", tt.x, tt.y, tt.wantIndex)
					return
				}
				if got != &monitors[tt.wantIndex] {
					t.Errorf("FindMonitorFromPhysicalPoint(%d, %d) got monitor at wrong index", tt.x, tt.y)
				}
			}
		})
	}
}

func TestFindMonitorFromLogicalRect(t *testing.T) {
	monitors := createTestMonitors()
	tests := []struct {
		name      string
		rect      Rect
		wantIndex int // -1 means no monitor should be found
	}{
		{
			name:      "fully inside first monitor",
			rect:      Rect{Left: 100, Top: 100, Right: 800, Bottom: 600},
			wantIndex: 0,
		},
		{
			name:      "spanning first and second monitors",
			rect:      Rect{Left: 1800, Top: 100, Right: 2100, Bottom: 600},
			wantIndex: 1, // should pick monitor with larger overlap (second monitor)
		},
		{
			name:      "fully inside second monitor",
			rect:      Rect{Left: 2000, Top: 100, Right: 3000, Bottom: 600},
			wantIndex: 1,
		},
		{
			name:      "spanning all monitors",
			rect:      Rect{Left: 0, Top: 0, Right: 3840, Bottom: 2160},
			wantIndex: 0, // should pick first monitor due to equal overlap
		},
		{
			name:      "outside all monitors",
			rect:      Rect{Left: 4000, Top: 4000, Right: 5000, Bottom: 5000},
			wantIndex: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromLogicalRect(monitors, tt.rect)
			if tt.wantIndex == -1 {
				if got != nil {
					t.Errorf("FindMonitorFromLogicalRect(%+v) = %+v, want nil", tt.rect, got)
				}
			} else {
				if got == nil {
					t.Errorf("FindMonitorFromLogicalRect(%+v) = nil, want monitor at index %d", tt.rect, tt.wantIndex)
					return
				}
				if got != &monitors[tt.wantIndex] {
					t.Errorf("FindMonitorFromLogicalRect(%+v) got monitor at wrong index", tt.rect)
				}
			}
		})
	}
}

func TestOverlapArea(t *testing.T) {
	tests := []struct {
		name string
		r1   Rect
		r2   Rect
		want int
	}{
		{
			name: "no overlap",
			r1:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			r2:   Rect{Left: 200, Top: 200, Right: 300, Bottom: 300},
			want: 0,
		},
		{
			name: "full overlap",
			r1:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			r2:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			want: 10000,
		},
		{
			name: "partial overlap",
			r1:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			r2:   Rect{Left: 50, Top: 50, Right: 150, Bottom: 150},
			want: 2500,
		},
		{
			name: "edge touch",
			r1:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			r2:   Rect{Left: 100, Top: 0, Right: 200, Bottom: 100},
			want: 0,
		},
		{
			name: "zero width first rect",
			r1:   Rect{Left: 0, Top: 0, Right: 0, Bottom: 100},
			r2:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			want: 0,
		},
		{
			name: "zero height second rect",
			r1:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			r2:   Rect{Left: 50, Top: 50, Right: 150, Bottom: 50},
			want: 0,
		},
		{
			name: "negative width",
			r1:   Rect{Left: 100, Top: 0, Right: 0, Bottom: 100},
			r2:   Rect{Left: 0, Top: 0, Right: 50, Bottom: 100},
			want: 0,
		},
		{
			name: "negative height",
			r1:   Rect{Left: 0, Top: 100, Right: 100, Bottom: 0},
			r2:   Rect{Left: 0, Top: 0, Right: 100, Bottom: 50},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOverlapArea(tt.r1, tt.r2)
			if got != tt.want {
				t.Errorf("getOverlapArea() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMonitorFromLogicalRect_EdgeCases(t *testing.T) {
	monitors := createTestMonitors()
	tests := []struct {
		name      string
		rect      Rect
		wantIndex int // -1 means no monitor should be found
	}{
		{
			name:      "zero size window",
			rect:      Rect{Left: 960, Top: 540, Right: 960, Bottom: 540},
			wantIndex: -1, // invalid rectangle should return nil
		},
		{
			name:      "single pixel window at monitor edge",
			rect:      Rect{Left: 1919, Top: 0, Right: 1920, Bottom: 1},
			wantIndex: 0,
		},
		{
			name:      "negative size window",
			rect:      Rect{Left: 1000, Top: 1000, Right: 900, Bottom: 900},
			wantIndex: -1, // invalid rectangle should return nil
		},
		{
			name:      "window exactly at monitor boundary",
			rect:      Rect{Left: 1920, Top: 0, Right: 1920, Bottom: 1080},
			wantIndex: -1, // zero-width rectangle should return nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromLogicalRect(monitors, tt.rect)
			if tt.wantIndex == -1 {
				if got != nil {
					t.Errorf("FindMonitorFromLogicalRect(%+v) = %+v, want nil", tt.rect, got)
				}
			} else {
				if got == nil {
					t.Errorf("FindMonitorFromLogicalRect(%+v) = nil, want monitor at index %d", tt.rect, tt.wantIndex)
					return
				}
				if got != &monitors[tt.wantIndex] {
					t.Errorf("FindMonitorFromLogicalRect(%+v) got monitor at wrong index", tt.rect)
				}
			}
		})
	}
}
