package multicam

import "testing"

func TestParseOneField(t *testing.T) {
	content := MetadataContentOneField
	data := []byte{16, 1, 2, 3}

	md, err := ParseMetadata(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentOneField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
}

func TestParseOneFieldFlipped(t *testing.T) {
	content := MetadataContentOneField
	data := []byte{1, 2, 3, 16}

	md, err := ParseMetadataFlipped(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentOneField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
}

func TestParseTwoField(t *testing.T) {
	content := MetadataContentTwoField
	data := []byte{16, 0x00, 72, 0x00, 0x00, 0x00, 1, 2, 3}

	md, err := ParseMetadata(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentTwoField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
	if md.Count() != 72 {
		t.Error("invalid count", md.Count(), "should be 72")
	}
}

func TestParseTwoFieldFlipped(t *testing.T) {
	content := MetadataContentTwoField
	data := []byte{3, 2, 1, 0x00, 0x00, 0x00, 72, 0x00, 16}

	md, err := ParseMetadataFlipped(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentTwoField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
	if md.Count() != 72 {
		t.Error("invalid count", md.Count(), "should be 72")
	}
}

func TestParseThreeField(t *testing.T) {
	content := MetadataContentThreeField
	data := []byte{16, 0x00, 72, 0x00, 0x00, 0x00, 42, 0x00, 0x00, 0x00, 1, 2, 3, 4}

	md, err := ParseMetadata(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentThreeField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
	if md.Count() != 42 {
		t.Error("invalid count", md.Count(), "should be 42")
	}

	if md.Qcount() != 72 {
		t.Error("invalid Qcount", md.Qcount(), "should be 72")
	}
}

func TestParseThreeFieldFlipped(t *testing.T) {
	content := MetadataContentThreeField
	data := []byte{3, 2, 1, 0x00, 0x00, 0x00, 42, 0x00, 0x00, 0x00, 72, 0x00, 16}

	md, err := ParseMetadataFlipped(content, data)
	if err != nil {
		t.Error("unable to parse metadata content")
	}
	if md.content != MetadataContentThreeField {
		t.Error("invalid metadata content")
	}
	if !md.DIN1() {
		t.Error("DIN1 should have been set")
	}
	if md.Count() != 42 {
		t.Error("invalid count", md.Count(), "should be 42")
	}

	if md.Qcount() != 72 {
		t.Error("invalid Qcount", md.Qcount(), "should be 72")
	}
}
