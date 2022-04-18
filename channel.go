package multicam

// #include <multicam.h>
// #include <stdlib.h>
/*
	extern void GoCallbackHandler(MCSIGNALINFO* SignalInfo);

	static void MCCallbackHandler(MCSIGNALINFO* SignalInfo)
	{
		GoCallbackHandler(SignalInfo);
	}

	static MCSTATUS SetCallbackHandler(MCHANDLE handle, PVOID context)
	{
		return McRegisterCallback(handle, MCCallbackHandler, context);
	}
*/
import "C"

const UninitializedChannel = 0

type SignalInfo C.MCSIGNALINFO

var channelCallbackHandlers map[int]func(*SignalInfo)

func initChannels() {
	channelCallbackHandlers = make(map[int]func(*SignalInfo))
}

type Channel struct {
	channel Handle
	handler func(*SignalInfo)
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
func (c *Channel) RegisterCallback(handler func(*SignalInfo)) error {
	c.handler = handler

	status := C.SetCallbackHandler(C.uint(c.channel), C.PVOID(nil))
	if status != C.MC_OK {
		return ErrCannotRegisterCallback
	}

	channelCallbackHandlers[int(c.channel)] = handler

	return nil
}

//export GoCallbackHandler
func GoCallbackHandler(info *SignalInfo) {
	if len(channelCallbackHandlers) == 0 {
		return
	}

	if cb, ok := channelCallbackHandlers[int(info.Instance)]; ok {
		cb(info)
		return
	}

	return
}
