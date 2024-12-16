//go:build windows
// +build windows

package platform

import (
	"syscall"
	"unsafe"

	"github.com/adnsv/multimon/types"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	shcore = syscall.NewLazyDLL("shcore.dll")

	procEnumDisplayMonitors    = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfo         = user32.NewProc("GetMonitorInfoW")
	procGetDpiForMonitor       = shcore.NewProc("GetDpiForMonitor")
	procGetDC                  = user32.NewProc("GetDC")
	procGetDeviceCaps          = user32.NewProc("GetDeviceCaps")
	procReleaseDC              = user32.NewProc("ReleaseDC")
	procSetProcessDPIAware     = user32.NewProc("SetProcessDPIAware")
	procSetProcessDpiAwareness = shcore.NewProc("SetProcessDpiAwareness")
)

func init() {
	// Try Windows 8.1+ API first
	if ret, _, _ := procSetProcessDpiAwareness.Call(2); ret != 0 { // PROCESS_PER_MONITOR_DPI_AWARE
		// Fall back to Windows Vista/7 API
		procSetProcessDPIAware.Call()
	}
}

type (
	HANDLE   uintptr
	HWND     HANDLE
	HDC      HANDLE
	HMONITOR HANDLE
)

type RECT struct {
	Left, Top, Right, Bottom int32
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

const (
	MONITORINFOF_PRIMARY = 0x1
	MDT_EFFECTIVE_DPI    = 0
	defaultWindowsDPI    = 96
)

func GetPlatformMonitors() []types.Monitor {
	var monitors []types.Monitor
	callback := func(hMonitor HMONITOR, hdcMonitor HDC, lprcMonitor *RECT, dwData uintptr) uintptr {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		ret, _, _ := procGetMonitorInfo.Call(
			uintptr(hMonitor),
			uintptr(unsafe.Pointer(&mi)),
		)

		if ret == 0 {
			return 1
		}

		var dpiX, dpiY uint32
		ret, _, _ = procGetDpiForMonitor.Call(
			uintptr(hMonitor),
			MDT_EFFECTIVE_DPI,
			uintptr(unsafe.Pointer(&dpiX)),
			uintptr(unsafe.Pointer(&dpiY)),
		)

		// If DPI call fails, try to get it from DC
		if ret != 0 && dpiX == 0 {
			dc, _, _ := procGetDC.Call(0)
			if dc != 0 {
				dpi, _, _ := procGetDeviceCaps.Call(dc, 88) // LOGPIXELSX = 88
				if dpi != 0 {
					dpiX = uint32(dpi)
					dpiY = uint32(dpi)
				}
				procReleaseDC.Call(0, dc)
			}
		}

		// Default to 96 if all methods fail
		if dpiX == 0 {
			dpiX = defaultWindowsDPI
		}

		// Calculate scale factor (1.0 = 100%, 1.5 = 150%, 2.0 = 200%, etc.)
		scale := float64(dpiX) / float64(defaultWindowsDPI)

		monitor := types.Monitor{
			Bounds: types.Rect{
				Left:   int(mi.RcMonitor.Left),
				Top:    int(mi.RcMonitor.Top),
				Right:  int(mi.RcMonitor.Right),
				Bottom: int(mi.RcMonitor.Bottom),
			},
			WorkArea: types.Rect{
				Left:   int(mi.RcWork.Left),
				Top:    int(mi.RcWork.Top),
				Right:  int(mi.RcWork.Right),
				Bottom: int(mi.RcWork.Bottom),
			},
			Scale: scale,
		}

		monitors = append(monitors, monitor)
		return 1
	}

	procEnumDisplayMonitors.Call(
		0,
		0,
		syscall.NewCallback(callback),
		0,
	)

	return monitors
}
