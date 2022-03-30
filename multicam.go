package multicam

// #include <multicam.h>
import "C"

type Handle uint32

const (
	ConfigurationHandle Handle = C.MC_CONFIGURATION
)

type ParamID uint32

const (
	ErrorHandlingParam ParamID = C.MC_ErrorHandling
	ErrorLogParam      ParamID = C.MC_ErrorLog
	BoardCountParam    ParamID = C.MC_BoardCount
)

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
