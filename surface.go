package multicam

// #include <multicam.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

const UninitializedSurface = 0

type Surface struct {
	h Handle
}

// NewSurface creates a new Multicam Surface.
func NewSurface() *Surface {
	return &Surface{}
}

// SurfaceForHandle returns a Multicam Surface for an existing Handle.
func SurfaceForHandle(h Handle) *Surface {
	return &Surface{h: h}
}

// Handle returns the Handle for this MultiCam surface object.
func (s *Surface) Handle() Handle {
	return s.h
}

// Create creates a new MultiCam Surface object.
func (s *Surface) Create() error {
	if s.h != UninitializedSurface {
		return ErrInvalidSurface
	}

	var h C.uint

	status := StatusCode(C.McCreate(C.MC_DEFAULT_SURFACE_HANDLE, &h))
	if status != StatusOK {
		return fmt.Errorf("%s: %w", status.String(), ErrCannotCreateSurface)
	}

	s.h = Handle(h)

	return nil
}

// Delete deletes an existing MultiCam surface object.
func (s *Surface) Delete() error {
	if s.h == UninitializedSurface {
		return ErrInvalidSurface
	}

	status := StatusCode(C.McDelete(C.uint(s.h)))
	if status != StatusOK {
		return fmt.Errorf("%s: %w", status.String(), ErrCannotDeleteSurface)
	}

	return nil
}

// SetParamStr sets a parameter string value for this surface.
func (s *Surface) SetParamStr(id ParamID, val string) error {
	return SetParamStr(s.h, id, val)
}

// GetParamStr gets a parameter string value for this surface.
func (s *Surface) GetParamStr(id ParamID) (string, error) {
	return GetParamStr(s.h, id)
}

// SetParamInt sets a parameter int value for this surface.
func (s *Surface) SetParamInt(id ParamID, val int) error {
	return SetParamInt(s.h, id, val)
}

// GetParamInt gets a parameter int value for this surface.
func (s *Surface) GetParamInt(id ParamID) (int, error) {
	return GetParamInt(s.h, id)
}

// SetParamPtr sets a parameter pointer value for this surface.
func (s *Surface) SetParamPtr(id ParamID, val unsafe.Pointer) error {
	return SetParamPtr(s.h, id, val)
}

// GetParamPtr gets a parameter pointer value for this surface.
func (s *Surface) GetParamPtr(id ParamID) (unsafe.Pointer, error) {
	return GetParamPtr(s.h, id)
}

// Ptr returns a slice of bytes of the data in the underlying Surface.
// Pass in the expected x and y dimensions of the surface.
// Note that the memory for this slice is under the control of Multicam and so
// may go away quickly, so copy the data elsewhere if you want to persist it.
func (s *Surface) Ptr(x, y int) ([]byte, error) {
	pimg, err := s.GetParamPtr(SurfaceAddrParam)
	if err != nil {
		return nil, err
	}
	h := &reflect.SliceHeader{
		Data: uintptr(pimg),
		Len:  int(x * y),
		Cap:  int(x * y),
	}
	ptr := *(*[]byte)(unsafe.Pointer(h))
	return ptr, nil
}

// ToBytes returns a slice of bytes which is a copy of the Surface data. This slice is safe to use as it is a normal Go slice.
func (s *Surface) ToBytes(x, y int) ([]byte, error) {
	ptr, err := s.Ptr(x, y)
	if err != nil {
		return nil, err
	}

	data := make([]byte, x*y)
	copy(data, ptr)
	return data, nil
}
