//go:build windows
// +build windows

package impl

import (
	"errors"
	"fmt"
	"image"
	"os"
	"syscall"
	"unsafe"

	"gitlab.com/256/Underbot/sys"
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")

	procGetForegroundWindow    = user32.NewProc("GetForegroundWindow")
	procEnumWindows            = user32.NewProc("EnumWindows")
	procIsWindow               = user32.NewProc("IsWindow")
	procGetWindowTextLengthW   = user32.NewProc("GetWindowTextLengthW")
	procGetWindowTextW         = user32.NewProc("GetWindowTextW")
	procGetClientRect          = user32.NewProc("GetClientRect")
	procClientToScreen         = user32.NewProc("ClientToScreen")
	procSetForegroundWindow    = user32.NewProc("SetForegroundWindow")
	procSetWindowPos           = user32.NewProc("SetWindowPos")
	procGetWindowThreadProcess = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess            = syscall.NewLazyDLL("kernel32.dll").NewProc("OpenProcess")
	procCloseHandle            = syscall.NewLazyDLL("kernel32.dll").NewProc("CloseHandle")
	procNtSuspendProcess       = syscall.NewLazyDLL("ntdll.dll").NewProc("NtSuspendProcess")
	procNtResumeProcess        = syscall.NewLazyDLL("ntdll.dll").NewProc("NtResumeProcess")
	procSendInput              = user32.NewProc("SendInput")
	procGetDC                  = user32.NewProc("GetDC")
	procReleaseDC              = user32.NewProc("ReleaseDC")
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procBitBlt                 = gdi32.NewProc("BitBlt")
	procGetDIBits              = gdi32.NewProc("GetDIBits")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
)

const (
	swRestore       = 9
	swpNoZOrder     = 0x0004
	swpNoActivate   = 0x0010
	srccopy         = 0x00CC0020
	dibRGBColors    = 0
	biRGB           = 0
	inputKeyboard   = 1
	keyeventKeyup   = 0x0002
	virtualKeyEnter = 0x0D
	processSuspendResume = 0x0800
)

type winRect struct{ Left, Top, Right, Bottom int32 }
type winPoint struct{ X, Y int32 }
type bitmapInfoHeader struct {
	Size uint32
	Width int32
	Height int32
	Planes uint16
	BitCount uint16
	Compression uint32
	SizeImage uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed uint32
	ClrImportant uint32
}
type bitmapInfo struct {
	Header bitmapInfoHeader
	Colors [1]uint32
}
type keybdInput struct {
	Vk uint16
	Scan uint16
	Flags uint32
	Time uint32
	ExtraInfo uintptr
}
type keyboardInput struct {
	Type uint32
	Ki keybdInput
}

// Server is the native Win32 implementation of sys.Server.
type Server struct{}

// NewServer initializes the native Windows backend.
func NewServer() (Server, error) { return Server{}, nil }

// ActiveWindow returns the current foreground window.
func (Server) ActiveWindow() (sys.Window, error) {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 { return nil, errors.New("no foreground window is available") }
	return window{hwnd: hwnd}, nil
}

type window struct { hwnd uintptr }

func (w window) valid() error {
	ok, _, _ := procIsWindow.Call(w.hwnd)
	if ok == 0 { return errors.New("the Windows window handle is no longer valid") }
	return nil
}

func (w window) Center() (image.Point, error) {
	x, y, err := w.originAndSize()
	if err != nil { return image.Point{}, err }
	width, height, err := w.dimensions()
	if err != nil { return image.Point{}, err }
	return image.Pt(x+width/2, y+height/2), nil
}

func (w window) originAndSize() (int, int, error) {
	var r winRect
	ok, _, callErr := procGetClientRect.Call(w.hwnd, uintptr(unsafe.Pointer(&r)))
	if ok == 0 { return 0, 0, fmt.Errorf("GetClientRect failed: %v", callErr) }
	p := winPoint{}
	ok, _, callErr = procClientToScreen.Call(w.hwnd, uintptr(unsafe.Pointer(&p)))
	if ok == 0 { return 0, 0, fmt.Errorf("ClientToScreen failed: %v", callErr) }
	return int(p.X), int(p.Y), nil
}
func (w window) dimensions() (int, int, error) {
	if err := w.valid(); err != nil { return 0, 0, err }
	var r winRect
	ok, _, callErr := procGetClientRect.Call(w.hwnd, uintptr(unsafe.Pointer(&r)))
	if ok == 0 { return 0, 0, fmt.Errorf("GetClientRect failed: %v", callErr) }
	width, height := int(r.Right-r.Left), int(r.Bottom-r.Top)
	if width <= 0 || height <= 0 { return 0, 0, errors.New("window has an empty client area") }
	return width, height, nil
}
func (w window) WxH() (int, int, error) { return w.dimensions() }

func (w window) GetImage() (image.RGBA, error) {
	width, height, err := w.dimensions()
	if err != nil { return image.RGBA{}, err }
	screenDC, _, _ := procGetDC.Call(0)
	if screenDC == 0 { return image.RGBA{}, errors.New("GetDC failed") }
	defer procReleaseDC.Call(0, screenDC)
	memDC, _, _ := procCreateCompatibleDC.Call(screenDC)
	if memDC == 0 { return image.RGBA{}, errors.New("CreateCompatibleDC failed") }
	defer procDeleteDC.Call(memDC)
	bitmap, _, _ := procCreateCompatibleBitmap.Call(screenDC, uintptr(width), uintptr(height))
	if bitmap == 0 { return image.RGBA{}, errors.New("CreateCompatibleBitmap failed") }
	defer procDeleteObject.Call(bitmap)
	old, _, _ := procSelectObject.Call(memDC, bitmap)
	x, y, err := w.originAndSize()
	if err != nil { return image.RGBA{}, err }
	ok, _, _ := procBitBlt.Call(memDC, 0, 0, uintptr(width), uintptr(height), screenDC, uintptr(x), uintptr(y), srccopy)
	if ok == 0 { return image.RGBA{}, errors.New("BitBlt failed while capturing the game window") }
	// GetDIBits requires the bitmap not to be selected into a device context.
	procSelectObject.Call(memDC, old)

	pixels := make([]byte, width*height*4)
	bmi := bitmapInfo{Header: bitmapInfoHeader{Size: uint32(unsafe.Sizeof(bitmapInfoHeader{})), Width: int32(width), Height: -int32(height), Planes: 1, BitCount: 32, Compression: biRGB}}
	lines, _, _ := procGetDIBits.Call(screenDC, bitmap, 0, uintptr(height), uintptr(unsafe.Pointer(&pixels[0])), uintptr(unsafe.Pointer(&bmi)), dibRGBColors)
	if lines == 0 { return image.RGBA{}, errors.New("GetDIBits failed while converting the window image") }
	for i := 0; i < len(pixels); i += 4 { pixels[i], pixels[i+2], pixels[i+3] = pixels[i+2], pixels[i], 255 }
	return image.RGBA{Pix: pixels, Stride: width*4, Rect: image.Rect(0, 0, width, height)}, nil
}

func (w window) Name() (string, error) {
	if err := w.valid(); err != nil { return "", err }
	n, _, _ := procGetWindowTextLengthW.Call(w.hwnd)
	buf := make([]uint16, int(n)+1)
	read, _, callErr := procGetWindowTextW.Call(w.hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if read == 0 { return "", fmt.Errorf("GetWindowTextW returned no title: %v", callErr) }
	return syscall.UTF16ToString(buf), nil
}

func (w window) Resize(width, height int) error {
	if err := w.valid(); err != nil { return err }
	if width <= 0 || height <= 0 { return errors.New("window dimensions must be positive") }
	ok, _, callErr := procSetWindowPos.Call(w.hwnd, 0, 0, 0, uintptr(width), uintptr(height), swpNoZOrder|swpNoActivate)
	if ok == 0 { return fmt.Errorf("SetWindowPos failed: %v", callErr) }
	return nil
}

func (w window) SetActive() error {
	if err := w.valid(); err != nil { return err }
	procShowWindow := user32.NewProc("ShowWindow")
	procShowWindow.Call(w.hwnd, swRestore)
	ok, _, callErr := procSetForegroundWindow.Call(w.hwnd)
	if ok == 0 { return fmt.Errorf("SetForegroundWindow failed: %v", callErr) }
	return nil
}

func (w window) Process() (*os.Process, error) {
	if err := w.valid(); err != nil { return nil, err }
	var pid uint32
	thread, _, callErr := procGetWindowThreadProcess.Call(w.hwnd, uintptr(unsafe.Pointer(&pid)))
	if thread == 0 || pid == 0 { return nil, fmt.Errorf("GetWindowThreadProcessId failed: %v", callErr) }
	return os.FindProcess(int(pid))
}

// suspendOrResume uses the process suspend/resume routines and always closes
// the process handle after the operation.
func (w window) suspendOrResume(resume bool) error {
	proc, err := w.Process()
	if err != nil { return err }
	handle, _, callErr := procOpenProcess.Call(processSuspendResume, 0, uintptr(proc.Pid))
	if handle == 0 { return fmt.Errorf("OpenProcess failed: %v", callErr) }
	defer procCloseHandle.Call(handle)
	var status uintptr
	if resume {
		status, _, _ = procNtResumeProcess.Call(handle)
	} else {
		status, _, _ = procNtSuspendProcess.Call(handle)
	}
	if int32(status) < 0 { return fmt.Errorf("Windows process suspend/resume failed with NTSTATUS 0x%X", uint32(status)) }
	return nil
}

func (w window) Pause() error { return w.suspendOrResume(false) }
func (w window) Resume() error { return w.suspendOrResume(true) }

func (w window) Press(key string) error {
	if err := w.valid(); err != nil { return err }
	var vk uint16
	switch key {
	case "z": vk = 'Z'
	case "x": vk = 'X'
	case "up": vk = 0x26
	case "down": vk = 0x28
	case "left": vk = 0x25
	case "right": vk = 0x27
	case "enter": vk = virtualKeyEnter
	default: return fmt.Errorf("unsupported key %q", key)
	}
	foreground, _, _ := procGetForegroundWindow.Call()
	if foreground != w.hwnd {
		if err := w.SetActive(); err != nil { return err }
	}
	inputs := []keyboardInput{{Type: inputKeyboard, Ki: keybdInput{Vk: vk}}, {Type: inputKeyboard, Ki: keybdInput{Vk: vk, Flags: keyeventKeyup}}}
	sent, _, callErr := procSendInput.Call(uintptr(len(inputs)), uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(inputs[0]))
	if sent != uintptr(len(inputs)) { return fmt.Errorf("SendInput sent %d of %d keyboard events: %v", sent, len(inputs), callErr) }
	return nil
}

func (w window) ID() (int, error) {
	if err := w.valid(); err != nil { return 0, err }
	return int(w.hwnd), nil
}

// FindWindowByTitle is a helper for the Windows backend and is also used by winmanage.
func FindWindowByTitle(title string) (sys.Window, error) {
	var found uintptr
	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		n, _, _ := procGetWindowTextLengthW.Call(hwnd)
		if n == 0 { return 1 }
		buf := make([]uint16, int(n)+1)
		procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if syscall.UTF16ToString(buf) == title {
			found = hwnd
			return 0
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	if found == 0 { return nil, fmt.Errorf("window with title %q was not found", title) }
	return window{hwnd: found}, nil
}
