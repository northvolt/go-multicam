package multicam

// #include <multicam.h>
// #include <stdlib.h>
import "C"
import "unsafe"

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

// Create creates a new MultiCam Surface object.
func (s *Surface) Create() error {
	if s.h != UninitializedSurface {
		return ErrInvalidSurface
	}

	var h C.uint

	status := C.McCreate(C.MC_DEFAULT_SURFACE_HANDLE, &h)
	if status != C.MC_OK {
		return ErrCannotCreateSurface
	}

	s.h = Handle(h)

	return nil
}

// Delete deletes an existing MultiCam surface object.
func (s *Surface) Delete() error {
	if s.h == UninitializedSurface {
		return ErrInvalidSurface
	}

	status := C.McDelete(C.uint(s.h))
	if status != C.MC_OK {
		return ErrCannotDeleteSurface
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
