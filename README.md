# Screen Capture Examples

This repository contains examples of screen capture implementations using different methods.

## Running the Examples

### main_legacy.go

To run the legacy Go implementation:

```
go run main_legacy.go
```

Note: This implementation is not compilable on MacOS 15 (beta).

### main_screencapturekit.go

To run the ScreenCaptureKit-based Go implementation:

```
go run main_screencapturekit.go
```

Note: While this runs without errors, it currently does not produce a screenshot. The capture method implementation is incomplete.

### demo.swift

To run the Swift demo:

```
swift demo.swift
```

## Status

All examples work as expected, except for the following:

1. `main_legacy.go` is not compatible with MacOS 15 (beta) and will not compile on this version.
2. `main_screencapturekit.go` runs without errors but does not produce a screenshot. This is due to an incomplete implementation of the capture method.

The Swift demo (`demo.swift`) works correctly and produces screenshots as intended.

