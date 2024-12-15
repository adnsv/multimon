//go:build darwin
// +build darwin

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework AppKit

#import <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <stdlib.h>

typedef struct Monitor {
    int x;
    int y;
    int width;
    int height;
    int workX;
    int workY;
    int workWidth;
    int workHeight;
    int isPrimary;
} Monitor;

int GetNumMonitors() {
    return [[NSScreen screens] count];
}

Monitor GetMonitorInfo(int nth) {
    Monitor result;
    NSArray<NSScreen *> *screens = [NSScreen screens];
    NSScreen* screen = [screens objectAtIndex:nth];

    // Get screen frame (bounds)
    NSRect frame = [screen frame];
    result.x = (int)frame.origin.x;
    result.y = (int)frame.origin.y;
    result.width = (int)frame.size.width;
    result.height = (int)frame.size.height;

    // Get visible frame (work area)
    NSRect visibleFrame = [screen visibleFrame];
    result.workX = (int)visibleFrame.origin.x;
    result.workY = (int)visibleFrame.origin.y;
    result.workWidth = (int)visibleFrame.size.width;
    result.workHeight = (int)visibleFrame.size.height;

    // First screen is primary
    result.isPrimary = (nth == 0);

    return result;
}
*/
import "C"
import (
	"github.com/adnsv/multimon/types"
)

// GetPlatformMonitors returns monitor information
func GetPlatformMonitors() []types.Monitor {
	var monitors []types.Monitor

	// Get number of monitors
	numMonitors := int(C.GetNumMonitors())

	// Iterate through monitors
	for i := 0; i < numMonitors; i++ {
		info := C.GetMonitorInfo(C.int(i))

		// Logical bounds are what we get directly from Cocoa
		logicalBounds := types.Rect{
			Left:   int(info.x),
			Top:    int(info.y),
			Right:  int(info.x + info.width),
			Bottom: int(info.y + info.height),
		}

		// Work area (visible frame)
		logicalWorkArea := types.Rect{
			Left:   int(info.workX),
			Top:    int(info.workY),
			Right:  int(info.workX + info.workWidth),
			Bottom: int(info.workY + info.workHeight),
		}

		// On macOS, logical and physical are the same since scaling is handled by the system
		m := types.Monitor{
			LogicalBounds:    logicalBounds,
			LogicalWorkArea:  logicalWorkArea,
			PhysicalBounds:   logicalBounds,   // Same as logical on macOS
			PhysicalWorkArea: logicalWorkArea, // Same as logical on macOS
		}

		monitors = append(monitors, m)
	}

	return monitors
}
