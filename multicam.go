package multicam

// #include <multicam.h>
import "C"

type Handle uint32

const (
	ConfigurationHandle  Handle = C.MC_CONFIGURATION
	DefaultSurfaceHandle Handle = C.MC_DEFAULT_SURFACE_HANDLE
	DefaultBoardHandle   Handle = C.MC_BOARD
)

type ParamID uint32

const (
	AcquisitionModeParam       ParamID = C.MC_AcquisitionMode
	BoardCountParam            ParamID = C.MC_BoardCount
	BoardNameParam             ParamID = C.MC_BoardName
	BoardNameChangeParam       ParamID = C.MC_NameBoard
	BoardIdentifierParam       ParamID = C.MC_BoardIdentifier
	BoardIndexParam            ParamID = C.MC_BoardIndex
	BoardPCIPositionParam      ParamID = C.MC_PciPosition
	BoardSerialNumberParam     ParamID = C.MC_SerialNumber
	BoardTypeParam             ParamID = C.MC_BoardType
	BufferPitchParam           ParamID = C.MC_BufferPitch
	BufferSizeParam            ParamID = C.MC_BufferSize
	CamFileParam               ParamID = C.MC_CamFile
	ChannelStateParam          ParamID = C.MC_ChannelState
	ClusterParam               ParamID = C.MC_Cluster
	ColorFormatParam           ParamID = C.MC_ColorFormat
	ConnectorParam             ParamID = C.MC_Connector
	DriverIndexParam           ParamID = C.MC_DriverIndex
	ElapsedPgParam             ParamID = C.MC_Elapsed_Pg
	ErrorHandlingParam         ParamID = C.MC_ErrorHandling
	ErrorLogParam              ParamID = C.MC_ErrorLog
	ForceTrigParam             ParamID = C.MC_ForceTrig
	ImageSizeXParam            ParamID = C.MC_ImageSizeX
	ImageSizeYParam            ParamID = C.MC_ImageSizeY
	MetadataContentParam       ParamID = C.MC_MetadataContent
	MetadataGPPCInputLineParam ParamID = C.MC_MetadataGPPCInputLine
	MetadataGPPCLocationParam  ParamID = C.MC_MetadataGPPCLocation
	MetadataGPPCResetLineParam ParamID = C.MC_MetadataGPPCResetLine
	MetadataInsertionParam     ParamID = C.MC_MetadataInsertion
	MetadataLocationParam      ParamID = C.MC_MetadataLocation
	MinBufferPitchParam        ParamID = C.MC_MinBufferPitch
	OutputConfigParam          ParamID = C.MC_OutputConfig
	OutputStateParam           ParamID = C.MC_OutputState
	SerialNumberParam          ParamID = C.MC_SerialNumber
	SeqLengthFrParam           ParamID = C.MC_SeqLength_Fr
	SeqLengthPgParam           ParamID = C.MC_SeqLength_Pg
	SeqLengthLnParam           ParamID = C.MC_SeqLength_Ln
	SignalEnableParam          ParamID = C.MC_SignalEnable
	SurfaceAddrParam           ParamID = C.MC_SurfaceAddr
	SurfaceCountParam          ParamID = C.MC_SurfaceCount
	SurfaceIndexParam          ParamID = C.MC_SurfaceIndex
	SurfacePitchParam          ParamID = C.MC_SurfacePitch
	SurfaceSizeParam           ParamID = C.MC_SurfaceSize
	SurfaceStateParam          ParamID = C.MC_SurfaceState
)

const (
	AnySignal                   ParamID = C.MC_SIG_ANY
	StartAcquisitionSignal      ParamID = C.MC_SIG_START_ACQUISITION_SEQUENCE
	EndAcquisitionSignal        ParamID = C.MC_SIG_END_ACQUISITION_SEQUENCE
	AcquisitionFailureSignal    ParamID = C.MC_SIG_ACQUISITION_FAILURE
	ClusterUnavailableSignal    ParamID = C.MC_SIG_CLUSTER_UNAVAILABLE
	EndChannelActivitySignal    ParamID = C.MC_SIG_END_CHANNEL_ACTIVITY
	FrameTriggerViolationSignal ParamID = C.MC_SIG_FRAMETRIGGER_VIOLATION
	SurfaceProcessingSignal     ParamID = C.MC_SIG_SURFACE_PROCESSING
	SurfaceFilledSignal         ParamID = C.MC_SIG_SURFACE_FILLED
	StartExposureSignal         ParamID = C.MC_SIG_START_EXPOSURE
	EndExposureSignal           ParamID = C.MC_SIG_END_EXPOSURE
	UnrecoverableOverrunSignal  ParamID = C.MC_SIG_UNRECOVERABLE_OVERRUN
	ReleaseSignal               ParamID = C.MC_SIG_RELEASE
)

const (
	SignalEnableOn  int = C.MC_SignalEnable_ON
	SignalEnableOff int = C.MC_SignalEnable_OFF
)

const (
	SurfaceStateFree       int = C.MC_SurfaceState_FREE
	SurfaceStateFilling        = C.MC_SurfaceState_FILLING
	SurfaceStateFilled         = C.MC_SurfaceState_FILLED
	SurfaceStateProcessing     = C.MC_SurfaceState_PROCESSING
	SurfaceStateReserved       = C.MC_SurfaceState_RESERVED
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

const (
	AcquisitionModeSnapshot = C.MC_AcquisitionMode_SNAPSHOT
	AcquisitionModeHFR      = C.MC_AcquisitionMode_HFR
	AcquisitionModePage     = C.MC_AcquisitionMode_PAGE
	AcquisitionModeWeb      = C.MC_AcquisitionMode_WEB
	AcquisitionModeLongPage = C.MC_AcquisitionMode_LONGPAGE
	AcquisitionModeInvalid  = C.MC_AcquisitionMode_INVALID
	AcquisitionModeVideo    = C.MC_AcquisitionMode_VIDEO
)

const IndeterminateLength = C.MC_INDETERMINATE

const (
	// LED is id of built-in LED
	LED = 25
)

const (
	MetadataContentNone       = C.MC_MetadataContent_NONE
	MetadataContentOneField   = C.MC_MetadataContent_ONE_FIELD
	MetadataContentTwoField   = C.MC_MetadataContent_TWO_FIELD
	MetadataContentThreeField = C.MC_MetadataContent_THREE_FIELD

	MetadataGPPCInputLineNone          = C.MC_MetadataGPPCInputLine_NONE
	MetadataGPPCInputLineIIN1          = C.MC_MetadataGPPCInputLine_IIN1
	MetadataGPPCLocationNone           = C.MC_MetadataGPPCLocation_NONE
	MetadataGPPCLocationInsteadLVALCNT = C.MC_MetadataGPPCLocation_INSTEAD_LVALCNT
	MetadataGPPCLocationInsteadQCNT    = C.MC_MetadataGPPCLocation_INSTEAD_QCNT
	MetadataGPPCResetLineNone          = C.MC_MetadataGPPCResetLine_NONE
	MetadataGPPCResetLineIIN4          = C.MC_MetadataGPPCResetLine_IIN4

	MetadataInsertionEnable  = C.MC_MetadataInsertion_ENABLE
	MetadataInsertionDisable = C.MC_MetadataInsertion_DISABLE

	MetadataLocationLeft        = C.MC_MetadataLocation_LEFT
	MetadataLocationSparse1     = C.MC_MetadataLocation_SPARSE_1
	MetadataLocationLeftBoarder = C.MC_MetadataLocation_LEFT_BOARDER
	MetadataLocationTap10       = C.MC_MetadataLocation_TAP10
	MetadataLocationLVALRISE    = C.MC_MetadataLocation_LVALRISE
)

// OpenDriver starts up the Multicam drivers.
func OpenDriver() error {
	status := C.McOpenDriver(nil)
	if status != C.MC_OK {
		return ErrCannotOpenDriver
	}

	initChannels()

	return nil
}

// CloseDriver closes the Multicam drivers. Call before exiting.
func CloseDriver() error {
	C.McCloseDriver()
	return nil
}
