package units

import "testing"

func TestParseDimension(t *testing.T) {
	tests := []struct {
		input    string
		expected Dimension
	}{
		// Pixel values (plain number)
		{"1024", Dimension{Value: 1024, Unit: Pixel}},
		{"768", Dimension{Value: 768, Unit: Pixel}},
		{"0", Dimension{Value: 0, Unit: Pixel}},
		{"100.5", Dimension{Value: 100.5, Unit: Pixel}},

		// Pixel values (with px suffix)
		{"1024px", Dimension{Value: 1024, Unit: Pixel}},
		{"768px", Dimension{Value: 768, Unit: Pixel}},
		{"0px", Dimension{Value: 0, Unit: Pixel}},
		{"100.5px", Dimension{Value: 100.5, Unit: Pixel}},

		// Em values
		{"60em", Dimension{Value: 60, Unit: Em}},
		{"40em", Dimension{Value: 40, Unit: Em}},
		{"1.5em", Dimension{Value: 1.5, Unit: Em}},
		{"0em", Dimension{Value: 0, Unit: Em}},

		// Percent values
		{"80%", Dimension{Value: 80, Unit: Percent}},
		{"50%", Dimension{Value: 50, Unit: Percent}},
		{"100%", Dimension{Value: 100, Unit: Percent}},
		{"33.3%", Dimension{Value: 33.3, Unit: Percent}},

		// Whitespace handling
		{"  1024  ", Dimension{Value: 1024, Unit: Pixel}},
		{"  1024px  ", Dimension{Value: 1024, Unit: Pixel}},
		{"  60em  ", Dimension{Value: 60, Unit: Em}},
		{"  80%  ", Dimension{Value: 80, Unit: Percent}},

		// Empty and invalid
		{"", Dimension{}},
		{"invalid", Dimension{}},
		{"em", Dimension{}},
		{"%", Dimension{}},
		{"px", Dimension{}},
		{"abc123", Dimension{}},
		{"12abc", Dimension{}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := ParseDimension(tt.input)
			if result != tt.expected {
				t.Errorf("ParseDimension(%q) = %+v, want %+v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseDimensionWithDefault(t *testing.T) {
	defaultVal := Pixels(800)

	tests := []struct {
		name     string
		input    string
		expected Dimension
	}{
		{"valid pixel", "1024", Dimension{Value: 1024, Unit: Pixel}},
		{"valid pixel with px", "1024px", Dimension{Value: 1024, Unit: Pixel}},
		{"valid em", "60em", Dimension{Value: 60, Unit: Em}},
		{"valid percent", "80%", Dimension{Value: 80, Unit: Percent}},
		{"empty uses default", "", defaultVal},
		{"invalid uses default", "invalid", defaultVal},
		{"zero returns default", "0", defaultVal},
		{"0px returns default", "0px", defaultVal},
		{"0em returns default", "0em", defaultVal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDimensionWithDefault(tt.input, defaultVal)
			if result != tt.expected {
				t.Errorf("ParseDimensionWithDefault(%q, %+v) = %+v, want %+v",
					tt.input, defaultVal, result, tt.expected)
			}
		})
	}
}

func TestDimensionString(t *testing.T) {
	tests := []struct {
		dim      Dimension
		expected string
	}{
		// Pixels
		{Dimension{Value: 1024, Unit: Pixel}, "1024"},
		{Dimension{Value: 0, Unit: Pixel}, "0"},

		// Em - integer values
		{Dimension{Value: 60, Unit: Em}, "60em"},
		{Dimension{Value: 1, Unit: Em}, "1em"},

		// Em - fractional values
		{Dimension{Value: 1.5, Unit: Em}, "1.5em"},
		{Dimension{Value: 2.25, Unit: Em}, "2.2em"},

		// Percent - integer values
		{Dimension{Value: 80, Unit: Percent}, "80%"},
		{Dimension{Value: 100, Unit: Percent}, "100%"},

		// Percent - fractional values
		{Dimension{Value: 33.3, Unit: Percent}, "33%"},
		{Dimension{Value: 50.5, Unit: Percent}, "50%"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.dim.String()
			if result != tt.expected {
				t.Errorf("Dimension{%v, %v}.String() = %q, want %q",
					tt.dim.Value, tt.dim.Unit, result, tt.expected)
			}
		})
	}
}

func TestDimensionIsZero(t *testing.T) {
	tests := []struct {
		dim      Dimension
		expected bool
	}{
		{Dimension{Value: 0, Unit: Pixel}, true},
		{Dimension{Value: 0, Unit: Em}, true},
		{Dimension{Value: 0, Unit: Percent}, true},
		{Dimension{}, true},
		{Dimension{Value: 1, Unit: Pixel}, false},
		{Dimension{Value: 0.1, Unit: Em}, false},
		{Dimension{Value: -1, Unit: Percent}, false},
	}

	for _, tt := range tests {
		t.Run(tt.dim.String(), func(t *testing.T) {
			result := tt.dim.IsZero()
			if result != tt.expected {
				t.Errorf("Dimension{%v, %v}.IsZero() = %v, want %v",
					tt.dim.Value, tt.dim.Unit, result, tt.expected)
			}
		})
	}
}

func TestDimensionConstructors(t *testing.T) {
	t.Run("Pixels", func(t *testing.T) {
		d := Pixels(1024)
		if d.Value != 1024 || d.Unit != Pixel {
			t.Errorf("Pixels(1024) = %+v, want {1024, Pixel}", d)
		}
	})

	t.Run("Ems", func(t *testing.T) {
		d := Ems(60.5)
		if d.Value != 60.5 || d.Unit != Em {
			t.Errorf("Ems(60.5) = %+v, want {60.5, Em}", d)
		}
	})

	t.Run("Pct", func(t *testing.T) {
		d := Pct(80)
		if d.Value != 80 || d.Unit != Percent {
			t.Errorf("Pct(80) = %+v, want {80, Percent}", d)
		}
	})
}
