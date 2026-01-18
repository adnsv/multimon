package units

import "testing"

func TestDimensionResolveWidth(t *testing.T) {
	ctx := ResolveContext{
		EmHeight: 16,
		WorkArea: WorkArea{Width: 1920, Height: 1080},
	}

	tests := []struct {
		name     string
		dim      Dimension
		expected int
	}{
		// Pixel values (returned as-is)
		{"pixel 1024", Dimension{Value: 1024, Unit: Pixel}, 1024},
		{"pixel 0", Dimension{Value: 0, Unit: Pixel}, 0},
		{"pixel fractional", Dimension{Value: 100.7, Unit: Pixel}, 100},

		// Em values (multiplied by EmHeight)
		{"em 60", Dimension{Value: 60, Unit: Em}, 960},   // 60 * 16 = 960
		{"em 1", Dimension{Value: 1, Unit: Em}, 16},      // 1 * 16 = 16
		{"em 1.5", Dimension{Value: 1.5, Unit: Em}, 24},  // 1.5 * 16 = 24
		{"em 0", Dimension{Value: 0, Unit: Em}, 0},

		// Percent values (percentage of WorkArea.Width)
		{"percent 100", Dimension{Value: 100, Unit: Percent}, 1920}, // 100% of 1920
		{"percent 50", Dimension{Value: 50, Unit: Percent}, 960},    // 50% of 1920
		{"percent 80", Dimension{Value: 80, Unit: Percent}, 1536},   // 80% of 1920
		{"percent 0", Dimension{Value: 0, Unit: Percent}, 0},
		{"percent 33.3", Dimension{Value: 33.3, Unit: Percent}, 639}, // 33.3% of 1920

		// Zero dimension
		{"zero dimension", Dimension{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dim.ResolveWidth(ctx)
			if result != tt.expected {
				t.Errorf("Dimension{%v, %v}.ResolveWidth(ctx) = %d, want %d",
					tt.dim.Value, tt.dim.Unit, result, tt.expected)
			}
		})
	}
}

func TestDimensionResolveHeight(t *testing.T) {
	ctx := ResolveContext{
		EmHeight: 16,
		WorkArea: WorkArea{Width: 1920, Height: 1080},
	}

	tests := []struct {
		name     string
		dim      Dimension
		expected int
	}{
		// Pixel values (returned as-is)
		{"pixel 768", Dimension{Value: 768, Unit: Pixel}, 768},
		{"pixel 0", Dimension{Value: 0, Unit: Pixel}, 0},
		{"pixel fractional", Dimension{Value: 100.9, Unit: Pixel}, 100},

		// Em values (multiplied by EmHeight)
		{"em 40", Dimension{Value: 40, Unit: Em}, 640},   // 40 * 16 = 640
		{"em 1", Dimension{Value: 1, Unit: Em}, 16},      // 1 * 16 = 16
		{"em 2.5", Dimension{Value: 2.5, Unit: Em}, 40},  // 2.5 * 16 = 40
		{"em 0", Dimension{Value: 0, Unit: Em}, 0},

		// Percent values (percentage of WorkArea.Height)
		{"percent 100", Dimension{Value: 100, Unit: Percent}, 1080}, // 100% of 1080
		{"percent 50", Dimension{Value: 50, Unit: Percent}, 540},    // 50% of 1080
		{"percent 70", Dimension{Value: 70, Unit: Percent}, 756},    // 70% of 1080
		{"percent 0", Dimension{Value: 0, Unit: Percent}, 0},
		{"percent 33.3", Dimension{Value: 33.3, Unit: Percent}, 359}, // 33.3% of 1080

		// Zero dimension
		{"zero dimension", Dimension{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.dim.ResolveHeight(ctx)
			if result != tt.expected {
				t.Errorf("Dimension{%v, %v}.ResolveHeight(ctx) = %d, want %d",
					tt.dim.Value, tt.dim.Unit, result, tt.expected)
			}
		})
	}
}

func TestResolveSize(t *testing.T) {
	ctx := ResolveContext{
		EmHeight: 16,
		WorkArea: WorkArea{Width: 1920, Height: 1080},
	}

	tests := []struct {
		name           string
		width          Dimension
		height         Dimension
		expectedWidth  int
		expectedHeight int
	}{
		{
			"pixel dimensions",
			Dimension{Value: 1024, Unit: Pixel},
			Dimension{Value: 768, Unit: Pixel},
			1024, 768,
		},
		{
			"em dimensions",
			Dimension{Value: 60, Unit: Em},
			Dimension{Value: 40, Unit: Em},
			960, 640,
		},
		{
			"percent dimensions",
			Dimension{Value: 80, Unit: Percent},
			Dimension{Value: 70, Unit: Percent},
			1536, 756,
		},
		{
			"mixed dimensions",
			Dimension{Value: 60, Unit: Em},
			Dimension{Value: 70, Unit: Percent},
			960, 756,
		},
		{
			"zero dimensions",
			Dimension{},
			Dimension{},
			0, 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := ResolveSize(tt.width, tt.height, ctx)
			if w != tt.expectedWidth || h != tt.expectedHeight {
				t.Errorf("ResolveSize(%+v, %+v, ctx) = (%d, %d), want (%d, %d)",
					tt.width, tt.height, w, h, tt.expectedWidth, tt.expectedHeight)
			}
		})
	}
}

func TestResolveWithDifferentContexts(t *testing.T) {
	// Test that resolution correctly uses different context values
	dim := Dimension{Value: 50, Unit: Percent}

	t.Run("small work area", func(t *testing.T) {
		ctx := ResolveContext{
			EmHeight: 12,
			WorkArea: WorkArea{Width: 800, Height: 600},
		}
		w := dim.ResolveWidth(ctx)
		h := dim.ResolveHeight(ctx)
		if w != 400 || h != 300 {
			t.Errorf("50%% of 800x600 = (%d, %d), want (400, 300)", w, h)
		}
	})

	t.Run("large work area", func(t *testing.T) {
		ctx := ResolveContext{
			EmHeight: 20,
			WorkArea: WorkArea{Width: 3840, Height: 2160},
		}
		w := dim.ResolveWidth(ctx)
		h := dim.ResolveHeight(ctx)
		if w != 1920 || h != 1080 {
			t.Errorf("50%% of 3840x2160 = (%d, %d), want (1920, 1080)", w, h)
		}
	})

	t.Run("different em heights", func(t *testing.T) {
		emDim := Dimension{Value: 10, Unit: Em}

		ctx12 := ResolveContext{EmHeight: 12, WorkArea: WorkArea{Width: 1920, Height: 1080}}
		ctx20 := ResolveContext{EmHeight: 20, WorkArea: WorkArea{Width: 1920, Height: 1080}}

		w12 := emDim.ResolveWidth(ctx12)
		w20 := emDim.ResolveWidth(ctx20)

		if w12 != 120 {
			t.Errorf("10em with EmHeight=12 = %d, want 120", w12)
		}
		if w20 != 200 {
			t.Errorf("10em with EmHeight=20 = %d, want 200", w20)
		}
	})
}
