//go:build darwin

package units

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

int getSystemFontEmHeight(void) {
    @autoreleasepool {
        // systemFontOfSize:0 returns the default system font at standard size
        NSFont *font = [NSFont systemFontOfSize:0];
        if (font == nil) {
            return 16;
        }
        // boundingRectForFont gives the em-square height
        NSRect rect = [font boundingRectForFont];
        int height = (int)ceil(rect.size.height);
        return height > 0 ? height : 16;
    }
}
*/
import "C"

// GetEmHeight returns the system font em-height in points.
// Uses the macOS system font (San Francisco).
func GetEmHeight() int {
	h := int(C.getSystemFontEmHeight())
	if h <= 0 {
		return 16
	}
	return h
}
