package multimon

import "testing"

func TestGetOverlapArea(t *testing.T) {
	tests := []struct {
		name string
		r1   Rect
		r2   Rect
		want int
	}{
		{
			name: "no overlap",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{200, 200, 300, 300},
			want: 0,
		},
		{
			name: "full overlap",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{0, 0, 100, 100},
			want: 10000, // 100x100
		},
		{
			name: "partial overlap",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{50, 50, 150, 150},
			want: 2500, // 50x50
		},
		{
			name: "edge touch",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{100, 100, 200, 200},
			want: 0,
		},
		{
			name: "one inside other",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{25, 25, 75, 75},
			want: 2500, // 50x50
		},
		{
			name: "overlap on one axis only",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{50, 200, 150, 300},
			want: 0,
		},
		// Edge cases
		{
			name: "invalid r1 - negative width",
			r1:   Rect{100, 0, 0, 100},
			r2:   Rect{0, 0, 100, 100},
			want: 0,
		},
		{
			name: "invalid r1 - negative height",
			r1:   Rect{0, 100, 100, 0},
			r2:   Rect{0, 0, 100, 100},
			want: 0,
		},
		{
			name: "invalid r2 - negative width",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{100, 0, 0, 100},
			want: 0,
		},
		{
			name: "invalid r2 - negative height",
			r1:   Rect{0, 0, 100, 100},
			r2:   Rect{0, 100, 100, 0},
			want: 0,
		},
		{
			name: "zero width rectangles",
			r1:   Rect{0, 0, 0, 100},
			r2:   Rect{0, 0, 0, 100},
			want: 0,
		},
		{
			name: "zero height rectangles",
			r1:   Rect{0, 0, 100, 0},
			r2:   Rect{0, 0, 100, 0},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOverlapArea(tt.r1, tt.r2)
			if got != tt.want {
				t.Errorf("getOverlapArea(%v, %v) = %v, want %v", tt.r1, tt.r2, got, tt.want)
			}
		})
	}
}

func TestGetEdgeDistance(t *testing.T) {
	bounds := Rect{0, 0, 1000, 1000}
	tests := []struct {
		name   string
		window Rect
		want   int
	}{
		{
			name:   "window inside bounds",
			window: Rect{100, 100, 200, 200},
			want:   0,
		},
		{
			name:   "window center left of bounds",
			window: Rect{-200, 100, -100, 200},
			want:   150, // distance from -150 (center) to 0
		},
		{
			name:   "window center right of bounds",
			window: Rect{1100, 100, 1200, 200},
			want:   150, // distance from 1150 (center) to 1000
		},
		{
			name:   "window center above bounds",
			window: Rect{100, -200, 200, -100},
			want:   150, // distance from -150 (center) to 0
		},
		{
			name:   "window center below bounds",
			window: Rect{100, 1100, 200, 1200},
			want:   150, // distance from 1150 (center) to 1000
		},
		{
			name:   "window center diagonal from bounds",
			window: Rect{1100, 1100, 1200, 1200},
			want:   300, // 150 horizontal + 150 vertical
		},
		// Edge cases
		{
			name:   "zero size window",
			window: Rect{500, 500, 500, 500},
			want:   0,
		},
		{
			name:   "window exactly at bounds edge",
			window: Rect{1000, 1000, 1100, 1100},
			want:   100, // center is at 1050,1050, so 50 pixels beyond both edges
		},
		{
			name:   "window with negative dimensions",
			window: Rect{200, 200, 100, 100},
			want:   0, // center calculation still works
		},
		{
			name:   "window far outside bounds",
			window: Rect{10000, 10000, 10100, 10100},
			want:   18100, // 9050 pixels beyond each edge
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEdgeDistance(tt.window, bounds)
			if got != tt.want {
				t.Errorf("getEdgeDistance(%v, %v) = %v, want %v", tt.window, bounds, got, tt.want)
			}
		})
	}
}

func TestFitRectDimension(t *testing.T) {
	tests := []struct {
		name      string
		pos, size int
		min, max  int
		wantPos   int
		wantSize  int
	}{
		{
			name:     "fits within bounds",
			pos:      100,
			size:     50,
			min:      0,
			max:      200,
			wantPos:  100,
			wantSize: 50,
		},
		{
			name:     "too large",
			pos:      100,
			size:     300,
			min:      0,
			max:      200,
			wantPos:  0,
			wantSize: 200,
		},
		{
			name:     "extends past max",
			pos:      150,
			size:     100,
			min:      0,
			max:      200,
			wantPos:  100,
			wantSize: 100,
		},
		{
			name:     "before min",
			pos:      -50,
			size:     100,
			min:      0,
			max:      200,
			wantPos:  0,
			wantSize: 100,
		},
		{
			name:     "zero size",
			pos:      100,
			size:     0,
			min:      0,
			max:      200,
			wantPos:  100,
			wantSize: 0,
		},
		// Edge cases
		{
			name:     "negative size",
			pos:      100,
			size:     -50,
			min:      0,
			max:      200,
			wantPos:  100,
			wantSize: 0,
		},
		{
			name:     "min equals max",
			pos:      100,
			size:     50,
			min:      100,
			max:      100,
			wantPos:  100,
			wantSize: 0,
		},
		{
			name:     "min greater than max",
			pos:      100,
			size:     50,
			min:      200,
			max:      100,
			wantPos:  200,
			wantSize: 0,
		},
		{
			name:     "position at int max",
			pos:      2147483647,
			size:     100,
			min:      0,
			max:      200,
			wantPos:  100,
			wantSize: 100,
		},
		{
			name:     "position at int min",
			pos:      -2147483648,
			size:     100,
			min:      0,
			max:      200,
			wantPos:  0,
			wantSize: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPos, gotSize := fitRectDimension(tt.pos, tt.size, tt.min, tt.max)
			if gotPos != tt.wantPos || gotSize != tt.wantSize {
				t.Errorf("fitRectDimension(%v, %v, %v, %v) = (%v, %v), want (%v, %v)",
					tt.pos, tt.size, tt.min, tt.max, gotPos, gotSize, tt.wantPos, tt.wantSize)
			}
		})
	}
}
