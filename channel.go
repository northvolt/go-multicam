package multicam

// #include <multicam.h>
// #include <stdlib.h>
import "C"

const UninitializedChannel = 0

type Channel struct {
	channel C.MCHANDLE
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

	status := C.McCreate(C.MC_CHANNEL, &c.channel)
	if status != C.MC_OK {
		return ErrCannotCreateChannel
	}

	return nil
}

// Delete deletes an existing MultiCam Channel object.
func (c *Channel) Delete() error {
	if c.channel == UninitializedChannel {
		return ErrInvalidChannel
	}

	status := C.McDelete(c.channel)
	if status != C.MC_OK {
		return ErrCannotDeleteChannel
	}

	return nil
}
