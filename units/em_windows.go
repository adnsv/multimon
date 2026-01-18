//go:build windows

package units

import (
	"syscall"
	"unsafe"
)

var (
	gdi32               = syscall.NewLazyDLL("gdi32.dll")
	user32              = syscall.NewLazyDLL("user32.dll")
	procCreateFontW     = gdi32.NewProc("CreateFontW")
	procSelectObject    = gdi32.NewProc("SelectObject")
	procGetTextMetricsW = gdi32.NewProc("GetTextMetricsW")
	procDeleteObject    = gdi32.NewProc("DeleteObject")
	procGetDC           = user32.NewProc("GetDC")
	procReleaseDC       = user32.NewProc("ReleaseDC")
)

// TEXTMETRICW structure for GetTextMetricsW
type textMetricW struct {
	TmHeight           int32
	TmAscent           int32
	TmDescent          int32
	TmInternalLeading  int32
	TmExternalLeading  int32
	TmAveCharWidth     int32
	TmMaxCharWidth     int32
	TmWeight           int32
	TmOverhang         int32
	TmDigitizedAspectX int32
	TmDigitizedAspectY int32
	TmFirstChar        uint16
	TmLastChar         uint16
	TmDefaultChar      uint16
	TmBreakChar        uint16
	TmItalic           uint8
	TmUnderlined       uint8
	TmStruckOut        uint8
	TmPitchAndFamily   uint8
	TmCharSet          uint8
}

const (
	DEFAULT_CHARSET = 1
	FW_NORMAL       = 400
)

// GetEmHeight returns the system font em-height in pixels.
// Uses Segoe UI (Windows system font) metrics.
func GetEmHeight() int {
	// Get screen DC
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 16 // fallback
	}
	defer procReleaseDC.Call(0, hdc)

	// Create Segoe UI font at default size (0 = system default)
	fontName, _ := syscall.UTF16PtrFromString("Segoe UI")
	hFont, _, _ := procCreateFontW.Call(
		0,               // height (0 = default)
		0,               // width
		0,               // escapement
		0,               // orientation
		FW_NORMAL,       // weight
		0,               // italic
		0,               // underline
		0,               // strikeout
		DEFAULT_CHARSET, // charset
		0,               // output precision
		0,               // clip precision
		0,               // quality
		0,               // pitch and family
		uintptr(unsafe.Pointer(fontName)),
	)
	if hFont == 0 {
		return 16
	}
	defer procDeleteObject.Call(hFont)

	// Select font and get metrics
	oldFont, _, _ := procSelectObject.Call(hdc, hFont)
	defer procSelectObject.Call(hdc, oldFont)

	var tm textMetricW
	ret, _, _ := procGetTextMetricsW.Call(hdc, uintptr(unsafe.Pointer(&tm)))
	if ret == 0 {
		return 16
	}

	return int(tm.TmHeight)
}
