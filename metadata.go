package multicam

import (
	"encoding/binary"
	"errors"
)

type Metadata struct {
	fields  int
	ioState uint8
	count   uint32
	qcount  uint32
}

// ParseMetadata returns the metadata for a given image frame, for the number of expected fields.
func ParseMetadata(fields int, data []byte) (Metadata, error) {
	switch fields {
	case 1:
		if len(data) < 1 {
			return Metadata{}, errors.New("invalid metadata length")
		}

		return Metadata{
			fields:  fields,
			ioState: data[0]}, nil

	case 2:
		if len(data) < 6 {
			return Metadata{}, errors.New("invalid metadata length")
		}

		return Metadata{
			fields:  fields,
			ioState: data[0],
			count:   binary.LittleEndian.Uint32(data[2:])}, nil

	case 3:
		if len(data) < 10 {
			return Metadata{}, errors.New("invalid metadata length")
		}

		return Metadata{
			fields:  fields,
			ioState: data[0],
			qcount:  binary.LittleEndian.Uint32(data[2:]),
			count:   binary.LittleEndian.Uint32(data[6:])}, nil

	default:
		return Metadata{}, errors.New("invalid metadata field count")
	}
}

func (m Metadata) Count() uint32 {
	return m.count
}

func (m Metadata) Qcount() uint32 {
	return m.qcount
}

func (m Metadata) IIN1() bool {
	return m.ioState&(1<<0) > 0
}

func (m Metadata) IIN2() bool {
	return m.ioState&(1<<1) > 0
}

func (m Metadata) IIN3() bool {
	return m.ioState&(1<<2) > 0
}

func (m Metadata) IIN4() bool {
	return m.ioState&(1<<3) > 0
}

func (m Metadata) DIN1() bool {
	return m.ioState&(1<<4) > 0
}

func (m Metadata) DIN2() bool {
	return m.ioState&(1<<5) > 0
}
