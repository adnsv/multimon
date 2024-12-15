package multimon

import (
	"fmt"
	"testing"

	"github.com/adnsv/multimon/types"
)

func checkMonitor(t *testing.T, i int, m types.Monitor) {
	if m.LogicalBounds.Right <= m.LogicalBounds.Left {
		t.Errorf("Monitor #%d: invalid logical bounds width: right(%d) <= left(%d)", i+1, m.LogicalBounds.Right, m.LogicalBounds.Left)
	}
	if m.LogicalBounds.Bottom <= m.LogicalBounds.Top {
		t.Errorf("Monitor #%d: invalid logical bounds height: bottom(%d) <= top(%d)", i+1, m.LogicalBounds.Bottom, m.LogicalBounds.Top)
	}
	if m.LogicalWorkArea.Right <= m.LogicalWorkArea.Left {
		t.Errorf("Monitor #%d: invalid logical work area width: right(%d) <= left(%d)", i+1, m.LogicalWorkArea.Right, m.LogicalWorkArea.Left)
	}
	if m.LogicalWorkArea.Bottom <= m.LogicalWorkArea.Top {
		t.Errorf("Monitor #%d: invalid logical work area height: bottom(%d) <= top(%d)", i+1, m.LogicalWorkArea.Bottom, m.LogicalWorkArea.Top)
	}
	if m.PhysicalBounds.Right <= m.PhysicalBounds.Left {
		t.Errorf("Monitor #%d: invalid physical bounds width: right(%d) <= left(%d)", i+1, m.PhysicalBounds.Right, m.PhysicalBounds.Left)
	}
	if m.PhysicalBounds.Bottom <= m.PhysicalBounds.Top {
		t.Errorf("Monitor #%d: invalid physical bounds height: bottom(%d) <= top(%d)", i+1, m.PhysicalBounds.Bottom, m.PhysicalBounds.Top)
	}
	if m.PhysicalWorkArea.Right <= m.PhysicalWorkArea.Left {
		t.Errorf("Monitor #%d: invalid physical work area width: right(%d) <= left(%d)", i+1, m.PhysicalWorkArea.Right, m.PhysicalWorkArea.Left)
	}
	if m.PhysicalWorkArea.Bottom <= m.PhysicalWorkArea.Top {
		t.Errorf("Monitor #%d: invalid physical work area height: bottom(%d) <= top(%d)", i+1, m.PhysicalWorkArea.Bottom, m.PhysicalWorkArea.Top)
	}
}

func logMonitor(t *testing.T, i int, m types.Monitor) {
	t.Logf("Monitor #%d:\n", i+1)
	t.Logf("  Logical Bounds:     (%d,%d)-(%d,%d)\n", m.LogicalBounds.Left, m.LogicalBounds.Top, m.LogicalBounds.Right, m.LogicalBounds.Bottom)
	t.Logf("  Logical WorkArea:   (%d,%d)-(%d,%d)\n", m.LogicalWorkArea.Left, m.LogicalWorkArea.Top, m.LogicalWorkArea.Right, m.LogicalWorkArea.Bottom)
	t.Logf("  Physical Bounds:    (%d,%d)-(%d,%d)\n", m.PhysicalBounds.Left, m.PhysicalBounds.Top, m.PhysicalBounds.Right, m.PhysicalBounds.Bottom)
	t.Logf("  Physical WorkArea:  (%d,%d)-(%d,%d)\n", m.PhysicalWorkArea.Left, m.PhysicalWorkArea.Top, m.PhysicalWorkArea.Right, m.PhysicalWorkArea.Bottom)
}

func TestGetMonitors(t *testing.T) {
	monitors := GetMonitors()
	if len(monitors) == 0 {
		t.Fatal("no monitors detected")
	}

	for i, m := range monitors {
		logMonitor(t, i, m)
		checkMonitor(t, i, m)
	}
}

func ExampleGetMonitors() {
	monitors := GetMonitors()
	for i, m := range monitors {
		fmt.Printf("Monitor #%d:\n", i+1)
		fmt.Printf("  Logical Bounds:     (%d,%d)-(%d,%d)\n", m.LogicalBounds.Left, m.LogicalBounds.Top, m.LogicalBounds.Right, m.LogicalBounds.Bottom)
		fmt.Printf("  Logical WorkArea:   (%d,%d)-(%d,%d)\n", m.LogicalWorkArea.Left, m.LogicalWorkArea.Top, m.LogicalWorkArea.Right, m.LogicalWorkArea.Bottom)
		fmt.Printf("  Physical Bounds:    (%d,%d)-(%d,%d)\n", m.PhysicalBounds.Left, m.PhysicalBounds.Top, m.PhysicalBounds.Right, m.PhysicalBounds.Bottom)
		fmt.Printf("  Physical WorkArea:  (%d,%d)-(%d,%d)\n", m.PhysicalWorkArea.Left, m.PhysicalWorkArea.Top, m.PhysicalWorkArea.Right, m.PhysicalWorkArea.Bottom)
	}
}
