# Screen Capture Examples

This repository contains examples of screen capture implementations using different methods for macOS.

## Running the Examples

### main_legacy.go

To run the legacy Go implementation:

```
go run main_legacy.go
```

### main_screencapturekit.go

To run the ScreenCaptureKit-based Go implementation:

```
go run main_screencapturekit.go
```

### demo.swift

To run the Swift demo:

```
swift demo.swift
```

## Status

All examples are now working as expected. This project demonstrates different ways of capturing screenshots on macOS:

1. `main_legacy.go`: A legacy implementation using older macOS APIs.
2. `main_screencapturekit.go`: An implementation using the newer ScreenCaptureKit framework, bridged with Objective-C.
3. `demo.swift`: A Swift implementation using ScreenCaptureKit.

Originally, this project was created to investigate issues with the Objective-C bridge in `main_screencapturekit.go`. These issues have now been resolved (thanks to this comment https://github.com/kbinani/screenshot/pull/23#issuecomment-2276493808), and all implementations are functioning correctly.
