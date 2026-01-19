//go:build linux && gtk4
// +build linux,gtk4

package platform

/*
#cgo linux pkg-config: gtk4 gtk4-x11

#include <gtk/gtk.h>
#include <gdk/gdk.h>
#include <gdk/x11/gdkx.h>

typedef struct Monitor {
    int x;
    int y;
    int width;
    int height;
    int workX;
    int workY;
    int workWidth;
    int workHeight;
    int scaleFactor;
} Monitor;

Monitor GetMonitorInfo(GdkDisplay *display, GdkMonitor *monitor) {
    Monitor result;
    GdkRectangle geometry;
    GdkRectangle workarea;

    gdk_monitor_get_geometry(monitor, &geometry);

    // GTK4: get_scale_factor returns integer scale
    result.scaleFactor = gdk_monitor_get_scale_factor(monitor);

    // Store screen coordinates
    result.x = geometry.x;
    result.y = geometry.y;
    result.width = geometry.width;
    result.height = geometry.height;

    // GTK4 removed gdk_monitor_get_workarea, but X11 backend still has it
    if (GDK_IS_X11_DISPLAY(display)) {
        gdk_x11_monitor_get_workarea(monitor, &workarea);
        result.workX = workarea.x;
        result.workY = workarea.y;
        result.workWidth = workarea.width;
        result.workHeight = workarea.height;
    } else {
        // Wayland: no standard work area API, fall back to geometry
        result.workX = geometry.x;
        result.workY = geometry.y;
        result.workWidth = geometry.width;
        result.workHeight = geometry.height;
    }

    return result;
}
*/
import "C"
import (
	"github.com/adnsv/multimon/types"
)

func init() {
	// Initialize GTK4
	C.gtk_init()
}

// GetPlatformMonitors returns monitor information
func GetPlatformMonitors() []types.Monitor {
	var monitors []types.Monitor

	// Get the default display
	display := C.gdk_display_get_default()
	if display == nil {
		return monitors
	}

	// GTK4: gdk_display_get_monitors returns a GListModel
	monitorList := C.gdk_display_get_monitors(display)
	if monitorList == nil {
		return monitors
	}

	n_monitors := int(C.g_list_model_get_n_items(monitorList))

	for i := 0; i < n_monitors; i++ {
		monitorPtr := C.g_list_model_get_item(monitorList, C.guint(i))
		if monitorPtr == nil {
			continue
		}
		monitor := (*C.GdkMonitor)(monitorPtr)

		info := C.GetMonitorInfo(display, monitor)
		scale := float64(info.scaleFactor)

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
		C.g_object_unref(C.gpointer(monitorPtr))
	}

	return monitors
}
