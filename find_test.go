package multimon

import "testing"

func TestFindPrimaryMonitor(t *testing.T) {
	tests := []struct {
		name     string
		monitors []Monitor
		want     *int // index in monitors array, nil for no monitor
	}{
		{
			name:     "empty monitors",
			monitors: []Monitor{},
			want:     nil,
		},
		{
			name: "single monitor at origin",
			monitors: []Monitor{
				{Bounds: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080}},
			},
			want: intPtr(0),
		},
		{
			name: "single monitor not at origin",
			monitors: []Monitor{
				{Bounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080}},
			},
			want: intPtr(0), // first monitor when no monitor contains origin
		},
		{
			name: "multiple monitors, one at origin",
			monitors: []Monitor{
				{Bounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080}},
				{Bounds: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080}},
			},
			want: intPtr(1),
		},
		{
			name: "multiple monitors, none at origin",
			monitors: []Monitor{
				{Bounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080}},
				{Bounds: Rect{Left: 3840, Top: 0, Right: 5760, Bottom: 1080}},
			},
			want: intPtr(0), // first monitor when no monitor contains origin
		},
		{
			name: "monitor partially containing origin",
			monitors: []Monitor{
				{Bounds: Rect{Left: -100, Top: -100, Right: 1820, Bottom: 980}},
			},
			want: intPtr(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindPrimaryMonitor(tt.monitors)
			if tt.want == nil {
				if got != nil {
					t.Errorf("FindPrimaryMonitor() = %v, want nil", got)
				}
			} else if got != &tt.monitors[*tt.want] {
				t.Errorf("FindPrimaryMonitor() = %v, want monitor %d", got, *tt.want)
			}
		})
	}
}

func TestFindMonitorFromScreenRect(t *testing.T) {
	baseMonitors := []Monitor{
		{Bounds: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080}},
		{Bounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080}},
	}

	tests := []struct {
		name      string
		monitors  []Monitor
		rect      Rect
		defaultTo DefaultMonitorMode
		want      *int // index in monitors array, nil for no monitor
	}{
		{
			name:      "empty monitors",
			monitors:  []Monitor{},
			rect:      Rect{Left: 0, Top: 0, Right: 100, Bottom: 100},
			defaultTo: DefaultMonitorNull,
			want:      nil,
		},
		{
			name:      "fully contained in first monitor",
			monitors:  baseMonitors,
			rect:      Rect{Left: 100, Top: 100, Right: 500, Bottom: 400},
			defaultTo: DefaultMonitorNull,
			want:      intPtr(0),
		},
		{
			name:      "fully contained in second monitor",
			monitors:  baseMonitors,
			rect:      Rect{Left: 2000, Top: 100, Right: 2500, Bottom: 400},
			defaultTo: DefaultMonitorNull,
			want:      intPtr(1),
		},
		{
			name:      "spanning both monitors, more in first",
			monitors:  baseMonitors,
			rect:      Rect{Left: 1820, Top: 100, Right: 2020, Bottom: 400},
			defaultTo: DefaultMonitorNull,
			want:      intPtr(0),
		},
		{
			name:      "spanning both monitors, more in second",
			monitors:  baseMonitors,
			rect:      Rect{Left: 1820, Top: 100, Right: 2220, Bottom: 400},
			defaultTo: DefaultMonitorNull,
			want:      intPtr(1),
		},
		{
			name:      "outside all monitors, default null",
			monitors:  baseMonitors,
			rect:      Rect{Left: 4000, Top: 100, Right: 4500, Bottom: 400},
			defaultTo: DefaultMonitorNull,
			want:      nil,
		},
		{
			name:      "outside all monitors, default primary",
			monitors:  baseMonitors,
			rect:      Rect{Left: 4000, Top: 100, Right: 4500, Bottom: 400},
			defaultTo: DefaultMonitorPrimary,
			want:      intPtr(0), // primary monitor (contains 0,0)
		},
		{
			name:      "outside all monitors, default nearest",
			monitors:  baseMonitors,
			rect:      Rect{Left: 4000, Top: 100, Right: 4500, Bottom: 400},
			defaultTo: DefaultMonitorNearest,
			want:      intPtr(1), // second monitor is closer
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromScreenRect(tt.monitors, tt.rect, tt.defaultTo)
			if tt.want == nil {
				if got != nil {
					t.Errorf("FindMonitorFromScreenRect() = %v, want nil", got)
				}
			} else if got != &tt.monitors[*tt.want] {
				t.Errorf("FindMonitorFromScreenRect() = %v, want monitor %d", got, *tt.want)
			}
		})
	}
}

func TestFindMonitorFromScreenPoint(t *testing.T) {
	baseMonitors := []Monitor{
		{Bounds: Rect{Left: 0, Top: 0, Right: 1920, Bottom: 1080}},
		{Bounds: Rect{Left: 1920, Top: 0, Right: 3840, Bottom: 1080}},
	}

	tests := []struct {
		name      string
		monitors  []Monitor
		x, y      int
		defaultTo DefaultMonitorMode
		want      *int // index in monitors array, nil for no monitor
	}{
		{
			name:      "empty monitors",
			monitors:  []Monitor{},
			x:         100,
			y:         100,
			defaultTo: DefaultMonitorNull,
			want:      nil,
		},
		{
			name:      "point in first monitor",
			monitors:  baseMonitors,
			x:         960,
			y:         540,
			defaultTo: DefaultMonitorNull,
			want:      intPtr(0),
		},
		{
			name:      "point in second monitor",
			monitors:  baseMonitors,
			x:         2880,
			y:         540,
			defaultTo: DefaultMonitorNull,
			want:      intPtr(1),
		},
		{
			name:      "point on edge of first monitor",
			monitors:  baseMonitors,
			x:         1919,
			y:         1079,
			defaultTo: DefaultMonitorNull,
			want:      intPtr(0),
		},
		{
			name:      "point outside all monitors, default null",
			monitors:  baseMonitors,
			x:         4000,
			y:         540,
			defaultTo: DefaultMonitorNull,
			want:      nil,
		},
		{
			name:      "point outside all monitors, default primary",
			monitors:  baseMonitors,
			x:         4000,
			y:         540,
			defaultTo: DefaultMonitorPrimary,
			want:      intPtr(0), // primary monitor (contains 0,0)
		},
		{
			name:      "point outside all monitors, default nearest",
			monitors:  baseMonitors,
			x:         4000,
			y:         540,
			defaultTo: DefaultMonitorNearest,
			want:      intPtr(1), // second monitor is closer
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindMonitorFromScreenPoint(tt.monitors, tt.x, tt.y, tt.defaultTo)
			if tt.want == nil {
				if got != nil {
					t.Errorf("FindMonitorFromScreenPoint() = %v, want nil", got)
				}
			} else if got != &tt.monitors[*tt.want] {
				t.Errorf("FindMonitorFromScreenPoint() = %v, want monitor %d", got, *tt.want)
			}
		})
	}
}

// Helper function to create pointer to int
func intPtr(i int) *int {
	return &i
}
