package multicam

import "errors"

var (
	ErrCannotOpenDriver       = errors.New("cannot open driver")
	ErrCannotSetParam         = errors.New("cannot set param")
	ErrCannotGetParam         = errors.New("cannot get param")
	ErrInvalidChannel         = errors.New("invalid channel")
	ErrCannotCreateChannel    = errors.New("cannot create MultiCam channel")
	ErrCannotDeleteChannel    = errors.New("cannot delete MultiCam channel")
	ErrCannotRegisterCallback = errors.New("cannot register callback for MultiCam channel")
	ErrInvalidSurface         = errors.New("invalid surface")
	ErrCannotCreateSurface    = errors.New("cannot create MultiCam surface")
	ErrCannotDeleteSurface    = errors.New("cannot delete MultiCam surface")
	ErrInvalidBoard           = errors.New("invalid board")
	ErrCannotWaitSignal       = errors.New("cannot wait for signal")
	ErrTimeoutOpenDriver      = errors.New("could not init driver, try restarting PC")
)
