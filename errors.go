package multicam

import "errors"

var (
	ErrCannotOpenDriver    = errors.New("cannot open driver")
	ErrCannotSetParam      = errors.New("cannot set param")
	ErrCannotGetParam      = errors.New("cannot get param")
	ErrInvalidChannel      = errors.New("invalid channel")
	ErrCannotCreateChannel = errors.New("cannot create MultiCam channel")
	ErrCannotDeleteChannel = errors.New("cannot delete MultiCam channel")
)
