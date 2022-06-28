package winapi

import (
	"image"
	"syscall"
	"time"
	"unsafe"
)

type (
	HANDLE  uintptr
	HWND    HANDLE
	HGDIOBJ HANDLE
	HDC     HANDLE
	HBITMAP HANDLE
)

const (
	HORZRES          = 8
	VERTRES          = 10
	BI_RGB           = 0
	InvalidParameter = 2
	DIB_RGB_COLORS   = 0
	SRCCOPY          = 0x00CC0020
)

var (
	modgdi32               = syscall.NewLazyDLL("gdi32.dll")
	moduser32              = syscall.NewLazyDLL("user32.dll")
	modkernel32            = syscall.NewLazyDLL("kernel32.dll")
	procGetDC              = moduser32.NewProc("GetDC")
	procReleaseDC          = moduser32.NewProc("ReleaseDC")
	procDeleteDC           = modgdi32.NewProc("DeleteDC")
	procBitBlt             = modgdi32.NewProc("BitBlt")
	procDeleteObject       = modgdi32.NewProc("DeleteObject")
	procSelectObject       = modgdi32.NewProc("SelectObject")
	procCreateDIBSection   = modgdi32.NewProc("CreateDIBSection")
	procCreateCompatibleDC = modgdi32.NewProc("CreateCompatibleDC")
	procGetDeviceCaps      = modgdi32.NewProc("GetDeviceCaps")
	procGetLastError       = modkernel32.NewProc("GetLastError")
	sendInput              = moduser32.NewProc("SendInput")
	setCursorPos           = moduser32.NewProc("SetCursorPos")

	my_hDC      HDC
	my_m_hDC    HDC
	initialized bool
	my_x        int
	my_y        int
	my_rect     image.Rectangle
)

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

func GetDC(hwnd HWND) HDC {
	ret, _, _ := procGetDC.Call(
		uintptr(hwnd))

	return HDC(ret)
}

func GetDeviceCaps(hdc HDC, index int) int {
	ret, _, _ := procGetDeviceCaps.Call(
		uintptr(hdc),
		uintptr(index))

	return int(ret)
}

func CreateCompatibleDC(hdc HDC) HDC {
	ret, _, _ := procCreateCompatibleDC.Call(
		uintptr(hdc))

	if ret == 0 {
		panic("Create compatible DC failed")
	}

	return HDC(ret)
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func DeleteObject(hObject HGDIOBJ) bool {
	ret, _, _ := procDeleteObject.Call(
		uintptr(hObject))

	return ret != 0
}

func CreateDIBSection(hdc HDC, pbmi *BITMAPINFO, iUsage uint, ppvBits *unsafe.Pointer, hSection HANDLE, dwOffset uint) HBITMAP {
	ret, _, _ := procCreateDIBSection.Call(
		uintptr(hdc),
		uintptr(unsafe.Pointer(pbmi)),
		uintptr(iUsage),
		uintptr(unsafe.Pointer(ppvBits)),
		uintptr(hSection),
		uintptr(dwOffset))

	return HBITMAP(ret)
}

func BitBlt(hdcDest HDC, nXDest, nYDest, nWidth, nHeight int, hdcSrc HDC, nXSrc, nYSrc int, dwRop uint) {
	ret, _, _ := procBitBlt.Call(
		uintptr(hdcDest),
		uintptr(nXDest),
		uintptr(nYDest),
		uintptr(nWidth),
		uintptr(nHeight),
		uintptr(hdcSrc),
		uintptr(nXSrc),
		uintptr(nYSrc),
		uintptr(dwRop))

	if ret == 0 {
		panic("BitBlt failed")
	}
}

func DeleteDC(hdc HDC) bool {
	ret, _, _ := procDeleteDC.Call(
		uintptr(hdc))

	return ret != 0
}

func SelectObject(hdc HDC, hgdiobj HGDIOBJ) HGDIOBJ {
	ret, _, _ := procSelectObject.Call(
		uintptr(hdc),
		uintptr(hgdiobj))

	if ret == 0 {
		panic("SelectObject failed")
	}

	return HGDIOBJ(ret)
}

type INPUT struct {
	Type  uint32
	mouse MOUSEINPUT
}

type MOUSEINPUT struct {
	Dx        int32
	Dy        int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	ExtraInfo uintptr
}

func SendInput(inputs ...INPUT) uint32 {
	if len(inputs) == 0 {
		return 0
	}
	ret, _, _ := sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	return uint32(ret)
}

func SetCursorPos(x, y int) bool {
	ret, _, _ := setCursorPos.Call(
		uintptr(x),
		uintptr(y),
	)
	return ret != 0
}

func MoveMouse(x, y int) {
	SetCursorPos(x, y)
}

func ClickMouse(t time.Duration, rightClick bool) {
	var flags = []uint32{0x0002, 0x0004}
	if rightClick {
		flags = []uint32{0x0008, 0x0010}
	}

	SendInput(INPUT{
		mouse: MOUSEINPUT{
			Flags: flags[0],
		},
	})
	time.Sleep(t)
	SendInput(INPUT{
		mouse: MOUSEINPUT{
			Flags: flags[1],
		},
	})
}
