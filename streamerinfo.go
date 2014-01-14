package rootio

import (
	"bytes"
	"fmt"
)

type StreamerInfo struct {
	named     named
	checksum  uint32
	classvers uint32
	elmts     []StreamerElement
}

func (si *StreamerInfo) Class() string {
	return "TStreamerInfo"
}

func (si *StreamerInfo) Name() string {
	return si.named.Name()
}

func (si *StreamerInfo) Title() string {
	return si.named.Title()
}

func (si *StreamerInfo) UnmarshalROOT(data *bytes.Buffer) error {
	var err error

	dec := NewDecoder(data)
	spos := dec.Pos()

	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		return err
	}
	myprintf("[streamerinfo] vers=%v pos=%v bcnt=%v\n", vers, pos, bcnt)

	err = si.named.UnmarshalROOT(data)
	if err != nil {
		return err
	}
	myprintf("name='%v' title='%v'\n", si.Name(), si.Title())
	if vers <= 1 {
		err = fmt.Errorf("too old version for StreamerInfo (v=%d)", vers)
		return err
	}

	err = dec.readBin(&si.checksum)
	if err != nil {
		return err
	}

	err = dec.readBin(&si.classvers)
	if err != nil {
		return err
	}
	/*
		elmts := b.read_elements()
		si.elmts = make([]StreamerElement, 0, len(elmts))
		for _, v := range elmts {
			switch vv := v.(type) {
			case StreamerElement:
				si.elmts = append(si.elmts, vv)
			default:
				si.elmts = append(si.elmts, nil)
			}
		}
	*/
	err = dec.checkByteCount(pos, bcnt, spos, "TStreamerInfo")
	if err != nil {
		return err
	}
	panic("not implemented")
	return err
}

func (si *StreamerInfo) MarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
	return err
}

// check interfaces
var _ Object = (*StreamerInfo)(nil)
var _ ROOTUnmarshaler = (*StreamerInfo)(nil)

// EOF
