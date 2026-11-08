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

// GetDeviceCaps lives in gdi32.dll; the GetDpiForMonitor fallback calls it.
func TestGetDeviceCapsResolves(t *testing.T) {
	if err := procGetDeviceCaps.Find(); err != nil {
		t.Fatalf("GetDeviceCaps not found: %v", err)
	}
	dc, _, _ := procGetDC.Call(0)
	if dc == 0 {
		t.Skip("no screen DC")
	}
	defer procReleaseDC.Call(0, dc)
	if dpi, _, _ := procGetDeviceCaps.Call(dc, 88); dpi == 0 { // LOGPIXELSX
		t.Fatal("GetDeviceCaps(LOGPIXELSX) returned 0")
	}
}
