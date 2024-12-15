# go-multimon

A Go package for handling window positioning and sizing across multiple monitors. Provides utilities for fitting windows to monitors while respecting work areas (taskbars/docks) and handling various edge cases.

[![Go Reference](https://pkg.go.dev/badge/github.com/adnsv/multimon.svg)](https://pkg.go.dev/github.com/adnsv/multimon)
[![Tests](https://github.com/adnsv/multimon/actions/workflows/test.yml/badge.svg)](https://github.com/adnsv/multimon/actions/workflows/test.yml)

## Features

- Monitor-aware window positioning and sizing
- Support for work areas (excluding taskbars/docks)
- Intelligent window fitting based on:
  - Overlap area with monitors
  - Edge distance when no overlap exists
  - Minimum size requirements
- Cross-platform support (Windows, Linux, macOS)
- Handles various edge cases:
  - Windows spanning multiple monitors
  - Windows outside all monitors
  - Zero or negative size windows
  - Windows larger than monitor bounds

## Installation

```bash
go get github.com/adnsv/multimon
```

## Platform Support

- **Windows**: Native support via Win32 API (pure Go)
- **macOS**: Support via Cocoa/AppKit (requires cgo)
- **Linux**: Support via GTK3/GDK (requires cgo, gtk3-dev package)

Each platform implementation provides:
- Monitor enumeration
- Physical and logical monitor bounds
- Work area detection (accounting for taskbars/docks)

### Dependencies

For non-Windows platforms, this package requires CGO and the appropriate development packages:
- **Linux**: `gtk3-dev` (or `libgtk-3-dev` on Debian/Ubuntu)
- **macOS**: Xcode Command Line Tools (provides Foundation, Cocoa, and AppKit frameworks)
