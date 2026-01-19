//go:build linux

package units

/*
#cgo linux,!gtk4 pkg-config: gtk+-3.0 pango pangocairo
#cgo linux,gtk4 pkg-config: gtk4 pango pangocairo

#include <gtk/gtk.h>
#include <pango/pango.h>
#include <pango/pangocairo.h>
#include <stdlib.h>

int getSystemFontEmHeight(void) {
    // Get the system font name from GtkSettings (e.g., "Cantarell 11")
    GtkSettings *settings = gtk_settings_get_default();
    if (settings == NULL) {
        return 16;
    }

    gchar *font_name = NULL;
    g_object_get(settings, "gtk-font-name", &font_name, NULL);
    if (font_name == NULL) {
        return 16;
    }

    // Parse the font description
    PangoFontDescription *desc = pango_font_description_from_string(font_name);
    g_free(font_name);
    if (desc == NULL) {
        return 16;
    }

    // Get font map and create context
    PangoFontMap *fontmap = pango_cairo_font_map_get_default();
    if (fontmap == NULL) {
        pango_font_description_free(desc);
        return 16;
    }

    PangoContext *context = pango_font_map_create_context(fontmap);
    if (context == NULL) {
        pango_font_description_free(desc);
        return 16;
    }

    // Load the font
    PangoFont *font = pango_font_map_load_font(fontmap, context, desc);
    if (font == NULL) {
        pango_font_description_free(desc);
        g_object_unref(context);
        return 16;
    }

    // Get font metrics
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
// Uses GtkSettings to get the actual system UI font.
func GetEmHeight() int {
	h := int(C.getSystemFontEmHeight())
	if h <= 0 {
		return 16
	}
	return h
}
