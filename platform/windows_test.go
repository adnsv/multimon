//go:build windows
// +build windows

package platform

import "testing"

// Go never frees a syscall.NewCallback and caps them at about 2000 per
// process, so enumeration must not create one per call.
func TestGetPlatformMonitorsRepeated(t *testing.T) {
	first := len(GetPlatformMonitors())
	for i := 0; i < 10000; i++ {
		if n := len(GetPlatformMonitors()); n != first {
			t.Fatalf("call %d: got %d monitors, want %d", i, n, first)
		}
	}
}
