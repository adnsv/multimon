//go:build linux

package units

/*
#cgo linux pkg-config: pango pangocairo

#include <pango/pango.h>
#include <pango/pangocairo.h>
#include <stdlib.h>

int getSystemFontEmHeight(void) {
    PangoFontMap *fontmap = pango_cairo_font_map_get_default();
    if (fontmap == NULL) {
        return 16;
    }

    PangoContext *context = pango_font_map_create_context(fontmap);
    if (context == NULL) {
        return 16;
    }

    // Get default font description (Sans is the default)
    PangoFontDescription *desc = pango_font_description_from_string("Sans");
    if (desc == NULL) {
        g_object_unref(context);
        return 16;
    }

    PangoFont *font = pango_font_map_load_font(fontmap, context, desc);
    if (font == NULL) {
        pango_font_description_free(desc);
        g_object_unref(context);
        return 16;
    }

    PangoFontMetrics *metrics = pango_font_get_metrics(font, NULL);
    if (metrics == NULL) {
        g_object_unref(font);
        pango_font_description_free(desc);
        g_object_unref(context);
        return 16;
    }

    int ascent = pango_font_metrics_get_ascent(metrics);
    int descent = pango_font_metrics_get_descent(metrics);
    int height = (ascent + descent) / PANGO_SCALE;

    pango_font_metrics_unref(metrics);
    g_object_unref(font);
    pango_font_description_free(desc);
    g_object_unref(context);

    return height > 0 ? height : 16;
}
*/
import "C"

// GetEmHeight returns the system font em-height in pixels.
// Uses Pango font metrics with the default Sans font.
func GetEmHeight() int {
	h := int(C.getSystemFontEmHeight())
	if h <= 0 {
		return 16
	}
	return h
}
