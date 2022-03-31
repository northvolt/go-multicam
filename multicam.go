package multicam

// #include <multicam.h>
import "C"

type Handle uint32

const (
	ConfigurationHandle Handle = C.MC_CONFIGURATION
)

type ParamID uint32

const (
	BoardCountParam    ParamID = C.MC_BoardCount
	BufferPitchParam   ParamID = C.MC_BufferPitch
	CamFileParam       ParamID = C.MC_CamFile
	ChannelStateParam  ParamID = C.MC_ChannelState
	ColorFormatParam   ParamID = C.MC_ColorFormat
	ConnectorParam     ParamID = C.MC_Connector
	DriverIndexParam   ParamID = C.MC_DriverIndex
	ErrorHandlingParam ParamID = C.MC_ErrorHandling
	ErrorLogParam      ParamID = C.MC_ErrorLog
	ImageSizeXParam    ParamID = C.MC_ImageSizeX
	ImageSizeYParam    ParamID = C.MC_ImageSizeY
	SeqLengthFrParam   ParamID = C.MC_SeqLength_Fr
	SignalEnableParam  ParamID = C.MC_SignalEnable
	SurfaceAddrParam   ParamID = C.MC_SurfaceAddr
)

type SignalID uint32

const (
	AcquisitionFailureSignal SignalID = C.MC_SIG_ACQUISITION_FAILURE
	SurfaceProcessingSignal  SignalID = C.MC_SIG_SURFACE_PROCESSING
)

type ChannelStateID uint32

const (
	ChannelStateActive ChannelStateID = C.MC_ChannelState_ACTIVE
	ChannelStateIdle   ChannelStateID = C.MC_ChannelState_IDLE
)

const (
	ColorFormatY8      = C.MC_ColorFormat_Y8
	ColorFormatY10     = C.MC_ColorFormat_Y10
	ColorFormatY10P    = C.MC_ColorFormat_Y10P
	ColorFormatY12     = C.MC_ColorFormat_Y12
	ColorFormatY14     = C.MC_ColorFormat_Y14
	ColorFormatY16     = C.MC_ColorFormat_Y16
	ColorFormatBayer8  = C.MC_ColorFormat_BAYER8
	ColorFormatBayer10 = C.MC_ColorFormat_BAYER10
	ColorFormatBayer12 = C.MC_ColorFormat_BAYER12
	ColorFormatBayer14 = C.MC_ColorFormat_BAYER14
	ColorFormatBayer16 = C.MC_ColorFormat_BAYER16
)

const IndeterminateLength = C.MC_INDETERMINATE

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
