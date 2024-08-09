//go:build cgo

package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"unsafe"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation -framework ScreenCaptureKit -framework Foundation -framework AppKit -framework CoreMedia -framework CoreVideo -framework CoreImage
#import <ScreenCaptureKit/ScreenCaptureKit.h>
#import <CoreGraphics/CoreGraphics.h>
#import <CoreFoundation/CoreFoundation.h>
#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import <CoreMedia/CoreMedia.h>
#import <CoreVideo/CoreVideo.h>
#import <CoreImage/CoreImage.h>

static void initializeIfNeeded() {
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        [NSApplication sharedApplication];
    });
}

static CGImageRef capture(CGDirectDisplayID id, CGRect diIntersectDisplayLocal, CGColorSpaceRef colorSpace) {
    dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
    __block CGImageRef result = nil;
    [SCShareableContent getShareableContentWithCompletionHandler:^(SCShareableContent* content, NSError* error) {
        @autoreleasepool {
            if (error) {
                dispatch_semaphore_signal(semaphore);
                return;
            }
            SCDisplay* target = nil;
            for (SCDisplay *display in content.displays) {
                if (display.displayID == id) {
                    target = display;
                    break;
                }
            }
            if (!target) {
                dispatch_semaphore_signal(semaphore);
                return;
            }
            SCContentFilter* filter = [[SCContentFilter alloc] initWithDisplay:target excludingWindows:@[]];
            SCStreamConfiguration* config = [[SCStreamConfiguration alloc] init];
            config.sourceRect = diIntersectDisplayLocal;
            config.width = diIntersectDisplayLocal.size.width;
            config.height = diIntersectDisplayLocal.size.height;
            [SCScreenshotManager captureImageWithFilter:filter
                                          configuration:config
                                      completionHandler:^(CGImageRef img, NSError* error) {
                if (!error) {
                    result = CGImageCreateCopyWithColorSpace(img, colorSpace);
                }
                dispatch_semaphore_signal(semaphore);
            }];
        }
    }];
    dispatch_semaphore_wait(semaphore, DISPATCH_TIME_FOREVER);
    return result;
}

void SCContentFilter_free(SCContentFilter* filter) {
    @autoreleasepool {
        [(id)filter release];
    }
}

SCStreamConfiguration* SCStreamConfiguration_init() {
    @autoreleasepool {
        initializeIfNeeded();
        SCStreamConfiguration* config = [[SCStreamConfiguration alloc] init];
        config.width = 1920;  // Default width, can be adjusted
        config.height = 1080; // Default height, can be adjusted
        config.showsCursor = NO;
        config.scalesToFit = NO;
        NSLog(@"SCStreamConfiguration initialized with width: %d, height: %d", (int)config.width, (int)config.height);
        return config;
    }
}

void SCStreamConfiguration_free(SCStreamConfiguration* config) {
    @autoreleasepool {
        [(id)config release];
    }
}

void SCStreamConfiguration_setWidth(SCStreamConfiguration* config, int width) {
    @autoreleasepool {
        [config setValue:@(width) forKey:@"width"];
        NSLog(@"SCStreamConfiguration width set to: %d", width);
    }
}

void SCStreamConfiguration_setHeight(SCStreamConfiguration* config, int height) {
    @autoreleasepool {
        [config setValue:@(height) forKey:@"height"];
        NSLog(@"SCStreamConfiguration height set to: %d", height);
    }
}

void SCStreamConfiguration_setShowsCursor(SCStreamConfiguration* config, bool showsCursor) {
    @autoreleasepool {
        config.showsCursor = showsCursor;
        NSLog(@"SCStreamConfiguration showsCursor set to: %d", showsCursor);
    }
}

void SCShareableContent_getDisplayCount(uint32_t* count) {
    @autoreleasepool {
        initializeIfNeeded();
        dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
        [SCShareableContent getShareableContentWithCompletionHandler:^(SCShareableContent * _Nullable shareableContent, NSError * _Nullable error) {
            if (error) {
                NSLog(@"Error getting shareable content: %@", error.localizedDescription);
                *count = 0;
            } else {
                *count = (uint32_t)[shareableContent.displays count];
                NSLog(@"Number of displays: %d", *count);
            }
            dispatch_semaphore_signal(semaphore);
        }];
        dispatch_semaphore_wait(semaphore, DISPATCH_TIME_FOREVER);
    }
}

typedef struct {
    int width;
    int height;
    int x;
    int y;
    uint32_t displayID;
} DisplayInfo;

void SCShareableContent_getDisplay(int index, DisplayInfo* display) {
    @autoreleasepool {
        initializeIfNeeded();
        dispatch_semaphore_t semaphore = dispatch_semaphore_create(0);
        [SCShareableContent getShareableContentWithCompletionHandler:^(SCShareableContent * _Nullable shareableContent, NSError * _Nullable error) {
            if (error) {
                NSLog(@"Error getting shareable content: %@", error.localizedDescription);
            } else if (index < [shareableContent.displays count]) {
                SCDisplay* scDisplay = shareableContent.displays[index];
                display->width = (int)scDisplay.width;
                display->height = (int)scDisplay.height;
                display->x = (int)scDisplay.frame.origin.x;
                display->y = (int)scDisplay.frame.origin.y;
                display->displayID = (uint32_t)scDisplay.displayID;
                NSLog(@"Display info - Index: %d, Width: %d, Height: %d, X: %d, Y: %d, ID: %u",
                      index, display->width, display->height, display->x, display->y, display->displayID);
            } else {
                NSLog(@"Invalid display index: %d", index);
            }
            dispatch_semaphore_signal(semaphore);
        }];
        dispatch_semaphore_wait(semaphore, DISPATCH_TIME_FOREVER);
    }
}
*/
import "C"

func CreateImage(rect image.Rectangle) (img *image.RGBA, e error) {
	img = nil
	e = errors.New("Cannot create image.RGBA")

	defer func() {
		err := recover()
		if err == nil {
			e = nil
		}
	}()
	// image.NewRGBA may panic if rect is too large.
	img = image.NewRGBA(rect)

	return img, e
}

func Capture(x, y, width, height int) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("width or height should be > 0")
	}

	rect := image.Rect(0, 0, width, height)
	img, err := CreateImage(rect)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	displayIndex := 0 // Assuming main display, adjust if needed
	var displayInfo C.DisplayInfo
	C.SCShareableContent_getDisplay(C.int(displayIndex), &displayInfo)

	colorSpace := C.CGColorSpaceCreateDeviceRGB()
	defer C.CGColorSpaceRelease(colorSpace)

	cgRect := C.CGRectMake(C.CGFloat(x), C.CGFloat(y), C.CGFloat(width), C.CGFloat(height))
	cgImage := C.capture(C.CGDirectDisplayID(displayInfo.displayID), cgRect, colorSpace)
	// if cgImage == nil {
	// 	return nil, errors.New("failed to capture screenshot")
	// }
	defer C.CGImageRelease(cgImage)

	context := C.CGBitmapContextCreate(
		unsafe.Pointer(&img.Pix[0]),
		C.size_t(width),
		C.size_t(height),
		8,
		C.size_t(img.Stride),
		colorSpace,
		C.kCGImageAlphaNoneSkipLast,
	)
	defer C.CGContextRelease(context)

	C.CGContextDrawImage(context, C.CGRectMake(0, 0, C.CGFloat(width), C.CGFloat(height)), cgImage)

	return img, nil
}

func NumActiveDisplays() int {
	var count C.uint32_t = 0
	C.SCShareableContent_getDisplayCount(&count)
	return int(count)
}

func GetDisplayBounds(displayIndex int) image.Rectangle {
	var display C.DisplayInfo
	C.SCShareableContent_getDisplay(C.int(displayIndex), &display)

	return image.Rect(
		int(display.x),
		int(display.y),
		int(display.x+display.width),
		int(display.y+display.height),
	)
}

// save *image.RGBA to filePath with PNG format.
func save(img *image.RGBA, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Capture each displays.
	n := NumActiveDisplays()
	if n <= 0 {
		panic("Active display not found")
	}

	var all image.Rectangle = image.Rect(0, 0, 0, 0)

	for i := 0; i < n; i++ {
		bounds := GetDisplayBounds(i)
		all = bounds.Union(all)

		img, err := Capture(bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy())
		if err != nil {
			panic(err)
		}
		fileName := fmt.Sprintf("%d_%dx%d.png", i, bounds.Dx(), bounds.Dy())
		save(img, fileName)

		fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)
	}

	// Capture all desktop region into an image.
	fmt.Printf("%v\n", all)
	img, err := Capture(all.Min.X, all.Min.Y, all.Dx(), all.Dy())
	if err != nil {
		panic(err)
	}
	save(img, "all.png")
}
