package multicam

import "errors"

var (
	ErrCannotOpenDriver = errors.New("cannot open driver")
	ErrCannotSetParam   = errors.New("cannot set param")
	ErrCannotGetParam   = errors.New("cannot get param")
)
