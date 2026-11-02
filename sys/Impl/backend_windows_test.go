//go:build windows
// +build windows

package impl

import (
	"testing"
	"time"
	"unsafe"
)

func TestVirtualKeyForSupportsOriginalUnderbotKeys(t *testing.T) {
	tests := map[string]uint16{
		"z": 0x5A,
		"x": 0x58,
		"up": 0x26,
		"down": 0x28,
		"left": 0x25,
		"right": 0x27,
		"enter": 0x0D,
		"Z": 0x5A,
		"LEFT": 0x25,
		"Enter": 0x0D,
	}
	for key, want := range tests {
		got, err := virtualKeyFor(key)
		if err != nil {
			t.Errorf("virtualKeyFor(%q) returned error: %v", key, err)
			continue
		}
		if got != want {
			t.Errorf("virtualKeyFor(%q) = %#x, want %#x", key, got, want)
		}
	}
}

func TestVirtualKeyForRejectsUnknownKeys(t *testing.T) {
	for _, key := range []string{"", "space", "escape", "mouse"} {
		if _, err := virtualKeyFor(key); err == nil {
			t.Errorf("virtualKeyFor(%q) unexpectedly succeeded", key)
		}
	}
}

func TestKeyboardInputHasTheNativeWin32InputSize(t *testing.T) {
	// sizeof(INPUT) is 40 bytes on 64-bit Windows and 28 bytes on 32-bit
	// Windows because INPUT is a union whose largest member is MOUSEINPUT.
	want := uintptr(28)
	if unsafe.Sizeof(uintptr(0)) == 8 {
		want = 40
	}
	if got := unsafe.Sizeof(keyboardInput{}); got != want {
		t.Fatalf("sizeof(keyboardInput) = %d, want native INPUT size %d", got, want)
	}
}

func TestBoolToUintptr(t *testing.T) {
	if boolToUintptr(true) != 1 || boolToUintptr(false) != 0 {
		t.Fatal("boolToUintptr did not return Win32 BOOL values")
	}
}

func TestDefaultKeyPressDurationMatchesOriginalBackend(t *testing.T) {
	if defaultKeyPressDuration != 40*time.Millisecond {
		t.Fatalf("default key press duration = %s, want 40ms to match the original backend", defaultKeyPressDuration)
	}
}
