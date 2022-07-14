package multicam

// #include <multicam.h>
// #include <stdlib.h>
import "C"
import "unsafe"

const UninitializedBoard = 0

type Board struct {
	h     Handle
	index int
}

// NewBoard creates a new Multicam Board.
func NewBoard() *Board {
	return &Board{}
}

// BoardForHandle returns a Multicam Board for an existing Handle.
func BoardForHandle(h Handle) *Board {
	return &Board{h: h}
}

// BoardForIndex returns a Multicam Board for a specific DriverIndex.
func BoardForIndex(index int) *Board {
	brd := BoardForHandle(DefaultBoardHandle + Handle(index))
	brd.index = index
	return brd
}

// Delete deletes an existing MultiCam Board object. Just here for consistency, no need to actually "delete" a board.
func (b *Board) Delete() error {
	return nil
}

// CreateChannel creates a new Channel for this Board.
func (b *Board) CreateChannel() (*Channel, error) {
	ch := NewChannel()
	err := ch.Create()
	if err != nil {
		return nil, err
	}

	// Link the channel to the board
	if err := ch.SetParamInt(DriverIndexParam, b.index); err != nil {
		return nil, err
	}

	return ch, nil
}

// SetParamStr sets a parameter string value for this Board.
func (b *Board) SetParamStr(id ParamID, val string) error {
	return SetParamStr(b.h, id, val)
}

// GetParamStr gets a parameter string value for this Board.
func (b *Board) GetParamStr(id ParamID) (string, error) {
	return GetParamStr(b.h, id)
}

// SetParamInt sets a parameter int value for this Board.
func (b *Board) SetParamInt(id ParamID, val int) error {
	return SetParamInt(b.h, id, val)
}

// GetParamInt gets a parameter int value for this Board.
func (b *Board) GetParamInt(id ParamID) (int, error) {
	return GetParamInt(b.h, id)
}

// SetParamPtr sets a parameter pointer value for this Board.
func (b *Board) SetParamPtr(id ParamID, val unsafe.Pointer) error {
	return SetParamPtr(b.h, id, val)
}

// GetParamPtr gets a parameter pointer value for this Board.
func (b *Board) GetParamPtr(id ParamID) (unsafe.Pointer, error) {
	return GetParamPtr(b.h, id)
}
