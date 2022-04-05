package multicam

// #include <multicam.h>
// #include <stdlib.h>
import "C"
import "unsafe"

const maxlength = 32

// SetParamStr sets a parameter string value.
func SetParamStr(handle Handle, id ParamID, val string) error {
	cval := C.CString(val)
	defer C.free(unsafe.Pointer(cval))

	status := C.McSetParamStr(C.MCHANDLE(handle), C.MCPARAMID(id), cval)
	if status != C.MC_OK {
		return ErrCannotSetParam
	}

	return nil
}

// GetParamStr gets a parameter string value.
func GetParamStr(handle Handle, id ParamID) (string, error) {
	data := [maxlength]byte{}
	val := C.CString(string(data[:]))
	defer C.free(unsafe.Pointer(val))

	status := C.McGetParamStr(C.MCHANDLE(handle), C.MCPARAMID(id), val, maxlength)
	if status != C.MC_OK {
		return "", ErrCannotGetParam
	}

	return C.GoString(val), nil
}

// SetParamInt sets a parameter int value.
func SetParamInt(handle Handle, id ParamID, val int) error {
	status := C.McSetParamInt(C.MCHANDLE(handle), C.MCPARAMID(id), C.int(val))
	if status != C.MC_OK {
		return ErrCannotSetParam
	}

	return nil
}

// GetParamInt gets a parameter int value.
func GetParamInt(handle Handle, id ParamID) (int, error) {
	var val C.INT32

	status := C.McGetParamInt(C.MCHANDLE(handle), C.MCPARAMID(id), &val)
	if status != C.MC_OK {
		return 0, ErrCannotGetParam
	}

	return int(val), nil
}

// SetParamPtr sets a parameter Ptr value.
func SetParamPtr(handle Handle, id ParamID, val unsafe.Pointer) error {
	status := C.McSetParamPtr(C.MCHANDLE(handle), C.MCPARAMID(id), C.PVOID(val))
	if status != C.MC_OK {
		return ErrCannotSetParam
	}

	return nil
}

// GetParamPtr gets a parameter Ptr value.
func GetParamPtr(handle Handle, id ParamID) (unsafe.Pointer, error) {
	var val C.PVOID

	status := C.McGetParamPtr(C.MCHANDLE(handle), C.MCPARAMID(id), &val)
	if status != C.MC_OK {
		return nil, ErrCannotGetParam
	}

	return unsafe.Pointer(val), nil
}

// SetParamInst sets a parameter Inst (Handle) value.
func SetParamInst(handle Handle, id ParamID, val Handle) error {
	status := C.McSetParamInst(C.MCHANDLE(handle), C.MCPARAMID(id), C.MCHANDLE(val))
	if status != C.MC_OK {
		return ErrCannotSetParam
	}

	return nil
}

// GetParamInst gets a parameter Inst (Handle) value.
func GetParamInst(handle Handle, id ParamID) (Handle, error) {
	var val C.MCHANDLE

	status := C.McGetParamInst(C.MCHANDLE(handle), C.MCPARAMID(id), &val)
	if status != C.MC_OK {
		return nil, ErrCannotGetParam
	}

	return Handle(val), nil
}
