//go:build darwin
// +build darwin

package platform

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa -framework AppKit

#import <Foundation/Foundation.h>
#include <AppKit/AppKit.h>
#include <stdlib.h>

typedef struct monitorInfo {
    int x;
    int y;
    int width;
    int height;
    int workX;
    int workY;
    int workWidth;
    int workHeight;
} monitorInfo;

int GetNumMonitors() {
    return [[NSScreen screens] count];
}

monitorInfo GetMonitorInfo(int nth) {
    monitorInfo result;
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

	// Get main screen height for Y-coordinate conversion
	mainScreen := C.GetMonitorInfo(0)
	mainHeight := int(mainScreen.height)

	// Iterate through monitors
	for i := 0; i < numMonitors; i++ {
		info := C.GetMonitorInfo(C.int(i))

		// Convert Y coordinates by subtracting from main screen height
		y := mainHeight - (int(info.y) + int(info.height))
		workY := mainHeight - (int(info.workY) + int(info.workHeight))

		// Create monitor with screen coordinates in points (macOS native units)
		m := types.Monitor{
			Bounds: types.Rect{
				Left:   int(info.x),
				Top:    y,
				Right:  int(info.x + info.width),
				Bottom: y + int(info.height),
			},
			WorkArea: types.Rect{
				Left:   int(info.workX),
				Top:    workY,
				Right:  int(info.workX + info.workWidth),
				Bottom: workY + int(info.workHeight),
			},
			Scale: 1.0, // Always 1.0 since we work with screen points
		}

		monitors = append(monitors, m)
	}

	return monitors
}
