package multimon

import (
	"github.com/adnsv/multimon/platform"
	"github.com/adnsv/multimon/types"
)

// Monitor represents a display monitor and its properties
type Monitor = types.Monitor

// Rect represents a rectangle with coordinates in screen space
type Rect = types.Rect

// GetMonitors returns monitor information
func GetMonitors() []Monitor {
	return platform.GetPlatformMonitors()
}
