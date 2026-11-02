//go:build windows
// +build windows

package winmanage

import (
	"time"
	"syscall"
)

var user32 = syscall.NewLazyDLL("user32.dll")
var procGetAsyncKeyState = user32.NewProc("GetAsyncKeyState")

// waitForSelect waits until Shift and the left mouse button are held together.
func waitForSelect() error {
	const (
		vkLButton = 0x01
		vkShift   = 0x10
	)
	for {
		shift, _, _ := procGetAsyncKeyState.Call(vkShift)
		left, _, _ := procGetAsyncKeyState.Call(vkLButton)
		if shift&0x8000 != 0 && left&0x8000 != 0 {
			return nil
		}
		time.Sleep(15 * time.Millisecond)
	}
}
