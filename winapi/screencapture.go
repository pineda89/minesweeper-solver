package winapi

import (
	"fmt"
	"image"
	"reflect"
	"unsafe"
)

func MyCaptureScreen() (*image.RGBA, error) {
	if !initialized {
		initialize()
	}

	x, y := my_rect.Dx(), my_rect.Dy()

	bt := BITMAPINFO{}
	bt.BmiHeader.BiSize = uint32(reflect.TypeOf(bt.BmiHeader).Size())
	bt.BmiHeader.BiWidth = int32(x)
	bt.BmiHeader.BiHeight = int32(-y)
	bt.BmiHeader.BiPlanes = 1
	bt.BmiHeader.BiBitCount = 32
	bt.BmiHeader.BiCompression = BI_RGB

	ptr := unsafe.Pointer(uintptr(0))

	my_m_hDC = CreateCompatibleDC(my_hDC)

	m_hBmp := CreateDIBSection(my_m_hDC, &bt, DIB_RGB_COLORS, &ptr, 0, 0)
	if m_hBmp == 0 {
		return nil, fmt.Errorf("Could not Create DIB Section err:%d.\n", GetLastError())
	}
	if m_hBmp == InvalidParameter {
		return nil, fmt.Errorf("One or more of the input parameters is invalid while calling CreateDIBSection.\n")
	}
	defer DeleteObject(HGDIOBJ(m_hBmp))

	obj := SelectObject(my_m_hDC, HGDIOBJ(m_hBmp))
	if obj == 0 {
		return nil, fmt.Errorf("error occurred and the selected object is not a region err:%d.\n", GetLastError())
	}
	if obj == 0xffffffff { //GDI_ERROR
		return nil, fmt.Errorf("GDI_ERROR while calling SelectObject err:%d.\n", GetLastError())
	}
	defer DeleteObject(obj)

	BitBlt(my_m_hDC, 0, 0, x, y, my_hDC, my_rect.Min.X, my_rect.Min.Y, SRCCOPY)

	DeleteDC(my_m_hDC)

	var slice []byte
	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	hdrp.Data = uintptr(ptr)
	hdrp.Len = x * y * 4
	hdrp.Cap = x * y * 4

	imageBytes := make([]byte, len(slice))

	for i := 0; i < len(imageBytes); i += 4 {
		imageBytes[i], imageBytes[i+2], imageBytes[i+1], imageBytes[i+3] = slice[i+2], slice[i], slice[i+1], slice[i+3]
	}

	img := &image.RGBA{Pix: imageBytes, Stride: 4 * x, Rect: image.Rect(0, 0, x, y)}
	return img, nil
}

func initialize() {
	my_hDC = GetDC(0)
	my_x = GetDeviceCaps(my_hDC, HORZRES)
	my_y = GetDeviceCaps(my_hDC, VERTRES)
	my_rect = image.Rect(0, 0, my_x, my_y)

	initialized = true
}
