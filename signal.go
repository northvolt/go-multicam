package multicam

// #include <multicam.h>
import "C"

// SignalInfo is a wrapper around the info returned to a Channel WaitSignal() function.
type SignalInfo struct {
	data C.struct__MC_CALLBACK_INFO
}

// PVOID Context;
func (s *SignalInfo) Context() {
	// TODO: implement this whenever actually needed.
}

// MCHANDLE Instance;
func (s *SignalInfo) Instance() Handle {
	return Handle(s.data.Instance)
}

// MCSIGNAL Signal;
func (s *SignalInfo) Signal() ParamID {
	return ParamID(s.data.Signal)
}

// UINT32 SignalInfo;
func (s *SignalInfo) SignalInfo() Handle {
	return Handle(s.data.SignalInfo)
}

// UINT32 SignalContext;
func (s *SignalInfo) SignalContext() uint32 {
	return uint32(s.data.SignalContext)
}
