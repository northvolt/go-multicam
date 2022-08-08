package multicam

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type Metadata struct {
	content int
	ioState uint8
	count   uint32
	qcount  uint32
}

// ParseMetadata returns the metadata for a given image frame, for the expected metadata type.
func ParseMetadata(content int, data []byte) (*Metadata, error) {
	switch content {
	case MetadataContentOneField:
		if len(data) < 1 {
			return nil, errors.New("invalid metadata length for 1 field")
		}

		return &Metadata{
			content: content,
			ioState: data[0]}, nil

	case MetadataContentTwoField:
		if len(data) < 6 {
			return nil, errors.New("invalid metadata length for 2 field")
		}

		return &Metadata{
			content: content,
			ioState: data[0],
			count:   binary.LittleEndian.Uint32(data[2:])}, nil

	case MetadataContentThreeField:
		if len(data) < 10 {
			return nil, errors.New("invalid metadata length for 3 field")
		}

		return &Metadata{
			content: content,
			ioState: data[0],
			qcount:  binary.LittleEndian.Uint32(data[2:]),
			count:   binary.LittleEndian.Uint32(data[6:])}, nil

	default:
		return nil, errors.New("invalid metadata field count")
	}
}

func (m *Metadata) Count() uint32 {
	return m.count
}

func (m *Metadata) Qcount() uint32 {
	return m.qcount
}

func (m *Metadata) IIN1() bool {
	return m.ioState&(1<<0) != 0
}

func (m *Metadata) IIN2() bool {
	return m.ioState&(1<<1) != 0
}

func (m *Metadata) IIN3() bool {
	return m.ioState&(1<<2) != 0
}

func (m *Metadata) IIN4() bool {
	return m.ioState&(1<<3) != 0
}

func (m *Metadata) DIN1() bool {
	return m.ioState&(1<<4) != 0
}

func (m *Metadata) DIN2() bool {
	return m.ioState&(1<<5) != 0
}

func (m *Metadata) String() string {
	s := ""
	switch m.content {
	case MetadataContentNone:
		return s
	case MetadataContentOneField:
		s = m.ioInfo()
	case MetadataContentTwoField:
		s = m.ioInfo()
		s += fmt.Sprintln("count:", m.Count())
	case MetadataContentThreeField:
		s = m.ioInfo()
		s += fmt.Sprintln("count:", m.Count())
		s += fmt.Sprintln("qcount:", m.Qcount())
	}

	return s
}

func (m *Metadata) ioInfo() string {
	s := fmt.Sprintf("\r\n=================================================>")
	s += fmt.Sprintf("\r\nIIN1    IIN2    IIN3    IIN4    DIN1    DIN2")
	s += fmt.Sprintf("\r\n %s      %s      %s      %s      %s      %s",
		showValue(m.IIN1()),
		showValue(m.IIN2()),
		showValue(m.IIN3()),
		showValue(m.IIN4()),
		showValue(m.DIN1()),
		showValue(m.DIN2()))
	s += fmt.Sprintf("\r\n------------------------------------------------->")
	return s
}

func showValue(b bool) string {
	if b {
		return "++"
	}
	return "--"
}
