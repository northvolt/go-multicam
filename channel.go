package multicam

// #include <multicam.h>
// #include <stdlib.h>
/*
	extern void GoCallbackHandler(MCCALLBACKINFO* info);

	static void MCCallbackHandler(MCCALLBACKINFO* info)
	{
		GoCallbackHandler(info);
	}

	static MCSTATUS SetCallbackHandler(MCHANDLE handle, PVOID context)
	{
		return McRegisterCallback(handle, MCCallbackHandler, context);
	}

	static MCSTATUS WaitSignal(MCHANDLE handle, MCSIGNAL signal, uint timeout, MCSIGNALINFO* info)
	{
		return McWaitSignal(handle, signal, timeout, info);
	}
*/
import "C"

const UninitializedChannel = 0

type CallbackInfo C.MCCALLBACKINFO

var channelCallbackHandlers map[int]func(*CallbackInfo)

func initChannels() {
	channelCallbackHandlers = make(map[int]func(*CallbackInfo))
}

type Channel struct {
	channel Handle
	handler func(*CallbackInfo)
}

// NewChannel creates a new Multicam Channel.
func NewChannel() *Channel {
	return &Channel{}
}

// Create creates a new MultiCam Channel object.
func (c *Channel) Create() error {
	if c.channel != UninitializedChannel {
		return ErrInvalidChannel
	}

	var ch C.uint

	status := C.McCreate(C.MC_CHANNEL, &ch)
	if status != C.MC_OK {
		return ErrCannotCreateChannel
	}

	c.channel = Handle(ch)

	return nil
}

// Delete deletes an existing MultiCam Channel object.
func (c *Channel) Delete() error {
	if c.channel == UninitializedChannel {
		return ErrInvalidChannel
	}

	status := C.McDelete(C.uint(c.channel))
	if status != C.MC_OK {
		return ErrCannotDeleteChannel
	}

	return nil
}

// SetParamStr sets a parameter string value for this channel.
func (c *Channel) SetParamStr(id ParamID, val string) error {
	return SetParamStr(c.channel, id, val)
}

// GetParamStr gets a parameter string value for this channel.
func (c *Channel) GetParamStr(id ParamID) (string, error) {
	return GetParamStr(c.channel, id)
}

// SetParamInt sets a parameter int value for this channel.
func (c *Channel) SetParamInt(id ParamID, val int) error {
	return SetParamInt(c.channel, id, val)
}

// GetParamInt gets a parameter int value for this channel.
func (c *Channel) GetParamInt(id ParamID) (int, error) {
	return GetParamInt(c.channel, id)
}

// SetParamInst sets a parameter instance value for this channel.
func (c *Channel) SetParamInst(id ParamID, val Handle) error {
	return SetParamInst(c.channel, id, val)
}

// GetParamInst gets a parameter instance value for this channel.
func (c *Channel) GetParamInst(id ParamID) (Handle, error) {
	return GetParamInst(c.channel, id)
}

// RegisterCallback allows setting a callback handler function for this channel.
// TODO(re): allow setting context data
func (c *Channel) RegisterCallback(handler func(*CallbackInfo)) error {
	c.handler = handler

	status := C.SetCallbackHandler(C.uint(c.channel), C.PVOID(nil))
	if status != C.MC_OK {
		return ErrCannotRegisterCallback
	}

	channelCallbackHandlers[int(c.channel)] = handler

	return nil
}

//export GoCallbackHandler
func GoCallbackHandler(info *CallbackInfo) {
	if len(channelCallbackHandlers) == 0 {
		return
	}

	if cb, ok := channelCallbackHandlers[int(info.Instance)]; ok {
		cb(info)
		return
	}

	return
}

// WaitSignal waits until a specific signal for this channel has occurred, or the timeout in ms is reached.
// On success, a pointer to the SignalInfo for this signal will be returned.
func (c *Channel) WaitSignal(signal ParamID, timeout int) (*SignalInfo, error) {
	var info SignalInfo

	status := C.WaitSignal(C.MCHANDLE(c.channel), C.MCSIGNAL(signal), C.uint(timeout), &(info.data))
	if status != C.MC_OK {
		return nil, ErrCannotWaitSignal
	}

	return &info, nil
}
