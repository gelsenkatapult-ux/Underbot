//go:build windows
// +build windows

package main

import (
	"fmt"
	"image/png"
	"os"

	impl "gitlab.com/256/Underbot/sys/Impl"
)

func main() {
	title := "UNDERTALE"
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	window, err := impl.FindWindowByTitle(title)
	if err != nil {
		fail(fmt.Errorf("could not find game window %q; pass its exact title as the first argument: %w", title, err))
	}
	name, err := window.Name()
	if err != nil {
		fail(err)
	}
	width, height, err := window.WxH()
	if err != nil {
		fail(err)
	}
	img, err := window.GetImage()
	if err != nil {
		fail(err)
	}
	const output = "underbot-window-test.png"
	file, err := os.Create(output)
	if err != nil {
		fail(err)
	}
	if err := png.Encode(file, &img); err != nil {
		_ = file.Close()
		fail(err)
	}
	if err := file.Close(); err != nil {
		fail(err)
	}
	fmt.Printf("Window: %s\nClient size: %dx%d\nScreenshot: %s\n", name, width, height, output)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "Windows backend check failed:", err)
	os.Exit(1)
}
