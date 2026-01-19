//go:build windows

package units

import (
	"syscall"
	"unsafe"
)

var (
	gdi32                    = syscall.NewLazyDLL("gdi32.dll")
	user32                   = syscall.NewLazyDLL("user32.dll")
	procCreateFontIndirectW  = gdi32.NewProc("CreateFontIndirectW")
	procSelectObject         = gdi32.NewProc("SelectObject")
	procGetTextMetricsW      = gdi32.NewProc("GetTextMetricsW")
	procDeleteObject         = gdi32.NewProc("DeleteObject")
	procGetDC                = user32.NewProc("GetDC")
	procReleaseDC            = user32.NewProc("ReleaseDC")
	procSystemParametersInfo = user32.NewProc("SystemParametersInfoW")
)

const (
	SPI_GETNONCLIENTMETRICS = 0x0029
	LF_FACESIZE             = 32
)

// LOGFONTW structure
type logFontW struct {
	LfHeight         int32
	LfWidth          int32
	LfEscapement     int32
	LfOrientation    int32
	LfWeight         int32
	LfItalic         uint8
	LfUnderline      uint8
	LfStrikeOut      uint8
	LfCharSet        uint8
	LfOutPrecision   uint8
	LfClipPrecision  uint8
	LfQuality        uint8
	LfPitchAndFamily uint8
	LfFaceName       [LF_FACESIZE]uint16
}

// NONCLIENTMETRICSW structure
type nonClientMetricsW struct {
	CbSize           uint32
	IBorderWidth     int32
	IScrollWidth     int32
	IScrollHeight    int32
	ICaptionWidth    int32
	ICaptionHeight   int32
	LfCaptionFont    logFontW
	ISmCaptionWidth  int32
	ISmCaptionHeight int32
	LfSmCaptionFont  logFontW
	IMenuWidth       int32
	IMenuHeight      int32
	LfMenuFont       logFontW
	LfStatusFont     logFontW
	LfMessageFont    logFontW
	IPaddedBorderWidth int32
}

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

// GetEmHeight returns the system font em-height in pixels.
// Uses SystemParametersInfo to get the actual system UI font (lfMessageFont).
func GetEmHeight() int {
	// Get non-client metrics to retrieve the system message font
	var ncm nonClientMetricsW
	ncm.CbSize = uint32(unsafe.Sizeof(ncm))

	ret, _, _ := procSystemParametersInfo.Call(
		SPI_GETNONCLIENTMETRICS,
		uintptr(ncm.CbSize),
		uintptr(unsafe.Pointer(&ncm)),
		0,
	)
	if ret == 0 {
		return 16 // fallback
	}

	// Create font from lfMessageFont (the system UI font)
	hFont, _, _ := procCreateFontIndirectW.Call(uintptr(unsafe.Pointer(&ncm.LfMessageFont)))
	if hFont == 0 {
		return 16
	}
	defer procDeleteObject.Call(hFont)

	// Get screen DC
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return 16
	}
	defer procReleaseDC.Call(0, hdc)

	// Select font and get metrics
	oldFont, _, _ := procSelectObject.Call(hdc, hFont)
	defer procSelectObject.Call(hdc, oldFont)

	var tm textMetricW
	ret, _, _ = procGetTextMetricsW.Call(hdc, uintptr(unsafe.Pointer(&tm)))
	if ret == 0 {
		return 16
	}

	return int(tm.TmHeight)
}
