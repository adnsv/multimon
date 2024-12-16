# multimon

A Go package for handling window positioning and sizing across multiple
monitors. Provides utilities for fitting windows to monitors while respecting
work areas (taskbars/docks) and handling various edge cases.

[![Go
Reference](https://pkg.go.dev/badge/github.com/adnsv/multimon.svg)](https://pkg.go.dev/github.com/adnsv/multimon)
[![Tests](https://github.com/adnsv/multimon/actions/workflows/test.yml/badge.svg)](https://github.com/adnsv/multimon/actions/workflows/test.yml)

## Overview

Window positioning and scaling is a common challenge in desktop GUI
applications, particularly when handling multiple monitors. Applications need to
properly restore window positions between sessions while gracefully handling
changes in monitor configurations. This includes cases where displays are added,
removed, or rearranged. The package helps ensure windows remain accessible and
properly positioned, avoiding issues like windows appearing outside viewable
areas or spanning multiple displays inappropriately.

## Features

- Monitor-aware window positioning and sizing
- Support for work areas (excluding taskbars/docks)
- Intelligent window fitting based on:
  - Overlap area with monitors
  - Edge distance when no overlap exists
  - Minimum size requirements
- Initial window placement with:
  - Margin support with minimum size guarantees
  - Automatic centering in work area
- Cross-platform support (Windows, Linux, macOS)

## Installation

```bash
go get github.com/adnsv/multimon
```

This package has no external Go dependencies - it only uses the standard library
and CGO bindings to system libraries.

## Terminology

Different operating systems handle screen coordinates and display scaling in
their own unique ways. Since there isn't standard terminology across platforms,
here are the terms we use in this package:

- **Physical Pixels**: The smallest addressable unit on the display.

- **Screen Units**: Units that are used by the display manager to position
  windows on the screen. All window coordinates and monitor bounds in this
  package use screen units.

- **Logical Units**: Units that provide resolution-independent way of describing
  window size and position. Used only when explicitly converting to/from screen
  units.

On Windows and Linux, screen units and physical pixels are the same. Window
positioning, monitor boundaries and mouse cursor movements are done in physical
pixel coordinates. Monitors may have display scale factors that provide a
mapping between screen units and logical units.

MacOS is different. In MacOS terminology, our screen units correspond to screen
points. On regular resolution displays, a screen unit is the same as physical
pixel. On Retina displays, a screen unit is 2x2 physical pixels.

Understanding "Effective" DPI terminology:

- The effective DPI set by display managers is a logical construct to ensure
  consistent UI scaling and does not directly correspond to the monitor's
  physical DPI
- On Windows and linux:
  - Monitors with 100% scaling factor have 96 effective DPI resolution.
  - Monitors with 200% scaling factor have 192 effective DPI resolution.
- On macOS:
  - It is assumed that non-retina displays have 72 effective DPI resolution.
  - Retina displays have 144 effective DPI resolution.

## Platform Support

- **Windows**: Native support via Win32 API (pure Go)
- **macOS**: Support via Cocoa/AppKit (requires cgo)
- **Linux**: Support via GTK3/GDK (requires cgo, gtk3-dev package)

Each platform implementation provides:
- Monitor enumeration
- Physical and logical monitor bounds
- Work area detection (accounting for taskbars/docks)

### Dependencies

For non-Windows platforms, this package requires CGO and the appropriate
development packages:
- **Linux**: `gtk3-dev` (or `libgtk-3-dev` on Debian/Ubuntu)
- **macOS**: Xcode Command Line Tools (provides Foundation, Cocoa, and AppKit
  frameworks)