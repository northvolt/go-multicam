package multicam

// #include <multicam.h>
import "C"

// OpenDriver starts up the Multicam drivers.
func OpenDriver() error {
	status := C.McOpenDriver(nil)
	if status != C.MC_OK {
		return ErrCannotOpenDriver
	}

	return nil
}

// CloseDriver closes the Multicam drivers. Call before exiting.
func CloseDriver() error {
	C.McCloseDriver()
	return nil
}
