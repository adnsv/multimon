//go:build linux
// +build linux

package platform

/*
#cgo linux pkg-config: gtk+-3.0

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <string.h>

typedef struct Monitor {
    int x;
    int y;
    int width;
    int height;
    int workX;
    int workY;
    int workWidth;
    int workHeight;
    int scaleFactor;     // from gdk_monitor_get_scale_factor
} Monitor;

Monitor GetMonitorInfo(GdkMonitor *monitor) {
    Monitor result;
    GdkRectangle geometry, workarea;

    gdk_monitor_get_geometry(monitor, &geometry);
    gdk_monitor_get_workarea(monitor, &workarea);

    // Get integer scale factor (GTK3 only supports integer scaling)
    result.scaleFactor = gdk_monitor_get_scale_factor(monitor);

    // Store screen coordinates
    result.x = geometry.x;
    result.y = geometry.y;
    result.width = geometry.width;
    result.height = geometry.height;
    result.workX = workarea.x;
    result.workY = workarea.y;
    result.workWidth = workarea.width;
    result.workHeight = workarea.height;

    return result;
}
*/
import "C"
import (
	"github.com/adnsv/multimon/types"
)

func init() {
	// Initialize GTK
	C.gtk_init_check(nil, nil)
}

// GetPlatformMonitors returns monitor information
func GetPlatformMonitors() []types.Monitor {
	var monitors []types.Monitor

	// Get the default display
	display := C.gdk_display_get_default()
	if display == nil {
		return monitors
	}

	// Get number of monitors
	n_monitors := int(C.gdk_display_get_n_monitors(display))

	// Iterate through monitors
	for i := 0; i < n_monitors; i++ {
		monitor := C.gdk_display_get_monitor(display, C.int(i))
		if monitor == nil {
			continue
		}

		info := C.GetMonitorInfo(monitor)
		scale := float64(info.scaleFactor)

		// Create monitor with screen coordinates and scale factor
		m := types.Monitor{
			Bounds: types.Rect{
				Left:   int(info.x),
				Top:    int(info.y),
				Right:  int(info.x + info.width),
				Bottom: int(info.y + info.height),
			},
			WorkArea: types.Rect{
				Left:   int(info.workX),
				Top:    int(info.workY),
				Right:  int(info.workX + info.workWidth),
				Bottom: int(info.workY + info.workHeight),
			},
			Scale: scale,
		}

		monitors = append(monitors, m)
	}

	return monitors
}
